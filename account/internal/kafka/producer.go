package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/rs/zerolog"
	kafkago "github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafkago.Writer
	logger *zerolog.Logger
}

type Message struct {
	Topic     string
	Key       []byte
	Value     []byte
	Headers   map[string]string
	Partition int
}

func NewProducer(cfg ProducerConfig, logger *zerolog.Logger) *Producer {
	return &Producer{
		writer: &kafkago.Writer{
			Addr: kafkago.TCP(cfg.Brokers...), Balancer: &kafkago.LeastBytes{},
			BatchSize: cfg.BatchSize, BatchTimeout: cfg.BatchTimeout,
			RequiredAcks: cfg.RequiredAcks, Async: false,
		},
		logger: logger,
	}
}

func (p *Producer) SendMessage(ctx context.Context, msg Message) error {
	kafkaMsg := kafkago.Message{Topic: msg.Topic, Key: msg.Key, Value: msg.Value, Partition: msg.Partition, Time: time.Now()}
	for key, value := range msg.Headers {
		kafkaMsg.Headers = append(kafkaMsg.Headers, kafkago.Header{Key: key, Value: []byte(value)})
	}
	if err := p.writer.WriteMessages(ctx, kafkaMsg); err != nil {
		p.logger.Error().Err(err).Str("topic", msg.Topic).Str("key", string(msg.Key)).Msg("failed to send message to kafka")
		return fmt.Errorf("failed to send message to kafka: %w", err)
	}
	p.logger.Debug().Str("topic", msg.Topic).Str("key", string(msg.Key)).Msg("message sent to kafka")
	return nil
}

func (p *Producer) SendJSONMessage(ctx context.Context, topic, key string, data interface{}) error {
	value, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal message to JSON: %w", err)
	}
	return p.SendMessage(ctx, Message{Topic: topic, Key: []byte(key), Value: value, Headers: map[string]string{"content-type": "application/json"}})
}

func (p *Producer) Close() error {
	if p.writer == nil {
		return nil
	}
	return p.writer.Close()
}
