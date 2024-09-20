package repository

import (
	"context"
	"github.com/abdelrahman146/digital-wallet/internal/model"
	"github.com/abdelrahman146/digital-wallet/internal/resource"
	"github.com/abdelrahman146/digital-wallet/pkg/api"
	"github.com/abdelrahman146/digital-wallet/pkg/logger"
)

type ExchangeRateRepo interface {
	// CreateExchangeRate Creates a new exchange rate
	CreateExchangeRate(ctx context.Context, exchangeRate *model.ExchangeRate) error
	// UpdateExchangeRate Updates an existing exchange rate
	UpdateExchangeRate(ctx context.Context, exchangeRate *model.ExchangeRate) error
	// DeleteExchangeRate Deletes an exchange rate by its ID
	DeleteExchangeRate(ctx context.Context, exchangeRate *model.ExchangeRate) error
	// FetchExchangeRateByID Retrieves an exchange rate by its ID
	FetchExchangeRateByID(ctx context.Context, exchangeRateId string) (*model.ExchangeRate, error)
	// FetchExchangeRate Retrieves an exchange rate by the source and destination wallet IDs, with an optional tier ID
	FetchExchangeRate(ctx context.Context, fromWalletId, toWalletId string, tierId *string) (*model.ExchangeRate, error)
	// FetchExchangeRates Retrieves a paginated list of exchange rates
	FetchExchangeRates(ctx context.Context, page int, limit int) ([]model.ExchangeRate, error)
	// CountExchangeRates Retrieves the total number of exchange rates
	CountExchangeRates(ctx context.Context) (int64, error)
	// FetchWalletExchangeRates Retrieves exchange rates by wallet ID with pagination
	FetchWalletExchangeRates(ctx context.Context, walletId string, page int, limit int) ([]model.ExchangeRate, error)
	// CountWalletExchangeRates Retrieves the total number of exchange rates for a specific wallet
	CountWalletExchangeRates(ctx context.Context, walletId string) (int64, error)
}

type exchangeRateRepo struct {
	resources *resource.Resources
}

// NewExchangeRateRepo initializes the exchange rate repository
func NewExchangeRateRepo(resources *resource.Resources) ExchangeRateRepo {
	return &exchangeRateRepo{resources: resources}
}

// FetchExchangeRateByID retrieves an exchange rate by its ID
func (r *exchangeRateRepo) FetchExchangeRateByID(ctx context.Context, exchangeRateId string) (*model.ExchangeRate, error) {
	var exchangeRate model.ExchangeRate
	err := r.resources.DB.Where("id = ?", exchangeRateId).First(&exchangeRate).Error
	if err != nil {
		api.GetLogger(ctx).Error("Failed to retrieve exchange rate by ID", logger.Field("error", err), logger.Field("exchangeRateId", exchangeRateId))
		return nil, err
	}
	return &exchangeRate, nil
}

// CreateExchangeRate creates a new exchange rate in the database
func (r *exchangeRateRepo) CreateExchangeRate(ctx context.Context, exchangeRate *model.ExchangeRate) error {
	if err := r.resources.DB.Create(exchangeRate).Error; err != nil {
		api.GetLogger(ctx).Error("Failed to create exchange rate", logger.Field("error", err), logger.Field("exchangeRate", exchangeRate))
		return err
	}
	return nil
}

// FetchExchangeRates retrieves a paginated list of exchange rates
func (r *exchangeRateRepo) FetchExchangeRates(ctx context.Context, page int, limit int) ([]model.ExchangeRate, error) {
	var exchangeRates []model.ExchangeRate
	err := r.resources.DB.Order("created_at desc").Offset((page - 1) * limit).Limit(limit).Find(&exchangeRates).Error
	if err != nil {
		api.GetLogger(ctx).Error("Failed to retrieve exchange rates", logger.Field("error", err))
		return nil, err
	}
	return exchangeRates, nil
}

// CountExchangeRates retrieves the total number of exchange rates
func (r *exchangeRateRepo) CountExchangeRates(ctx context.Context) (int64, error) {
	var total int64
	err := r.resources.DB.Model(&model.ExchangeRate{}).Count(&total).Error
	if err != nil {
		api.GetLogger(ctx).Error("Failed to retrieve total exchange rates count", logger.Field("error", err))
		return 0, err
	}
	return total, nil
}

// FetchWalletExchangeRates retrieves exchange rates by wallet ID with pagination
func (r *exchangeRateRepo) FetchWalletExchangeRates(ctx context.Context, walletId string, page int, limit int) ([]model.ExchangeRate, error) {
	var exchangeRates []model.ExchangeRate
	err := r.resources.DB.Where("from_wallet_id = ?", walletId).Order("created_at desc").Offset((page - 1) * limit).Limit(limit).Find(&exchangeRates).Error
	if err != nil {
		api.GetLogger(ctx).Error("Failed to retrieve exchange rates by wallet ID", logger.Field("error", err), logger.Field("walletId", walletId))
		return nil, err
	}
	return exchangeRates, nil
}

// CountWalletExchangeRates retrieves the total number of exchange rates for a specific wallet
func (r *exchangeRateRepo) CountWalletExchangeRates(ctx context.Context, walletId string) (int64, error) {
	var total int64
	err := r.resources.DB.Model(&model.ExchangeRate{}).Where("from_wallet_id = ?", walletId).Count(&total).Error
	if err != nil {
		api.GetLogger(ctx).Error("Failed to retrieve total exchange rates by wallet ID", logger.Field("error", err), logger.Field("walletId", walletId))
		return 0, err
	}
	return total, nil
}

// UpdateExchangeRate updates an existing exchange rate in the database
func (r *exchangeRateRepo) UpdateExchangeRate(ctx context.Context, exchangeRate *model.ExchangeRate) error {
	if err := r.resources.DB.Save(exchangeRate).Error; err != nil {
		api.GetLogger(ctx).Error("Failed to update exchange rate", logger.Field("error", err), logger.Field("exchangeRate", exchangeRate))
		return err
	}
	return nil
}

// DeleteExchangeRate deletes an exchange rate
func (r *exchangeRateRepo) DeleteExchangeRate(ctx context.Context, exchangeRate *model.ExchangeRate) error {
	if err := r.resources.DB.Delete(exchangeRate).Error; err != nil {
		api.GetLogger(ctx).Error("Failed to delete exchange rate", logger.Field("error", err), logger.Field("exchangeRate", exchangeRate))
		return err
	}
	return nil
}

// FetchExchangeRate retrieves an exchange rate by source wallet ID, destination wallet ID, and optionally a tier ID
func (r *exchangeRateRepo) FetchExchangeRate(ctx context.Context, fromWalletId, toWalletId string, tierId *string) (*model.ExchangeRate, error) {
	var exchangeRate model.ExchangeRate
	err := r.resources.DB.Where("from_wallet_id = ? AND to_wallet_id = ? AND tier_id = ?", fromWalletId, toWalletId, tierId).First(&exchangeRate).Error
	if err != nil {
		api.GetLogger(ctx).Error("Failed to retrieve exchange rate by wallet IDs and tier ID", logger.Field("error", err), logger.Field("fromWalletId", fromWalletId), logger.Field("toWalletId", toWalletId), logger.Field("tierId", tierId))
		return nil, err
	}
	return &exchangeRate, nil
}
