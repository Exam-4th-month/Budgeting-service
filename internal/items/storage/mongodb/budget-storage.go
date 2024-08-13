package mongodb

import (
	pb "budgeting-service/genproto/budget"
	"budgeting-service/internal/items/config"
	"budgeting-service/internal/items/redisservice"
	"budgeting-service/internal/items/repository"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"log/slog"
)

type BudgetStorage struct {
	redis   *redisservice.RedisService
	mongodb *mongo.Database
	cfg     *config.Config
	logger  *slog.Logger
}

func NewBudgetStorage(redis *redisservice.RedisService, mongodb *mongo.Database, cfg *config.Config, logger *slog.Logger) repository.BudgetI {
	return &BudgetStorage{
		redis:   redis,
		mongodb: mongodb,
		cfg:     cfg,
		logger:  logger,
	}
}

func (s *BudgetStorage) CreateBudget(ctx context.Context, req *pb.CreateBudgetRequest) (*pb.BudgetResponse, error) {
	s.logger.Info("CreateBudget", slog.String("req", req.String()))
	budgetCollection := s.mongodb.Collection("budgets")
	created_at := time.Now()

	budgetDoc := bson.D{
		{Key: "user_id", Value: req.UserId},
		{Key: "category_id", Value: req.CategoryId},
		{Key: "amount", Value: req.Amount},
		{Key: "period", Value: req.Period},
		{Key: "start_date", Value: req.StartDate},
		{Key: "end_date", Value: req.EndDate},
		{Key: "created_at", Value: created_at},
		{Key: "updated_at", Value: created_at},
		{Key: "deleted_at", Value: nil},
	}

	res, err := budgetCollection.InsertOne(ctx, budgetDoc)
	if err != nil {
		s.logger.Error("Error while inserting budget", slog.Any("error", err))
		return nil, err
	}

	budgetID := res.InsertedID.(primitive.ObjectID)

	return &pb.BudgetResponse{
		Id:         budgetID.Hex(),
		UserId:     req.UserId,
		CategoryId: req.CategoryId,
		Amount:     req.Amount,
		Period:     req.Period,
		StartDate:  req.StartDate,
		EndDate:    req.EndDate,
		CreatedAt:  created_at.String(),
	}, nil
}

func (s *BudgetStorage) GetBudgets(ctx context.Context, req *pb.GetBudgetsRequest) (*pb.BudgetsResponse, error) {
	s.logger.Info("GetBudgets", slog.String("req", req.String()))
	budgetCollection := s.mongodb.Collection("budgets")

	filter := bson.D{{Key: "user_id", Value: req.UserId}}

	cursor, err := budgetCollection.Find(ctx, filter)
	if err != nil {
		s.logger.Error("Error while retrieving budgets", slog.Any("error", err))
		return nil, err
	}
	defer cursor.Close(ctx)

	var budgets []*pb.BudgetResponse
	for cursor.Next(ctx) {
		var budget bson.M
		if err = cursor.Decode(&budget); err != nil {
			s.logger.Error("Error while decoding budget", slog.Any("error", err))
			return nil, err
		}

		budgets = append(budgets, &pb.BudgetResponse{
			Id:         budget["_id"].(primitive.ObjectID).Hex(),
			UserId:     budget["user_id"].(string),
			CategoryId: budget["category_id"].(string),
			Amount:     budget["amount"].(float32),
			Period:     budget["period"].(string),
			StartDate:  budget["start_date"].(string),
			EndDate:    budget["end_date"].(string),
			CreatedAt:  budget["created_at"].(primitive.DateTime).Time().String(),
			UpdatedAt:  budget["updated_at"].(primitive.DateTime).Time().String(),
		})
	}

	if err := cursor.Err(); err != nil {
		s.logger.Error("Error while iterating over cursor", slog.Any("error", err))
		return nil, err
	}

	return &pb.BudgetsResponse{Budgets: budgets}, nil
}

func (s *BudgetStorage) GetBudgetById(ctx context.Context, req *pb.GetBudgetByIdRequest) (*pb.BudgetResponse, error) {
	s.logger.Info("GetBudgetById", slog.String("req", req.Id))
	budgetCollection := s.mongodb.Collection("budgets")

	objID, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		s.logger.Error("Error while converting ID", slog.Any("error", err))
		return nil, err
	}

	filter := bson.D{{Key: "_id", Value: objID}}

	var budget bson.M
	err = budgetCollection.FindOne(ctx, filter).Decode(&budget)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			s.logger.Info("Budget not found")
			return nil, nil
		}
		s.logger.Error("Error while retrieving budget", slog.Any("error", err))
		return nil, err
	}

	return &pb.BudgetResponse{
		Id:         budget["_id"].(primitive.ObjectID).Hex(),
		UserId:     budget["user_id"].(string),
		CategoryId: budget["category_id"].(string),
		Amount:     budget["amount"].(float32),
		Period:     budget["period"].(string),
		StartDate:  budget["start_date"].(string),
		EndDate:    budget["end_date"].(string),
		CreatedAt:  budget["created_at"].(primitive.DateTime).Time().String(),
		UpdatedAt:  budget["updated_at"].(primitive.DateTime).Time().String(),
	}, nil
}

func (s *BudgetStorage) UpdateBudget(ctx context.Context, req *pb.UpdateBudgetRequest) (*pb.BudgetResponse, error) {
	budgetCollection := s.mongodb.Collection("budgets")

	objID, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		s.logger.Error("Error while converting ID", slog.Any("error", err))
		return nil, err
	}

	filter := bson.D{
		{Key: "_id", Value: objID},
		{Key: "deleted_at", Value: bson.D{{Key: "$exists", Value: false}}},
	}

	updateFields := bson.D{}
	if req.Amount != 0 {
		updateFields = append(updateFields, bson.E{Key: "amount", Value: req.Amount})
	}
	if req.Period != "" {
		updateFields = append(updateFields, bson.E{Key: "period", Value: req.Period})
	}
	if req.StartDate != "" {
		updateFields = append(updateFields, bson.E{Key: "start_date", Value: req.StartDate})
	}
	if req.EndDate != "" {
		updateFields = append(updateFields, bson.E{Key: "end_date", Value: req.EndDate})
	}
	if len(updateFields) > 0 {
		updateFields = append(updateFields, bson.E{Key: "updated_at", Value: time.Now()})
	}

	if len(updateFields) == 0 {
		s.logger.Info("No fields to update")
		return nil, nil
	}

	update := bson.D{{Key: "$set", Value: updateFields}}

	res := budgetCollection.FindOneAndUpdate(ctx, filter, update, options.FindOneAndUpdate().SetReturnDocument(options.After))
	if res.Err() != nil {
		if res.Err() == mongo.ErrNoDocuments {
			s.logger.Info("Budget not found")
			return nil, nil
		}
		s.logger.Error("Error while updating budget", slog.Any("error", res.Err()))
		return nil, res.Err()
	}

	var updatedBudget bson.M
	if err = res.Decode(&updatedBudget); err != nil {
		s.logger.Error("Error while decoding updated budget", slog.Any("error", err))
		return nil, err
	}

	return &pb.BudgetResponse{
		Id:         updatedBudget["_id"].(primitive.ObjectID).Hex(),
		UserId:     updatedBudget["user_id"].(string),
		CategoryId: updatedBudget["category_id"].(string),
		Amount:     updatedBudget["amount"].(float32),
		Period:     updatedBudget["period"].(string),
		StartDate:  updatedBudget["start_date"].(string),
		EndDate:    updatedBudget["end_date"].(string),
		CreatedAt:  updatedBudget["created_at"].(primitive.DateTime).Time().String(),
		UpdatedAt:  updatedBudget["updated_at"].(primitive.DateTime).Time().String(),
	}, nil
}

func (s *BudgetStorage) DeleteBudget(ctx context.Context, req *pb.DeleteBudgetRequest) (*pb.Empty, error) {
	budgetCollection := s.mongodb.Collection("budgets")

	objID, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		s.logger.Error("Error while converting ID", slog.Any("error", err))
		return nil, err
	}

	filter := bson.D{{Key: "_id", Value: objID}}
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "deleted_at", Value: time.Now()},
		}},
	}

	_, err = budgetCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		s.logger.Error("Error while deleting budget", slog.Any("error", err))
		return nil, err
	}

	return &pb.Empty{}, nil
}
