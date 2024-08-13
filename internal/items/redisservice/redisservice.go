package redisservice

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	pb "budgeting-service/genproto/account"

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

func (r *RedisService) StoreAccountInRedis(ctx context.Context, account *pb.AccountResponse) (*pb.AccountResponse, error) {
	key := fmt.Sprintf("account:%s", account.Id)
	athleteJSON, err := json.Marshal(account)
	if err != nil {
		r.logger.Error("Error marshalling account:", slog.String("err: ", err.Error()))
		return nil, err
	}

	if err := r.redisDb.Set(ctx, key, athleteJSON, 10*time.Minute).Err(); err != nil {
		r.logger.Error("Error setting account in Redis:", slog.String("err: ", err.Error()))
		return nil, err
	}

	return account, nil
}

func (r *RedisService) GetAccountFromRedis(ctx context.Context, id string) (*pb.AccountResponse, error) {
	key := fmt.Sprintf("account:%s", id)
	val, err := r.redisDb.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		r.logger.Error("Error getting account from Redis:", slog.String("err: ", err.Error()))
		return nil, err
	}

	var athlete pb.AccountResponse
	if err := json.Unmarshal([]byte(val), &athlete); err != nil {
		r.logger.Error("Error unmarshalling account:", slog.String("err: ", err.Error()))
		return nil, err
	}

	return &athlete, nil
}