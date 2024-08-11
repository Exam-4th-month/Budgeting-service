package service

import (
	pb "budgeting-service/genproto/goal"
	"budgeting-service/internal/items/repository"
	"log/slog"
)

type GoalService struct {
	pb.UnimplementedGoalServiceServer
	goalstorage repository.GoalI
	logger      *slog.Logger
}

func NewGoalService(goalstorage repository.GoalI, logger *slog.Logger) *GoalService {
	return &GoalService{
		goalstorage: goalstorage,
		logger:      logger,
	}
}
