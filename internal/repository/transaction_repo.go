package repository

import (
	"digital-wallet/internal/model"
	"digital-wallet/pkg/errs"
	"digital-wallet/pkg/logger"
	"fmt"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TransactionRepo interface {
	GetTransactionsByWalletID(walletId string, page int, limit int) ([]model.Transaction, error)
	GetTotalTransactionsByWalletID(walletId string) (int64, error)
	Create(transaction *model.Transaction, walletVersion int64) error
	Transfer(from *model.Transaction, fromWalletVersion int64, to *model.Transaction, toWalletVersion int64) error
	GetTransactionsSumByWalletID(walletId string) (decimal.Decimal, error)
	GetTransactionsSum() (decimal.Decimal, error)
}

type transactionRepo struct {
	db *gorm.DB
}

func NewTransactionRepo(db *gorm.DB) TransactionRepo {
	return &transactionRepo{db: db}
}

func (r *transactionRepo) GetTransactionsByWalletID(walletId string, page int, limit int) ([]model.Transaction, error) {
	var transactions []model.Transaction
	err := r.db.Where("wallet_id = ?", walletId).Order("version desc").Offset((page - 1) * limit).Limit(limit).Find(&transactions).Error
	if err != nil {
		logger.GetLogger().Error("Error while fetching transactions by wallet id", logger.Field("error", err))
		return nil, err
	}
	return transactions, nil
}

func (r *transactionRepo) GetTotalTransactionsByWalletID(walletId string) (int64, error) {
	var total int64
	err := r.db.Model(&model.Transaction{}).Where("wallet_id = ?", walletId).Count(&total).Error
	if err != nil {
		logger.GetLogger().Error("Error while fetching total transactions by wallet id", logger.Field("error", err))
		return 0, err
	}
	return total, nil
}

func (r *transactionRepo) Create(transaction *model.Transaction, walletVersion int64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var wallet model.Wallet

		// Lock the wallet
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", transaction.WalletID).First(&wallet).Error; err != nil {
			logger.GetLogger().Error("Error while fetching wallet by id", logger.Field("error", err))
			return err
		}

		// optimistic locking
		if wallet.Version != walletVersion {
			return errs.NewConflictError(fmt.Sprintf("wallet %s has been modified by another transaction", wallet.ID), nil)
		}

		// check if wallet has sufficient balance
		if wallet.Balance.Add(transaction.Amount).LessThan(decimal.NewFromFloat(0)) {
			return errs.NewNotAcceptableError("insufficient balance", nil)
		}

		transaction.PreviousBalance = wallet.Balance
		wallet.Balance = wallet.Balance.Add(transaction.Amount)
		wallet.Version++
		transaction.NewBalance = wallet.Balance
		transaction.Version = wallet.Version

		if err := tx.Create(transaction).Error; err != nil {
			logger.GetLogger().Error("Error while creating transaction", logger.Field("error", err))
			return err
		}

		if err := tx.Save(&wallet).Error; err != nil {
			logger.GetLogger().Error("Error while saving wallet", logger.Field("error", err))
			return err
		}

		return nil
	})
}

func (r *transactionRepo) Transfer(from *model.Transaction, fromWalletVersion int64, to *model.Transaction, toWalletVersion int64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var fromWallet, toWallet model.Wallet

		// Lock the wallets
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", from.WalletID).First(&fromWallet).Error; err != nil {
			return err
		}
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", to.WalletID).First(&toWallet).Error; err != nil {
			return err
		}

		// optimistic locking
		if fromWallet.Version != fromWalletVersion {
			return errs.NewConflictError(fmt.Sprintf("sender wallet %s has been modified by another transaction", fromWallet.ID), nil)
		}
		if toWallet.Version != toWalletVersion {
			return errs.NewConflictError(fmt.Sprintf("receiver wallet %s has been modified by another transaction", toWallet.ID), nil)
		}

		// check if from wallet has sufficient balance
		if fromWallet.Balance.Add(from.Amount).LessThan(decimal.NewFromFloat(0)) {
			return errs.NewNotAcceptableError("insufficient balance", nil)
		}

		from.PreviousBalance = fromWallet.Balance
		fromWallet.Balance = fromWallet.Balance.Add(from.Amount)
		fromWallet.Version++
		from.NewBalance = fromWallet.Balance

		to.PreviousBalance = toWallet.Balance
		toWallet.Balance = toWallet.Balance.Add(to.Amount)
		toWallet.Version++
		to.NewBalance = toWallet.Balance

		if err := tx.Create(from).Error; err != nil {
			return err
		}

		to.ReferenceID = &from.ID
		if err := tx.Create(to).Error; err != nil {
			logger.GetLogger().Error("Error while creating transaction", logger.Field("error", err))
			return err
		}

		if err := tx.Save(&fromWallet).Error; err != nil {
			logger.GetLogger().Error("Error while saving from wallet", logger.Field("error", err))
			return err
		}

		if err := tx.Save(&toWallet).Error; err != nil {
			logger.GetLogger().Error("Error while saving to wallet", logger.Field("error", err))
			return err
		}
		return nil
	})
}

func (r *transactionRepo) GetTransactionsSumByWalletID(walletId string) (decimal.Decimal, error) {
	var sum float64
	err := r.db.Model(&model.Transaction{}).Select("COALESCE(SUM(amount), 0)").Where("wallet_id = ?", walletId).Scan(&sum).Error
	if err != nil {
		return decimal.Zero, err
	}
	return decimal.NewFromFloat(sum), nil
}

func (r *transactionRepo) GetTransactionsSum() (decimal.Decimal, error) {
	var sum float64
	err := r.db.Model(&model.Transaction{}).Select("COALESCE(SUM(amount), 0)").Scan(&sum).Error
	if err != nil {
		return decimal.Zero, err
	}
	return decimal.NewFromFloat(sum), nil
}
