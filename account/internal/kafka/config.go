package kafka

import (
	"time"

	kafkago "github.com/segmentio/kafka-go"
)

type ProducerConfig struct {
	Brokers      []string
	BatchSize    int
	BatchTimeout time.Duration
	RequiredAcks kafkago.RequiredAcks
}

type ConsumerConfig struct {
	Brokers           []string
	GroupID           string
	Topics            []string
	MinBytes          int
	MaxBytes          int
	MaxWait           time.Duration
	ReadBatchTimeout  time.Duration
	HeartbeatInterval time.Duration
	CommitInterval    time.Duration
}

func DefaultProducerConfig(brokers []string) ProducerConfig {
	return ProducerConfig{Brokers: brokers, BatchSize: 100, BatchTimeout: 10 * time.Millisecond, RequiredAcks: kafkago.RequireOne}
}

func DefaultConsumerConfig(brokers []string, groupID string, topics []string) ConsumerConfig {
	return ConsumerConfig{
		Brokers: brokers, GroupID: groupID, Topics: topics,
		MinBytes: 10e3, MaxBytes: 10e6, MaxWait: time.Second,
		ReadBatchTimeout: time.Second, HeartbeatInterval: time.Second, CommitInterval: time.Second,
	}
}
