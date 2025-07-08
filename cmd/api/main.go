package main

import (
	"context"
	"log"
	"net/http"
	"time"
	"wb-tech/internal/cache"
	"wb-tech/internal/common/config"
	"wb-tech/internal/common/kafka"
	"wb-tech/internal/handler"
	"wb-tech/internal/pkg/logger"
	"wb-tech/internal/repository"
	"wb-tech/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	logger.Init()
	defer logger.Sync()

	if err := godotenv.Load("./.env"); err != nil {
		logger.Log.Panic("No .env file found")
	}

	cfg := config.MustLoad()

	db, err := repository.Connect(cfg.Database)
	if err != nil {
		log.Fatal("failed to connect to DB:", err)
	}
	repo := repository.NewRepository(db.DB())

	redisClient := cache.NewRedis(cfg.Redis)
	cacheService := cache.NewCacheService(repo.Order, redisClient)
	orderService := service.NewOrderService(repo, cacheService)

	ctx := context.Background()
	if err := orderService.PreloadCache(ctx); err != nil {
		log.Println("Cache preload error:", err)
	}

	go func() {
		orderHandler := service.NewOrderKafkaHandler(orderService, 1)
		consumer, err := kafka.NewConsumer(orderHandler, cfg.Kafka.Address, cfg.Kafka.Topic, cfg.Kafka.Group, 1)
		if err != nil {
			log.Fatal("kafka consumer error:", err)
		}
		consumer.Start()
	}()

	router := gin.Default()

	orderHandler := handler.NewOrderHandler(orderService)
	router.GET("/api/order/:uid", orderHandler.GetOrder)

	router.GET("/order/:uid", func(c *gin.Context) {
		c.File("./cmd/api/static/index.html")
	})

	srv := &http.Server{
		Addr:           ":8081",
		Handler:        router,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Println("HTTP server started on :8081")
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
