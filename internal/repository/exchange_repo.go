package repository

import (
	"digital-wallet/internal/model"
	"digital-wallet/pkg/errs"
	"digital-wallet/pkg/logger"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ExchangeRateRepo interface {
	GetExchangeRateByID(exchangeRateId string) (*model.ExchangeRate, error)
	GetExchangeRate(fromWalletId, toWalletId, tierId string) (*model.ExchangeRate, error)
	CreateExchangeRate(exchangeRate *model.ExchangeRate) error
	GetExchangeRates(page int, limit int) ([]model.ExchangeRate, error)
	GetTotalExchangeRates() (int64, error)
	GetExchangeRatesByWalletID(walletId string, page int, limit int) ([]model.ExchangeRate, error)
	GetTotalExchangeRatesByWalletID(walletId string) (int64, error)
	UpdateExchangeRate(exchangeRate *model.ExchangeRate) error
	DeleteExchangeRate(exchangeRateId string) error
	Exchange(from *ExchangeRequest, to *ExchangeRequest) error
}

type exchangeRateRepo struct {
	db *gorm.DB
}

func NewExchangeRateRepo(db *gorm.DB) ExchangeRateRepo {
	return &exchangeRateRepo{db: db}
}

func (r *exchangeRateRepo) GetExchangeRateByID(exchangeRateId string) (*model.ExchangeRate, error) {
	var exchangeRate model.ExchangeRate
	err := r.db.Where("id = ?", exchangeRateId).First(&exchangeRate).Error
	if err != nil {
		return nil, err
	}
	return &exchangeRate, nil
}

func (r *exchangeRateRepo) CreateExchangeRate(exchangeRate *model.ExchangeRate) error {
	return r.db.Create(exchangeRate).Error
}

func (r *exchangeRateRepo) GetExchangeRates(page int, limit int) ([]model.ExchangeRate, error) {
	var exchangeRates []model.ExchangeRate
	err := r.db.Order("created_at desc").Offset((page - 1) * limit).Limit(limit).Find(&exchangeRates).Error
	if err != nil {
		return nil, err
	}
	return exchangeRates, nil
}

func (r *exchangeRateRepo) GetTotalExchangeRates() (int64, error) {
	var total int64
	err := r.db.Model(&model.ExchangeRate{}).Count(&total).Error
	if err != nil {
		return 0, err
	}
	return total, nil
}

func (r *exchangeRateRepo) GetExchangeRatesByWalletID(walletId string, page int, limit int) ([]model.ExchangeRate, error) {
	var exchangeRates []model.ExchangeRate
	err := r.db.Where("from_wallet_id = ?", walletId).Order("created_at desc").Offset((page - 1) * limit).Limit(limit).Find(&exchangeRates).Error
	if err != nil {
		return nil, err
	}
	return exchangeRates, nil
}

func (r *exchangeRateRepo) GetTotalExchangeRatesByWalletID(walletId string) (int64, error) {
	var total int64
	err := r.db.Model(&model.ExchangeRate{}).Where("from_wallet_id = ?", walletId).Count(&total).Error
	if err != nil {
		return 0, err
	}
	return total, nil
}

func (r *exchangeRateRepo) UpdateExchangeRate(exchangeRate *model.ExchangeRate) error {
	return r.db.Save(exchangeRate).Error
}

func (r *exchangeRateRepo) DeleteExchangeRate(exchangeRateId string) error {
	return r.db.Where("id = ?", exchangeRateId).Delete(&model.ExchangeRate{}).Error
}

func (r *exchangeRateRepo) GetExchangeRate(fromWalletId, toWalletId, tierId string) (*model.ExchangeRate, error) {
	var exchangeRate model.ExchangeRate
	err := r.db.Where("from_wallet_id = ? AND to_wallet_id = ? AND tier_id = ?", fromWalletId, toWalletId, tierId).First(&exchangeRate).Error
	if err != nil {
		return nil, err
	}
	return &exchangeRate, nil
}

func (r *exchangeRateRepo) Exchange(from *ExchangeRequest, to *ExchangeRequest) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		fromAccount := &model.Account{}

		// set schema to from wallet
		if err := tx.Exec("SET search_path TO ?", from.WalletID).Error; err != nil {
			logger.GetLogger().Error("Error while setting search path", logger.Field("error", err))
			return err
		}

		// lock update on from account
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", from.Transaction.AccountID).First(fromAccount).Error; err != nil {
			logger.GetLogger().Error("Error while fetching account by id", logger.Field("error", err))
			return err
		}

		// optimistic locking
		if fromAccount.Version != from.AccountVersion {
			return errs.NewConflictError(fmt.Sprintf("account %s has been modified by another transaction", fromAccount.ID), nil)
		}

		if fromAccount.Balance < from.Transaction.Amount {
			return errs.NewNotAcceptableError("insufficient balance", nil)
		}

		from.Transaction.PreviousBalance = fromAccount.Balance
		fromAccount.Balance = fromAccount.Balance - from.Transaction.Amount
		fromAccount.Version++
		from.Transaction.NewBalance = fromAccount.Balance
		from.Transaction.Version = fromAccount.Version

		// create transaction id
		if err := tx.Raw("SELECT generate_transaction_id(?);", from.WalletID).Scan(&from.Transaction.ID).Error; err != nil {
			logger.GetLogger().Error("Error while generating transaction id", logger.Field("error", err))
			return err
		}

		// create transaction id
		if err := tx.Raw("SELECT generate_transaction_id(?);", to.WalletID).Scan(&to.Transaction.ID).Error; err != nil {
			logger.GetLogger().Error("Error while generating transaction id", logger.Field("error", err))
			return err
		}

		if err := tx.Create(from.Transaction).Error; err != nil {
			logger.GetLogger().Error("Error while creating transaction", logger.Field("error", err))
			return err
		}

		if err := tx.Save(fromAccount).Error; err != nil {
			logger.GetLogger().Error("Error while saving account", logger.Field("error", err))
			return err
		}

		// switch schema
		if err := tx.Exec("SET search_path TO ?", to.WalletID).Error; err != nil {
			logger.GetLogger().Error("Error while setting search path", logger.Field("error", err))
			return err
		}

		toAccount := &model.Account{}

		// lock update on to account
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", to.Transaction.AccountID).First(toAccount).Error; err != nil {
			logger.GetLogger().Error("Error while fetching account by id", logger.Field("error", err))
			return err
		}

		// optimistic locking
		if toAccount.Version != to.AccountVersion {
			return errs.NewConflictError(fmt.Sprintf("account %s has been modified by another transaction", toAccount.ID), nil)
		}

		to.Transaction.PreviousBalance = toAccount.Balance
		toAccount.Balance = toAccount.Balance + to.Transaction.Amount
		toAccount.Version++
		to.Transaction.NewBalance = toAccount.Balance
		to.Transaction.Version = toAccount.Version

		if err := tx.Create(to.Transaction).Error; err != nil {
			logger.GetLogger().Error("Error while creating transaction", logger.Field("error", err))
			return err
		}

		if err := tx.Save(toAccount).Error; err != nil {
			logger.GetLogger().Error("Error while saving account", logger.Field("error", err))
			return err
		}

		return nil

	})
}
