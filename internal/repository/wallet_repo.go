package repository

import (
	"context"
	"github.com/abdelrahman146/digital-wallet/internal/model"
	"github.com/abdelrahman146/digital-wallet/internal/resource"
	"github.com/abdelrahman146/digital-wallet/pkg/api"
	"github.com/abdelrahman146/digital-wallet/pkg/logger"
)

type WalletRepo interface {
	// CreateWallet Creates a new wallet
	CreateWallet(ctx context.Context, wallet *model.Wallet) error
	// UpdateWallet Updates an existing wallet
	UpdateWallet(ctx context.Context, wallet *model.Wallet) error
	// DeleteWallet Deletes a wallet by its wallet ID
	DeleteWallet(ctx context.Context, wallet *model.Wallet) error
	// FetchWalletByID Retrieves a wallet by its wallet ID
	FetchWalletByID(ctx context.Context, walletId string) (*model.Wallet, error)
	// FetchWallets Retrieves a paginated list of wallets
	FetchWallets(ctx context.Context, page int, limit int) ([]model.Wallet, error)
	// CountTotalWallets Retrieves the total number of wallets
	CountTotalWallets(ctx context.Context) (int64, error)
}

type walletRepo struct {
	resources *resource.Resources
}

// NewWalletRepo initializes the wallet repository
func NewWalletRepo(resources *resource.Resources) WalletRepo {
	return &walletRepo{resources: resources}
}

// CreateWallet creates a new wallet and its associated schema
func (r *walletRepo) CreateWallet(ctx context.Context, wallet *model.Wallet) error {
	// Create the wallet in the database
	if err := r.resources.DB.Create(wallet).Error; err != nil {
		api.GetLogger(ctx).Error("Failed to create wallet", logger.Field("error", err), logger.Field("wallet", wallet))
		return err
	}
	return nil
}

// FetchWalletByID retrieves a wallet by its wallet ID
func (r *walletRepo) FetchWalletByID(ctx context.Context, walletId string) (*model.Wallet, error) {
	var wallet model.Wallet
	err := r.resources.DB.Where("id = ?", walletId).First(&wallet).Error
	if err != nil {
		api.GetLogger(ctx).Error("Failed to fetch wallet by ID", logger.Field("error", err), logger.Field("walletId", walletId))
		return nil, err
	}
	return &wallet, nil
}

// UpdateWallet updates an existing wallet in the database
func (r *walletRepo) UpdateWallet(ctx context.Context, wallet *model.Wallet) error {
	if err := r.resources.DB.Save(wallet).Error; err != nil {
		api.GetLogger(ctx).Error("Failed to update wallet", logger.Field("error", err), logger.Field("wallet", wallet))
		return err
	}
	return nil
}

// FetchWallets retrieves a paginated list of wallets from the database
func (r *walletRepo) FetchWallets(ctx context.Context, page int, limit int) ([]model.Wallet, error) {
	var wallets []model.Wallet
	err := r.resources.DB.Order("created_at desc").Offset((page - 1) * limit).Limit(limit).Find(&wallets).Error
	if err != nil {
		api.GetLogger(ctx).Error("Failed to fetch wallets", logger.Field("error", err))
		return nil, err
	}
	return wallets, nil
}

// CountTotalWallets retrieves the total number of wallets in the database
func (r *walletRepo) CountTotalWallets(ctx context.Context) (int64, error) {
	var total int64
	err := r.resources.DB.Model(&model.Wallet{}).Count(&total).Error
	if err != nil {
		api.GetLogger(ctx).Error("Failed to count total wallets", logger.Field("error", err))
		return 0, err
	}
	return total, nil
}

// DeleteWallet deletes a wallet from the database by its wallet ID
func (r *walletRepo) DeleteWallet(ctx context.Context, wallet *model.Wallet) error {
	if err := r.resources.DB.Delete(wallet).Error; err != nil {
		api.GetLogger(ctx).Error("Failed to delete wallet", logger.Field("error", err), logger.Field("wallet", wallet))
		return err
	}
	return nil
}
