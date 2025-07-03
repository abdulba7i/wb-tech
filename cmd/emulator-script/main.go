package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"
	"wb-tech/internal/common/kafka"
	"wb-tech/internal/model"
)

func main() {
	producer, err := kafka.NewProducer(os.Getenv("KAFKA_ADDRESS"))
	if err != nil {
		log.Fatalf("Error with create producer: %v", err)
	}
	defer producer.Close()

	for i := 0; i < 10; i++ {
		order := generateTestOrder(i)
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

func generateTestOrder(id int) model.Order {
	now := time.Now().UTC()

	return model.Order{
		OrderUID:    fmt.Sprintf("b563feb7b2b84b6test%d", id),
		TrackNumber: fmt.Sprintf("WBILMTESTTRACK%d", id),
		Entry:       "WBIL",
		Delivery: model.Delivery{
			Name:    "Test Testov",
			Phone:   "+9720000000",
			Zip:     "2639809",
			City:    "Kiryat Mozkin",
			Address: "Ploshad Mira 15",
			Region:  "Kraiot",
			Email:   "test@gmail.com",
		},
		Payment: model.Payment{
			Transaction:  fmt.Sprintf("b563feb7b2b84b6test%d", id),
			Currency:     "USD",
			Provider:     "wbpay",
			Amount:       rand.Intn(1000) + 500,
			PaymentDT:    now.Unix(),
			Bank:         "alpha",
			DeliveryCost: 1500,
			GoodsTotal:   rand.Intn(500) + 100,
		},
		Items: []model.Item{
			{
				ChrtID:      9934930 + id,
				TrackNumber: fmt.Sprintf("WBILMTESTTRACK%d", id),
				Price:       453,
				RID:         fmt.Sprintf("ab4219087a764ae0btest%d", id),
				Name:        "Mascaras",
				Sale:        30,
				Size:        "0",
				TotalPrice:  317,
				NMID:        2389212 + id,
				Brand:       "Vivienne Sabo",
				Status:      202,
			},
		},
		DateCreated: now.Format(time.RFC3339),
	}
}
