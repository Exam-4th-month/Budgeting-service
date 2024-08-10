package mongodb

import (
	"budgeting-service/internal/items/config"
	"budgeting-service/internal/items/redisservice"
	"budgeting-service/internal/items/repository"

	"go.mongodb.org/mongo-driver/mongo"

	"log/slog"
)

type TransactionStorage struct {
	redis   *redisservice.RedisService
	mongodb *mongo.Database
	cfg     *config.Config
	logger  *slog.Logger
}

func NewTransactionStorage(redis *redisservice.RedisService, mongodb *mongo.Database, cfg *config.Config, logger *slog.Logger) repository.TransactionI {
	return &TransactionStorage{
		redis:   redis,
		mongodb: mongodb,
		cfg:     cfg,
		logger:  logger,
	}
}
