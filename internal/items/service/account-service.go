package service

import (
	pb "budgeting-service/genproto/account"
	"budgeting-service/internal/items/repository"
)

type AccountService struct {
	pb.UnimplementedAccountServiceServer
	accountstorage repository.AccountI
}

func NewAccountService(accountstorage repository.AccountI) *AccountService {
	return &AccountService{
		accountstorage: accountstorage,
	}
}
