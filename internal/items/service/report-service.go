package service

import (
	pb "budgeting-service/genproto/report"
	"budgeting-service/internal/items/repository"
)

type ReportService struct {
	pb.UnimplementedReportServiceServer
	reportstorage repository.ReportI
}

func NewReportService(reportstorage repository.ReportI) *ReportService {
	return &ReportService{
		reportstorage: reportstorage,
	}
}
