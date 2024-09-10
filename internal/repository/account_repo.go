package repository

import (
	"digital-wallet/internal/model"
	"fmt"
	"gorm.io/gorm"
)

type AccountRepo interface {
	CreateAccount(walletId string, account *model.Account) error
	GetAccountByID(walletId string, accountId string) (*model.Account, error)
	GetAccountByUserID(walletId string, userId string) (*model.Account, error)
	GetAccounts(walletId string, page int, limit int) ([]model.Account, error)
	GetTotalAccounts(walletId string) (int64, error)
	GetAccountsSum(walletId string) (uint64, error)
	DeleteAccount(walletId string, accountId string) error
}

type accountRepo struct {
	db *gorm.DB
}

func NewAccountRepo(db *gorm.DB) AccountRepo {
	return &accountRepo{db: db}
}

func (r *accountRepo) CreateAccount(walletId string, account *model.Account) error {
	if err := r.db.Raw("SELECT generate_account_id(?);", walletId).Scan(&account.ID).Error; err != nil {
		return err
	}
	return r.db.Transaction(func(tx *gorm.DB) error {
		return tx.Exec(fmt.Sprintf("SET search_path TO %s_wallet", walletId)).Create(account).Error
	})
}

func (r *accountRepo) GetAccountByUserID(walletId string, userId string) (*model.Account, error) {
	var wallet model.Account
	err := r.db.Transaction(func(tx *gorm.DB) error {
		return tx.Exec(fmt.Sprintf("SET search_path TO %s_wallet", walletId)).Where("user_id = ?", userId).First(&wallet).Error
	})
	if err != nil {
		return nil, err
	}
	return &wallet, nil
}

func (r *accountRepo) GetAccountByID(walletId, accountId string) (*model.Account, error) {
	var wallet model.Account
	err := r.db.Transaction(func(tx *gorm.DB) error {
		return tx.Exec(fmt.Sprintf("SET search_path TO %s_wallet", walletId)).Where("id = ?", accountId).First(&wallet).Error
	})
	if err != nil {
		return nil, err
	}
	return &wallet, nil
}

func (r *accountRepo) GetAccountsSum(walletId string) (uint64, error) {
	var sum uint64
	err := r.db.Transaction(func(tx *gorm.DB) error {
		return tx.Exec(fmt.Sprintf("SET search_path TO %s_wallet", walletId)).Model(&model.Account{}).Select("COALESCE(SUM(balance), 0)").Row().Scan(&sum)
	})
	if err != nil {
		return 0, err
	}
	return sum, nil
}

func (r *accountRepo) GetAccounts(walletId string, page int, limit int) ([]model.Account, error) {
	var wallets []model.Account
	err := r.db.Transaction(func(tx *gorm.DB) error {
		return tx.Exec(fmt.Sprintf("SET search_path TO %s_wallet", walletId)).Order("created_at desc").Offset((page - 1) * limit).Limit(limit).Find(&wallets).Error
	})
	if err != nil {
		return nil, err
	}
	return wallets, nil
}

func (r *accountRepo) GetTotalAccounts(walletId string) (int64, error) {
	var total int64
	err := r.db.Transaction(func(tx *gorm.DB) error {
		return tx.Exec(fmt.Sprintf("SET search_path TO %s_wallet", walletId)).Model(&model.Account{}).Count(&total).Error
	})
	if err != nil {
		return 0, err
	}
	return total, nil
}

func (r *accountRepo) DeleteAccount(walletId, accountId string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		return tx.Exec(fmt.Sprintf("SET search_path TO %s_wallet", walletId)).Where("id = ?", accountId).Delete(&model.Account{}).Error
	})
}
