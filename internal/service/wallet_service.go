package service

import (
	"context"
	"digital-wallet/internal/model"
	"digital-wallet/internal/repository"
	"digital-wallet/pkg/api"
	"digital-wallet/pkg/errs"
	"digital-wallet/pkg/logger"
	"digital-wallet/pkg/validator"
	"time"
)

type WalletService interface {
	CreateWallet(ctx context.Context, req *CreateWalletRequest) (*model.Wallet, error)
	GetAccountsSum(ctx context.Context, walletId string) (uint64, error)
	GetTransactionsSum(ctx context.Context, walletId string) (uint64, error)
	UpdateWallet(ctx context.Context, walletId string, req *UpdateWalletRequest) (*model.Wallet, error)
	GetWalletByID(ctx context.Context, walletId string) (*model.Wallet, error)
	GetWallets(ctx context.Context, page int, limit int) (*api.List[model.Wallet], error)
	DeleteWallet(ctx context.Context, walletId string) error
}

type walletService struct {
	repos *repository.Repos
}

func NewWalletService(repos *repository.Repos) WalletService {
	return &walletService{repos: repos}
}

func (s *walletService) CreateWallet(ctx context.Context, req *CreateWalletRequest) (*model.Wallet, error) {
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
	if req.PointsExpireAfter != nil {
		pointsExpireAfter := time.Duration(*req.PointsExpireAfter) * time.Millisecond
		wallet.PointsExpireAfter = &pointsExpireAfter
	}
	if err := s.repos.Wallet.CreateWallet(ctx, wallet); err != nil {
		return nil, err
	}
	return wallet, nil
}

func (s *walletService) UpdateWallet(ctx context.Context, walletId string, req *UpdateWalletRequest) (*model.Wallet, error) {
	if err := validator.GetValidator().ValidateStruct(req); err != nil {
		fields := validator.GetValidator().GetValidationErrors(err)
		api.GetLogger(ctx).Error("Invalid transaction request", logger.Field("fields", fields), logger.Field("request", req))
		return nil, errs.NewValidationError("Invalid transaction request", "", fields)
	}
	wallet, err := s.repos.Wallet.GetWalletByID(ctx, walletId)
	if err != nil {
		return nil, err
	}
	wallet.Name = req.Name
	wallet.Description = req.Description
	wallet.Currency = req.Currency
	wallet.LimitPerUser = req.LimitPerUser
	wallet.LimitGlobal = req.LimitGlobal
	if req.IsMonetary != nil {
		wallet.IsMonetary = *req.IsMonetary
	}
	if req.PointsExpireAfter != nil {
		pointsExpireAfter := time.Duration(*req.PointsExpireAfter) * time.Millisecond
		wallet.PointsExpireAfter = &pointsExpireAfter
	}
	if err := s.repos.Wallet.UpdateWallet(ctx, wallet); err != nil {
		return nil, err
	}
	return wallet, nil
}

func (s *walletService) GetWalletByID(ctx context.Context, walletId string) (*model.Wallet, error) {
	wallet, err := s.repos.Wallet.GetWalletByID(ctx, walletId)
	if wallet == nil {
		return nil, errs.NewNotFoundError("wallet not found", "WALLET_NOT_FOUND", err)
	}
	return wallet, nil
}

func (s *walletService) GetWallets(ctx context.Context, page int, limit int) (*api.List[model.Wallet], error) {
	wallets, err := s.repos.Wallet.GetWallets(ctx, page, limit)
	if err != nil {
		return nil, err
	}
	total, err := s.repos.Wallet.GetTotalWallets(ctx)
	if err != nil {
		return nil, err
	}
	return &api.List[model.Wallet]{Items: wallets, Page: page, Limit: limit, Total: total}, nil
}

func (s *walletService) DeleteWallet(ctx context.Context, walletId string) error {
	return s.repos.Wallet.DeleteWallet(ctx, walletId)
}

func (s *walletService) GetAccountsSum(ctx context.Context, walletId string) (uint64, error) {
	return s.repos.Account.GetAccountsSum(ctx, walletId)
}

func (s *walletService) GetTransactionsSum(ctx context.Context, walletId string) (uint64, error) {
	return s.repos.Transaction.GetTransactionsSum(ctx, walletId)
}
