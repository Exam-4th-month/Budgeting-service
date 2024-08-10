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
		log.Fatalln(err)
	}

	logFile, err := os.OpenFile("application.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()

	logger := slog.New(slog.NewJSONHandler(logFile, nil))

	db, err := mdb.ConnectDB(config)
	if err != nil {
		logger.Error("error while connecting postgres:", slog.String("err:", err.Error()))
	}

	redis, err := redisCl.NewRedisDB(config)
	if err != nil {
		logger.Error("error while connecting redis:", slog.String("err:", err.Error()))
	}

	service := service.New(storage.New(
		redisservice.New(redis, logger),
		db,
		config,
		logger,
	), logger)

	msgs := msgbroker.InitMessageBroker(config)

	msgbroker := msgbroker.New(service, logger, msgs, &sync.WaitGroup{}, 4)

	api := api.New(service)

	go func() {
		log.Fatalln(api.RUN(config, service))
	}()

	msgbroker.StartToConsume(context.Background(), "application/json")

}
