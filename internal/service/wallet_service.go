package service

import (
	"context"
	"github.com/abdelrahman146/digital-wallet/internal/model"
	"github.com/abdelrahman146/digital-wallet/internal/repository"
	"github.com/abdelrahman146/digital-wallet/pkg/api"
	"github.com/abdelrahman146/digital-wallet/pkg/errs"
	"github.com/abdelrahman146/digital-wallet/pkg/logger"
	"github.com/abdelrahman146/digital-wallet/pkg/types"
	"github.com/abdelrahman146/digital-wallet/pkg/validator"
	"time"
)

type WalletService interface {
	// CreateWallet creates a new wallet
	CreateWallet(ctx context.Context, req *CreateWalletRequest) (*model.Wallet, error)
	// GetAccountsSum fetches the sum of all accounts for a wallet
	GetAccountsSum(ctx context.Context, walletId string) (uint64, error)
	// GetTransactionsSum fetches the sum of all transactions for a wallet
	GetTransactionsSum(ctx context.Context, walletId string) (uint64, error)
	// UpdateWallet updates a wallet
	UpdateWallet(ctx context.Context, walletId string, req *UpdateWalletRequest) (*model.Wallet, error)
	// GetWalletByID fetches a wallet by ID
	GetWalletByID(ctx context.Context, walletId string) (*model.Wallet, error)
	// GetWallets fetches all wallets
	GetWallets(ctx context.Context, page int, limit int) (*api.List[model.Wallet], error)
	// DeleteWallet deletes a wallet by ID
	DeleteWallet(ctx context.Context, walletId string) error
}

type walletService struct {
	repos *repository.Repos
}

func NewWalletService(repos *repository.Repos) WalletService {
	return &walletService{repos: repos}
}

func (s *walletService) CreateWallet(ctx context.Context, req *CreateWalletRequest) (*model.Wallet, error) {
	if err := api.IsAdmin(ctx); err != nil {
		api.GetLogger(ctx).Error("User not authorized")
		return nil, err
	}
	if err := validator.GetValidator().ValidateStruct(req); err != nil {
		fields := validator.GetValidator().GetValidationErrors(err)
		api.GetLogger(ctx).Error("Invalid request", logger.Field("fields", fields), logger.Field("request", req))
		return nil, errs.NewValidationError("Invalid request", "", fields)
	}
	wallet := &model.Wallet{
		ID:           req.ID,
		Name:         req.Name,
		Description:  req.Description,
		Currency:     req.Currency,
		LimitPerUser: req.LimitPerUser,
		IsMonetary:   req.IsMonetary,
		LimitGlobal:  req.LimitGlobal,
	}
	wallet.SetActor(api.GetActor(ctx), api.GetActorID(ctx))
	wallet.SetRemarks("Wallet created")
	if req.PointsExpireAfter != nil {
		pointsExpireAfter := types.Interval(time.Duration(*req.PointsExpireAfter) * time.Millisecond)
		wallet.PointsExpireAfter = &pointsExpireAfter
	}
	if err := s.repos.Wallet.CreateWallet(ctx, wallet); err != nil {
		return nil, err
	}
	return wallet, nil
}

func (s *walletService) UpdateWallet(ctx context.Context, walletId string, req *UpdateWalletRequest) (*model.Wallet, error) {
	if err := api.IsAdmin(ctx); err != nil {
		api.GetLogger(ctx).Error("User not authorized")
		return nil, err
	}
	if err := validator.GetValidator().ValidateStruct(req); err != nil {
		fields := validator.GetValidator().GetValidationErrors(err)
		api.GetLogger(ctx).Error("Invalid transaction request", logger.Field("fields", fields), logger.Field("request", req))
		return nil, errs.NewValidationError("Invalid transaction request", "", fields)
	}
	wallet, err := s.repos.Wallet.FetchWalletByID(ctx, walletId)
	if err != nil {
		return nil, err
	}
	wallet.SetActor(api.GetActor(ctx), api.GetActorID(ctx))
	wallet.SetRemarks("Wallet updated")
	wallet.SetOldRecord(*wallet)
	wallet.Name = req.Name
	wallet.Description = req.Description
	wallet.Currency = req.Currency
	wallet.LimitPerUser = req.LimitPerUser
	wallet.LimitGlobal = req.LimitGlobal
	if req.IsMonetary != nil {
		wallet.IsMonetary = *req.IsMonetary
	}
	if req.PointsExpireAfter != nil {
		pointsExpireAfter := types.Interval(time.Duration(*req.PointsExpireAfter) * time.Millisecond)
		wallet.PointsExpireAfter = &pointsExpireAfter
	}

	if err := s.repos.Wallet.UpdateWallet(ctx, wallet); err != nil {
		return nil, err
	}
	return wallet, nil
}

func (s *walletService) GetWalletByID(ctx context.Context, walletId string) (*model.Wallet, error) {
	wallet, err := s.repos.Wallet.FetchWalletByID(ctx, walletId)
	if wallet == nil {
		return nil, errs.NewNotFoundError("wallet not found", "WALLET_NOT_FOUND", err)
	}
	return wallet, nil
}

func (s *walletService) GetWallets(ctx context.Context, page int, limit int) (*api.List[model.Wallet], error) {
	wallets, err := s.repos.Wallet.FetchWallets(ctx, page, limit)
	if err != nil {
		return nil, err
	}
	total, err := s.repos.Wallet.CountTotalWallets(ctx)
	if err != nil {
		return nil, err
	}
	return &api.List[model.Wallet]{Items: wallets, Page: page, Limit: limit, Total: total}, nil
}

func (s *walletService) DeleteWallet(ctx context.Context, walletId string) error {
	if err := api.IsAdmin(ctx); err != nil {
		api.GetLogger(ctx).Error("User not authorized")
		return err
	}
	wallet, err := s.repos.Wallet.FetchWalletByID(ctx, walletId)
	if wallet == nil {
		api.GetLogger(ctx).Error("Wallet not found", logger.Field("walletId", walletId), logger.Field("error", err))
		return errs.NewNotFoundError("wallet not found", "WALLET_NOT_FOUND", err)
	}
	wallet.SetActor(api.GetActor(ctx), api.GetActorID(ctx))
	wallet.SetRemarks("Wallet deleted")
	wallet.SetOldRecord(*wallet)
	return s.repos.Wallet.DeleteWallet(ctx, wallet)
}

func (s *walletService) GetAccountsSum(ctx context.Context, walletId string) (uint64, error) {
	if err := api.IsAdmin(ctx); err != nil {
		api.GetLogger(ctx).Error("User not authorized")
		return 0, err
	}
	return s.repos.Account.SumWalletAccounts(ctx, walletId)
}

func (s *walletService) GetTransactionsSum(ctx context.Context, walletId string) (uint64, error) {
	if err := api.IsAdmin(ctx); err != nil {
		api.GetLogger(ctx).Error("User not authorized")
		return 0, err
	}
	return s.repos.Transaction.SumWalletTransactions(ctx, walletId)
}
