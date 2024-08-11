package service

import (
	pb "budgeting-service/genproto/notification"
	"budgeting-service/internal/items/repository"
	"log/slog"
)

type NotificationService struct {
	pb.UnimplementedNotificationServiceServer
	notificationstorage repository.NotificationI
	logger              *slog.Logger
}

func NewNotificationService(notificationstorage repository.NotificationI, logger *slog.Logger) *NotificationService {
	return &NotificationService{
		notificationstorage: notificationstorage,
		logger:              logger,
	}
}
