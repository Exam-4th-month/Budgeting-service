package storage

import (
	"log/slog"

	"budgeting-service/internal/items/config"
	"budgeting-service/internal/items/redisservice"
	"budgeting-service/internal/items/repository"

	"go.mongodb.org/mongo-driver/mongo"

	mdb "budgeting-service/internal/items/storage/mongodb"
)

type StrorageI interface {
	Account() repository.AccountI
	Budget() repository.BudgetI
	Category() repository.CategoryI
	Goal() repository.GoalI
	Notification() repository.NotificationI
	Report() repository.ReportI
	Transaction() repository.TransactionI
}

type Storage struct {
	accountRepo      repository.AccountI
	budgetRepo       repository.BudgetI
	categoryRepo     repository.CategoryI
	goalRepo         repository.GoalI
	notificationRepo repository.NotificationI
	reportRepo       repository.ReportI
	transactionRepo  repository.TransactionI
}

func New(redis *redisservice.RedisService, mongodb *mongo.Database, cfg *config.Config, logger *slog.Logger) StrorageI {
	return &Storage{
		accountRepo:      mdb.NewAccountStorage(redis, mongodb, cfg, logger),
		budgetRepo:       mdb.NewBudgetStorage(redis, mongodb, cfg, logger),
		categoryRepo:     mdb.NewCategoryStorage(redis, mongodb, cfg, logger),
		goalRepo:         mdb.NewGoalStorage(redis, mongodb, cfg, logger),
		notificationRepo: mdb.NewNotificationStorage(redis, mongodb, cfg, logger),
		reportRepo:       mdb.NewReportStorage(redis, mongodb, cfg, logger),
		transactionRepo:  mdb.NewTransactionStorage(redis, mongodb, cfg, logger),
	}
}

func (s *Storage) Account() repository.AccountI {
	return s.accountRepo
}

func (s *Storage) Budget() repository.BudgetI {
	return s.budgetRepo
}

func (s *Storage) Category() repository.CategoryI {
	return s.categoryRepo
}

func (s *Storage) Goal() repository.GoalI {
	return s.goalRepo
}

func (s *Storage) Notification() repository.NotificationI {
	return s.notificationRepo
}

func (s *Storage) Report() repository.ReportI {
	return s.reportRepo
}

func (s *Storage) Transaction() repository.TransactionI {
	return s.transactionRepo
}
