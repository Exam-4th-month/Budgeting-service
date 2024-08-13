package service

import (
	pb "budgeting-service/genproto/report"
	"budgeting-service/internal/items/repository"
	"context"
	"log/slog"
)

type ReportService struct {
	pb.UnimplementedReportServiceServer
	reportstorage repository.ReportI
	logger        *slog.Logger
}

func NewReportService(reportstorage repository.ReportI, logger *slog.Logger) *ReportService {
	return &ReportService{
		reportstorage: reportstorage,
		logger:        logger,
	}
}

func (s *ReportService) GetSpendingReport(ctx context.Context, req *pb.GetSpendingReportRequest) (*pb.SpendingReportResponse, error) {
	s.logger.Info("GetSpendingReport")
	return s.reportstorage.GetSpendingReport(ctx, req)
}

func (s *ReportService) GetIncomeReport(ctx context.Context, req *pb.GetIncomeReportRequest) (*pb.IncomeReportResponse, error) {
	s.logger.Info("GetIncomeReport")
	return s.reportstorage.GetIncomeReport(ctx, req)
}
