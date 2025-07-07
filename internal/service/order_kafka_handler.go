package service

import (
	"context"
	"encoding/json"
	"fmt"
	"wb-tech/internal/model"
	"wb-tech/internal/pkg/logger"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type OrderKafkaHandler struct {
	service        *OrderService
	consumerNumber int
}

func NewOrderKafkaHandler(service *OrderService, consumerNumber int) *OrderKafkaHandler {
	return &OrderKafkaHandler{
		service:        service,
		consumerNumber: consumerNumber,
	}
}

func (h *OrderKafkaHandler) HandleOrderMessage(message []byte, topic kafka.TopicPartition, cn int) error {
	ctx := context.Background()

	var order model.Order
	if err := json.Unmarshal(message, &order); err != nil {
		logger.Log.Errorf("[Consumer %d] Ошибка парсинга сообщения из Kafka: %v", cn, err)
	}

	if order.OrderUID == "" {
		err := fmt.Errorf("[Consumer %d] order_uid пустой", cn)
		logger.Log.Error(err)
		return err
	}

	if err := h.service.CreateOrder(ctx, &order); err != nil {
		logger.Log.Errorf("[Consumer %d] Ошибка сохранения заказа: %v", cn, err)
		return err
	}

	logger.Log.Infof("[Consumer %d] Успешно обработан заказ %s", cn, order.OrderUID)
	return nil
}
