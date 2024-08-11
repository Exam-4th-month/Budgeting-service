package mongodb

import (
	"budgeting-service/internal/items/config"
	"budgeting-service/internal/items/redisservice"
	"budgeting-service/internal/items/repository"
	"context"

	pb "budgeting-service/genproto/notification"

	"go.mongodb.org/mongo-driver/mongo"

	"log/slog"
)

type NotificationStorage struct {
	redis   *redisservice.RedisService
	mongodb *mongo.Database
	cfg     *config.Config
	logger  *slog.Logger
}

func NewNotificationStorage(redis *redisservice.RedisService, mongodb *mongo.Database, cfg *config.Config, logger *slog.Logger) repository.NotificationI {
	return &NotificationStorage{
		redis:   redis,
		mongodb: mongodb,
		cfg:     cfg,
		logger:  logger,
	}
}

func (s *NotificationStorage) CreateNotification(ctx context.Context, req *pb.GetNotificationsRequest) (*pb.NotificationsResponse, error) {
	s.logger.Info("CreateNotification", slog.String("req", req.String()))

	// notificationCollection := s.mongodb.Collection("notifications")

	return nil, nil
}
