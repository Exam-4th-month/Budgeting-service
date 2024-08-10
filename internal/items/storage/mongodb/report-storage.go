package mongodb

import (
	"budgeting-service/internal/items/config"
	"budgeting-service/internal/items/redisservice"
	"budgeting-service/internal/items/repository"

	"go.mongodb.org/mongo-driver/mongo"

	"log/slog"
)

type ReportStorage struct {
	redis   *redisservice.RedisService
	mongodb *mongo.Database
	cfg     *config.Config
	logger  *slog.Logger
}

func NewReportStorage(redis *redisservice.RedisService, mongodb *mongo.Database, cfg *config.Config, logger *slog.Logger) repository.ReportI {
	return &ReportStorage{
		redis:   redis,
		mongodb: mongodb,
		cfg:     cfg,
		logger:  logger,
	}
}
