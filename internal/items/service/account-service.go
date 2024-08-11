package service

import (
	pb "budgeting-service/genproto/account"
	"budgeting-service/internal/items/repository"
	"context"
	"log/slog"
)

type AccountService struct {
	pb.UnimplementedAccountServiceServer
	accountstorage repository.AccountI
	logger         *slog.Logger
}

func NewAccountService(accountstorage repository.AccountI, logger *slog.Logger) *AccountService {
	return &AccountService{
		accountstorage: accountstorage,
		logger:         logger,
	}
}

func (s *AccountService) CreateAccount(ctx context.Context, req *pb.CreateAccountRequest) (*pb.AccountResponse, error) {
	s.logger.Info("CreateAccount", "req", req)
	return s.accountstorage.CreateAccount(ctx, req)
}

func (s *AccountService) GetAccounts(ctx context.Context, req *pb.GetAccountsRequest) (*pb.AccountsResponse, error) {
	s.logger.Info("GetAccounts", "req", req)
	return s.accountstorage.GetAccounts(ctx, req)
}

func (s *AccountService) GetAccountById(ctx context.Context, req *pb.GetAccountByIdRequest) (*pb.AccountResponse, error) {
	s.logger.Info("GetAccountById", "req", req)
	return s.accountstorage.GetAccountById(ctx, req)
}

func (s *AccountService) UpdateAccount(ctx context.Context, req *pb.UpdateAccountRequest) (*pb.AccountResponse, error) {
	s.logger.Info("UpdateAccount", "req", req)
	return s.accountstorage.UpdateAccount(ctx, req)
}

func (s *AccountService) DeleteAccount(ctx context.Context, req *pb.DeleteAccountRequest) (*pb.Empty, error) {
	s.logger.Info("DeleteAccount", "req", req)
	return s.accountstorage.DeleteAccount(ctx, req)
}
