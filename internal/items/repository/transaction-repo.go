package repository

import (
	pb "budgeting-service/genproto/transaction"
	"context"
)

type TransactionI interface {
	CreateTransaction(ctx context.Context, req *pb.CreateTransactionRequest) (*pb.TransactionResponse, error)
	GetTransactions(ctx context.Context, req *pb.GetTransactionsRequest) (*pb.TransactionsResponse, error)
	GetTransactionById(ctx context.Context, req *pb.GetTransactionByIdRequest) (*pb.TransactionResponse, error)
	UpdateTransaction(ctx context.Context, req *pb.UpdateTransactionRequest) (*pb.TransactionResponse, error)
	DeleteTransaction(ctx context.Context, req *pb.DeleteTransactionRequest) (*pb.Empty, error)
}
