package service

import (
	pb "budgeting-service/genproto/notification"
	"budgeting-service/internal/items/repository"
)

type NotificationService struct {
	pb.UnimplementedNotificationServiceServer
	notificationstorage repository.NotificationI
}

func NewNotificationService(notificationstorage repository.NotificationI) *NotificationService {
	return &NotificationService{
		notificationstorage: notificationstorage,
	}
}
