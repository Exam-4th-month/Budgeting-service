package service

import (
	pb "budgeting-service/genproto/budget"
	"budgeting-service/internal/items/repository"
)

type BudgetService struct {
	pb.UnimplementedBudgetServiceServer
	budgetstorage repository.BudgetI
}

func NewBudgetService(budgetstorage repository.BudgetI) *BudgetService {
	return &BudgetService{
		budgetstorage: budgetstorage,
	}
}
