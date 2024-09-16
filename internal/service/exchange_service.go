package service

import (
	"context"
	"digital-wallet/internal/model"
	"digital-wallet/internal/repository"
	"digital-wallet/pkg/api"
	"digital-wallet/pkg/errs"
	"digital-wallet/pkg/logger"
	"github.com/shopspring/decimal"
)

type ExchangeRateService interface {
	// CreateExchangeRate creates a new exchange rate
	CreateExchangeRate(ctx context.Context, req *CreateExchangeRateRequest) (*model.ExchangeRate, error)
	// GetExchangeRates fetches all exchange rates
	GetExchangeRates(ctx context.Context, page int, limit int) (*api.List[model.ExchangeRate], error)
	// GetExchangeRatesByWalletID fetches all exchange rates for a wallet
	GetExchangeRatesByWalletID(ctx context.Context, walletId string, page int, limit int) (*api.List[model.ExchangeRate], error)
	// UpdateExchangeRate updates an exchange rate
	UpdateExchangeRate(ctx context.Context, exchangeRateId string, exchangeRate decimal.Decimal) (*model.ExchangeRate, error)
	// DeleteExchangeRate deletes an exchange rate
	DeleteExchangeRate(ctx context.Context, exchangeRateId string) error
}

type exchangeRateService struct {
	repos *repository.Repos
}

func NewExchangeRateService(repos *repository.Repos) ExchangeRateService {
	return &exchangeRateService{repos: repos}
}

func (s *exchangeRateService) CreateExchangeRate(ctx context.Context, req *CreateExchangeRateRequest) (*model.ExchangeRate, error) {
	if err := api.IsAdmin(ctx); err != nil {
		api.GetLogger(ctx).Error("User not authorized")
		return nil, err
	}
	exchangeRate := &model.ExchangeRate{
		FromWalletID: req.FromWalletID,
		ToWalletID:   req.ToWalletID,
		TierID:       req.TierID,
		ExchangeRate: req.ExchangeRate,
	}
	exchangeRate.SetActor(*api.GetActor(ctx), api.GetActorID(ctx))
	exchangeRate.SetRemarks("Exchange rate created")
	if err := s.repos.ExchangeRate.CreateExchangeRate(ctx, exchangeRate); err != nil {
		return nil, err
	}
	return exchangeRate, nil
}

func (s *exchangeRateService) GetExchangeRates(ctx context.Context, page int, limit int) (*api.List[model.ExchangeRate], error) {
	if err := api.IsAdmin(ctx); err != nil {
		api.GetLogger(ctx).Error("User not authorized")
		return nil, err
	}
	exchangeRates, err := s.repos.ExchangeRate.FetchExchangeRates(ctx, page, limit)
	if err != nil {
		return nil, err
	}
	total, err := s.repos.ExchangeRate.CountExchangeRates(ctx)
	return &api.List[model.ExchangeRate]{Items: exchangeRates, Total: total, Page: page, Limit: limit}, nil
}

func (s *exchangeRateService) GetExchangeRatesByWalletID(ctx context.Context, walletId string, page int, limit int) (*api.List[model.ExchangeRate], error) {
	if err := api.IsAdmin(ctx); err != nil {
		api.GetLogger(ctx).Error("User not authorized")
		return nil, err
	}
	exchangeRates, err := s.repos.ExchangeRate.FetchWalletExchangeRates(ctx, walletId, page, limit)
	if err != nil {
		return nil, err
	}
	total, err := s.repos.ExchangeRate.CountWalletExchangeRates(ctx, walletId)
	return &api.List[model.ExchangeRate]{Items: exchangeRates, Total: total, Page: page, Limit: limit}, nil
}

func (s *exchangeRateService) UpdateExchangeRate(ctx context.Context, exchangeRateId string, newRate decimal.Decimal) (*model.ExchangeRate, error) {
	if err := api.IsAdmin(ctx); err != nil {
		api.GetLogger(ctx).Error("User not authorized")
		return nil, err
	}
	exchangeRate, err := s.repos.ExchangeRate.FetchExchangeRateByID(ctx, exchangeRateId)
	if err != nil {
		return nil, err
	}
	exchangeRate.SetActor(*api.GetActor(ctx), api.GetActorID(ctx))
	exchangeRate.SetRemarks("Exchange rate updated")
	exchangeRate.SetOldRecord(*exchangeRate)
	exchangeRate.ExchangeRate = newRate
	if err := s.repos.ExchangeRate.UpdateExchangeRate(ctx, exchangeRate); err != nil {
		return nil, err
	}
	return exchangeRate, nil
}

func (s *exchangeRateService) DeleteExchangeRate(ctx context.Context, exchangeRateId string) error {
	if err := api.IsAdmin(ctx); err != nil {
		return err
	}
	exchangeRate, err := s.repos.ExchangeRate.FetchExchangeRateByID(ctx, exchangeRateId)
	if exchangeRate == nil {
		api.GetLogger(ctx).Error("Exchange Rate not found", logger.Field("exchangeRateId", exchangeRateId))
		return errs.NewNotFoundError("Exchange Rate not found", "EXCHANGE_RATE_NOT_FOUND", err)
	}
	exchangeRate.SetActor(*api.GetActor(ctx), api.GetActorID(ctx))
	exchangeRate.SetRemarks("Exchange rate deleted")
	exchangeRate.SetOldRecord(exchangeRate)
	return s.repos.ExchangeRate.RemoveExchangeRate(ctx, exchangeRate)
}
