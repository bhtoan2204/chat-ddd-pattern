package pubsub

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
)

var (
	ErrBusClosed           = errors.New("localpubsub: bus closed")
	ErrSubscriptionClosed  = errors.New("localpubsub: subscription closed")
	ErrSubscriberQueueFull = errors.New("localpubsub: subscriber queue full")
)

type Message struct {
	Topic string
	Data  any
}

type PublishMode int

const (
	PublishBlocking PublishMode = iota
	PublishNonBlocking
)

type Bus struct {
	mu          sync.RWMutex
	closed      bool
	topics      map[string]map[uint64]*Subscription
	nextSubID   uint64
	defaultMode PublishMode
	bufferSize  int
}

type Subscription struct {
	id    uint64
	topic string
	bus   *Bus
	ch    chan Message

	stateMu sync.RWMutex
	closed  bool
}

type Config struct {
	BufferSize  int
	PublishMode PublishMode
}

func New(cfg Config) *Bus {
	bufferSize := cfg.BufferSize
	if bufferSize <= 0 {
		bufferSize = 16
	}

	return &Bus{
		topics:      make(map[string]map[uint64]*Subscription),
		defaultMode: cfg.PublishMode,
		bufferSize:  bufferSize,
	}
}

func (b *Bus) Subscribe(topic string) (*Subscription, error) {
	if topic == "" {
		return nil, errors.New("localpubsub: topic is required")
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	if b.closed {
		return nil, ErrBusClosed
	}

	id := atomic.AddUint64(&b.nextSubID, 1)

	sub := &Subscription{
		id:    id,
		topic: topic,
		bus:   b,
		ch:    make(chan Message, b.bufferSize),
	}

	if b.topics[topic] == nil {
		b.topics[topic] = make(map[uint64]*Subscription)
	}
	b.topics[topic][sub.id] = sub

	return sub, nil
}

func (b *Bus) Publish(ctx context.Context, topic string, data any) error {
	return b.publish(ctx, topic, data, b.defaultMode)
}

func (b *Bus) PublishWithMode(ctx context.Context, topic string, data any, mode PublishMode) error {
	return b.publish(ctx, topic, data, mode)
}

func (b *Bus) publish(ctx context.Context, topic string, data any, mode PublishMode) error {
	if topic == "" {
		return errors.New("localpubsub: topic is required")
	}

	b.mu.RLock()
	if b.closed {
		b.mu.RUnlock()
		return ErrBusClosed
	}

	subsMap := b.topics[topic]
	if len(subsMap) == 0 {
		b.mu.RUnlock()
		return nil
	}

	subs := make([]*Subscription, 0, len(subsMap))
	for _, sub := range subsMap {
		subs = append(subs, sub)
	}
	b.mu.RUnlock()

	msg := Message{
		Topic: topic,
		Data:  data,
	}

	for _, sub := range subs {
		if err := sub.deliver(ctx, msg, mode); err != nil {
			if errors.Is(err, ErrSubscriptionClosed) {
				continue
			}
			if mode == PublishBlocking {
				return err
			}
		}
	}

	return nil
}

func (b *Bus) Close() {
	b.mu.Lock()
	if b.closed {
		b.mu.Unlock()
		return
	}
	b.closed = true

	topics := b.topics
	b.topics = make(map[string]map[uint64]*Subscription)
	b.mu.Unlock()

	for _, subs := range topics {
		for _, sub := range subs {
			sub.close()
		}
	}
}

func (b *Bus) removeSubscription(topic string, id uint64) {
	b.mu.Lock()
	defer b.mu.Unlock()

	subs := b.topics[topic]
	if len(subs) == 0 {
		return
	}

	delete(subs, id)
	if len(subs) == 0 {
		delete(b.topics, topic)
	}
}

func (s *Subscription) C() <-chan Message {
	return s.ch
}

func (s *Subscription) Topic() string {
	return s.topic
}

func (s *Subscription) Unsubscribe() {
	s.bus.removeSubscription(s.topic, s.id)
	s.close()
}

func (s *Subscription) close() {
	s.stateMu.Lock()
	defer s.stateMu.Unlock()

	if s.closed {
		return
	}

	s.closed = true
	close(s.ch)
}

func (s *Subscription) deliver(ctx context.Context, msg Message, mode PublishMode) error {
	s.stateMu.RLock()
	defer s.stateMu.RUnlock()

	if s.closed {
		return ErrSubscriptionClosed
	}

	switch mode {
	case PublishNonBlocking:
		select {
		case s.ch <- msg:
			return nil
		default:
			return ErrSubscriberQueueFull
		}

	case PublishBlocking:
		select {
		case s.ch <- msg:
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}

	default:
		return errors.New("localpubsub: unknown publish mode")
	}
}
