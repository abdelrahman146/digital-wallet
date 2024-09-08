package service

import (
	"digital-wallet/internal/model"
	"digital-wallet/internal/repository"
	"digital-wallet/pkg/api"
)

type AccountService interface {
	CreateAccount(userId string) (*model.Account, error)
	GetAccountByUserID(userId string) (*model.Account, error)
	GetAccountByID(accountId string) (*model.Account, error)
	GetAccounts(page int, limit int) (*api.List[model.Account], error)
	GetAccountsSum() (uint64, error)
}

type accountService struct {
	repos *repository.Repos
}

func NewAccountService(repos *repository.Repos) AccountService {
	return &accountService{repos: repos}
}

func (s *accountService) CreateAccount(userId string) (*model.Account, error) {
	account := &model.Account{
		UserID: userId,
	}
	err := s.repos.Account.CreateAccount(account)
	if err != nil {
		return nil, err
	}
	return account, nil
}

func (s *accountService) GetAccountByUserID(userId string) (*model.Account, error) {
	account, err := s.repos.Account.GetAccountByUserID(userId)
	if err != nil {
		return nil, err
	}
	return account, nil
}

func (s *accountService) GetAccountByID(accountId string) (*model.Account, error) {
	account, err := s.repos.Account.GetAccountByID(accountId)
	if err != nil {
		return nil, err
	}
	return account, nil
}

func (s *accountService) GetAccountsSum() (uint64, error) {
	sum, err := s.repos.Account.GetAccountsSum()
	if err != nil {
		return 0, err
	}
	return sum, nil
}

func (s *accountService) GetAccounts(page int, limit int) (*api.List[model.Account], error) {
	accounts, err := s.repos.Account.GetAccounts(page, limit)
	if err != nil {
		return nil, err
	}
	total, err := s.repos.Account.GetTotalAccounts()
	if err != nil {
		return nil, err
	}
	return &api.List[model.Account]{Items: accounts, Page: page, Limit: limit, Total: total}, nil
}
