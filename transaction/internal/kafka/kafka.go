package kafka

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/rs/zerolog"
	kafkago "github.com/segmentio/kafka-go"
)

type MessageHandler func(ctx context.Context, topic, key string, data []byte) error

type Client interface {
	Publish(ctx context.Context, topic, key string, data interface{}) error
	Subscribe(ctx context.Context, topic string, handler MessageHandler) error
	Close() error
}

type Kafka struct {
	producer *Producer
	brokers  []string
	groupID  string
	logger   *zerolog.Logger

	mu        sync.Mutex
	consumers []*Consumer
}

func New(producer *Producer, brokers []string, groupID string, logger *zerolog.Logger) *Kafka {
	return &Kafka{producer: producer, brokers: brokers, groupID: groupID, logger: logger}
}

func (k *Kafka) Publish(ctx context.Context, topic, key string, data interface{}) error {
	return k.producer.SendJSONMessage(ctx, topic, key, data)
}

func (k *Kafka) Subscribe(ctx context.Context, topic string, handler MessageHandler) error {
	groupID := k.groupID + "_" + topic
	consumer, err := NewConsumer(DefaultConsumerConfig(k.brokers, groupID, []string{topic}), k.logger)
	if err != nil {
		return err
	}

	k.mu.Lock()
	k.consumers = append(k.consumers, consumer)
	k.mu.Unlock()

	go func() {
		k.logger.Info().Str("topic", topic).Str("group_id", groupID).Msg("subscribing to topic")
		if err := consumer.Start(ctx, func(ctx context.Context, msg kafkago.Message) error {
			return handler(ctx, msg.Topic, string(msg.Key), msg.Value)
		}); err != nil {
			k.logger.Error().Err(err).Str("topic", topic).Msg("kafka consumer stopped")
		}
	}()

	return nil
}

func (k *Kafka) Close() error {
	k.mu.Lock()
	consumers := append([]*Consumer(nil), k.consumers...)
	k.consumers = nil
	k.mu.Unlock()

	var errs []error
	for _, consumer := range consumers {
		if err := consumer.Close(); err != nil {
			errs = append(errs, fmt.Errorf("failed to close consumer: %w", err))
		}
	}
	if k.producer != nil {
		if err := k.producer.Close(); err != nil {
			errs = append(errs, fmt.Errorf("failed to close producer: %w", err))
		}
	}
	return errors.Join(errs...)
}
