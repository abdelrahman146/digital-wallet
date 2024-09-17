package repository

import (
	"context"
	"github.com/abdelrahman146/digital-wallet/internal/model"
	"github.com/abdelrahman146/digital-wallet/pkg/api"
	"github.com/abdelrahman146/digital-wallet/pkg/logger"
	"gorm.io/gorm"
)

type AccountRepo interface {
	// CreateAccount Creates a new account
	CreateAccount(ctx context.Context, account *model.Account) error
	// FetchAccountByID Retrieves an account by its account ID
	FetchAccountByID(ctx context.Context, accountId string) (*model.Account, error)
	// FetchAccountByUserID Retrieves an account by wallet ID and user ID
	FetchAccountByUserID(ctx context.Context, walletId, userId string) (*model.Account, error)
	// FetchWalletAccounts Retrieves a paginated list of accounts for a wallet
	FetchWalletAccounts(ctx context.Context, walletId string, page int, limit int) ([]model.Account, error)
	// CountWalletAccounts Retrieves the total number of accounts for a wallet
	CountWalletAccounts(ctx context.Context, walletId string) (int64, error)
	// SumWalletAccounts Retrieves the sum of account balances for a wallet
	SumWalletAccounts(ctx context.Context, walletId string) (uint64, error)
	// RemoveAccount Deletes an account
	RemoveAccount(ctx context.Context, account *model.Account) error
}

type accountRepo struct {
	db *gorm.DB
}

// NewAccountRepo initializes the account repository
func NewAccountRepo(db *gorm.DB) AccountRepo {
	return &accountRepo{db: db}
}

// CreateAccount creates a new account and generates an account ID
func (r *accountRepo) CreateAccount(ctx context.Context, account *model.Account) error {
	// Generate account ID based on the wallet ID
	if err := r.db.Raw("SELECT generate_account_id(?);", account.WalletID).Scan(&account.ID).Error; err != nil {
		api.GetLogger(ctx).Error("Failed to generate account ID", logger.Field("error", err), logger.Field("account", account))
		return err
	}
	// Create the account
	if err := r.db.Create(account).Error; err != nil {
		api.GetLogger(ctx).Error("Failed to create account", logger.Field("error", err), logger.Field("account", account))
		return err
	}
	return nil
}

// FetchAccountByUserID retrieves an account by wallet ID and user ID
func (r *accountRepo) FetchAccountByUserID(ctx context.Context, walletId, userId string) (*model.Account, error) {
	var account model.Account
	err := r.db.Where("wallet_id = ? AND user_id = ?", walletId, userId).First(&account).Error
	if err != nil {
		api.GetLogger(ctx).Error("Failed to retrieve account by user ID", logger.Field("error", err), logger.Field("userId", userId), logger.Field("walletId", walletId))
		return nil, err
	}
	return &account, nil
}

// FetchAccountByID retrieves an account by its account ID
func (r *accountRepo) FetchAccountByID(ctx context.Context, accountId string) (*model.Account, error) {
	var account model.Account
	err := r.db.Where("id = ?", accountId).First(&account).Error
	if err != nil {
		api.GetLogger(ctx).Error("Failed to retrieve account by ID", logger.Field("error", err), logger.Field("accountId", accountId))
		return nil, err
	}
	return &account, nil
}

// SumWalletAccounts retrieves the sum of balances for accounts in a wallet
func (r *accountRepo) SumWalletAccounts(ctx context.Context, walletId string) (uint64, error) {
	var sum uint64
	err := r.db.Model(&model.Account{}).Select("COALESCE(SUM(balance), 0)").Where("wallet_id = ?", walletId).Scan(&sum).Error
	if err != nil {
		api.GetLogger(ctx).Error("Failed to retrieve account balance sum", logger.Field("error", err), logger.Field("walletId", walletId))
		return 0, err
	}
	return sum, nil
}

// FetchWalletAccounts retrieves a paginated list of accounts for a wallet
func (r *accountRepo) FetchWalletAccounts(ctx context.Context, walletId string, page int, limit int) ([]model.Account, error) {
	var accounts []model.Account
	err := r.db.Where("wallet_id = ?", walletId).Order("created_at desc").
		Offset((page - 1) * limit).Limit(limit).Find(&accounts).Error
	if err != nil {
		api.GetLogger(ctx).Error("Failed to retrieve accounts", logger.Field("error", err), logger.Field("walletId", walletId))
		return nil, err
	}
	return accounts, nil
}

// CountWalletAccounts retrieves the total number of accounts in a wallet
func (r *accountRepo) CountWalletAccounts(ctx context.Context, walletId string) (int64, error) {
	var total int64
	err := r.db.Model(&model.Account{}).Where("wallet_id = ?", walletId).Count(&total).Error
	if err != nil {
		api.GetLogger(ctx).Error("Failed to retrieve total accounts", logger.Field("error", err), logger.Field("walletId", walletId))
		return 0, err
	}
	return total, nil
}

// RemoveAccount deletes an account from the database
func (r *accountRepo) RemoveAccount(ctx context.Context, account *model.Account) error {
	if err := r.db.Delete(account).Error; err != nil {
		api.GetLogger(ctx).Error("Failed to delete account", logger.Field("error", err), logger.Field("account", account))
		return err
	}
	return nil
}
