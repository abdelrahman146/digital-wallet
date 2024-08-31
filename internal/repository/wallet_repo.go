package repository

import (
	"digital-wallet/internal/model"
	"gorm.io/gorm"
)

type WalletRepo interface {
	CreateWallet(wallet *model.Wallet) error
	GetWalletByUserID(userID string) (*model.Wallet, error)
}

type walletRepo struct {
	db *gorm.DB
}

func NewWalletRepo(db *gorm.DB) WalletRepo {
	return &walletRepo{db: db}
}

func (r *walletRepo) CreateWallet(wallet *model.Wallet) error {
	return r.db.Create(wallet).Error
}

func (r *walletRepo) GetWalletByUserID(userID string) (*model.Wallet, error) {
	var wallet model.Wallet
	err := r.db.Where("user_id = ?", userID).First(&wallet).Error
	if err != nil {
		return nil, err
	}
	return &wallet, nil
}
