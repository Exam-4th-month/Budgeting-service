package service

import (
	pb "budgeting-service/genproto/category"
	"budgeting-service/internal/items/repository"
)

type CategoryService struct {
	pb.UnimplementedCategoryServiceServer
	categorystorage repository.CategoryI
}

func NewCategoryService(categorystorage repository.CategoryI) *CategoryService {
	return &CategoryService{
		categorystorage: categorystorage,
	}
}
