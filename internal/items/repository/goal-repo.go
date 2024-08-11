package repository

import (
	pb "budgeting-service/genproto/goal"
	"context"
)

type GoalI interface {
	CreateGoal(ctx context.Context, req *pb.CreateGoalRequest) (*pb.GoalResponse, error)
	GetGoals(ctx context.Context, req *pb.GetGoalsRequest) (*pb.GoalsResponse, error)
	GetGoalById(ctx context.Context, req *pb.GetGoalByIdRequest) (*pb.GoalResponse, error)
	UpdateGoal(ctx context.Context, req *pb.UpdateGoalRequest) (*pb.GoalResponse, error)
	DeleteGoal(ctx context.Context, req *pb.DeleteGoalRequest) (*pb.Empty, error)
}
