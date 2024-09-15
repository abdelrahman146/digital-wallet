package service

import (
	"context"
	"digital-wallet/internal/model"
	"digital-wallet/internal/repository"
	"digital-wallet/pkg/api"
	"digital-wallet/pkg/errs"
	"digital-wallet/pkg/logger"
	"github.com/shopspring/decimal"
	"time"
)

type ExchangeRateService interface {
	CreateExchangeRate(ctx context.Context, req *CreateExchangeRateRequest) (*model.ExchangeRate, error)
	GetExchangeRates(ctx context.Context, page int, limit int) (*api.List[model.ExchangeRate], error)
	GetExchangeRatesByWalletID(ctx context.Context, walletId string, page int, limit int) (*api.List[model.ExchangeRate], error)
	UpdateExchangeRate(ctx context.Context, exchangeRateId string, exchangeRate decimal.Decimal) (*model.ExchangeRate, error)
	DeleteExchangeRate(ctx context.Context, exchangeRateId string) error
	Exchange(ctx context.Context, fromWalletId, toWalletId, userId string, amount uint64) (*ExchangeResponse, error)
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
		api.GetLogger(ctx).Error("PerformExchange Rate not found", logger.Field("exchangeRateId", exchangeRateId))
		return errs.NewNotFoundError("PerformExchange Rate not found", "EXCHANGE_RATE_NOT_FOUND", err)
	}
	return s.repos.ExchangeRate.RemoveExchangeRate(ctx, exchangeRateId)
}

func (s *exchangeRateService) Exchange(ctx context.Context, fromWalletId, toWalletId, userId string, amount uint64) (*ExchangeResponse, error) {
	// Get User
	user, err := s.repos.User.FetchUserByID(ctx, userId)
	if user == nil {
		api.GetLogger(ctx).Error("User not found", logger.Field("userId", userId))
		return nil, errs.NewNotFoundError("User not found", "USER_NOT_FOUND", err)
	}

	// Check if user is authorized
	if err := api.IsAuthorizedUser(ctx, userId); err != nil {
		api.GetLogger(ctx).Error("User not authorized", logger.Field("userId", userId))
		return nil, err
	}

	// Get Wallets
	fromWallet, err := s.repos.Wallet.FetchWalletByID(ctx, fromWalletId)
	if err != nil {
		return nil, err
	}
	if fromWallet == nil {
		api.GetLogger(ctx).Error("From Wallet not found", logger.Field("fromWalletId", fromWalletId))
		return nil, errs.NewNotFoundError("fromWallet not found", "FROM_WALLET_NOT_FOUND", nil)
	}
	toWallet, err := s.repos.Wallet.FetchWalletByID(ctx, toWalletId)
	if err != nil {
		return nil, err
	}
	if toWallet == nil {
		api.GetLogger(ctx).Error("To Wallet not found", logger.Field("toWalletId", toWalletId))
		return nil, errs.NewNotFoundError("toWallet not found", "TO_WALLET_NOT_FOUND", nil)
	}

	// Get PerformExchange Rate
	exchangeRate, err := s.repos.ExchangeRate.FetchExchangeRate(ctx, fromWalletId, toWalletId, user.TierID)
	if exchangeRate == nil {
		api.GetLogger(ctx).Error("PerformExchange Rate not found", logger.Field("fromWalletId", fromWalletId), logger.Field("toWalletId", toWalletId), logger.Field("tierId", user.TierID))
		return nil, errs.NewNotFoundError("PerformExchange Rate not found", "EXCHANGE_RATE_NOT_FOUND", err)
	}

	// Get Accounts
	fromAccount, err := s.repos.Account.FetchAccountByUserID(ctx, fromWalletId, userId)
	if fromAccount == nil {
		api.GetLogger(ctx).Error("From Account not found", logger.Field("fromWalletId", fromWalletId), logger.Field("userId", userId))
		return nil, errs.NewNotFoundError("From Account not found", "FROM_ACCOUNT_NOT_FOUND", err)
	}
	toAccount, err := s.repos.Account.FetchAccountByUserID(ctx, toWalletId, userId)
	if toAccount == nil {
		api.GetLogger(ctx).Error("From Account not found", logger.Field("toWalletId", fromWalletId), logger.Field("userId", userId))
		return nil, errs.NewNotFoundError("From Account not found", "TO_ACCOUNT_NOT_FOUND", err)
	}

	// Check if balance is sufficient
	if amount > fromAccount.Balance {
		api.GetLogger(ctx).Error("Insufficient balance", logger.Field("amount", amount), logger.Field("balance", fromAccount.Balance))
		return nil, errs.NewPaymentRequiredError("Insufficient balance", "INSUFFICIENT_BALANCE", nil)
	}
	// Check toWallet limit per user is not exceeded
	if toWallet.LimitPerUser != nil && amount+toAccount.Balance > *toWallet.LimitPerUser {
		api.GetLogger(ctx).Error("Limit per user exceeded", logger.Field("limit", *toWallet.LimitPerUser), logger.Field("amount", amount))
		return nil, errs.NewForbiddenError("Limit per user exceeded", "LIMIT_PER_USER_EXCEEDED", nil)
	}

	// Check toWallet total limit is not exceeded
	sum, err := s.repos.Account.SumWalletAccounts(ctx, toWalletId)
	if err != nil {
		api.GetLogger(ctx).Error("Error while getting wallet accounts sum", logger.Field("toWalletId", toWalletId))
		return nil, err
	}
	if toWallet.LimitGlobal != nil && amount+sum > *toWallet.LimitGlobal {
		api.GetLogger(ctx).Error("Limit global exceeded", logger.Field("limit", *toWallet.LimitGlobal), logger.Field("totalWalletBalance", sum), logger.Field("amount", amount))
		return nil, errs.NewForbiddenError("Limit global exceeded", "LIMIT_GLOBAL_EXCEEDED", nil)
	}

	// CreateTransaction From Transaction
	fromTransaction := &model.Transaction{
		AccountID: fromAccount.ID,
		Amount:    amount,
		Type:      model.TransactionTypeDebit,
	}

	// Setup To Transaction
	amountToCredit := decimal.NewFromUint64(amount)
	if toWallet.IsMonetary {
		amountToCredit = amountToCredit.Div(decimal.NewFromInt(100))
	}
	amountToCredit = amountToCredit.Mul(exchangeRate.ExchangeRate)
	if toWallet.IsMonetary {
		amountToCredit = amountToCredit.Mul(decimal.NewFromInt(100))
	}
	toTransaction := &model.Transaction{
		AccountID: toAccount.ID,
		Amount:    amountToCredit.BigInt().Uint64(),
		Type:      model.TransactionTypeCredit,
	}
	if toWallet.PointsExpireAfter != nil {
		expireAt := time.Now().Add(toWallet.PointsExpireAfter.Duration())
		toTransaction.ExpireAt = &expireAt
	}

	from := &repository.ExchangeRequest{
		WalletID:       fromWalletId,
		Transaction:    fromTransaction,
		AccountVersion: fromAccount.Version,
	}

	to := &repository.ExchangeRequest{
		WalletID:       toWalletId,
		Transaction:    toTransaction,
		AccountVersion: toAccount.Version,
	}

	if err := s.repos.ExchangeRate.Exchange(ctx, from, to); err != nil {
		return nil, err
	}
	resp := &ExchangeResponse{
		FromTransaction: *fromTransaction,
		ToTransaction:   *toTransaction,
	}
	return resp, nil
}
