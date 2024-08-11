package service

import (
	pb "budgeting-service/genproto/report"
	"budgeting-service/internal/items/repository"
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
