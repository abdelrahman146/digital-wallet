package repository

import (
	"context"
	"digital-wallet/internal/model"
	"digital-wallet/pkg/api"
	"digital-wallet/pkg/logger"
	"gorm.io/gorm"
)

type WalletRepo interface {
	CreateWallet(ctx context.Context, wallet *model.Wallet) error
	GetWalletByID(ctx context.Context, walletId string) (*model.Wallet, error)
	UpdateWallet(ctx context.Context, wallet *model.Wallet) error
	GetWallets(ctx context.Context, page int, limit int) ([]model.Wallet, error)
	GetTotalWallets(ctx context.Context) (int64, error)
	DeleteWallet(ctx context.Context, walletId string) error
}

type walletRepo struct {
	db *gorm.DB
}

func NewWalletRepo(db *gorm.DB) WalletRepo {
	return &walletRepo{db: db}
}

func (r *walletRepo) CreateWallet(ctx context.Context, wallet *model.Wallet) error {
	if err := r.db.Create(wallet).Error; err != nil {
		api.GetLogger(ctx).Error("Error while creating wallet", logger.Field("error", err), logger.Field("wallet", wallet))
		return err
	}
	err := r.db.Exec("SELECT create_wallet_schema(?);", wallet.ID).Error
	if err != nil {
		api.GetLogger(ctx).Error("Error while creating wallet schema", logger.Field("error", err), logger.Field("wallet", wallet))
		return err
	}
	return nil
}

func (r *walletRepo) GetWalletByID(ctx context.Context, walletId string) (*model.Wallet, error) {
	var wallet model.Wallet
	err := r.db.Where("id = ?", walletId).First(&wallet).Error
	if err != nil {
		api.GetLogger(ctx).Error("Error while fetching wallet by id", logger.Field("error", err), logger.Field("walletId", walletId))
		return nil, err
	}
	return &wallet, nil
}

func (r *walletRepo) UpdateWallet(ctx context.Context, wallet *model.Wallet) error {
	if err := r.db.Save(wallet).Error; err != nil {
		api.GetLogger(ctx).Error("Error while updating wallet", logger.Field("error", err), logger.Field("wallet", wallet))
		return err
	}
	return nil
}

func (r *walletRepo) GetWallets(ctx context.Context, page int, limit int) ([]model.Wallet, error) {
	var wallets []model.Wallet
	err := r.db.Order("created_at desc").Offset((page - 1) * limit).Limit(limit).Find(&wallets).Error
	if err != nil {
		api.GetLogger(ctx).Error("Error while fetching wallets", logger.Field("error", err))
		return nil, err
	}
	return wallets, nil
}

func (r *walletRepo) DeleteWallet(ctx context.Context, walletId string) error {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ?", walletId).Delete(&model.Wallet{}).Error; err != nil {
			api.GetLogger(ctx).Error("Error while deleting wallet", logger.Field("error", err), logger.Field("walletId", walletId))
			return err
		}
		err := tx.Exec("DROP SCHEMA IF EXISTS " + walletId + "_wallet CASCADE;").Error
		if err != nil {
			api.GetLogger(ctx).Error("Error while dropping wallet schema", logger.Field("error", err), logger.Field("walletId", walletId))
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *walletRepo) GetTotalWallets(ctx context.Context) (int64, error) {
	var total int64
	err := r.db.Model(&model.Wallet{}).Count(&total).Error
	if err != nil {
		api.GetLogger(ctx).Error("Error while fetching total wallets", logger.Field("error", err))
		return 0, err
	}
	return total, nil
}
