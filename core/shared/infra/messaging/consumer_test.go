package messaging

import (
	"context"
	"errors"
	"testing"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

type fakeKafkaConsumer struct {
	committed []*kafka.Message
	seeked    []kafka.TopicPartition
}

func (f *fakeKafkaConsumer) SubscribeTopics([]string, kafka.RebalanceCb) error {
	return nil
}

func (f *fakeKafkaConsumer) Poll(int) kafka.Event {
	return nil
}

func (f *fakeKafkaConsumer) Assign([]kafka.TopicPartition) error {
	return nil
}

func (f *fakeKafkaConsumer) Unassign() error {
	return nil
}

func (f *fakeKafkaConsumer) Unsubscribe() error {
	return nil
}

func (f *fakeKafkaConsumer) Close() error {
	return nil
}

func (f *fakeKafkaConsumer) CommitMessage(msg *kafka.Message) ([]kafka.TopicPartition, error) {
	f.committed = append(f.committed, msg)
	return []kafka.TopicPartition{msg.TopicPartition}, nil
}

func (f *fakeKafkaConsumer) Seek(partition kafka.TopicPartition, _ int) error {
	f.seeked = append(f.seeked, partition)
	return nil
}

func TestHandleMessageCommitsOffsetOnlyAfterSuccessfulProcessing(t *testing.T) {
	topic := "CHATAPP.APPUSER.ACCOUNT_OUTBOX_EVENTS"
	msg := &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: 0, Offset: 12},
		Value:          []byte(`{"event_name":"EventAccountCreated"}`),
	}

	fake := &fakeKafkaConsumer{}
	consumer := &consumer{
		instance:     fake,
		chanStop:     make(chan bool, 1),
		retryBackoff: 0,
	}

	consumer.handleMessage(zap.NewNop().Sugar(), func(context.Context, string, []byte) error {
		return nil
	}, msg)

	if got := len(fake.committed); got != 1 {
		t.Fatalf("expected 1 committed offset, got %d", got)
	}
	if got := len(fake.seeked); got != 0 {
		t.Fatalf("expected no rewind on success, got %d", got)
	}
}

func TestHandleMessageRewindsFailedOffsetWithoutCommitting(t *testing.T) {
	topic := "CHATAPP.APPUSER.ACCOUNT_OUTBOX_EVENTS"
	msg := &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: 0, Offset: 27},
		Value:          []byte(`{"event_name":"EventAccountCreated"}`),
	}

	fake := &fakeKafkaConsumer{}
	consumer := &consumer{
		instance:     fake,
		chanStop:     make(chan bool, 1),
		retryBackoff: 0,
	}

	callCount := 0
	consumer.handleMessage(zap.NewNop().Sugar(), func(context.Context, string, []byte) error {
		callCount++
		return errors.New("boom")
	}, msg)

	if callCount == 0 {
		t.Fatalf("expected callback to be invoked")
	}
	if got := len(fake.committed); got != 0 {
		t.Fatalf("expected no committed offsets on failure, got %d", got)
	}
	if got := len(fake.seeked); got != 1 {
		t.Fatalf("expected 1 rewind on failure, got %d", got)
	}
	if fake.seeked[0].Offset != msg.TopicPartition.Offset {
		t.Fatalf("expected rewind to offset %v, got %v", msg.TopicPartition.Offset, fake.seeked[0].Offset)
	}
}

func TestStopClosesConsumerWithoutUnsubscribe(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockClient := NewMockkafkaConsumerClient(ctrl)

	mockClient.EXPECT().Close().Return(nil)

	consumer := &consumer{
		instance: mockClient,
		chanStop: make(chan bool, 1),
		loopDone: make(chan struct{}),
	}
	close(consumer.loopDone)

	consumer.Stop()
}
