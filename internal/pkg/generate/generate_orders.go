package generate

import (
	"fmt"
	"math/rand"
	"time"
	"wb-tech/internal/model"
)

func GenerateTestOrder() model.Order {
	id := time.Now().Unix()

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
				ChrtID:      9934930 + int(id),
				TrackNumber: fmt.Sprintf("WBILMTESTTRACK%d", id),
				Price:       453,
				RID:         fmt.Sprintf("ab4219087a764ae0btest%d", id),
				Name:        "Mascaras",
				Sale:        30,
				Size:        "0",
				TotalPrice:  317,
				NMID:        2389212 + int(id),
				Brand:       "Vivienne Sabo",
				Status:      202,
			},
		},
		Locale:            "en",
		InternalSignature: "",
		CustomerID:        "test",
		DeliveryService:   "meest",
		ShardKey:          "9",
		SMID:              99,
		DateCreated:       now,
		OOFShard:          "1",
	}
}
