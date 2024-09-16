package service

import (
	"context"
	"digital-wallet/internal/model"
	"digital-wallet/internal/repository"
	"digital-wallet/pkg/api"
	"digital-wallet/pkg/errs"
	"digital-wallet/pkg/logger"
	"fmt"
)

type AccountService interface {
	// CreateAccount creates a new account for a user
	CreateAccount(ctx context.Context, walletId, userId string) (*model.Account, error)
	// GetAccount fetches an account by ID
	GetAccount(ctx context.Context, accountId string) (*model.Account, error)
	// GetWalletAccounts fetches all accounts for a wallet
	GetWalletAccounts(ctx context.Context, walletId string, page int, limit int) (*api.List[model.Account], error)
	// GetWalletAccountsSum fetches the sum of all accounts for a wallet
	GetWalletAccountsSum(ctx context.Context, walletId string) (uint64, error)
	// DeleteAccount deletes an account by ID
	DeleteAccount(ctx context.Context, accountId string) error
}

type accountService struct {
	repos *repository.Repos
}

func NewAccountService(repos *repository.Repos) AccountService {
	return &accountService{repos: repos}
}

func (s *accountService) CreateAccount(ctx context.Context, walletId, userId string) (*model.Account, error) {
	user, err := s.repos.User.FetchUserByID(ctx, userId)
	if user == nil {
		return nil, errs.NewNotFoundError("User not found", "USER_NOT_FOUND", err)
	}
	if err := api.IsAuthorizedUser(ctx, user.ID); err != nil {
		api.GetLogger(ctx).Error("User not authorized", logger.Field("userId", user.ID))
		return nil, err
	}
	wallet, err := s.repos.Wallet.FetchWalletByID(ctx, walletId)
	if wallet == nil {
		return nil, errs.NewNotFoundError("Wallet not found", "WALLET_NOT_FOUND", err)
	}
	account := &model.Account{
		UserID: userId,
	}
	account.SetActor(*api.GetActor(ctx), api.GetActorID(ctx))
	account.SetRemarks(fmt.Sprintf("Account created for user %s", userId))
	err = s.repos.Account.CreateAccount(ctx, account)
	if err != nil {
		return nil, err
	}
	return account, nil
}

func (s *accountService) GetAccount(ctx context.Context, accountId string) (*model.Account, error) {
	account, err := s.repos.Account.FetchAccountByID(ctx, accountId)
	if account == nil {
		return nil, errs.NewNotFoundError("Account not found", "ACCOUNT_NOT_FOUND", err)
	}
	if err := api.IsAuthorizedUser(ctx, account.UserID); err != nil {
		api.GetLogger(ctx).Error("User not authorized", logger.Field("userId", account.UserID))
		return nil, err
	}
	return account, nil
}

func (s *accountService) GetWalletAccountsSum(ctx context.Context, walletId string) (uint64, error) {
	if err := api.IsAdmin(ctx); err != nil {
		api.GetLogger(ctx).Error("User not authorized")
		return 0, err
	}
	wallet, err := s.repos.Wallet.FetchWalletByID(ctx, walletId)
	if wallet == nil {
		return 0, errs.NewNotFoundError("Wallet not found", "WALLET_NOT_FOUND", err)
	}
	sum, err := s.repos.Account.SumWalletAccounts(ctx, wallet.ID)
	if err != nil {
		return 0, err
	}
	return sum, nil
}

func (s *accountService) GetWalletAccounts(ctx context.Context, walletId string, page int, limit int) (*api.List[model.Account], error) {
	if err := api.IsAdmin(ctx); err != nil {
		api.GetLogger(ctx).Error("User not authorized")
		return nil, err
	}
	wallet, err := s.repos.Wallet.FetchWalletByID(ctx, walletId)
	if wallet == nil {
		return nil, errs.NewNotFoundError("Wallet not found", "WALLET_NOT_FOUND", err)
	}
	accounts, err := s.repos.Account.FetchWalletAccounts(ctx, wallet.ID, page, limit)
	if err != nil {
		return nil, err
	}
	total, err := s.repos.Account.CountWalletAccounts(ctx, wallet.ID)
	if err != nil {
		return nil, err
	}
	return &api.List[model.Account]{Items: accounts, Page: page, Limit: limit, Total: total}, nil
}

func (s *accountService) DeleteAccount(ctx context.Context, accountId string) error {
	if err := api.IsAdmin(ctx); err != nil {
		api.GetLogger(ctx).Error("User not authorized")
		return err
	}
	account, err := s.repos.Account.FetchAccountByID(ctx, accountId)
	if account == nil {
		return errs.NewNotFoundError("Account not found", "ACCOUNT_NOT_FOUND", err)
	}
	return s.repos.Account.RemoveAccount(ctx, account)
}
