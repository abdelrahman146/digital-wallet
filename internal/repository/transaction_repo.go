package repository

import (
	"digital-wallet/internal/model"
	"digital-wallet/pkg/errs"
	"digital-wallet/pkg/logger"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TransactionRepo interface {
	GetTransactionsByAccountID(walletId, accountId string, page int, limit int) ([]model.Transaction, error)
	GetTotalTransactionsByAccountID(walletId, accountId string) (int64, error)
	GetTransactions(walletId string, page int, limit int) ([]model.Transaction, error)
	GetTotalTransactions(walletId string) (int64, error)
	Create(walletId string, transaction *model.Transaction, accountVersion uint64) error
	GetTransactionsSumByAccountID(walletId string, accountId string) (uint64, error)
	GetTransactionsSum(walletId string) (uint64, error)
}

type transactionRepo struct {
	db *gorm.DB
}

func NewTransactionRepo(db *gorm.DB) TransactionRepo {
	return &transactionRepo{db: db}
}

func (r *transactionRepo) GetTransactionsByAccountID(walletId, accountId string, page int, limit int) ([]model.Transaction, error) {
	var transactions []model.Transaction
	var transaction model.Transaction
	//err := r.db.Transaction(func(tx *gorm.DB) error {})
	err := r.db.Table(fmt.Sprintf("%s_wallet.%s", walletId, transaction.TableName())).Where("account_id = ?", accountId).Order("created_at desc").Offset((page - 1) * limit).Limit(limit).Find(&transactions).Error
	if err != nil {
		logger.GetLogger().Error("Error while fetching transactions by account id", logger.Field("error", err), logger.Field("accountId", accountId))
		return nil, err
	}
	return transactions, nil
}

func (r *transactionRepo) GetTotalTransactionsByAccountID(walletId, accountId string) (int64, error) {
	var total int64
	transaction := &model.Transaction{}
	err := r.db.Table(fmt.Sprintf("%s_wallet.%s"), walletId, transaction.TableName()).Model(transaction).Where("account_id = ?", accountId).Count(&total).Error
	if err != nil {
		logger.GetLogger().Error("Error while fetching total transactions by account id", logger.Field("error", err), logger.Field("walletId", walletId), logger.Field("accountId", accountId))
		return 0, err
	}
	return total, nil
}

func (r *transactionRepo) Create(walletId string, transaction *model.Transaction, accountVersion uint64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var account model.Account

		// Lock the account
		if err := tx.Table(fmt.Sprintf("%s_wallet.%s", walletId, account.TableName())).Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", transaction.AccountID).First(&account).Error; err != nil {
			logger.GetLogger().Error("Error while fetching account by id", logger.Field("error", err), logger.Field("accountId", transaction.AccountID))
			return err
		}

		// optimistic locking
		if account.Version != accountVersion {
			logger.GetLogger().Error("Account has been modified by another transaction", logger.Field("account", account), logger.Field("accountVersion", accountVersion), logger.Field("transaction", transaction))
			return errs.NewConflictError(fmt.Sprintf("account %s has been modified by another transaction", account.ID), "ACCOUNT_VERSION_MODIFIED", nil)
		}

		// check if account has sufficient balance
		if transaction.Type == model.TransactionTypeDebit && account.Balance < transaction.Amount {
			logger.GetLogger().Error("Insufficient balance", logger.Field("account", account), logger.Field("transaction", transaction))
			return errs.NewPaymentRequiredError("insufficient balance", "INSUFFICIENT_BALANCE", nil)
		}

		transaction.PreviousBalance = account.Balance
		switch transaction.Type {
		case model.TransactionTypeDebit:
			account.Balance = account.Balance - transaction.Amount
		case model.TransactionTypeCredit:
			account.Balance = account.Balance + transaction.Amount
		}
		account.Version++
		transaction.NewBalance = account.Balance
		transaction.Version = account.Version

		// create transaction id
		if err := tx.Raw("SELECT generate_transaction_id(?);", walletId).Scan(&transaction.ID).Error; err != nil {
			logger.GetLogger().Error("Error while generating transaction id", logger.Field("error", err))
			return err
		}

		if err := tx.Table(fmt.Sprintf("%s_wallet.%s", walletId, transaction.TableName())).Create(transaction).Error; err != nil {
			logger.GetLogger().Error("Error while creating transaction", logger.Field("error", err))
			return err
		}

		if err := tx.Table(fmt.Sprintf("%s_wallet.%s", walletId, account.TableName())).Save(&account).Error; err != nil {
			logger.GetLogger().Error("Error while saving account", logger.Field("error", err))
			return err
		}

		return nil
	})
}

func (r *transactionRepo) GetTransactionsSumByAccountID(walletId, accountId string) (uint64, error) {
	var sum float64
	transaction := &model.Transaction{}
	err := r.db.Table(fmt.Sprintf("%s_wallet.%s", walletId, transaction.TableName())).Model(transaction).Select("COALESCE(SUM(amount), 0)").Where("account_id = ?", accountId).Scan(&sum).Error
	if err != nil {
		logger.GetLogger().Error("Error while fetching transactions sum by account id", logger.Field("error", err), logger.Field("walletId", walletId), logger.Field("accountId", accountId))
		return 0, err
	}
	return 0, nil
}

func (r *transactionRepo) GetTransactionsSum(walletId string) (uint64, error) {
	var sum uint64
	var res []struct {
		Type string
		Sum  uint64
	}
	transaction := &model.Transaction{}
	err := r.db.Table(fmt.Sprintf("%s_wallet.%s", walletId, transaction.TableName())).Model(transaction).Select("type, COALESCE(SUM(amount), 0) as sum").Group("type").Scan(&res).Error
	if err != nil {
		logger.GetLogger().Error("Error while fetching transactions sum", logger.Field("error", err), logger.Field("walletId", walletId))
		return 0, err
	}
	mappedRes := make(map[string]uint64)
	for _, r := range res {
		mappedRes[r.Type] = r.Sum
	}
	sum = mappedRes[model.TransactionTypeCredit] - mappedRes[model.TransactionTypeDebit]
	return sum, nil
}

func (r *transactionRepo) GetTransactions(walletId string, page int, limit int) ([]model.Transaction, error) {
	var transactions []model.Transaction
	var transaction model.Transaction
	err := r.db.Table(fmt.Sprintf("%s_wallet.%s", walletId, transaction.TableName())).Order("created_at desc").Offset((page - 1) * limit).Limit(limit).Find(&transactions).Error
	if err != nil {
		logger.GetLogger().Error("Error while fetching transactions", logger.Field("error", err), logger.Field("walletId", walletId))
		return nil, err
	}
	return transactions, nil
}

func (r *transactionRepo) GetTotalTransactions(walletId string) (int64, error) {
	var total int64
	transaction := &model.Transaction{}
	err := r.db.Table(fmt.Sprintf("%s_wallet.%s", walletId, transaction.TableName())).Model(transaction).Count(&total).Error
	if err != nil {
		logger.GetLogger().Error("Error while fetching total transactions", logger.Field("error", err), logger.Field("walletId", walletId))
		return 0, err
	}
	return total, nil
}
