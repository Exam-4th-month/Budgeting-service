package service

import (
	pb "budgeting-service/genproto/transaction"
	"budgeting-service/internal/items/repository"
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
		logger: logger,
	}
}
