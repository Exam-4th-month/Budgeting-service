package mongodb

import (
	"budgeting-service/internal/items/config"
	"budgeting-service/internal/items/redisservice"
	"budgeting-service/internal/items/repository"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	pb "budgeting-service/genproto/goal"

	"log/slog"
)

type GoalStorage struct {
	redis   *redisservice.RedisService
	mongodb *mongo.Database
	cfg     *config.Config
	logger  *slog.Logger
}

func NewGoalStorage(redis *redisservice.RedisService, mongodb *mongo.Database, cfg *config.Config, logger *slog.Logger) repository.GoalI {
	return &GoalStorage{
		redis:   redis,
		mongodb: mongodb,
		cfg:     cfg,
		logger:  logger,
	}
}

func (s *GoalStorage) CreateGoal(ctx context.Context, req *pb.CreateGoalRequest) (*pb.GoalResponse, error) {
	s.logger.Info("CreateGoal", slog.String("req", req.String()))

	goalCollecton := s.mongodb.Collection("goals")
	created_at := time.Now()

	goalDoc := bson.D{
		{Key: "user_id", Value: req.UserId},
		{Key: "name", Value: req.Name},
		{Key: "target_amount", Value: req.TargetAmount},
		{Key: "current_amount", Value: req.CurrentAmount},
		{Key: "deadline", Value: req.Deadline},
		{Key: "status", Value: req.Status},
		{Key: "created_at", Value: created_at},
	}

	res, err := goalCollecton.InsertOne(ctx, goalDoc)
	if err != nil {
		s.logger.Error(err.Error())
		return nil, err
	}

	goalID := res.InsertedID.(primitive.ObjectID).Hex()

	return &pb.GoalResponse{
		Id:            goalID,
		UserId:        req.UserId,
		Name:          req.Name,
		TargetAmount:  req.TargetAmount,
		CurrentAmount: req.CurrentAmount,
		Deadline:      req.Deadline,
		Status:        req.Status,
		CreatedAt:     created_at.String(),
	}, nil
}

func (s *GoalStorage) GetGoals(ctx context.Context, req *pb.GetGoalsRequest) (*pb.GoalsResponse, error) {
	s.logger.Info("GetGoals", slog.String("req", req.String()))
	goalCollection := s.mongodb.Collection("goals")

	filter := bson.D{{Key: "user_id", Value: req.UserId}}

	cursor, err := goalCollection.Find(ctx, filter)
	if err != nil {
		s.logger.Error(err.Error())
		return nil, err
	}
	defer cursor.Close(ctx)

	var goals []*pb.GoalResponse
	for cursor.Next(ctx) {
		var goal bson.M
		if err := cursor.Decode(&goal); err != nil {
			s.logger.Error(err.Error())
			return nil, err
		}

		goals = append(goals, &pb.GoalResponse{
			Id:            goal["_id"].(primitive.ObjectID).Hex(),
			UserId:        goal["user_id"].(string),
			Name:          goal["name"].(string),
			TargetAmount:  float32(goal["target_amount"].(float64)),
			CurrentAmount: float32(goal["current_amount"].(float64)),
			Deadline:      goal["deadline"].(string),
			Status:        goal["status"].(string),
			CreatedAt:     goal["created_at"].(primitive.DateTime).Time().String(),
			UpdatedAt:     goal["updated_at"].(primitive.DateTime).Time().String(),
		})
	}

	if err := cursor.Err(); err != nil {
		s.logger.Error(err.Error())
		return nil, err
	}

	return &pb.GoalsResponse{Goals: goals}, nil
}

func (s *GoalStorage) GetGoalById(ctx context.Context, req *pb.GetGoalByIdRequest) (*pb.GoalResponse, error) {
	s.logger.Info("GetGoalById", slog.String("req", req.Id))
	goalCollection := s.mongodb.Collection("goals")

	objID, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		s.logger.Error(err.Error())
		return nil, err
	}

	filter := bson.D{{Key: "_id", Value: objID}}

	var goal bson.M
	err = goalCollection.FindOne(ctx, filter).Decode(&goal)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			s.logger.Error(err.Error())
			return nil, nil
		}
		s.logger.Error(err.Error())
		return nil, err
	}

	return &pb.GoalResponse{
		Id:            goal["_id"].(primitive.ObjectID).Hex(),
		UserId:        goal["user_id"].(string),
		Name:          goal["name"].(string),
		TargetAmount:  float32(goal["target_amount"].(float64)),
		CurrentAmount: float32(goal["current_amount"].(float64)),
		Deadline:      goal["deadline"].(string),
		Status:        goal["status"].(string),
		CreatedAt:     goal["created_at"].(primitive.DateTime).Time().String(),
		UpdatedAt:     goal["updated_at"].(primitive.DateTime).Time().String(),
	}, nil
}

func (s *GoalStorage) UpdateGoal(ctx context.Context, req *pb.UpdateGoalRequest) (*pb.GoalResponse, error) {
	s.logger.Info("UpdateGoal", slog.String("req", req.String()))
	goalCollection := s.mongodb.Collection("goals")

	objID, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		s.logger.Error(err.Error())
		return nil, err
	}

	filter := bson.D{
		{Key: "_id", Value: objID},
		{Key: "deleted_at", Value: bson.D{{Key: "$exists", Value: false}}},
	}

	updateFields := bson.D{}
	if req.Name != "" {
		updateFields = append(updateFields, bson.E{Key: "name", Value: req.Name})
	}
	if req.TargetAmount != 0 {
		updateFields = append(updateFields, bson.E{Key: "target_amount", Value: req.TargetAmount})
	}
	if req.CurrentAmount != 0 {
		updateFields = append(updateFields, bson.E{Key: "current_amount", Value: req.CurrentAmount})
	}
	if req.Deadline != "" {
		updateFields = append(updateFields, bson.E{Key: "deadline", Value: req.Deadline})
	}
	if req.Status != "" {
		updateFields = append(updateFields, bson.E{Key: "status", Value: req.Status})
	}
	if len(updateFields) > 0 {
		updateFields = append(updateFields, bson.E{Key: "updated_at", Value: time.Now()})
	}

	if len(updateFields) == 0 {
		s.logger.Info("No fields to update")
		return nil, nil
	}

	update := bson.D{{Key: "$set", Value: updateFields}}

	res := goalCollection.FindOneAndUpdate(ctx, filter, update, options.FindOneAndUpdate().SetReturnDocument(options.After))
	if res.Err() != nil {
		if res.Err() == mongo.ErrNoDocuments {
			s.logger.Error(res.Err().Error())
			return nil, nil
		}
		s.logger.Error(res.Err().Error())
		return nil, res.Err()
	}

	var updatedGoal bson.M
	if err = res.Decode(&updatedGoal); err != nil {
		s.logger.Error(err.Error())
		return nil, err
	}

	return &pb.GoalResponse{
		Id:            updatedGoal["_id"].(primitive.ObjectID).Hex(),
		UserId:        updatedGoal["user_id"].(string),
		Name:          updatedGoal["name"].(string),
		TargetAmount:  float32(updatedGoal["target_amount"].(float64)),
		CurrentAmount: float32(updatedGoal["current_amount"].(float64)),
		Deadline:      updatedGoal["deadline"].(string),
		Status:        updatedGoal["status"].(string),
		CreatedAt:     updatedGoal["created_at"].(primitive.DateTime).Time().String(),
		UpdatedAt:     updatedGoal["updated_at"].(primitive.DateTime).Time().String(),
	}, nil
}

func (s *GoalStorage) DeleteGoal(ctx context.Context, req *pb.DeleteGoalRequest) (*pb.Empty, error) {
	s.logger.Info("DeleteGoal", slog.String("req", req.Id))
	goalCollection := s.mongodb.Collection("goals")

	objID, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		s.logger.Error(err.Error())
		return nil, err
	}

	filter := bson.D{{Key: "_id", Value: objID}}
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "deleted_at", Value: time.Now()},
		}},
	}

	_, err = goalCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		s.logger.Error(err.Error())
		return nil, err
	}

	return &pb.Empty{}, nil
}
