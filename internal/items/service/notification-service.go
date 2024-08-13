package service

import (
	pb "budgeting-service/genproto/notification"
	"budgeting-service/internal/items/repository"
	"context"
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

func (s *NotificationService) CreateNotification(ctx context.Context, req *pb.CreateNotificationRequest) (*pb.NotificationResponse, error) {
	s.logger.Info("CreateNotification", slog.String("req", req.String()))
	return s.notificationstorage.CreateNotification(ctx, req)
}

func (s *NotificationService) GetNotifications(ctx context.Context, req *pb.GetNotificationsRequest) (*pb.NotificationsResponse, error) {
	s.logger.Info("GetNotification", slog.String("id", req.UserId))
	return s.notificationstorage.GetNotifications(ctx, req)
}

func (s *NotificationService) MarkNotificationAsRead(ctx context.Context, req *pb.MarkNotificationAsReadRequest) (*pb.Empty, error) {
	s.logger.Info("MarkNotificationAsRead", slog.String("id", req.Id))
	return s.notificationstorage.MarkNotificationAsRead(ctx, req)
}
