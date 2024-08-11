package repository

import (
	pb "budgeting-service/genproto/budget"
	"context"
)

type BudgetI interface {
	CreateBudget(ctx context.Context, req *pb.CreateBudgetRequest) (*pb.BudgetResponse, error)
	GetBudgets(ctx context.Context, req *pb.GetBudgetsRequest) (*pb.BudgetsResponse, error)
	GetBudgetById(ctx context.Context, req *pb.GetBudgetByIdRequest) (*pb.BudgetResponse, error)
	UpdateBudget(ctx context.Context, req *pb.UpdateBudgetRequest) (*pb.BudgetResponse, error)
	DeleteBudget(ctx context.Context, req *pb.DeleteBudgetRequest) (*pb.Empty, error)
}
