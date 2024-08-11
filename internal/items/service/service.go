package service

import (
	"log/slog"

	"budgeting-service/internal/items/storage"
)

type Service struct {
	AccountService      *AccountService
	BudgetService       *BudgetService
	CategoryService     *CategoryService
	GoalService         *GoalService
	NotificationService *NotificationService
	ReportService       *ReportService
	TransactionService  *TransactionService
}

func New(storage storage.StrorageI, logger *slog.Logger) *Service {
	return &Service{
		AccountService:      NewAccountService(storage.Account(), logger),
		BudgetService:       NewBudgetService(storage.Budget(), logger),
		CategoryService:     NewCategoryService(storage.Category(), logger),
		GoalService:         NewGoalService(storage.Goal(), logger),
		NotificationService: NewNotificationService(storage.Notification(), logger),
		ReportService:       NewReportService(storage.Report(), logger),
		TransactionService:  NewTransactionService(storage.Transaction(), logger),
	}

}
