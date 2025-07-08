package main

import (
	"log"
	"os"
	"time"
	"wb-tech/internal/common/kafka"
	"wb-tech/internal/pkg/generate"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	topic := os.Getenv("TOPIC")
	producer, err := kafka.NewProducer("localhost:9092")

	if err != nil {
		log.Fatalf("Error with create producer: %v", err)
	}
	defer producer.Close()

	for i := 0; i < 10; i++ {
		order := generate.GenerateTestOrder()
		currentTime := time.Now().UTC()

		err := producer.Produce(order, topic, order.OrderUID, currentTime)

		if err != nil {
			log.Printf("Error sending order: %v", err)
			continue
		}

		log.Printf("Order %s successfully send", order.OrderUID)
		time.Sleep(2 * time.Second)
	}
}
