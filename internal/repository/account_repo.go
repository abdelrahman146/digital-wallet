package repository

import (
	"context"
	"digital-wallet/internal/model"
	"digital-wallet/pkg/api"
	"digital-wallet/pkg/logger"
	"fmt"
	"gorm.io/gorm"
)

type AccountRepo interface {
	CreateAccount(ctx context.Context, walletId string, account *model.Account) error
	GetAccountByID(ctx context.Context, walletId string, accountId string) (*model.Account, error)
	GetAccountByUserID(ctx context.Context, walletId string, userId string) (*model.Account, error)
	GetAccounts(ctx context.Context, walletId string, page int, limit int) ([]model.Account, error)
	GetTotalAccounts(ctx context.Context, walletId string) (int64, error)
	GetAccountsSum(ctx context.Context, walletId string) (uint64, error)
	DeleteAccount(ctx context.Context, walletId string, accountId string) error
}

type accountRepo struct {
	db *gorm.DB
}

func NewAccountRepo(db *gorm.DB) AccountRepo {
	return &accountRepo{db: db}
}

func (r *accountRepo) CreateAccount(ctx context.Context, walletId string, account *model.Account) error {
	if err := r.db.Raw("SELECT generate_account_id(?);", walletId).Scan(&account.ID).Error; err != nil {
		return err
	}
	if err := r.db.Table(fmt.Sprintf("%s_wallet.%s", walletId, account.TableName())).Create(account).Error; err != nil {
		api.GetLogger(ctx).Error("failed to create account", logger.Field("error", err), logger.Field("account", account), logger.Field("walletId", walletId))
		return err
	}
	return nil
}

func (r *accountRepo) GetAccountByUserID(ctx context.Context, walletId string, userId string) (*model.Account, error) {
	var account model.Account
	err := r.db.Table(fmt.Sprintf("%s_wallet.%s", walletId, account.TableName())).Where("user_id = ?", userId).First(&account).Error
	if err != nil {
		api.GetLogger(ctx).Error("failed to get account by user id", logger.Field("error", err), logger.Field("userId", userId), logger.Field("walletId", walletId))
		return nil, err
	}
	return &account, nil
}

func (r *accountRepo) GetAccountByID(ctx context.Context, walletId, accountId string) (*model.Account, error) {
	var account model.Account
	err := r.db.Table(fmt.Sprintf("%s_wallet.%s", walletId, account.TableName())).Where("id = ?", accountId).First(&account).Error
	if err != nil {
		api.GetLogger(ctx).Error("failed to get account by id", logger.Field("error", err), logger.Field("accountId", accountId), logger.Field("walletId", walletId))
		return nil, err
	}
	return &account, nil
}

func (r *accountRepo) GetAccountsSum(ctx context.Context, walletId string) (uint64, error) {
	var sum uint64
	var account model.Account
	err := r.db.Table(fmt.Sprintf("%s_wallet.%s", walletId, account.TableName())).Model(&model.Account{}).Select("COALESCE(SUM(balance), 0)").Row().Scan(&sum)
	if err != nil {
		api.GetLogger(ctx).Error("failed to get accounts sum", logger.Field("error", err), logger.Field("walletId", walletId))
		return 0, err
	}
	return sum, nil
}

func (r *accountRepo) GetAccounts(ctx context.Context, walletId string, page int, limit int) ([]model.Account, error) {
	var accounts []model.Account
	var account model.Account
	err := r.db.Table(fmt.Sprintf("%s_wallet.%s", walletId, account.TableName())).Order("created_at desc").Offset((page - 1) * limit).Limit(limit).Find(&accounts).Error
	if err != nil {
		api.GetLogger(ctx).Error("failed to get accounts", logger.Field("error", err), logger.Field("walletId", walletId))
		return nil, err
	}
	return accounts, nil
}

func (r *accountRepo) GetTotalAccounts(ctx context.Context, walletId string) (int64, error) {
	var total int64
	account := &model.Account{}
	err := r.db.Table(fmt.Sprintf("%s_wallet.%s", walletId, account.TableName())).Model(&model.Account{}).Count(&total).Error
	if err != nil {
		api.GetLogger(ctx).Error("failed to get total accounts", logger.Field("error", err), logger.Field("walletId", walletId))
		return 0, err
	}
	return total, nil
}

func (r *accountRepo) DeleteAccount(ctx context.Context, walletId, accountId string) error {
	account := &model.Account{}
	if err := r.db.Table(fmt.Sprintf("%s_wallet.%s", walletId, account.TableName())).Where("id = ?", accountId).Delete(account).Error; err != nil {
		api.GetLogger(ctx).Error("failed to delete account", logger.Field("error", err), logger.Field("accountId", accountId), logger.Field("walletId", walletId))
		return err
	}
	return nil
}
