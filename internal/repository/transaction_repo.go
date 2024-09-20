package repository

import (
	"context"
	"fmt"
	"github.com/abdelrahman146/digital-wallet/internal/model"
	"github.com/abdelrahman146/digital-wallet/internal/resource"
	"github.com/abdelrahman146/digital-wallet/pkg/api"
	"github.com/abdelrahman146/digital-wallet/pkg/errs"
	"github.com/abdelrahman146/digital-wallet/pkg/logger"
	"github.com/abdelrahman146/digital-wallet/pkg/types"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type TransactionRepo interface {
	// FetchAccountTransactions Retrieves a list of transactions by account ID with pagination
	FetchAccountTransactions(ctx context.Context, accountId string, page int, limit int) ([]model.Transaction, error)
	// CountAccountTransactions Retrieves the total count of transactions by account ID
	CountAccountTransactions(ctx context.Context, accountId string) (int64, error)
	// FetchWalletTransactions Retrieves a list of transactions by wallet ID with pagination
	FetchWalletTransactions(ctx context.Context, walletId string, page int, limit int) ([]model.Transaction, error)
	// CountWalletTransactions Retrieves the total count of transactions by wallet ID
	CountWalletTransactions(ctx context.Context, walletId string) (int64, error)
	// SumAccountTransactions Retrieves the sum of transactions for a specific account ID
	SumAccountTransactions(ctx context.Context, accountId string) (uint64, error)
	// SumWalletTransactions Retrieves the sum of transactions for a specific wallet ID
	SumWalletTransactions(ctx context.Context, walletId string) (uint64, error)
	// FetchExpiredWalletTransactions Retrieves transactions for a wallet that have expired
	FetchExpiredWalletTransactions(ctx context.Context, walletId string) ([]model.Transaction, error)
	// SumExpiringAccountTransactions Retrieves the sum of transactions about to expire for a specific account ID
	SumExpiringAccountTransactions(ctx context.Context, accountId string, expireInterval types.Interval) (uint64, error)
	// CreateTransaction Creates a new transaction
	CreateTransaction(ctx context.Context, transaction *model.Transaction, accountVersion uint64) error
	// PerformExchange Performs an exchange between two accounts
	PerformExchange(ctx context.Context, from *ExchangeRequest, to *ExchangeRequest) error
}

type transactionRepo struct {
	resources *resource.Resources
}

// NewTransactionRepo initializes the transaction repository
func NewTransactionRepo(resources *resource.Resources) TransactionRepo {
	return &transactionRepo{resources: resources}
}

// FetchAccountTransactions retrieves transactions by account ID with pagination
func (r *transactionRepo) FetchAccountTransactions(ctx context.Context, accountId string, page int, limit int) ([]model.Transaction, error) {
	var transactions []model.Transaction
	err := r.resources.DB.Where("account_id = ?", accountId).Order("created_at desc").
		Offset((page - 1) * limit).Limit(limit).Find(&transactions).Error
	if err != nil {
		api.GetLogger(ctx).Error("Error fetching transactions by account ID", logger.Field("error", err), logger.Field("accountId", accountId))
		return nil, err
	}
	return transactions, nil
}

// CountAccountTransactions retrieves the total number of transactions by account ID
func (r *transactionRepo) CountAccountTransactions(ctx context.Context, accountId string) (int64, error) {
	var total int64
	err := r.resources.DB.Model(&model.Transaction{}).Where("account_id = ?", accountId).Count(&total).Error
	if err != nil {
		api.GetLogger(ctx).Error("Error fetching total transactions by account ID", logger.Field("error", err), logger.Field("accountId", accountId))
		return 0, err
	}
	return total, nil
}

// SumAccountTransactions retrieves the sum of transaction amounts for a specific account ID
func (r *transactionRepo) SumAccountTransactions(ctx context.Context, accountId string) (uint64, error) {
	var sum uint64
	err := r.resources.DB.Model(&model.Transaction{}).Select("COALESCE(SUM(amount), 0)").Where("account_id = ?", accountId).Scan(&sum).Error
	if err != nil {
		api.GetLogger(ctx).Error("Error fetching transaction sum by account ID", logger.Field("error", err), logger.Field("accountId", accountId))
		return 0, err
	}
	return sum, nil
}

// SumExpiringAccountTransactions retrieves the sum of transactions that are about to expire for an account ID
func (r *transactionRepo) SumExpiringAccountTransactions(ctx context.Context, accountId string, expireInterval types.Interval) (uint64, error) {
	var sum uint64
	err := r.resources.DB.Model(&model.Transaction{}).
		Select("COALESCE(SUM(amount), 0)").
		Where("account_id = ? AND expire_at < ?", accountId, time.Now().Add(expireInterval.Duration()).Format(time.DateOnly)).
		Scan(&sum).Error
	if err != nil {
		api.GetLogger(ctx).Error("Error fetching sum of expiring transactions", logger.Field("error", err), logger.Field("accountId", accountId))
		return 0, err
	}
	return sum, nil
}

// FetchWalletTransactions retrieves transactions by wallet ID with pagination
func (r *transactionRepo) FetchWalletTransactions(ctx context.Context, walletId string, page int, limit int) ([]model.Transaction, error) {
	var transactions []model.Transaction
	err := r.resources.DB.Where("wallet_id = ?", walletId).Order("created_at desc").
		Offset((page - 1) * limit).Limit(limit).Find(&transactions).Error
	if err != nil {
		api.GetLogger(ctx).Error("Error fetching transactions by wallet ID", logger.Field("error", err))
		return nil, err
	}
	return transactions, nil
}

// CountWalletTransactions retrieves the total number of transactions by wallet ID
func (r *transactionRepo) CountWalletTransactions(ctx context.Context, walletId string) (int64, error) {
	var total int64
	err := r.resources.DB.Model(&model.Transaction{}).Where("wallet_id = ?", walletId).Count(&total).Error
	if err != nil {
		api.GetLogger(ctx).Error("Error fetching total transactions by wallet ID", logger.Field("error", err), logger.Field("walletId", walletId))
		return 0, err
	}
	return total, nil
}

// SumWalletTransactions retrieves the sum of transaction amounts for a specific wallet ID
func (r *transactionRepo) SumWalletTransactions(ctx context.Context, walletId string) (uint64, error) {
	var res []struct {
		Type string
		Sum  uint64
	}
	err := r.resources.DB.Model(&model.Transaction{}).Where("wallet_id = ?", walletId).
		Select("type, COALESCE(SUM(amount), 0) as sum").Group("type").Scan(&res).Error
	if err != nil {
		api.GetLogger(ctx).Error("Error fetching transaction sums by wallet ID", logger.Field("error", err), logger.Field("walletId", walletId))
		return 0, err
	}

	// Calculate the net sum (credits - debits)
	var credit, debit uint64
	for _, r := range res {
		if r.Type == model.TransactionTypeCredit {
			credit = r.Sum
		} else if r.Type == model.TransactionTypeDebit {
			debit = r.Sum
		}
	}
	return credit - debit, nil
}

// FetchExpiredWalletTransactions retrieves expired transactions by wallet ID
func (r *transactionRepo) FetchExpiredWalletTransactions(ctx context.Context, walletId string) ([]model.Transaction, error) {
	var transactions []model.Transaction
	err := r.resources.DB.Where("wallet_id = ? AND expire_at < ? AND available_amount > 0", walletId, time.Now().Format(time.DateOnly)).Find(&transactions).Error
	if err != nil {
		api.GetLogger(ctx).Error("Error fetching expired transactions", logger.Field("error", err))
		return nil, err
	}
	return transactions, nil
}

// CreateTransaction creates a new transaction for the specified account
func (r *transactionRepo) CreateTransaction(ctx context.Context, transaction *model.Transaction, accountVersion uint64) error {
	return r.resources.DB.Transaction(func(tx *gorm.DB) error {
		account, err := r.lockAndFetchAccount(ctx, tx, transaction.AccountID, accountVersion)
		if err != nil {
			return err
		}
		return r.createTransaction(ctx, tx, transaction, account)
	})
}

// PerformExchange performs a transaction exchange between two accounts
func (r *transactionRepo) PerformExchange(ctx context.Context, from *ExchangeRequest, to *ExchangeRequest) error {
	return r.resources.DB.Transaction(func(tx *gorm.DB) error {
		fromAccount, err := r.lockAndFetchAccount(ctx, tx, from.Transaction.AccountID, from.AccountVersion)
		if err != nil {
			return err
		}
		if err := r.createTransaction(ctx, tx, from.Transaction, fromAccount); err != nil {
			return err
		}
		toAccount, err := r.lockAndFetchAccount(ctx, tx, to.Transaction.AccountID, to.AccountVersion)
		if err != nil {
			return err
		}
		return r.createTransaction(ctx, tx, to.Transaction, toAccount)
	})
}

// createTransaction handles the actual transaction creation logic
func (r *transactionRepo) createTransaction(ctx context.Context, tx *gorm.DB, transaction *model.Transaction, account *model.Account) error {
	// Check if account has sufficient balance for debit transactions
	if transaction.Type == model.TransactionTypeDebit && account.Balance < transaction.Amount {
		api.GetLogger(ctx).Error("Insufficient balance", logger.Field("account", account), logger.Field("transaction", transaction))
		return errs.NewPaymentRequiredError("insufficient balance", "INSUFFICIENT_BALANCE", nil)
	}

	transaction.PreviousBalance = account.Balance
	switch transaction.Type {
	case model.TransactionTypeDebit:
		account.Balance -= transaction.Amount
		if err := r.applyDebitFIFO(ctx, tx, transaction); err != nil {
			return err
		}
	case model.TransactionTypeCredit:
		account.Balance += transaction.Amount
		transaction.AvailableAmount = transaction.Amount
	}

	account.Version++
	transaction.NewBalance = account.Balance
	transaction.Version = account.Version

	// Generate transaction ID
	if err := tx.Raw("SELECT generate_transaction_id();").Scan(&transaction.ID).Error; err != nil {
		api.GetLogger(ctx).Error("Error generating transaction ID", logger.Field("error", err))
		return err
	}

	// Create transaction and update account balance
	if err := tx.Create(transaction).Error; err != nil {
		api.GetLogger(ctx).Error("Error creating transaction", logger.Field("error", err))
		return err
	}
	if err := tx.Save(&account).Error; err != nil {
		api.GetLogger(ctx).Error("Error saving account", logger.Field("error", err))
		return err
	}

	return nil
}

// lockAndFetchAccount locks and retrieves an account by ID with version checking for optimistic concurrency
func (r *transactionRepo) lockAndFetchAccount(ctx context.Context, tx *gorm.DB, accountId string, accountVersion uint64) (*model.Account, error) {
	var account model.Account
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", accountId).First(&account).Error; err != nil {
		api.GetLogger(ctx).Error("Error fetching account by ID", logger.Field("error", err), logger.Field("accountId", accountId))
		return nil, err
	}
	// Check for version conflicts (optimistic locking)
	if account.Version != accountVersion {
		api.GetLogger(ctx).Error("Account version conflict", logger.Field("account", account), logger.Field("accountVersion", accountVersion))
		return nil, errs.NewConflictError(fmt.Sprintf("Account %s has been modified by another transaction", accountId), "ACCOUNT_VERSION_MODIFIED", nil)
	}
	return &account, nil
}

// applyDebitFIFO applies FIFO (first-in, first-out) logic for debiting transactions from available balance
func (r *transactionRepo) applyDebitFIFO(ctx context.Context, tx *gorm.DB, transaction *model.Transaction) error {
	var transactions []model.Transaction
	var modifiedTransactions []model.Transaction

	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("account_id = ? AND type = 'CREDIT' AND available_amount > 0", transaction.AccountID).
		Order("created_at asc").Find(&transactions).Error
	if err != nil {
		api.GetLogger(ctx).Error("Error fetching credit transactions", logger.Field("error", err))
		return err
	}

	amount := transaction.Amount
	for _, t := range transactions {
		modifiedTransactions = append(modifiedTransactions, t)
		if t.AvailableAmount >= amount {
			t.AvailableAmount -= amount
			amount = 0
			break
		} else {
			amount -= t.AvailableAmount
			t.AvailableAmount = 0
		}
	}

	if amount > 0 {
		api.GetLogger(ctx).Error("Insufficient balance for debit", logger.Field("account_id", transaction.AccountID))
		return errs.NewPaymentRequiredError("insufficient balance", "INSUFFICIENT_BALANCE", nil)
	}

	if len(modifiedTransactions) > 0 {
		if err := tx.Save(&modifiedTransactions).Error; err != nil {
			api.GetLogger(ctx).Error("Error saving modified transactions", logger.Field("error", err))
			return err
		}
	}
	return nil
}
