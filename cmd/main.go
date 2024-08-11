package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"sync"

	"budgeting-service/api"
	"budgeting-service/internal/items/config"
	"budgeting-service/internal/items/msgbroker"
	"budgeting-service/internal/items/redisservice"
	"budgeting-service/internal/items/service"
	"budgeting-service/internal/items/storage"
	mdb "budgeting-service/internal/items/storage/mongodb"
	redisCl "budgeting-service/internal/pkg/redis"
)

func main() {
	config, err := config.New()
	if err != nil {
		log.Fatalln("Error loading config:", err)
	}

	logFile, err := os.OpenFile("application.log", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatalln("Error opening log file:", err)
	}
	defer logFile.Close()

	logger := slog.New(slog.NewJSONHandler(logFile, nil))

	db, err := mdb.ConnectDB(config)
	if err != nil {
		logger.Error("Error connecting to MongoDB", slog.String("err", err.Error()))
	}

	redis, err := redisCl.NewRedisDB(config)
	if err != nil {
		logger.Error("Error connecting to Redis", slog.String("err", err.Error()))
	}

	service := service.New(storage.New(
		redisservice.New(redis, logger),
		db,
		config,
		logger,
	), logger)

	msgs, conn, ch := msgbroker.InitMessageBroker(config)
	defer conn.Close()
	defer ch.Close()

	msgbroker := msgbroker.New(service, logger, msgs, &sync.WaitGroup{}, 4)

	api := api.New(service)

	go func() {
		if err := api.RUN(config, service); err != nil {
			logger.Error("API run error", slog.String("err", err.Error()))
		}
	}()

	msgbroker.StartToConsume(context.Background(), "application/json")
}
