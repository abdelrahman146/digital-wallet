package repository

import (
	"digital-wallet/internal/model"
	"gorm.io/gorm"
)

type AccountRepo interface {
	CreateAccount(account *model.Account) error
	GetAccountByUserID(userId string) (*model.Account, error)
	GetAccountByID(accountId string) (*model.Account, error)
	GetAccounts(page int, limit int) ([]model.Account, error)
	GetTotalAccounts() (int64, error)
	GetAccountsSum() (uint64, error)
}

type accountRepo struct {
	db *gorm.DB
}

func NewAccountRepo(db *gorm.DB) AccountRepo {
	return &accountRepo{db: db}
}

func (r *accountRepo) CreateAccount(account *model.Account) error {
	return r.db.Create(account).Error
}

func (r *accountRepo) GetAccountByUserID(userId string) (*model.Account, error) {
	var wallet model.Account
	err := r.db.Where("user_id = ?", userId).First(&wallet).Error
	if err != nil {
		return nil, err
	}
	return &wallet, nil
}

func (r *accountRepo) GetAccountByID(accountId string) (*model.Account, error) {
	var wallet model.Account
	err := r.db.Where("id = ?", accountId).First(&wallet).Error
	if err != nil {
		return nil, err
	}
	return &wallet, nil
}

func (r *accountRepo) GetAccountsSum() (uint64, error) {
	var sum uint64
	err := r.db.Model(&model.Account{}).Select("COALESCE(SUM(balance), 0)").Row().Scan(&sum)
	if err != nil {
		return 0, err
	}
	return sum, nil
}

func (r *accountRepo) GetAccounts(page int, limit int) ([]model.Account, error) {
	var wallets []model.Account
	err := r.db.Order("created_at desc").Offset((page - 1) * limit).Limit(limit).Find(&wallets).Error
	if err != nil {
		return nil, err
	}
	return wallets, nil
}

func (r *accountRepo) GetTotalAccounts() (int64, error) {
	var total int64
	err := r.db.Model(&model.Account{}).Count(&total).Error
	if err != nil {
		return 0, err
	}
	return total, nil
}
