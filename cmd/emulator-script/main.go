package main

import (
	"log"
	"os"
	"time"
	"wb-tech/internal/common/kafka"
	"wb-tech/internal/pkg/generate"
)

func main() {
	producer, err := kafka.NewProducer(os.Getenv("KAFKA_ADDRESS"))
	if err != nil {
		log.Fatalf("Error with create producer: %v", err)
	}
	defer producer.Close()

	for i := 0; i < 10; i++ {
		order := generate.GenerateTestOrder(i)
		currentTime := time.Now().UTC()

		err := producer.Produce(order, os.Getenv("TOPIC"), order.OrderUID, currentTime)
		if err != nil {
			log.Printf("Error sending order: %v", err)
			continue
		}

		log.Printf("Order %s successfully send", order.OrderUID)
		time.Sleep(2 * time.Second)
	}
}
