package mongodb

import (
	"budgeting-service/internal/items/config"
	"budgeting-service/internal/items/repository"
	"context"
	"time"

	pb "budgeting-service/genproto/category"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"log/slog"
)

type CategoryStorage struct {
	mongodb *mongo.Database
	cfg     *config.Config
	logger  *slog.Logger
}

func NewCategoryStorage(mongodb *mongo.Database, cfg *config.Config, logger *slog.Logger) repository.CategoryI {
	return &CategoryStorage{
		mongodb: mongodb,
		cfg:     cfg,
		logger:  logger,
	}
}

func (s *CategoryStorage) CreateCategory(ctx context.Context, req *pb.CreateCategoryRequest) (*pb.CategoryResponse, error) {
	s.logger.Info("CreateCategory", slog.String("req", req.String()))

	categoryCollection := s.mongodb.Collection("categories")
	created_at := time.Now()

	categoryDoc := bson.D{
		{Key: "user_id", Value: req.UserId},
		{Key: "name", Value: req.Name},
		{Key: "type", Value: req.Type},
		{Key: "created_at", Value: created_at},
		{Key: "updated_at", Value: created_at},
		{Key: "deleted_at", Value: nil},
	}

	res, err := categoryCollection.InsertOne(ctx, categoryDoc)
	if err != nil {
		s.logger.Error(err.Error())
		return nil, err
	}

	categoryID := res.InsertedID.(primitive.ObjectID).Hex()

	return &pb.CategoryResponse{
		Id:        categoryID,
		UserId:    req.UserId,
		Name:      req.Name,
		Type:      req.Type,
		CreatedAt: created_at.String(),
	}, nil
}

func (s *CategoryStorage) GetCategories(ctx context.Context, req *pb.GetCategoriesRequest) (*pb.CategoriesResponse, error) {
	s.logger.Info("GetCategories", slog.String("req", req.String()))
	categoryCollection := s.mongodb.Collection("categories")

	filter := bson.D{{Key: "user_id", Value: req.UserId}}

	cursor, err := categoryCollection.Find(ctx, filter)
	if err != nil {
		s.logger.Error(err.Error())
		return nil, err
	}
	defer cursor.Close(ctx)

	var categories []*pb.CategoryResponse
	for cursor.Next(ctx) {
		var category bson.M
		if err := cursor.Decode(&category); err != nil {
			s.logger.Error(err.Error())
			return nil, err
		}

		categories = append(categories, &pb.CategoryResponse{
			Id:        category["_id"].(primitive.ObjectID).Hex(),
			UserId:    category["user_id"].(string),
			Name:      category["name"].(string),
			Type:      category["type"].(string),
			CreatedAt: category["created_at"].(primitive.DateTime).Time().String(),
			UpdatedAt: category["updated_at"].(primitive.DateTime).Time().String(),
		})
	}

	if err := cursor.Err(); err != nil {
		s.logger.Error(err.Error())
		return nil, err
	}

	return &pb.CategoriesResponse{Categories: categories}, nil
}

func (s *CategoryStorage) GetCategoryById(ctx context.Context, req *pb.GetCategoryByIdRequest) (*pb.CategoryResponse, error) {
	s.logger.Info("GetCategoryById", slog.String("req", req.Id))
	categoryCollection := s.mongodb.Collection("categories")

	objID, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		s.logger.Error(err.Error())
		return nil, err
	}

	filter := bson.D{{Key: "_id", Value: objID}}

	var category bson.M
	err = categoryCollection.FindOne(ctx, filter).Decode(&category)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			s.logger.Error(err.Error())
			return nil, nil
		}
		s.logger.Error(err.Error())
		return nil, err
	}

	return &pb.CategoryResponse{
		Id:        category["_id"].(primitive.ObjectID).Hex(),
		UserId:    category["user_id"].(string),
		Name:      category["name"].(string),
		Type:      category["type"].(string),
		CreatedAt: category["created_at"].(primitive.DateTime).Time().String(),
		UpdatedAt: category["updated_at"].(primitive.DateTime).Time().String(),
	}, nil
}

func (s *CategoryStorage) UpdateCategory(ctx context.Context, req *pb.UpdateCategoryRequest) (*pb.CategoryResponse, error) {
	categoryCollection := s.mongodb.Collection("categories")

	objID, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		s.logger.Error(err.Error())
		return nil, err
	}

	filter := bson.D{
		{Key: "_id", Value: objID},
		{Key: "deleted_at", Value: bson.D{{Key: "$eq", Value: nil}}},
	}

	updateFields := bson.D{}
	if req.Name != "" {
		updateFields = append(updateFields, bson.E{Key: "name", Value: req.Name})
	}
	if req.Type != "" {
		updateFields = append(updateFields, bson.E{Key: "type", Value: req.Type})
	}
	if len(updateFields) > 0 {
		updateFields = append(updateFields, bson.E{Key: "updated_at", Value: time.Now()})
	}

	if len(updateFields) == 0 {
		s.logger.Info("No fields to update")
		return nil, nil
	}

	update := bson.D{{Key: "$set", Value: updateFields}}

	res := categoryCollection.FindOneAndUpdate(ctx, filter, update, options.FindOneAndUpdate().SetReturnDocument(options.After))
	if res.Err() != nil {
		if res.Err() == mongo.ErrNoDocuments {
			s.logger.Error(res.Err().Error())
			return nil, nil
		}
		s.logger.Error(res.Err().Error())
		return nil, res.Err()
	}

	var updatedCategory bson.M
	if err = res.Decode(&updatedCategory); err != nil {
		s.logger.Error(err.Error())
		return nil, err
	}

	return &pb.CategoryResponse{
		Id:        updatedCategory["_id"].(primitive.ObjectID).Hex(),
		UserId:    updatedCategory["user_id"].(string),
		Name:      updatedCategory["name"].(string),
		Type:      updatedCategory["type"].(string),
		CreatedAt: updatedCategory["created_at"].(primitive.DateTime).Time().String(),
		UpdatedAt: updatedCategory["updated_at"].(primitive.DateTime).Time().String(),
	}, nil
}

func (s *CategoryStorage) DeleteCategory(ctx context.Context, req *pb.DeleteCategoryRequest) (*pb.Empty, error) {
	categoryCollection := s.mongodb.Collection("categories")

	objID, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		s.logger.Error(err.Error())
		return nil, err
	}

	filter := bson.D{{Key: "_id", Value: objID}}
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "deleted_at", Value: time.Now()},
		}},
	}

	_, err = categoryCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		s.logger.Error(err.Error())
		return nil, err
	}

	return &pb.Empty{}, nil
}
