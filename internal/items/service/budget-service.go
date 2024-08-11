package service

import (
	pb "budgeting-service/genproto/budget"
	"budgeting-service/internal/items/repository"
	"context"
	"log/slog"
)

type BudgetService struct {
	pb.UnimplementedBudgetServiceServer
	budgetstorage repository.BudgetI
	logger        *slog.Logger
}

func NewBudgetService(budgetstorage repository.BudgetI, logger *slog.Logger) *BudgetService {
	return &BudgetService{
		budgetstorage: budgetstorage,
		logger:        logger,
	}
}

func (s *BudgetService) CreateBudget(ctx context.Context, req *pb.CreateBudgetRequest) (*pb.BudgetResponse, error) {
	s.logger.Info("CreateBudget", "req", req)
	return s.budgetstorage.CreateBudget(ctx, req)
}
func (s *BudgetService) GetBudgets(ctx context.Context, req *pb.GetBudgetsRequest) (*pb.BudgetsResponse, error) {
	s.logger.Info("GetBudgets", "req", req)
	return s.budgetstorage.GetBudgets(ctx, req)
}
func (s *BudgetService) GetBudgetById(ctx context.Context, req *pb.GetBudgetByIdRequest) (*pb.BudgetResponse, error) {
	s.logger.Info("GetBudgetById", "req", req)
	return s.budgetstorage.GetBudgetById(ctx, req)
}
func (s *BudgetService) UpdateBudget(ctx context.Context, req *pb.UpdateBudgetRequest) (*pb.BudgetResponse, error) {
	s.logger.Info("UpdateBudget", "req", req)
	return s.budgetstorage.UpdateBudget(ctx, req)
}
func (s *BudgetService) DeleteBudget(ctx context.Context, req *pb.DeleteBudgetRequest) (*pb.Empty, error) {
	s.logger.Info("DeleteBudget", "req", req)
	return s.budgetstorage.DeleteBudget(ctx, req)
}
