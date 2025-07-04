package kafka

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
	"wb-tech/internal/model"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

const (
	flushTimeout   = 5000
	messageTimeout = 5000
)

var errUnknownType = errors.New("unknown event type")

type Producer struct {
	producer *kafka.Producer
}

func NewProducer(address string) (*Producer, error) {
	conf := kafka.ConfigMap{
		"bootstrap.servers":  address,
		"message.timeout.ms": messageTimeout,
		"acks":               "all",
	}

	p, err := kafka.NewProducer(&conf)
	if err != nil {
		return nil, fmt.Errorf("error with producer: %w", err)
	}

	return &Producer{producer: p}, nil
}

func (p *Producer) Produce(message model.Order, topic, key string, tn time.Time) error {
	messageBytes, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("error marshaling message to JSON: %w", err)
	}

	kafkaMsg := &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Value:     messageBytes,
		Key:       []byte(key),
		Timestamp: tn,
	}

	kafkaChan := make(chan kafka.Event)

	if err := p.producer.Produce(kafkaMsg, kafkaChan); err != nil {
		return fmt.Errorf("error with chan: %w", err)
	}

	e := <-kafkaChan
	switch ev := e.(type) {
	case *kafka.Message:
		return nil
	case kafka.Error:
		return fmt.Errorf("kafka error: %v", ev)
	default:
		return errUnknownType
	}
}

func (p *Producer) Close() {
	p.producer.Flush(flushTimeout)
	p.producer.Close()
}
