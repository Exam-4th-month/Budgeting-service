package mongodb

import (
	"budgeting-service/internal/items/config"
	"budgeting-service/internal/items/repository"
	"context"
	"time"

	pb "budgeting-service/genproto/notification"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"log/slog"
)

type NotificationStorage struct {
	mongodb *mongo.Database
	cfg     *config.Config
	logger  *slog.Logger
}

func NewNotificationStorage(mongodb *mongo.Database, cfg *config.Config, logger *slog.Logger) repository.NotificationI {
	return &NotificationStorage{
		mongodb: mongodb,
		cfg:     cfg,
		logger:  logger,
	}
}

func (s *NotificationStorage) CreateNotification(ctx context.Context, req *pb.CreateNotificationRequest) (*pb.NotificationResponse, error) {
	s.logger.Info("CreateNotification", slog.String("req", req.String()))

	notificationCollection := s.mongodb.Collection("notifications")
	created_at := time.Now()

	notificationDoc := bson.D{
		{Key: "user_id", Value: req.UserId},
		{Key: "message", Value: req.Message},
		{Key: "is_read", Value: false},
		{Key: "created_at", Value: created_at},
	}

	res, err := notificationCollection.InsertOne(ctx, notificationDoc)
	if err != nil {
		s.logger.Error("Error creating notification", slog.Any("error", err))
		return nil, err
	}

	notificationID := res.InsertedID.(primitive.ObjectID).Hex()

	return &pb.NotificationResponse{
		Id:        notificationID,
		UserId:    req.UserId,
		Message:   req.Message,
		IsRead:    false,
		CreatedAt: created_at.String(),
	}, nil
}

func (s *NotificationStorage) GetNotifications(ctx context.Context, req *pb.GetNotificationsRequest) (*pb.NotificationsResponse, error) {
	s.logger.Info("GetNotification", slog.String("id", req.UserId))

	notificationCollection := s.mongodb.Collection("notifications")

	filter := bson.D{
		{Key: "user_id", Value: req.UserId},
		{Key: "is_read", Value: false},
	}

	cursor, err := notificationCollection.Find(ctx, filter)
	if err != nil {
		s.logger.Error("Error finding notifications", slog.Any("error", err))
		return nil, err
	}
	defer cursor.Close(ctx)

	var notifications []*pb.NotificationResponse
	for cursor.Next(ctx) {
		var notification bson.M
		if err := cursor.Decode(&notification); err != nil {
			s.logger.Error("Error decoding notification", slog.Any("error", err))
			return nil, err
		}
		notifications = append(notifications, &pb.NotificationResponse{
			Id:        notification["_id"].(primitive.ObjectID).Hex(),
			UserId:    notification["user_id"].(string),
			Message:   notification["message"].(string),
			IsRead:    notification["is_read"].(bool),
			CreatedAt: notification["created_at"].(primitive.DateTime).Time().String(),
		})
	}

	if err := cursor.Err(); err != nil {
		s.logger.Error("Cursor error", slog.Any("error", err))
		return nil, err
	}

	return &pb.NotificationsResponse{Notifications: notifications}, nil
}

func (s *NotificationStorage) MarkNotificationAsRead(ctx context.Context, req *pb.MarkNotificationAsReadRequest) (*pb.Empty, error) {
	s.logger.Info("MarkNotificationAsRead", slog.String("id", req.Id))

	notificationCollection := s.mongodb.Collection("notifications")

	objID, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		s.logger.Error("Invalid ObjectID", slog.Any("error", err))
		return nil, err
	}

	update := bson.D{{Key: "$set", Value: bson.D{{Key: "is_read", Value: true}}}}

	_, err = notificationCollection.UpdateByID(ctx, objID, update)
	if err != nil {
		s.logger.Error("Error marking notification as read", slog.Any("error", err))
		return nil, err
	}

	return &pb.Empty{}, nil
}
