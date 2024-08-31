package repository

import (
	"digital-wallet/internal/model"
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TransactionRepo interface {
	GetTransactionsByWalletID(walletID string, page int, limit int) ([]model.Transaction, error)
	GetTotalTransactionsByWalletID(walletID string) (int64, error)
	Create(transaction *model.Transaction, walletVersion int64) error
	Transfer(from *model.Transaction, fromWalletVersion int64, to *model.Transaction, toWalletVersion int64) error
	GetTransactionsSumByWalletID(walletID string) (float64, error)
	GetTransactionsSum() (float64, error)
}

type transactionRepo struct {
	db *gorm.DB
}

func NewTransactionRepo(db *gorm.DB) TransactionRepo {
	return &transactionRepo{db: db}
}

func (r *transactionRepo) GetTransactionsByWalletID(walletID string, page int, limit int) ([]model.Transaction, error) {
	var transactions []model.Transaction
	err := r.db.Where("wallet_id = ?", walletID).Order("created_at desc").Offset((page - 1) * limit).Limit(limit).Find(&transactions).Error
	if err != nil {
		return nil, err
	}
	return transactions, nil
}

func (r *transactionRepo) GetTotalTransactionsByWalletID(walletID string) (int64, error) {
	var total int64
	err := r.db.Model(&model.Transaction{}).Where("wallet_id = ?", walletID).Count(&total).Error
	if err != nil {
		return 0, err
	}
	return total, nil
}

func (r *transactionRepo) Create(transaction *model.Transaction, walletVersion int64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var wallet model.Wallet

		// Lock the wallet
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", transaction.WalletID).First(&wallet).Error; err != nil {
			return err
		}

		// optimistic locking
		if wallet.Version != walletVersion {
			return errors.New("wallet has been modified by another transaction")
		}

		// check if wallet has sufficient balance
		if wallet.Balance+transaction.Amount < 0 {
			return errors.New("insufficient balance")
		}

		transaction.PreviousBalance = wallet.Balance
		wallet.Balance = wallet.Balance + transaction.Amount
		wallet.Version++
		transaction.NewBalance = wallet.Balance
		transaction.Version = wallet.Version

		if err := tx.Create(transaction).Error; err != nil {
			return err
		}

		return tx.Save(&wallet).Error
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
			return errors.New("from wallet has been modified by another transaction")
		}
		if toWallet.Version != toWalletVersion {
			return errors.New("to wallet has been modified by another transaction")
		}

		// check if from wallet has sufficient balance
		if fromWallet.Balance+from.Amount < 0 {
			return errors.New("insufficient balance")
		}

		from.PreviousBalance = fromWallet.Balance
		fromWallet.Balance = fromWallet.Balance + from.Amount
		fromWallet.Version++
		from.NewBalance = fromWallet.Balance

		to.PreviousBalance = toWallet.Balance
		toWallet.Balance = toWallet.Balance + to.Amount
		toWallet.Version++
		to.NewBalance = toWallet.Balance

		if err := tx.Create(from).Error; err != nil {
			return err
		}

		if err := tx.Create(to).Error; err != nil {
			return err
		}

		if err := tx.Save(&fromWallet).Error; err != nil {
			return err
		}

		return tx.Save(&toWallet).Error
	})
}

func (r *transactionRepo) GetTransactionsSumByWalletID(walletID string) (float64, error) {
	var sum float64
	err := r.db.Model(&model.Transaction{}).Select("COALESCE(SUM(amount), 0)").Where("wallet_id = ?", walletID).Scan(&sum).Error
	if err != nil {
		return 0, err
	}
	return sum, nil
}

func (r *transactionRepo) GetTransactionsSum() (float64, error) {
	var sum float64
	err := r.db.Model(&model.Transaction{}).Select("COALESCE(SUM(amount), 0)").Scan(&sum).Error
	if err != nil {
		return 0, err
	}
	return sum, nil
}
