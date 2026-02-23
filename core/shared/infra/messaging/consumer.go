package messaging

import (
	"context"
	"fmt"
	"go-socket/core/shared/pkg/logging"
	"strings"
	"sync"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"go.opentelemetry.io/otel/codes"
	"go.uber.org/zap"
)

type CallBack func(ctx context.Context, topic string, value []byte) error
type Handler func(ctx context.Context, value []byte) error

type Consumer interface {
	Read(callback CallBack)
	Stop()
	SetHandler(h Handler)
	GetHandler() Handler
	GetHandlerName() string
}

type consumer struct {
	instance *kafka.Consumer

	startSingleton sync.Once
	stopSingleton  sync.Once

	chanStop    chan bool
	handler     Handler
	handlerName string

	dlq bool
}

func NewConsumer(config *kafka.ConfigMap) Consumer {
	return &consumer{}
}

func (c *consumer) Read(f CallBack) {
	c.startSingleton.Do(func() {
		go func() {
			c.start(f)
		}()
	})
}

func (c *consumer) SetHandler(f Handler) {
	if c.handler == nil {
		c.handler = f
	}
}

func (c *consumer) GetHandler() Handler {
	return c.handler
}

func (c *consumer) GetHandlerName() string {
	return c.handlerName
}

func (c *consumer) start(f CallBack) {
	log := logging.DefaultLogger()

loop:
	for {
		select {
		case <-c.chanStop:
			log.Infow("Caught signal stop kafa consumer, terminating ...")
			break loop
		default:
			ev := c.instance.Poll(100)
			if ev == nil {
				continue
			}

			switch e := ev.(type) {
			case *kafka.Message:
				var topic string
				if e.TopicPartition.Topic != nil {
					topic = *e.TopicPartition.Topic
				}

				ctx, span := c.startSpan(e)

				err := processMessageWithRetry(ctx, f, e)
				if err != nil {
					span.RecordError(err)
					span.SetStatus(codes.Error, err.Error())
					span.End()

					log.Errorw("consumer process got error",
						zap.String("topic", topic), zap.Error(err),
						zap.ByteString("val", e.Value),
						zap.ByteString("key", e.Key),
						zap.Int64("offset", int64(e.TopicPartition.Offset)))

					if c.dlq {
						c.StoreDLQ(ctx, e)
					}
				}

				_, err = c.instance.StoreMessage(e)
				if err != nil {
					span.RecordError(err)
					span.SetStatus(codes.Error, err.Error())
					log.Warnw("Error storing offset after message", zap.Any("topic_partition", e.TopicPartition))
				}
				span.End()
			case kafka.Error:
				if !e.IsTimeout() {
					log.Errorw("Consume kafka got error", zap.Error(e))
				}
			case kafka.AssignedPartitions:
				log.Warnw("Partitions assigned", zap.Any("partitions", e.Partitions))
				err := c.instance.Assign(e.Partitions)
				if err != nil {
					log.Errorw("Failed to assign partitions", zap.Error(err))
				}
			case kafka.RevokedPartitions:
				log.Warnw("Partitions revoked", zap.Any("partitions", e.Partitions))

				if _, err := c.instance.Commit(); err != nil {
					log.Warnw("Failed to commit offsets on revoke", zap.Error(err))
				}
				err := c.instance.Unassign()
				if err != nil {
					log.Errorw("Failed to unassign partitions", zap.Error(err))
				}
			default:
			}
		}
	}
}

func (c *consumer) Stop() {
	c.stopSingleton.Do(func() {
		log := logging.DefaultLogger()
		log.Infow("Stopping Kafka consumer gracefully...")

		c.chanStop <- true

		time.Sleep(500 * time.Millisecond)

		if _, err := c.instance.Commit(); err != nil {
			log.Warnw("Failed to commit offsets on shutdown", zap.Error(err))
		}

		if err := c.instance.Unsubscribe(); err != nil {
			log.Warnw("Failed to unsubscribe", zap.Error(err))
		}

		if err := c.instance.Close(); err != nil {
			log.Warnw("Failed to close consumer", zap.Error(err))
		}

		log.Infow("Kafka consumer stopped successfully")
	})
}

const DLQSuffix = "dlq"

func (c *consumer) StoreDLQ(ctx context.Context, msg *kafka.Message) {
	// topic := *msg.TopicPartition.Topic
	// c.producer.ProduceRawWithKey(ctx, GetDLQTopic(topic), msg.Key, msg.Value)
}

func GetDLQTopic(topic string) string {
	if !strings.HasSuffix(topic, DLQSuffix) {
		topic = fmt.Sprintf("%s.%s", topic, DLQSuffix)
	}

	return topic
}

func processMessageWithRetry(ctx context.Context, f CallBack, msg *kafka.Message) error {
	retryTimes := uint(3)
	topic := *msg.TopicPartition.Topic
	if strings.HasSuffix(topic, DLQSuffix) {
		retryTimes = 0
	}

	options := []retry.Option{
		retry.Attempts(retryTimes),
		retry.DelayType(retry.BackOffDelay),
		retry.LastErrorOnly(true),
		retry.Context(ctx),
		retry.MaxDelay(time.Second * 5),
	}

	return retry.Do(func() error {
		return f(ctx, topic, msg.Value)
	}, options...)
}
