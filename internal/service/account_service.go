package service

import (
	"digital-wallet/internal/model"
	"digital-wallet/internal/repository"
	"digital-wallet/pkg/api"
	"digital-wallet/pkg/errs"
)

type AccountService interface {
	CreateAccount(walletId, userId string) (*model.Account, error)
	GetAccountByID(walletId, accountId string) (*model.Account, error)
	GetAccounts(walletId string, page int, limit int) (*api.List[model.Account], error)
	GetAccountsSum(walletId string) (uint64, error)
	DeleteAccount(walletId, accountId string) error
}

type accountService struct {
	repos *repository.Repos
}

func NewAccountService(repos *repository.Repos) AccountService {
	return &accountService{repos: repos}
}

func (s *accountService) CreateAccount(walletId, userId string) (*model.Account, error) {
	user, err := s.repos.User.GetUserByID(userId)
	if user == nil {
		return nil, errs.NewNotFoundError("User not found", "USER_NOT_FOUND", err)
	}
	wallet, err := s.repos.Wallet.GetWalletByID(walletId)
	if wallet == nil {
		return nil, errs.NewNotFoundError("Wallet not found", "WALLET_NOT_FOUND", err)
	}

	account := &model.Account{
		UserID: userId,
	}
	err = s.repos.Account.CreateAccount(wallet.ID, account)
	if err != nil {
		return nil, err
	}
	return account, nil
}

func (s *accountService) GetAccountByID(walletId, accountId string) (*model.Account, error) {
	wallet, err := s.repos.Wallet.GetWalletByID(walletId)
	if wallet == nil {
		return nil, errs.NewNotFoundError("Wallet not found", "WALLET_NOT_FOUND", err)
	}
	account, err := s.repos.Account.GetAccountByID(wallet.ID, accountId)
	if err != nil {
		return nil, err
	}
	return account, nil
}

func (s *accountService) GetAccountsSum(walletId string) (uint64, error) {
	wallet, err := s.repos.Wallet.GetWalletByID(walletId)
	if wallet == nil {
		return 0, errs.NewNotFoundError("Wallet not found", "WALLET_NOT_FOUND", err)
	}
	sum, err := s.repos.Account.GetAccountsSum(wallet.ID)
	if err != nil {
		return 0, err
	}
	return sum, nil
}

func (s *accountService) GetAccounts(walletId string, page int, limit int) (*api.List[model.Account], error) {
	wallet, err := s.repos.Wallet.GetWalletByID(walletId)
	if wallet == nil {
		return nil, errs.NewNotFoundError("Wallet not found", "WALLET_NOT_FOUND", err)
	}
	accounts, err := s.repos.Account.GetAccounts(wallet.ID, page, limit)
	if err != nil {
		return nil, err
	}
	total, err := s.repos.Account.GetTotalAccounts(wallet.ID)
	if err != nil {
		return nil, err
	}
	return &api.List[model.Account]{Items: accounts, Page: page, Limit: limit, Total: total}, nil
}

func (s *accountService) DeleteAccount(walletId, accountId string) error {
	wallet, err := s.repos.Wallet.GetWalletByID(walletId)
	if wallet == nil {
		return errs.NewNotFoundError("Wallet not found", "WALLET_NOT_FOUND", err)
	}
	account, err := s.repos.Account.GetAccountByID(wallet.ID, accountId)
	if account == nil {
		return errs.NewNotFoundError("Account not found", "ACCOUNT_NOT_FOUND", err)
	}
	return s.repos.Account.DeleteAccount(wallet.ID, accountId)
}
