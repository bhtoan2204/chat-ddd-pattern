package messaging

import (
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"go.opentelemetry.io/otel/propagation"
)

type textMapCarrier struct {
	msg *kafka.Message
}

var _ propagation.TextMapCarrier = (*textMapCarrier)(nil)

func NewMessageCarrier(message *kafka.Message) propagation.TextMapCarrier {
	return &textMapCarrier{message}
}

func (c *textMapCarrier) Get(key string) string {
	for _, h := range c.msg.Headers {
		if h.Key == key {
			return string(h.Value)
		}
	}
	return ""
}

func (c *textMapCarrier) Set(key, value string) {
	for i := len(c.msg.Headers) - 1; i >= 0; i-- {
		if c.msg.Headers[i].Key == key {
			c.msg.Headers = append(c.msg.Headers[:i], c.msg.Headers[i+1:]...)
		}
	}
	c.msg.Headers = append(c.msg.Headers, kafka.Header{
		Key:   key,
		Value: []byte(value),
	})
}

func (c *textMapCarrier) Keys() []string {
	out := make([]string, len(c.msg.Headers))
	for i, h := range c.msg.Headers {
		out[i] = h.Key
	}
	return out
}
