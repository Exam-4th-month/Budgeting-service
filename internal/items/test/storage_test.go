package test

import (
	account_pb "budgeting-service/genproto/account"
	budget_pb "budgeting-service/genproto/budget"
	category_pb "budgeting-service/genproto/category"
	goal_pb "budgeting-service/genproto/goal"
	transaction_pb "budgeting-service/genproto/transaction"

	"budgeting-service/internal/items/config"
	"budgeting-service/internal/items/redisservice"
	"budgeting-service/internal/items/storage"
	"budgeting-service/internal/items/storage/mongodb"
	redisCl "budgeting-service/internal/pkg/redis"

	"context"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"log"
	"log/slog"
	"os"
)

func setupStorage() (storage.StrorageI, *mongo.Database) {
	config, err := config.New()
	if err != nil {
		log.Fatalln(err)
	}

	logFile, err := os.OpenFile("test.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()

	logger := slog.New(slog.NewJSONHandler(logFile, nil))

	db, err := mongodb.ConnectDB(config)
	if err != nil {
		logger.Error("error while connecting postgres:", slog.String("err:", err.Error()))
	}

	redis, err := redisCl.NewRedisDB(config)
	if err != nil {
		logger.Error("error while connecting redis:", slog.String("err:", err.Error()))
	}

	return storage.New(
		redisservice.New(redis, logger),
		db,
		config,
		logger,
	), db
}

func TestCreateAccount(t *testing.T) {
	storage, db := setupStorage()
	ctx := context.Background()

	test := account_pb.CreateAccountRequest{
		UserId:   "68819df6-1db1-447a-837e-4f4bd6ec577f",
		Name:     "test",
		Type:     "test",
		Balance:  1000,
		Currency: "UZS",
	}

	res, err := storage.Account().CreateAccount(ctx, &test)
	if err != nil {
		t.Error(err)
	}

	_, err = db.Collection("accounts").DeleteOne(ctx, bson.M{"_id": res.Id})
	if err != nil {
		t.Error(err)
	}
}

func TestCreateBudget(t *testing.T) {
	storage, db := setupStorage()
	ctx := context.Background()

	test := budget_pb.CreateBudgetRequest{
		UserId:     "4ed1de6b-d3de-4811-aac2-7d86bd544659",
		CategoryId: "68819df6-1db1-447a-837e-4f4bd6ec577f",
		Amount:     1000,
		Period:     "monthly",
		StartDate:  "2023-01-01",
		EndDate:    "2023-12-31",
	}

	res, err := storage.Budget().CreateBudget(ctx, &test)
	if err != nil {
		t.Error(err)
	}

	_, err = db.Collection("budgets").DeleteOne(ctx, bson.M{"_id": res.Id})
	if err != nil {
		t.Error(err)
	}
}

func TestCategory(t *testing.T) {
	storage, db := setupStorage()
	ctx := context.Background()

	test := category_pb.CreateCategoryRequest{
		UserId: "4ed1de6b-d3de-4811-aac2-7d86bd544659",
		Name:   "test",
		Type:   "test",
	}

	res, err := storage.Category().CreateCategory(ctx, &test)
	if err != nil {
		t.Error(err)
	}
	_, err = db.Collection("categories").DeleteOne(ctx, bson.M{"_id": res.Id})
	if err != nil {
		t.Error(err)
	}
}

func TestCreateGoal(t *testing.T) {
	storage, db := setupStorage()
	ctx := context.Background()

	test := goal_pb.CreateGoalRequest{
		UserId:        "4ed1de6b-d3de-4811-aac2-7d86bd544659",
		Name:          "test",
		TargetAmount:  1000,
		CurrentAmount: 0,
		Deadline:      "2023-12-31",
		Status:        "active",
	}

	res, err := storage.Goal().CreateGoal(ctx, &test)
	if err != nil {
		t.Error(err)
	}

	_, err = db.Collection("goals").DeleteOne(ctx, bson.M{"_id": res.Id})
	if err != nil {
		t.Error(err)
	}
}

func TestCreateTransaction(t *testing.T) {
	storage, db := setupStorage()
	ctx := context.Background()

	test := transaction_pb.CreateTransactionRequest{
		UserId:      "4ed1de6b-d3de-4811-aac2-7d86bd544659",
		AccountId:   "68819df6-1db1-447a-837e-4f4bd6ec577f",
		CategoryId:  "68819df6-1db1-447a-837e-4f4bd6ec577f",
		Amount:      1000,
		Type:        "test",
		Description: "test",
		Date:        "2023-12-31",
	}

	res, err := storage.Transaction().CreateTransaction(ctx, &test)
	if err != nil {
		t.Error(err)
	}

	_, err = db.Collection("transactions").DeleteOne(ctx, bson.M{"_id": res.Id})
	if err != nil {
		t.Error(err)
	}
}
