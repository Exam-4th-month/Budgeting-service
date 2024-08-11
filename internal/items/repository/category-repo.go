package repository

import (
	pb "budgeting-service/genproto/category"
	"context"
)

type CategoryI interface {
	CreateCategory(ctx context.Context, req *pb.CreateCategoryRequest) (*pb.CategoryResponse, error)
	GetCategories(ctx context.Context, req *pb.GetCategoriesRequest) (*pb.CategoriesResponse, error)
	GetCategoryById(ctx context.Context, req *pb.GetCategoryByIdRequest) (*pb.CategoryResponse, error)
	UpdateCategory(ctx context.Context, req *pb.UpdateCategoryRequest) (*pb.CategoryResponse, error)
	DeleteCategory(ctx context.Context, req *pb.DeleteCategoryRequest) (*pb.Empty, error)
}
