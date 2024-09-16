package repository

import (
	"context"
	"digital-wallet/internal/model"
	"digital-wallet/pkg/api"
	"digital-wallet/pkg/logger"
	"gorm.io/gorm"
)

type WalletRepo interface {
	// CreateWallet Creates a new wallet
	CreateWallet(ctx context.Context, wallet *model.Wallet) error
	// FetchWalletByID Retrieves a wallet by its wallet ID
	FetchWalletByID(ctx context.Context, walletId string) (*model.Wallet, error)
	// UpdateWallet Updates an existing wallet
	UpdateWallet(ctx context.Context, wallet *model.Wallet) error
	// FetchWallets Retrieves a paginated list of wallets
	FetchWallets(ctx context.Context, page int, limit int) ([]model.Wallet, error)
	// CountTotalWallets Retrieves the total number of wallets
	CountTotalWallets(ctx context.Context) (int64, error)
	// RemoveWallet Deletes a wallet by its wallet ID
	RemoveWallet(ctx context.Context, wallet *model.Wallet) error
}

type walletRepo struct {
	db *gorm.DB
}

// NewWalletRepo initializes the wallet repository
func NewWalletRepo(db *gorm.DB) WalletRepo {
	return &walletRepo{db: db}
}

// CreateWallet creates a new wallet and its associated schema
func (r *walletRepo) CreateWallet(ctx context.Context, wallet *model.Wallet) error {
	// Create the wallet in the database
	if err := r.db.Create(wallet).Error; err != nil {
		api.GetLogger(ctx).Error("Failed to create wallet", logger.Field("error", err), logger.Field("wallet", wallet))
		return err
	}

	// Create the wallet schema in the database
	if err := r.db.Exec("SELECT create_wallet_schema(?);", wallet.ID).Error; err != nil {
		api.GetLogger(ctx).Error("Failed to create wallet schema", logger.Field("error", err), logger.Field("wallet", wallet))
		return err
	}

	return nil
}

// FetchWalletByID retrieves a wallet by its wallet ID
func (r *walletRepo) FetchWalletByID(ctx context.Context, walletId string) (*model.Wallet, error) {
	var wallet model.Wallet
	err := r.db.Where("id = ?", walletId).First(&wallet).Error
	if err != nil {
		api.GetLogger(ctx).Error("Failed to fetch wallet by ID", logger.Field("error", err), logger.Field("walletId", walletId))
		return nil, err
	}
	return &wallet, nil
}

// UpdateWallet updates an existing wallet in the database
func (r *walletRepo) UpdateWallet(ctx context.Context, wallet *model.Wallet) error {
	if err := r.db.Save(wallet).Error; err != nil {
		api.GetLogger(ctx).Error("Failed to update wallet", logger.Field("error", err), logger.Field("wallet", wallet))
		return err
	}
	return nil
}

// FetchWallets retrieves a paginated list of wallets from the database
func (r *walletRepo) FetchWallets(ctx context.Context, page int, limit int) ([]model.Wallet, error) {
	var wallets []model.Wallet
	err := r.db.Order("created_at desc").Offset((page - 1) * limit).Limit(limit).Find(&wallets).Error
	if err != nil {
		api.GetLogger(ctx).Error("Failed to fetch wallets", logger.Field("error", err))
		return nil, err
	}
	return wallets, nil
}

// CountTotalWallets retrieves the total number of wallets in the database
func (r *walletRepo) CountTotalWallets(ctx context.Context) (int64, error) {
	var total int64
	err := r.db.Model(&model.Wallet{}).Count(&total).Error
	if err != nil {
		api.GetLogger(ctx).Error("Failed to count total wallets", logger.Field("error", err))
		return 0, err
	}
	return total, nil
}

// RemoveWallet deletes a wallet from the database by its wallet ID
func (r *walletRepo) RemoveWallet(ctx context.Context, wallet *model.Wallet) error {
	if err := r.db.Delete(wallet).Error; err != nil {
		api.GetLogger(ctx).Error("Failed to delete wallet", logger.Field("error", err), logger.Field("wallet", wallet))
		return err
	}
	return nil
}
