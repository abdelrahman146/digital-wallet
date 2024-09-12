package repository

import (
	"context"
	"digital-wallet/internal/model"
	"digital-wallet/pkg/api"
	"digital-wallet/pkg/errs"
	"digital-wallet/pkg/logger"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ExchangeRateRepo interface {
	GetExchangeRateByID(ctx context.Context, exchangeRateId string) (*model.ExchangeRate, error)
	GetExchangeRate(ctx context.Context, fromWalletId, toWalletId, tierId string) (*model.ExchangeRate, error)
	CreateExchangeRate(ctx context.Context, exchangeRate *model.ExchangeRate) error
	GetExchangeRates(ctx context.Context, page int, limit int) ([]model.ExchangeRate, error)
	GetTotalExchangeRates(ctx context.Context) (int64, error)
	GetExchangeRatesByWalletID(ctx context.Context, walletId string, page int, limit int) ([]model.ExchangeRate, error)
	GetTotalExchangeRatesByWalletID(ctx context.Context, walletId string) (int64, error)
	UpdateExchangeRate(ctx context.Context, exchangeRate *model.ExchangeRate) error
	DeleteExchangeRate(ctx context.Context, exchangeRateId string) error
	Exchange(ctx context.Context, from *ExchangeRequest, to *ExchangeRequest) error
}

type exchangeRateRepo struct {
	db *gorm.DB
}

func NewExchangeRateRepo(db *gorm.DB) ExchangeRateRepo {
	return &exchangeRateRepo{db: db}
}

func (r *exchangeRateRepo) GetExchangeRateByID(ctx context.Context, exchangeRateId string) (*model.ExchangeRate, error) {
	var exchangeRate model.ExchangeRate
	err := r.db.Where("id = ?", exchangeRateId).First(&exchangeRate).Error
	if err != nil {
		api.GetLogger(ctx).Error("failed to get exchange rate by id", logger.Field("error", err), logger.Field("exchangeRateId", exchangeRateId))
		return nil, err
	}
	return &exchangeRate, nil
}

func (r *exchangeRateRepo) CreateExchangeRate(ctx context.Context, exchangeRate *model.ExchangeRate) error {
	if err := r.db.Create(exchangeRate).Error; err != nil {
		api.GetLogger(ctx).Error("failed to create exchange rate", logger.Field("error", err), logger.Field("exchangeRate", exchangeRate))
		return err
	}
	return nil
}

func (r *exchangeRateRepo) GetExchangeRates(ctx context.Context, page int, limit int) ([]model.ExchangeRate, error) {
	var exchangeRates []model.ExchangeRate
	err := r.db.Order("created_at desc").Offset((page - 1) * limit).Limit(limit).Find(&exchangeRates).Error
	if err != nil {
		api.GetLogger(ctx).Error("failed to get exchange rates", logger.Field("error", err))
		return nil, err
	}
	return exchangeRates, nil
}

func (r *exchangeRateRepo) GetTotalExchangeRates(ctx context.Context) (int64, error) {
	var total int64
	err := r.db.Model(&model.ExchangeRate{}).Count(&total).Error
	if err != nil {
		api.GetLogger(ctx).Error("failed to get total exchange rates", logger.Field("error", err))
		return 0, err
	}
	return total, nil
}

func (r *exchangeRateRepo) GetExchangeRatesByWalletID(ctx context.Context, walletId string, page int, limit int) ([]model.ExchangeRate, error) {
	var exchangeRates []model.ExchangeRate
	err := r.db.Where("from_wallet_id = ?", walletId).Order("created_at desc").Offset((page - 1) * limit).Limit(limit).Find(&exchangeRates).Error
	if err != nil {
		api.GetLogger(ctx).Error("failed to get exchange rates by wallet id", logger.Field("error", err), logger.Field("walletId", walletId))
		return nil, err
	}
	return exchangeRates, nil
}

func (r *exchangeRateRepo) GetTotalExchangeRatesByWalletID(ctx context.Context, walletId string) (int64, error) {
	var total int64
	err := r.db.Model(&model.ExchangeRate{}).Where("from_wallet_id = ?", walletId).Count(&total).Error
	if err != nil {
		api.GetLogger(ctx).Error("failed to get total exchange rates by wallet id", logger.Field("error", err), logger.Field("walletId", walletId))
		return 0, err
	}
	return total, nil
}

func (r *exchangeRateRepo) UpdateExchangeRate(ctx context.Context, exchangeRate *model.ExchangeRate) error {
	if err := r.db.Save(exchangeRate).Error; err != nil {
		api.GetLogger(ctx).Error("failed to update exchange rate", logger.Field("error", err), logger.Field("exchangeRate", exchangeRate))
		return err
	}
	return nil
}

func (r *exchangeRateRepo) DeleteExchangeRate(ctx context.Context, exchangeRateId string) error {
	if err := r.db.Where("id = ?", exchangeRateId).Delete(&model.ExchangeRate{}).Error; err != nil {
		api.GetLogger(ctx).Error("failed to delete exchange rate", logger.Field("error", err), logger.Field("exchangeRateId", exchangeRateId))
		return err
	}
	return nil
}

func (r *exchangeRateRepo) GetExchangeRate(ctx context.Context, fromWalletId, toWalletId, tierId string) (*model.ExchangeRate, error) {
	var exchangeRate model.ExchangeRate
	err := r.db.Where("from_wallet_id = ? AND to_wallet_id = ? AND tier_id = ?", fromWalletId, toWalletId, tierId).First(&exchangeRate).Error
	if err != nil {
		api.GetLogger(ctx).Error("failed to get exchange rate by wallet id and tier id", logger.Field("error", err), logger.Field("fromWalletId", fromWalletId), logger.Field("toWalletId", toWalletId), logger.Field("tierId", tierId))
		return nil, err
	}
	return &exchangeRate, nil
}

func (r *exchangeRateRepo) Exchange(ctx context.Context, from *ExchangeRequest, to *ExchangeRequest) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		fromAccount := &model.Account{}

		// lock update on from account
		if err := tx.Table(fmt.Sprintf("%s_wallet.%s", from.WalletID, fromAccount.TableName())).Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", from.Transaction.AccountID).First(fromAccount).Error; err != nil {
			api.GetLogger(ctx).Error("Error while fetching account by id", logger.Field("error", err), logger.Field("accountId", from.Transaction.AccountID))
			return err
		}

		// optimistic locking
		if fromAccount.Version != from.AccountVersion {
			api.GetLogger(ctx).Error("account has been modified by another transaction", logger.Field("accountId", fromAccount.ID), logger.Field("accountVersion", from.AccountVersion), logger.Field("account", fromAccount))
			return errs.NewConflictError(fmt.Sprintf("account %s has been modified by another transaction", fromAccount.ID), "ACCOUNT_VERSION_MODIfIED", nil)
		}

		if fromAccount.Balance < from.Transaction.Amount {
			api.GetLogger(ctx).Error("insufficient balance in account", logger.Field("accountId", fromAccount.ID))
			return errs.NewPaymentRequiredError(fmt.Sprintf("insufficient balance in account %s", fromAccount.ID), "INSUFFICIENT_BALANCE", nil)
		}

		from.Transaction.PreviousBalance = fromAccount.Balance
		fromAccount.Balance = fromAccount.Balance - from.Transaction.Amount
		fromAccount.Version++
		from.Transaction.NewBalance = fromAccount.Balance
		from.Transaction.Version = fromAccount.Version

		// create transaction id
		if err := tx.Raw("SELECT generate_transaction_id(?);", from.WalletID).Scan(&from.Transaction.ID).Error; err != nil {
			api.GetLogger(ctx).Error("Error while generating transaction id", logger.Field("error", err))
			return err
		}

		if err := tx.Table(fmt.Sprintf("%s_wallet.%s", from.WalletID, from.Transaction.TableName())).Create(from.Transaction).Error; err != nil {
			api.GetLogger(ctx).Error("Error while creating transaction", logger.Field("error", err), logger.Field("transaction", from.Transaction))
			return err
		}

		if err := tx.Table(fmt.Sprintf("%s_wallet.%s", from.WalletID, fromAccount.TableName())).Save(fromAccount).Error; err != nil {
			api.GetLogger(ctx).Error("Error while saving account", logger.Field("error", err), logger.Field("account", fromAccount))
			return err
		}

		toAccount := &model.Account{}

		// lock update on to account
		if err := tx.Table(fmt.Sprintf("%s_wallet.%s", to.WalletID, toAccount.TableName())).Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", to.Transaction.AccountID).First(toAccount).Error; err != nil {
			api.GetLogger(ctx).Error("Error while fetching account by id", logger.Field("error", err))
			return err
		}

		// optimistic locking
		if toAccount.Version != to.AccountVersion {
			api.GetLogger(ctx).Error("account has been modified by another transaction", logger.Field("accountId", toAccount.ID), logger.Field("accountVersion", to.AccountVersion), logger.Field("account", toAccount))
			return errs.NewConflictError(fmt.Sprintf("account %s has been modified by another transaction", toAccount.ID), "ACCOUNT_VERSION_MODIFIED", nil)
		}

		to.Transaction.PreviousBalance = toAccount.Balance
		toAccount.Balance = toAccount.Balance + to.Transaction.Amount
		toAccount.Version++
		to.Transaction.NewBalance = toAccount.Balance
		to.Transaction.Version = toAccount.Version

		// create transaction id
		if err := tx.Raw("SELECT generate_transaction_id(?);", to.WalletID).Scan(&to.Transaction.ID).Error; err != nil {
			api.GetLogger(ctx).Error("Error while generating transaction id", logger.Field("error", err), logger.Field("walletId", to.WalletID))
			return err
		}

		if err := tx.Table(fmt.Sprintf("%s_wallet.%s", to.WalletID, to.Transaction.TableName())).Create(to.Transaction).Error; err != nil {
			api.GetLogger(ctx).Error("Error while creating transaction", logger.Field("error", err), logger.Field("walletId", to.WalletID), logger.Field("transaction", to.Transaction))
			return err
		}

		if err := tx.Table(fmt.Sprintf("%s_wallet.%s", to.WalletID, toAccount.TableName())).Save(toAccount).Error; err != nil {
			api.GetLogger(ctx).Error("Error while saving account", logger.Field("error", err), logger.Field("account", toAccount))
			return err
		}

		return nil

	})
}
