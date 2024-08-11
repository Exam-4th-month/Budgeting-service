package mongodb

import (
	"budgeting-service/internal/items/config"
	"budgeting-service/internal/items/redisservice"
	"budgeting-service/internal/items/repository"
	"context"
	"time"

	pb "budgeting-service/genproto/transaction"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"log/slog"
)

type TransactionStorage struct {
	redis   *redisservice.RedisService
	mongodb *mongo.Database
	cfg     *config.Config
	logger  *slog.Logger
}

func NewTransactionStorage(redis *redisservice.RedisService, mongodb *mongo.Database, cfg *config.Config, logger *slog.Logger) repository.TransactionI {
	return &TransactionStorage{
		redis:   redis,
		mongodb: mongodb,
		cfg:     cfg,
		logger:  logger,
	}
}

func (s *TransactionStorage) CreateTransaction(ctx context.Context, req *pb.CreateTransactionRequest) (*pb.TransactionResponse, error) {
	s.logger.Info("CreateTransaction", slog.Any("req", req))

	transactionCollection := s.mongodb.Collection("transactions")
	created_at := time.Now()

	transactionDoc := bson.D{
		{Key: "user_id", Value: req.UserId},
		{Key: "account_id", Value: req.AccountId},
		{Key: "category_id", Value: req.CategoryId},
		{Key: "amount", Value: req.Amount},
		{Key: "type", Value: req.Type},
		{Key: "description", Value: req.Description},
		{Key: "date", Value: req.Date},
		{Key: "created_at", Value: created_at},
	}

	res, err := transactionCollection.InsertOne(ctx, transactionDoc)
	if err != nil {
		s.logger.Error("Error while creating transaction", slog.Any("error", err))
		return nil, err
	}

	transactionID := res.InsertedID.(primitive.ObjectID).Hex()

	return &pb.TransactionResponse{
		Id:          transactionID,
		UserId:      req.UserId,
		AccountId:   req.AccountId,
		CategoryId:  req.CategoryId,
		Amount:      req.Amount,
		Type:        req.Type,
		Description: req.Description,
		Date:        req.Date,
		CreatedAt:   created_at.String(),
	}, nil
}

func (s *TransactionStorage) GetTransactions(ctx context.Context, req *pb.GetTransactionsRequest) (*pb.TransactionsResponse, error) {
	s.logger.Info("GetTransactions", slog.Any("req", req))

	transactionCollection := s.mongodb.Collection("transactions")

	filter := bson.D{}
	if req.UserId != "" {
		filter = append(filter, bson.E{Key: "user_id", Value: req.UserId})
	}
	if req.AccountId != "" {
		filter = append(filter, bson.E{Key: "account_id", Value: req.AccountId})
	}
	if req.CategoryId != "" {
		filter = append(filter, bson.E{Key: "category_id", Value: req.CategoryId})
	}

	cursor, err := transactionCollection.Find(ctx, filter)
	if err != nil {
		s.logger.Error("Error while fetching transactions", slog.Any("error", err))
		return nil, err
	}
	defer cursor.Close(ctx)

	var transactions []*pb.TransactionResponse
	for cursor.Next(ctx) {
		var transaction bson.M
		if err := cursor.Decode(&transaction); err != nil {
			s.logger.Error("Error while decoding transaction", slog.Any("error", err))
			return nil, err
		}

		transactions = append(transactions, &pb.TransactionResponse{
			Id:          transaction["_id"].(primitive.ObjectID).Hex(),
			UserId:      transaction["user_id"].(string),
			AccountId:   transaction["account_id"].(string),
			CategoryId:  transaction["category_id"].(string),
			Amount:      float32(transaction["amount"].(float64)),
			Type:        transaction["type"].(string),
			Description: transaction["description"].(string),
			Date:        transaction["date"].(string),
			CreatedAt:   transaction["created_at"].(primitive.DateTime).Time().String(),
			UpdatedAt:   transaction["updated_at"].(primitive.DateTime).Time().String(),
		})
	}

	if err := cursor.Err(); err != nil {
		s.logger.Error("Cursor error", slog.Any("error", err))
		return nil, err
	}

	return &pb.TransactionsResponse{Transactions: transactions}, nil
}

func (s *TransactionStorage) GetTransactionById(ctx context.Context, req *pb.GetTransactionByIdRequest) (*pb.TransactionResponse, error) {
	s.logger.Info("GetTransactionById", slog.String("id", req.Id))

	transactionCollection := s.mongodb.Collection("transactions")

	objID, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		s.logger.Error("Invalid ObjectID", slog.Any("error", err))
		return nil, err
	}

	filter := bson.D{{Key: "_id", Value: objID}}

	var transaction bson.M
	err = transactionCollection.FindOne(ctx, filter).Decode(&transaction)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			s.logger.Info("Transaction not found", slog.String("id", req.Id))
			return nil, nil
		}
		s.logger.Error("Error finding transaction", slog.Any("error", err))
		return nil, err
	}

	return &pb.TransactionResponse{
		Id:          transaction["_id"].(primitive.ObjectID).Hex(),
		UserId:      transaction["user_id"].(string),
		AccountId:   transaction["account_id"].(string),
		CategoryId:  transaction["category_id"].(string),
		Amount:      float32(transaction["amount"].(float64)),
		Type:        transaction["type"].(string),
		Description: transaction["description"].(string),
		Date:        transaction["date"].(string),
		CreatedAt:   transaction["created_at"].(primitive.DateTime).Time().String(),
		UpdatedAt:   transaction["updated_at"].(primitive.DateTime).Time().String(),
	}, nil
}

func (s *TransactionStorage) UpdateTransaction(ctx context.Context, req *pb.UpdateTransactionRequest) (*pb.TransactionResponse, error) {
	s.logger.Info("UpdateTransaction", slog.Any("req", req))

	transactionCollection := s.mongodb.Collection("transactions")

	objID, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		s.logger.Error("Invalid ObjectID", slog.Any("error", err))
		return nil, err
	}

	filter := bson.D{{Key: "_id", Value: objID}}
	updateFields := bson.D{}
	if req.Amount != 0 {
		updateFields = append(updateFields, bson.E{Key: "amount", Value: req.Amount})
	}
	if req.Type != "" {
		updateFields = append(updateFields, bson.E{Key: "type", Value: req.Type})
	}
	if req.Description != "" {
		updateFields = append(updateFields, bson.E{Key: "description", Value: req.Description})
	}
	if req.Date != "" {
		updateFields = append(updateFields, bson.E{Key: "date", Value: req.Date})
	}
	if len(updateFields) > 0 {
		updateFields = append(updateFields, bson.E{Key: "updated_at", Value: time.Now()})
	}

	if len(updateFields) == 0 {
		s.logger.Info("No fields to update")
		return nil, nil
	}

	update := bson.D{{Key: "$set", Value: updateFields}}

	res := transactionCollection.FindOneAndUpdate(ctx, filter, update, options.FindOneAndUpdate().SetReturnDocument(options.After))
	if res.Err() != nil {
		if res.Err() == mongo.ErrNoDocuments {
			s.logger.Info("Transaction not found", slog.String("id", req.Id))
			return nil, nil
		}
		s.logger.Error("Error updating transaction", slog.Any("error", res.Err()))
		return nil, res.Err()
	}

	var updatedTransaction bson.M
	if err = res.Decode(&updatedTransaction); err != nil {
		s.logger.Error("Error decoding updated transaction", slog.Any("error", err))
		return nil, err
	}

	return &pb.TransactionResponse{
		Id:          updatedTransaction["_id"].(primitive.ObjectID).Hex(),
		UserId:      updatedTransaction["user_id"].(string),
		AccountId:   updatedTransaction["account_id"].(string),
		CategoryId:  updatedTransaction["category_id"].(string),
		Amount:      float32(updatedTransaction["amount"].(float64)),
		Type:        updatedTransaction["type"].(string),
		Description: updatedTransaction["description"].(string),
		Date:        updatedTransaction["date"].(string),
		CreatedAt:   updatedTransaction["created_at"].(primitive.DateTime).Time().String(),
		UpdatedAt:   updatedTransaction["updated_at"].(primitive.DateTime).Time().String(),
	}, nil
}

func (s *TransactionStorage) DeleteTransaction(ctx context.Context, req *pb.DeleteTransactionRequest) (*pb.Empty, error) {
	s.logger.Info("DeleteTransaction", slog.String("id", req.Id))

	transactionCollection := s.mongodb.Collection("transactions")

	objID, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		s.logger.Error("Invalid ObjectID", slog.Any("error", err))
		return nil, err
	}

	filter := bson.D{{Key: "_id", Value: objID}}
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "deleted_at", Value: time.Now()},
		}},
	}

	_, err = transactionCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		s.logger.Error("Error deleting transaction", slog.Any("error", err))
		return nil, err
	}

	return &pb.Empty{}, nil
}
