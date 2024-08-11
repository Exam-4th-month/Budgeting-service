package service

import (
	pb "budgeting-service/genproto/goal"
	"budgeting-service/internal/items/repository"
	"context"
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

func (s *GoalService) CreateGoal(ctx context.Context, req *pb.CreateGoalRequest) (*pb.GoalResponse, error) {
	s.logger.Info("CreateGoal", "req", req)
	return s.goalstorage.CreateGoal(ctx, req)
}

func (s *GoalService) GetGoals(ctx context.Context, req *pb.GetGoalsRequest) (*pb.GoalsResponse, error) {
	s.logger.Info("GetGoals", "req", req)
	return s.goalstorage.GetGoals(ctx, req)
}

func (s *GoalService) GetGoalById(ctx context.Context, req *pb.GetGoalByIdRequest) (*pb.GoalResponse, error) {
	s.logger.Info("GetGoalById", "req", req)
	return s.goalstorage.GetGoalById(ctx, req)
}

func (s *GoalService) UpdateGoal(ctx context.Context, req *pb.UpdateGoalRequest) (*pb.GoalResponse, error) {
	s.logger.Info("UpdateGoal", "req", req)
	return s.goalstorage.UpdateGoal(ctx, req)
}

func (s *GoalService) DeleteGoal(ctx context.Context, req *pb.DeleteGoalRequest) (*pb.Empty, error) {
	s.logger.Info("DeleteGoal", "req", req)
	return s.goalstorage.DeleteGoal(ctx, req)
}
