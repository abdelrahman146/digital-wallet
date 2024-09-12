package service

import (
	"context"
	"digital-wallet/internal/model"
	"digital-wallet/internal/repository"
	"digital-wallet/pkg/api"
	"digital-wallet/pkg/errs"
)

type AccountService interface {
	CreateAccount(ctx context.Context, walletId, userId string) (*model.Account, error)
	GetAccountByID(ctx context.Context, walletId, accountId string) (*model.Account, error)
	GetAccounts(ctx context.Context, walletId string, page int, limit int) (*api.List[model.Account], error)
	GetAccountsSum(ctx context.Context, walletId string) (uint64, error)
	DeleteAccount(ctx context.Context, walletId, accountId string) error
}

type accountService struct {
	repos *repository.Repos
}

func NewAccountService(repos *repository.Repos) AccountService {
	return &accountService{repos: repos}
}

func (s *accountService) CreateAccount(ctx context.Context, walletId, userId string) (*model.Account, error) {
	user, err := s.repos.User.GetUserByID(ctx, userId)
	if user == nil {
		return nil, errs.NewNotFoundError("User not found", "USER_NOT_FOUND", err)
	}
	wallet, err := s.repos.Wallet.GetWalletByID(ctx, walletId)
	if wallet == nil {
		return nil, errs.NewNotFoundError("Wallet not found", "WALLET_NOT_FOUND", err)
	}

	account := &model.Account{
		UserID: userId,
	}
	err = s.repos.Account.CreateAccount(ctx, wallet.ID, account)
	if err != nil {
		return nil, err
	}
	return account, nil
}

func (s *accountService) GetAccountByID(ctx context.Context, walletId, accountId string) (*model.Account, error) {
	wallet, err := s.repos.Wallet.GetWalletByID(ctx, walletId)
	if wallet == nil {
		return nil, errs.NewNotFoundError("Wallet not found", "WALLET_NOT_FOUND", err)
	}
	account, err := s.repos.Account.GetAccountByID(ctx, wallet.ID, accountId)
	if err != nil {
		return nil, err
	}
	return account, nil
}

func (s *accountService) GetAccountsSum(ctx context.Context, walletId string) (uint64, error) {
	wallet, err := s.repos.Wallet.GetWalletByID(ctx, walletId)
	if wallet == nil {
		return 0, errs.NewNotFoundError("Wallet not found", "WALLET_NOT_FOUND", err)
	}
	sum, err := s.repos.Account.GetAccountsSum(ctx, wallet.ID)
	if err != nil {
		return 0, err
	}
	return sum, nil
}

func (s *accountService) GetAccounts(ctx context.Context, walletId string, page int, limit int) (*api.List[model.Account], error) {
	wallet, err := s.repos.Wallet.GetWalletByID(ctx, walletId)
	if wallet == nil {
		return nil, errs.NewNotFoundError("Wallet not found", "WALLET_NOT_FOUND", err)
	}
	accounts, err := s.repos.Account.GetAccounts(ctx, wallet.ID, page, limit)
	if err != nil {
		return nil, err
	}
	total, err := s.repos.Account.GetTotalAccounts(ctx, wallet.ID)
	if err != nil {
		return nil, err
	}
	return &api.List[model.Account]{Items: accounts, Page: page, Limit: limit, Total: total}, nil
}

func (s *accountService) DeleteAccount(ctx context.Context, walletId, accountId string) error {
	wallet, err := s.repos.Wallet.GetWalletByID(ctx, walletId)
	if wallet == nil {
		return errs.NewNotFoundError("Wallet not found", "WALLET_NOT_FOUND", err)
	}
	account, err := s.repos.Account.GetAccountByID(ctx, wallet.ID, accountId)
	if account == nil {
		return errs.NewNotFoundError("Account not found", "ACCOUNT_NOT_FOUND", err)
	}
	return s.repos.Account.DeleteAccount(ctx, wallet.ID, accountId)
}
