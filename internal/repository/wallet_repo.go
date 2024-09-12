package repository

import (
	"digital-wallet/internal/model"
	"digital-wallet/pkg/logger"
	"gorm.io/gorm"
)

type WalletRepo interface {
	CreateWallet(wallet *model.Wallet) error
	GetWalletByID(walletId string) (*model.Wallet, error)
	UpdateWallet(wallet *model.Wallet) error
	GetWallets(page int, limit int) ([]model.Wallet, error)
	GetTotalWallets() (int64, error)
	DeleteWallet(walletId string) error
}

type walletRepo struct {
	db *gorm.DB
}

func NewWalletRepo(db *gorm.DB) WalletRepo {
	return &walletRepo{db: db}
}

func (r *walletRepo) CreateWallet(wallet *model.Wallet) error {
	if err := r.db.Create(wallet).Error; err != nil {
		logger.GetLogger().Error("Error while creating wallet", logger.Field("error", err), logger.Field("wallet", wallet))
		return err
	}
	err := r.db.Exec("SELECT create_wallet_schema(?);", wallet.ID).Error
	if err != nil {
		logger.GetLogger().Error("Error while creating wallet schema", logger.Field("error", err), logger.Field("wallet", wallet))
		return err
	}
	return nil
}

func (r *walletRepo) GetWalletByID(walletId string) (*model.Wallet, error) {
	var wallet model.Wallet
	err := r.db.Where("id = ?", walletId).First(&wallet).Error
	if err != nil {
		logger.GetLogger().Error("Error while fetching wallet by id", logger.Field("error", err), logger.Field("walletId", walletId))
		return nil, err
	}
	return &wallet, nil
}

func (r *walletRepo) UpdateWallet(wallet *model.Wallet) error {
	if err := r.db.Save(wallet).Error; err != nil {
		logger.GetLogger().Error("Error while updating wallet", logger.Field("error", err), logger.Field("wallet", wallet))
		return err
	}
	return nil
}

func (r *walletRepo) GetWallets(page int, limit int) ([]model.Wallet, error) {
	var wallets []model.Wallet
	err := r.db.Order("created_at desc").Offset((page - 1) * limit).Limit(limit).Find(&wallets).Error
	if err != nil {
		logger.GetLogger().Error("Error while fetching wallets", logger.Field("error", err))
		return nil, err
	}
	return wallets, nil
}

func (r *walletRepo) DeleteWallet(walletId string) error {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ?", walletId).Delete(&model.Wallet{}).Error; err != nil {
			logger.GetLogger().Error("Error while deleting wallet", logger.Field("error", err), logger.Field("walletId", walletId))
			return err
		}
		err := tx.Exec("DROP SCHEMA IF EXISTS " + walletId + "_wallet CASCADE;").Error
		if err != nil {
			logger.GetLogger().Error("Error while dropping wallet schema", logger.Field("error", err), logger.Field("walletId", walletId))
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *walletRepo) GetTotalWallets() (int64, error) {
	var total int64
	err := r.db.Model(&model.Wallet{}).Count(&total).Error
	if err != nil {
		logger.GetLogger().Error("Error while fetching total wallets", logger.Field("error", err))
		return 0, err
	}
	return total, nil
}
