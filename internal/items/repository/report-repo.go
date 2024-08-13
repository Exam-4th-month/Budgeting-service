package repository

import (
	pb "budgeting-service/genproto/report"
	"context"
)

type ReportI interface {
	GetSpendingReport(ctx context.Context, req *pb.GetSpendingReportRequest) (*pb.SpendingReportResponse, error)
	GetIncomeReport(ctx context.Context, req *pb.GetIncomeReportRequest) (*pb.IncomeReportResponse, error)
	GetBudgetPerformanceReport(ctx context.Context, req *pb.GetBudgetPerformanceReportRequest) (*pb.BudgetPerformanceReportResponse, error)
	GetGoalProgressReport(ctx context.Context, req *pb.GetGoalProgressReportRequest) (*pb.GoalProgressReportResponse, error)
}
