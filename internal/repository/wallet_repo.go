package repository

import (
	"digital-wallet/internal/model"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type WalletRepo interface {
	CreateWallet(wallet *model.Wallet) error
	GetWalletByUserID(userId string) (*model.Wallet, error)
	GetWalletByID(walletId string) (*model.Wallet, error)
	GetWallets(page int, limit int) ([]model.Wallet, error)
	GetTotalWallets() (int64, error)
	GetWalletsSum() (decimal.Decimal, error)
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

func (r *walletRepo) GetWalletByUserID(userId string) (*model.Wallet, error) {
	var wallet model.Wallet
	err := r.db.Where("user_id = ?", userId).First(&wallet).Error
	if err != nil {
		return nil, err
	}
	return &wallet, nil
}

func (r *walletRepo) GetWalletByID(walletId string) (*model.Wallet, error) {
	var wallet model.Wallet
	err := r.db.Where("id = ?", walletId).First(&wallet).Error
	if err != nil {
		return nil, err
	}
	return &wallet, nil
}

func (r *walletRepo) GetWalletsSum() (decimal.Decimal, error) {
	var sum float64
	err := r.db.Model(&model.Wallet{}).Select("COALESCE(SUM(balance), 0)").Row().Scan(&sum)
	if err != nil {
		return decimal.Zero, err
	}
	return decimal.NewFromFloat(sum), nil
}

func (r *walletRepo) GetWallets(page int, limit int) ([]model.Wallet, error) {
	var wallets []model.Wallet
	err := r.db.Order("created_at desc").Offset((page - 1) * limit).Limit(limit).Find(&wallets).Error
	if err != nil {
		return nil, err
	}
	return wallets, nil
}

func (r *walletRepo) GetTotalWallets() (int64, error) {
	var total int64
	err := r.db.Model(&model.Wallet{}).Count(&total).Error
	if err != nil {
		return 0, err
	}
	return total, nil
}
