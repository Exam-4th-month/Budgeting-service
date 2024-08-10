package service

import (
	pb "budgeting-service/genproto/goal"
	"budgeting-service/internal/items/repository"
)

type GoalService struct {
	pb.UnimplementedGoalServiceServer
	goalstorage repository.GoalI
}

func NewGoalService(goalstorage repository.GoalI) *GoalService {
	return &GoalService{
		goalstorage: goalstorage,
	}
}
