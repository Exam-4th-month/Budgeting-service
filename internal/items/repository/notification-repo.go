package repository

import (
	pb "budgeting-service/genproto/notification"
	"context"
)

type NotificationI interface {
	CreateNotification(ctx context.Context, req *pb.CreateNotificationRequest) (*pb.NotificationResponse, error)
	GetNotifications(ctx context.Context, req *pb.GetNotificationsRequest) (*pb.NotificationsResponse, error)
	MarkNotificationAsRead(ctx context.Context, req *pb.MarkNotificationAsReadRequest) (*pb.Empty, error)
}
