package service

import (
	pb "budgeting-service/genproto/transaction"
	"budgeting-service/internal/items/repository"
)

type TransactionService struct {
	pb.UnimplementedTransactionServiceServer
	transactionstorage repository.TransactionI
}

func NewTransactionService(transactionstorage repository.TransactionI) *TransactionService {
	return &TransactionService{
		transactionstorage: transactionstorage,
	}
}
