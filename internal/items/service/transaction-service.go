package service

import (
	pb "budgeting-service/genproto/transaction"
	"budgeting-service/internal/items/repository"
	"context"
	"log/slog"
)

type TransactionService struct {
	pb.UnimplementedTransactionServiceServer
	transactionstorage repository.TransactionI
	logger             *slog.Logger
}

func NewTransactionService(transactionstorage repository.TransactionI, logger *slog.Logger) *TransactionService {
	return &TransactionService{
		transactionstorage: transactionstorage,
		logger:             logger,
	}
}

func (s *TransactionService) CreateTransaction(ctx context.Context, req *pb.CreateTransactionRequest) (*pb.TransactionResponse, error) {
	s.logger.Info("CreateTransaction", slog.Any("req", req))
	return s.transactionstorage.CreateTransaction(ctx, req)
}

func (s *TransactionService) GetTransactions(ctx context.Context, req *pb.GetTransactionsRequest) (*pb.TransactionsResponse, error) {
	s.logger.Info("GetTransactions", slog.Any("req", req))
	return s.transactionstorage.GetTransactions(ctx, req)
}

func (s *TransactionService) GetTransactionById(ctx context.Context, req *pb.GetTransactionByIdRequest) (*pb.TransactionResponse, error) {
	s.logger.Info("GetTransactionById", slog.String("id", req.Id))
	return s.transactionstorage.GetTransactionById(ctx, req)
}

func (s *TransactionService) UpdateTransaction(ctx context.Context, req *pb.UpdateTransactionRequest) (*pb.TransactionResponse, error) {
	s.logger.Info("UpdateTransaction", slog.Any("req", req))
	return s.transactionstorage.UpdateTransaction(ctx, req)
}

func (s *TransactionService) DeleteTransaction(ctx context.Context, req *pb.DeleteTransactionRequest) (*pb.Empty, error) {
	s.logger.Info("DeleteTransaction", slog.String("id", req.Id))
	return s.transactionstorage.DeleteTransaction(ctx, req)
}
