package redisservice

import (
	"log/slog"

	"github.com/go-redis/redis/v8"
)

type (
	RedisService struct {
		redisDb *redis.Client
		logger  *slog.Logger
	}
)

func New(redisDb *redis.Client, logger *slog.Logger) *RedisService {
	return &RedisService{
		logger:  logger,
		redisDb: redisDb,
	}
}
