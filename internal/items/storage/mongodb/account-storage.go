package mongodb

import (
	pb "budgeting-service/genproto/account"
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

type AccountStorage struct {
	redis   *redisservice.RedisService
	mongodb *mongo.Database
	cfg     *config.Config
	logger  *slog.Logger
}

func NewAccountStorage(redis *redisservice.RedisService, mongodb *mongo.Database, cfg *config.Config, logger *slog.Logger) repository.AccountI {
	return &AccountStorage{
		redis:   redis,
		mongodb: mongodb,
		cfg:     cfg,
		logger:  logger,
	}
}

func (s *AccountStorage) CreateAccount(ctx context.Context, req *pb.CreateAccountRequest) (*pb.AccountResponse, error) {
	s.logger.Info("CreateAccount", "req", req)
	accountCollection := s.mongodb.Collection("accounts")
	created_at := time.Now()

	accountDoc := bson.D{
		{Key: "user_id", Value: req.UserId},
		{Key: "name", Value: req.Name},
		{Key: "type", Value: req.Type},
		{Key: "balance", Value: req.Balance},
		{Key: "currency", Value: req.Currency},
		{Key: "created_at", Value: time.Now()},
	}

	// Hujjatni MongoDB ga kiriting
	res, err := accountCollection.InsertOne(ctx, accountDoc)
	if err != nil {
		s.logger.Error("error while inserting account", slog.Any("error", err))
		return nil, err
	}

	// Javob qaytarish (InsertOne natijasidan _id ni olish)
	accountID := res.InsertedID.(primitive.ObjectID)

	return &pb.AccountResponse{
		Id:        accountID.Hex(),
		UserId:    req.UserId,
		Name:      req.Name,
		Type:      req.Type,
		Balance:   req.Balance,
		Currency:  req.Currency,
		CreatedAt: created_at.String(),
	}, nil
}

func (s *AccountStorage) GetAccounts(ctx context.Context, req *pb.GetAccountsRequest) (*pb.AccountsResponse, error) {
	s.logger.Info("GetAccounts", "req", req)
	accountCollection := s.mongodb.Collection("accounts")

	// Foydalanuvchining barcha hisoblarini filtr qilish
	filter := bson.D{{Key: "user_id", Value: req.UserId}}

	cursor, err := accountCollection.Find(ctx, filter)
	if err != nil {
		s.logger.Error(err.Error())
		return nil, err
	}
	defer cursor.Close(ctx)

	var accounts []*pb.AccountResponse
	for cursor.Next(ctx) {
		var account bson.M
		if err = cursor.Decode(&account); err != nil {
			s.logger.Error(err.Error())
			return nil, err
		}

		accounts = append(accounts, &pb.AccountResponse{
			Id:        account["_id"].(primitive.ObjectID).Hex(),
			UserId:    account["user_id"].(string),
			Name:      account["name"].(string),
			Type:      account["type"].(string),
			Balance:   account["balance"].(float32),
			Currency:  account["currency"].(string),
			CreatedAt: account["created_at"].(primitive.DateTime).Time().String(),
			UpdatedAt: account["updated_at"].(primitive.DateTime).Time().String(),
		})
	}

	if err := cursor.Err(); err != nil {
		s.logger.Error(err.Error())
		return nil, err
	}

	return &pb.AccountsResponse{Accounts: accounts}, nil
}

func (s *AccountStorage) GetAccountById(ctx context.Context, req *pb.GetAccountByIdRequest) (*pb.AccountResponse, error) {
	s.logger.Info("GetAccountById", "req: ",req.Id)
	accountCollection := s.mongodb.Collection("accounts")

	objID, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		s.logger.Error(err.Error())
		return nil, err
	}

	// Hisobni ID bo'yicha qidirish
	filter := bson.D{{Key: "_id", Value: objID}}

	var account bson.M
	err = accountCollection.FindOne(ctx, filter).Decode(&account)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			s.logger.Error(err.Error())
			return nil, nil // Hisob topilmadi
		}
		s.logger.Error(err.Error())
		return nil, err
	}

	return &pb.AccountResponse{
		Id:        account["_id"].(primitive.ObjectID).Hex(),
		UserId:    account["user_id"].(string),
		Name:      account["name"].(string),
		Type:      account["type"].(string),
		Balance:   account["balance"].(float32),
		Currency:  account["currency"].(string),
		CreatedAt: account["created_at"].(primitive.DateTime).Time().String(),
		UpdatedAt: account["updated_at"].(primitive.DateTime).Time().String(),
	}, nil
}

func (s *AccountStorage) UpdateAccount(ctx context.Context, req *pb.UpdateAccountRequest) (*pb.AccountResponse, error) {
	accountCollection := s.mongodb.Collection("accounts")

	objID, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		s.logger.Error(err.Error())
		return nil, err
	}

	// Hisobni yangilash uchun ID va `deleted_at` bo'sh bo'lishi filtr
	filter := bson.D{
		{Key: "_id", Value: objID},
		{Key: "deleted_at", Value: bson.D{{Key: "$exists", Value: false}}},
	}

	// Yangilanish uchun maydonlarni dinamik ravishda qo'shish
	updateFields := bson.D{}
	if req.Name != "" {
		updateFields = append(updateFields, bson.E{Key: "name", Value: req.Name})
	}
	if req.Type != "" {
		updateFields = append(updateFields, bson.E{Key: "type", Value: req.Type})
	}
	if req.Balance != 0 {
		updateFields = append(updateFields, bson.E{Key: "balance", Value: req.Balance})
	}
	if req.Currency != "" {
		updateFields = append(updateFields, bson.E{Key: "currency", Value: req.Currency})
	}
	if len(updateFields) > 0 {
		updateFields = append(updateFields, bson.E{Key: "updated_at", Value: time.Now()})
	}

	if len(updateFields) == 0 {
		// Yangilanish uchun hech qanday maydon mavjud emas
		s.logger.Info("No fields to update")
		return nil, nil
	}

	update := bson.D{{Key: "$set", Value: updateFields}}

	res := accountCollection.FindOneAndUpdate(ctx, filter, update, options.FindOneAndUpdate().SetReturnDocument(options.After))
	if res.Err() != nil {
		if res.Err() == mongo.ErrNoDocuments {
			s.logger.Error(res.Err().Error())
			return nil, nil // Hisob topilmadi
		}
		s.logger.Error(res.Err().Error())
		return nil, res.Err()
	}

	var updatedAccount bson.M
	if err = res.Decode(&updatedAccount); err != nil {
		s.logger.Error(err.Error())
		return nil, err
	}

	return &pb.AccountResponse{
		Id:        updatedAccount["_id"].(primitive.ObjectID).Hex(),
		UserId:    updatedAccount["user_id"].(string),
		Name:      updatedAccount["name"].(string),
		Type:      updatedAccount["type"].(string),
		Balance:   updatedAccount["balance"].(float32),
		Currency:  updatedAccount["currency"].(string),
		CreatedAt: updatedAccount["created_at"].(primitive.DateTime).Time().String(),
		UpdatedAt: updatedAccount["updated_at"].(primitive.DateTime).Time().String(),
	}, nil
}

func (s *AccountStorage) DeleteAccount(ctx context.Context, req *pb.DeleteAccountRequest) (*pb.Empty, error) {
	accountCollection := s.mongodb.Collection("accounts")

	objID, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		s.logger.Error(err.Error())
		return nil, err
	}

	// Hisobni ID bo'yicha topish va `deleted_at` maydonini yangilash
	filter := bson.D{{Key: "_id", Value: objID}}
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "deleted_at", Value: time.Now()},
		}},
	}

	_, err = accountCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		s.logger.Error(err.Error())
		return nil, err
	}

	return &pb.Empty{}, nil
}
