package mongodb

import (
	pb "budgeting-service/genproto/report"
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

type ReportStorage struct {
	redis   *redisservice.RedisService
	mongodb *mongo.Database
	cfg     *config.Config
	logger  *slog.Logger
}

func NewReportStorage(redis *redisservice.RedisService, mongodb *mongo.Database, cfg *config.Config, logger *slog.Logger) repository.ReportI {
	return &ReportStorage{
		redis:   redis,
		mongodb: mongodb,
		cfg:     cfg,
		logger:  logger,
	}
}

func (s *ReportStorage) GetSpendingReport(ctx context.Context, req *pb.GetSpendingReportRequest) (*pb.SpendingReportResponse, error) {
	s.logger.Info("GetSpendingReport")

	transactionCollection := s.mongodb.Collection("transactions")
	categoryCollection := s.mongodb.Collection("categories")

	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		s.logger.Error("error while parsing start date:", slog.String("err", err.Error()))
		return nil, err
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		s.logger.Error("error while parsing end date:", slog.String("err", err.Error()))
		return nil, err
	}

	filter := bson.M{
		"user_id": req.UserId,
		"date": bson.M{
			"$gte": startDate,
			"$lte": endDate,
		},
		"type": "expense",
	}

	projection := bson.M{
		"category_id": 1,
		"amount":      1,
		"_id":         0,
	}

	var totalSpending float32
	categorySpending := make(map[string]float32)

	cursor, err := transactionCollection.Find(ctx, filter, options.Find().SetProjection(projection))
	if err != nil {
		s.logger.Error("error while querying transactions:", slog.String("err", err.Error()))
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var transaction struct {
			CategoryID primitive.ObjectID `bson:"category_id"`
			Amount     float64            `bson:"amount"`
		}
		if err := cursor.Decode(&transaction); err != nil {
			s.logger.Error("error while decoding transaction:", slog.String("err", err.Error()))
			return nil, err
		}

		filter := bson.M{"_id": transaction.CategoryID}
		projection := bson.M{"name": 1, "_id": 0}

		var categoryName struct {
			Name string `bson:"name"`
		}

		err = categoryCollection.FindOne(ctx, filter, options.FindOne().SetProjection(projection)).Decode(&categoryName)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				s.logger.Error("Category not found", slog.String("category_id", transaction.CategoryID.Hex()))
			} else {
				s.logger.Error("error while querying category:", slog.String("err", err.Error()))
				return nil, err
			}
		}

		if _, exists := categorySpending[categoryName.Name]; !exists {
			categorySpending[categoryName.Name] = 0
		}
		categorySpending[categoryName.Name] += float32(transaction.Amount)

		totalSpending += float32(transaction.Amount)
	}

	if err := cursor.Err(); err != nil {
		s.logger.Error("error while iterating cursor:", slog.String("err", err.Error()))
		return nil, err
	}

	return &pb.SpendingReportResponse{
		TotalSpending:    totalSpending,
		CategorySpending: categorySpending,
	}, nil
}

func (s *ReportStorage) GetIncomeReport(ctx context.Context, req *pb.GetIncomeReportRequest) (*pb.IncomeReportResponse, error) {
	s.logger.Info("GetIncomeReport")

	transactionCollection := s.mongodb.Collection("transactions")
	categoryCollection := s.mongodb.Collection("categories")

	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		s.logger.Error("error while parsing start date:", slog.String("err", err.Error()))
		return nil, err
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		s.logger.Error("error while parsing end date:", slog.String("err", err.Error()))
		return nil, err
	}

	filter := bson.M{
		"user_id": req.UserId,
		"date": bson.M{
			"$gte": startDate,
			"$lte": endDate,
		},
		"type": "income",
	}

	projection := bson.M{
		"category_id": 1,
		"amount":      1,
		"_id":         0,
	}

	var totalIncoming float32
	categoryIncoming := make(map[string]float32)

	cursor, err := transactionCollection.Find(ctx, filter, options.Find().SetProjection(projection))
	if err != nil {
		s.logger.Error("error while querying transactions:", slog.String("err", err.Error()))
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var transaction struct {
			CategoryID primitive.ObjectID `bson:"category_id"`
			Amount     float64            `bson:"amount"`
		}
		if err := cursor.Decode(&transaction); err != nil {
			s.logger.Error("error while decoding transaction:", slog.String("err", err.Error()))
			return nil, err
		}

		filter := bson.M{"_id": transaction.CategoryID}
		projection := bson.M{"name": 1, "_id": 0}

		var categoryName struct {
			Name string `bson:"name"`
		}

		err = categoryCollection.FindOne(ctx, filter, options.FindOne().SetProjection(projection)).Decode(&categoryName)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				s.logger.Error("Category not found", slog.String("category_id", transaction.CategoryID.Hex()))
			} else {
				s.logger.Error("error while querying category:", slog.String("err", err.Error()))
				return nil, err
			}
		}

		if _, exists := categoryIncoming[categoryName.Name]; !exists {
			categoryIncoming[categoryName.Name] = 0
		}
		categoryIncoming[categoryName.Name] += float32(transaction.Amount)

		totalIncoming += float32(transaction.Amount)
	}

	if err := cursor.Err(); err != nil {
		s.logger.Error("error while iterating cursor:", slog.String("err", err.Error()))
		return nil, err
	}

	return &pb.IncomeReportResponse{
		TotalIncome:    totalIncoming,
		CategoryIncome: categoryIncoming,
	}, nil
}

func (s *ReportStorage) GetBudgetPerformanceReport(ctx context.Context, req *pb.GetBudgetPerformanceReportRequest) (*pb.BudgetPerformanceReportResponse, error) {
	s.logger.Info("GetBudgetPerformanceReport")

	transactionCollection := s.mongodb.Collection("transactions")
	budgetCollection := s.mongodb.Collection("budgets")
	categoryCollection := s.mongodb.Collection("categories")

	budgetId, err := primitive.ObjectIDFromHex(req.BudgetId)
	if err != nil {
		return nil, err
	}

	filter1 := bson.D{{Key: "_id", Value: budgetId}}
	projection1 := bson.D{
		{Key: "start_date", Value: 1},
		{Key: "end_date", Value: 1},
		{Key: "amount", Value: 1},
	}

	var budget struct {
		StartDate time.Time `bson:"start_date"`
		EndDate   time.Time `bson:"end_date"`
		Amount    float64   `bson:"amount"`
	}

	err = budgetCollection.FindOne(ctx, filter1, options.FindOne().SetProjection(projection1)).Decode(&budget)
	if err != nil {
		s.logger.Error("error while querying budget:", slog.String("err", err.Error()))
		return nil, err
	}

	filter2 := bson.M{
		"user_id": req.UserId,
		"date": bson.M{
			"$gte": budget.StartDate,
			"$lte": budget.EndDate,
		},
		"type": "expense",
	}

	projection2 := bson.M{
		"category_id": 1,
		"amount":      1,
		"_id":         0,
	}

	var totalSpending float32
	categorySpending := make(map[string]float32)

	cursor, err := transactionCollection.Find(ctx, filter2, options.Find().SetProjection(projection2))
	if err != nil {
		s.logger.Error("error while querying transactions:", slog.String("err", err.Error()))
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var transaction struct {
			CategoryID primitive.ObjectID `bson:"category_id"`
			Amount     float64            `bson:"amount"`
		}
		if err := cursor.Decode(&transaction); err != nil {
			s.logger.Error("error while decoding transaction:", slog.String("err", err.Error()))
			return nil, err
		}

		filter := bson.M{"_id": transaction.CategoryID}
		projection := bson.M{"name": 1, "_id": 0}

		var categoryName struct {
			Name string `bson:"name"`
		}

		err = categoryCollection.FindOne(ctx, filter, options.FindOne().SetProjection(projection)).Decode(&categoryName)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				s.logger.Error("Category not found", slog.String("category_id", transaction.CategoryID.Hex()))
			} else {
				s.logger.Error("error while querying category:", slog.String("err", err.Error()))
				return nil, err
			}
		}

		if _, exists := categorySpending[categoryName.Name]; !exists {
			categorySpending[categoryName.Name] = 0
		}
		categorySpending[categoryName.Name] += float32(transaction.Amount)

		totalSpending += float32(transaction.Amount)
	}

	if err := cursor.Err(); err != nil {
		s.logger.Error("error while iterating cursor:", slog.String("err", err.Error()))
		return nil, err
	}

	return &pb.BudgetPerformanceReportResponse{
		TotalBudget:         float32(budget.Amount),
		TotalSpent:          totalSpending,
		CategoryPerformance: categorySpending,
	}, nil
}

func (s *ReportStorage) GetGoalProgressReport(ctx context.Context, req *pb.GetGoalProgressReportRequest) (*pb.GoalProgressReportResponse, error) {
	s.logger.Info("GetGoalProgressReport")

	transactionCollection := s.mongodb.Collection("transactions")
	goalCollection := s.mongodb.Collection("goals")
	categoryCollection := s.mongodb.Collection("categories")

	goalId, err := primitive.ObjectIDFromHex(req.GoalId)
	if err != nil {
		return nil, err
	}

	filter1 := bson.D{{Key: "_id", Value: goalId}}
	projection1 := bson.D{
		{Key: "created_at", Value: 1},
		{Key: "deadline", Value: 1},
		{Key: "target_amount", Value: 1},
		{Key: "current_amount", Value: 1},
	}

	var goal struct {
		CreatedAt     time.Time `bson:"created_at"`
		Deadline      time.Time `bson:"deadline"`
		TargetAmount  float64   `bson:"target_amount"`
		CurrentAmount float64   `bson:"current_amount"`
	}

	err = goalCollection.FindOne(ctx, filter1, options.FindOne().SetProjection(projection1)).Decode(&goal)
	if err != nil {
		s.logger.Error("error while querying budget:", slog.String("err", err.Error()))
		return nil, err
	}

	filter2 := bson.M{
		"user_id": req.UserId,
		"date": bson.M{
			"$gte": goal.CreatedAt,
			"$lte": goal.Deadline,
		},
		"type": "income",
	}

	projection2 := bson.M{
		"category_id": 1,
		"amount":      1,
		"_id":         0,
	}

	var totalIncoming float32
	categoryIncoming := make(map[string]float32)

	cursor, err := transactionCollection.Find(ctx, filter2, options.Find().SetProjection(projection2))
	if err != nil {
		s.logger.Error("error while querying transactions:", slog.String("err", err.Error()))
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var transaction struct {
			CategoryID primitive.ObjectID `bson:"category_id"`
			Amount     float64            `bson:"amount"`
		}
		if err := cursor.Decode(&transaction); err != nil {
			s.logger.Error("error while decoding transaction:", slog.String("err", err.Error()))
			return nil, err
		}

		filter := bson.M{"_id": transaction.CategoryID}
		projection := bson.M{"name": 1, "_id": 0}

		var categoryName struct {
			Name string `bson:"name"`
		}

		err = categoryCollection.FindOne(ctx, filter, options.FindOne().SetProjection(projection)).Decode(&categoryName)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				s.logger.Error("Category not found", slog.String("category_id", transaction.CategoryID.Hex()))
			} else {
				s.logger.Error("error while querying category:", slog.String("err", err.Error()))
				return nil, err
			}
		}

		if _, exists := categoryIncoming[categoryName.Name]; !exists {
			categoryIncoming[categoryName.Name] = 0
		}
		categoryIncoming[categoryName.Name] += float32(transaction.Amount)

		totalIncoming += float32(transaction.Amount)
	}

	if err := cursor.Err(); err != nil {
		s.logger.Error("error while iterating cursor:", slog.String("err", err.Error()))
		return nil, err
	}

	return &pb.GoalProgressReportResponse{
		TotalProgress:       totalIncoming,
		TargetAmount:        float32(goal.TargetAmount) - float32(goal.CurrentAmount),
		CategoryPerformance: categoryIncoming,
	}, nil
}
