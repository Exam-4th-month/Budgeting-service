package repository

import (
	pb "budgeting-service/genproto/account"
	"context"
)

type AccountI interface {
	CreateAccount(ctx context.Context, req *pb.CreateAccountRequest) (*pb.AccountResponse, error)
	GetAccounts(ctx context.Context, req *pb.GetAccountsRequest) (*pb.AccountsResponse, error)
	GetAccountById(ctx context.Context, req *pb.GetAccountByIdRequest) (*pb.AccountResponse, error)
	UpdateAccount(ctx context.Context, req *pb.UpdateAccountRequest) (*pb.AccountResponse, error)
	DeleteAccount(ctx context.Context, req *pb.DeleteAccountRequest) (*pb.Empty, error)
}
