package service

import (
	pb "budgeting-service/genproto/category"
	"budgeting-service/internal/items/repository"
	"context"
	"log/slog"
)

type CategoryService struct {
	pb.UnimplementedCategoryServiceServer
	categorystorage repository.CategoryI
	logger          *slog.Logger
}

func NewCategoryService(categorystorage repository.CategoryI, logger *slog.Logger) *CategoryService {
	return &CategoryService{
		categorystorage: categorystorage,
		logger:          logger,
	}
}

func (s *CategoryService) CreateCategory(ctx context.Context, req *pb.CreateCategoryRequest) (*pb.CategoryResponse, error) {
	s.logger.Info("CreateCategory", "req", req)
	return s.categorystorage.CreateCategory(ctx, req)
}

func (s *CategoryService) GetCategories(ctx context.Context, req *pb.GetCategoriesRequest) (*pb.CategoriesResponse, error) {
	s.logger.Info("GetCategories", "req", req)
	return s.categorystorage.GetCategories(ctx, req)
}

func (s *CategoryService) GetCategoryById(ctx context.Context, req *pb.GetCategoryByIdRequest) (*pb.CategoryResponse, error) {
	s.logger.Info("GetCategoryById", "req", req)
	return s.categorystorage.GetCategoryById(ctx, req)
}

func (s *CategoryService) UpdateCategory(ctx context.Context, req *pb.UpdateCategoryRequest) (*pb.CategoryResponse, error) {
	s.logger.Info("UpdateCategory", "req", req)
	return s.categorystorage.UpdateCategory(ctx, req)
}

func (s *CategoryService) DeleteCategory(ctx context.Context, req *pb.DeleteCategoryRequest) (*pb.Empty, error) {
	s.logger.Info("DeleteCategory", "req", req)
	return s.categorystorage.DeleteCategory(ctx, req)
}
