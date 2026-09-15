package kafka

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"
	kafkago "github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader *kafkago.Reader
	logger *zerolog.Logger
}

type ConsumerMessageHandler func(context.Context, kafkago.Message) error

func NewConsumer(cfg ConsumerConfig, logger *zerolog.Logger) (*Consumer, error) {
	if len(cfg.Topics) == 0 || cfg.Topics[0] == "" {
		return nil, fmt.Errorf("at least one kafka topic is required")
	}

	reader := kafkago.NewReader(kafkago.ReaderConfig{
		Brokers:           cfg.Brokers,
		GroupID:           cfg.GroupID,
		Topic:             cfg.Topics[0],
		MinBytes:          cfg.MinBytes,
		MaxBytes:          cfg.MaxBytes,
		MaxWait:           cfg.MaxWait,
		ReadBatchTimeout:  cfg.ReadBatchTimeout,
		HeartbeatInterval: cfg.HeartbeatInterval,
		CommitInterval:    cfg.CommitInterval,
	})

	return &Consumer{reader: reader, logger: logger}, nil
}

func (c *Consumer) Start(ctx context.Context, handler ConsumerMessageHandler) error {
	cfg := c.reader.Config()
	c.logger.Info().Str("topic", cfg.Topic).Str("group_id", cfg.GroupID).Msg("starting kafka consumer")

	for {
		msg, err := c.reader.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				c.logger.Info().Msg("stopping kafka consumer")
				return nil
			}
			return fmt.Errorf("failed to read kafka message: %w", err)
		}

		if err := handler(ctx, msg); err != nil {
			c.logger.Error().Err(err).Str("topic", msg.Topic).Int("partition", msg.Partition).Int64("offset", msg.Offset).Msg("failed to process message")
			continue
		}

		c.logger.Debug().Str("topic", msg.Topic).Int("partition", msg.Partition).Int64("offset", msg.Offset).Msg("message processed")
	}
}

func (c *Consumer) Close() error {
	if c.reader == nil {
		return nil
	}
	return c.reader.Close()
}
