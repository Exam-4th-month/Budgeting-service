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
		AccountService:      NewAccountService(storage.Account()),
		BudgetService:       NewBudgetService(storage.Budget()),
		CategoryService:     NewCategoryService(storage.Category()),
		GoalService:         NewGoalService(storage.Goal()),
		NotificationService: NewNotificationService(storage.Notification()),
		ReportService:       NewReportService(storage.Report()),
		TransactionService:  NewTransactionService(storage.Transaction()),
	}

}
