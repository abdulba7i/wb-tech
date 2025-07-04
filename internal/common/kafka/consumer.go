package kafka

import (
	"fmt"
	"wb-tech/internal/pkg/logger"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

const (
	sessionTimeout = 7000
	noTimeout      = -1
)

type OrderHandler interface {
	HandleOrderMessage(message []byte, topic kafka.TopicPartition, cn int) error
}

type Consumer struct {
	consumer       *kafka.Consumer
	handler        OrderHandler
	stop           bool
	consumerNumber int
}

func NewConsumer(handler OrderHandler, address, topic, consumerGroup string, consumerNumber int) (*Consumer, error) {
	cfg := &kafka.ConfigMap{
		"bootstrap.servers":        address,
		"group.id":                 consumerGroup,
		"session.timeout.ms":       sessionTimeout,
		"enable.auto.offset.store": false,
		"enable.auto.commit":       true,
		"auto.commit.interval.ms":  5000,
		"auto.offset.reset":        "earliest",
	}

	c, err := kafka.NewConsumer(cfg)
	if err != nil {
		return nil, fmt.Errorf("error with new consumer: %w", err)
	}

	if err := c.Subscribe(topic, nil); err != nil {
		return nil, fmt.Errorf("error with subscribe topic: %w", err)
	}

	return &Consumer{
		consumer:       c,
		handler:        handler,
		consumerNumber: consumerNumber,
	}, nil
}

func (c *Consumer) Start() {
	for {
		if c.stop {
			break
		}

		kafkaMsg, err := c.consumer.ReadMessage(noTimeout)
		if err != nil {
			logger.Log.Errorf("Kafka read error: %v", err)
			continue
		}

		if kafkaMsg == nil {
			continue
		}

		logger.Log.Infof("Received raw message: %s", string(kafkaMsg.Value))

		if err = c.handler.HandleOrderMessage(kafkaMsg.Value, kafkaMsg.TopicPartition, c.consumerNumber); err != nil {
			logger.Log.Error(err)
			continue
		}
		if _, err = c.consumer.StoreMessage(kafkaMsg); err != nil {
			logger.Log.Error(err)
			continue
		}
	}

}

func (c *Consumer) Stop() error {
	c.stop = true
	if _, err := c.consumer.Commit(); err != nil {
		return err
	}

	logger.Log.Info("Commited offset")
	return c.consumer.Close()
}
