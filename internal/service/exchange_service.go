package service

import (
	"digital-wallet/internal/model"
	"digital-wallet/internal/repository"
	"digital-wallet/pkg/api"
	"digital-wallet/pkg/errs"
	"digital-wallet/pkg/logger"
	"github.com/shopspring/decimal"
)

type ExchangeRateService interface {
	CreateExchangeRate(req *CreateExchangeRateRequest) (*model.ExchangeRate, error)
	GetExchangeRates(page int, limit int) (*api.List[model.ExchangeRate], error)
	GetExchangeRatesByWalletID(walletId string, page int, limit int) (*api.List[model.ExchangeRate], error)
	UpdateExchangeRate(exchangeRateId string, exchangeRate decimal.Decimal) (*model.ExchangeRate, error)
	DeleteExchangeRate(exchangeRateId string) error
	Exchange(fromWalletId, toWalletId, userId, actorType, actorId string, amount uint64) (*ExchangeResponse, error)
}

type exchangeRateService struct {
	repos *repository.Repos
}

func NewExchangeRateService(repos *repository.Repos) ExchangeRateService {
	return &exchangeRateService{repos: repos}
}

func (s *exchangeRateService) CreateExchangeRate(req *CreateExchangeRateRequest) (*model.ExchangeRate, error) {
	exchangeRate := &model.ExchangeRate{
		FromWalletID: req.FromWalletID,
		ToWalletID:   req.ToWalletID,
		TierID:       req.TierID,
		ExchangeRate: req.ExchangeRate,
	}
	if err := s.repos.ExchangeRate.CreateExchangeRate(exchangeRate); err != nil {
		return nil, err
	}
	return exchangeRate, nil
}

func (s *exchangeRateService) GetExchangeRates(page int, limit int) (*api.List[model.ExchangeRate], error) {
	exchangeRates, err := s.repos.ExchangeRate.GetExchangeRates(page, limit)
	if err != nil {
		return nil, err
	}
	total, err := s.repos.ExchangeRate.GetTotalExchangeRates()
	return &api.List[model.ExchangeRate]{Items: exchangeRates, Total: total, Page: page, Limit: limit}, nil
}

func (s *exchangeRateService) GetExchangeRatesByWalletID(walletId string, page int, limit int) (*api.List[model.ExchangeRate], error) {
	exchangeRates, err := s.repos.ExchangeRate.GetExchangeRatesByWalletID(walletId, page, limit)
	if err != nil {
		return nil, err
	}
	total, err := s.repos.ExchangeRate.GetTotalExchangeRatesByWalletID(walletId)
	return &api.List[model.ExchangeRate]{Items: exchangeRates, Total: total, Page: page, Limit: limit}, nil
}

func (s *exchangeRateService) UpdateExchangeRate(exchangeRateId string, newRate decimal.Decimal) (*model.ExchangeRate, error) {
	exchangeRate, err := s.repos.ExchangeRate.GetExchangeRateByID(exchangeRateId)
	if err != nil {
		return nil, err
	}
	exchangeRate.ExchangeRate = newRate
	if err := s.repos.ExchangeRate.UpdateExchangeRate(exchangeRate); err != nil {
		return nil, err
	}
	return exchangeRate, nil
}

func (s *exchangeRateService) DeleteExchangeRate(exchangeRateId string) error {
	exchangeRate, err := s.repos.ExchangeRate.GetExchangeRateByID(exchangeRateId)
	if exchangeRate == nil {
		logger.GetLogger().Error("Exchange Rate not found", logger.Field("exchangeRateId", exchangeRateId))
		return errs.NewNotFoundError("Exchange Rate not found", "EXCHANGE_RATE_NOT_FOUND", err)
	}
	return s.repos.ExchangeRate.DeleteExchangeRate(exchangeRateId)
}

func (s *exchangeRateService) Exchange(fromWalletId, toWalletId, userId, actorType, actorId string, amount uint64) (*ExchangeResponse, error) {
	fromWallet, err := s.repos.Wallet.GetWalletByID(fromWalletId)
	if err != nil {
		return nil, err
	}
	if fromWallet == nil {
		logger.GetLogger().Error("From Wallet not found", logger.Field("fromWalletId", fromWalletId))
		return nil, errs.NewNotFoundError("fromWallet not found", "FROM_WALLET_NOT_FOUND", nil)
	}
	toWallet, err := s.repos.Wallet.GetWalletByID(toWalletId)
	if err != nil {
		return nil, err
	}
	if toWallet == nil {
		logger.GetLogger().Error("To Wallet not found", logger.Field("toWalletId", toWalletId))
		return nil, errs.NewNotFoundError("toWallet not found", "TO_WALLET_NOT_FOUND", nil)
	}
	user, err := s.repos.User.GetUserByID(userId)
	if err != nil {
		return nil, err
	}
	if user == nil {
		logger.GetLogger().Error("User not found", logger.Field("userId", userId))
		return nil, errs.NewNotFoundError("User not found", "USER_NOT_FOUND", nil)
	}
	exchangeRate, err := s.repos.ExchangeRate.GetExchangeRate(fromWalletId, toWalletId, user.TierID)
	if err != nil {
		return nil, err
	}
	if exchangeRate == nil {
		logger.GetLogger().Error("Exchange Rate not found", logger.Field("fromWalletId", fromWalletId), logger.Field("toWalletId", toWalletId), logger.Field("tierId", user.TierID))
		return nil, errs.NewNotFoundError("Exchange Rate not found", "EXCHANGE_RATE_NOT_FOUND", nil)
	}
	fromAccount, err := s.repos.Account.GetAccountByUserID(fromWalletId, userId)
	if err != nil {
		return nil, err
	}
	if fromAccount == nil {
		logger.GetLogger().Error("From Account not found", logger.Field("fromWalletId", fromWalletId), logger.Field("userId", userId))
		return nil, errs.NewNotFoundError("From Account not found", "FROM_ACCOUNT_NOT_FOUND", nil)
	}
	toAccount, err := s.repos.Account.GetAccountByUserID(toWalletId, userId)
	if err != nil {
		return nil, err
	}
	if toAccount == nil {
		toAccount = &model.Account{
			UserID: userId,
		}
		if err = s.repos.Account.CreateAccount(toWalletId, toAccount); err != nil {
			return nil, err
		}
	}
	if amount > fromAccount.Balance {
		logger.GetLogger().Error("Insufficient balance", logger.Field("amount", amount), logger.Field("balance", fromAccount.Balance))
		return nil, errs.NewPaymentRequiredError("Insufficient balance", "INSUFFICIENT_BALANCE", nil)
	}
	fromTransaction := &model.Transaction{
		AccountID: fromAccount.ID,
		Amount:    amount,
		Type:      model.TransactionTypeDebit,
		ActorType: actorType,
		ActorID:   actorId,
	}
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
		ActorType: actorType,
		ActorID:   actorId,
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

	if err := s.repos.ExchangeRate.Exchange(from, to); err != nil {
		return nil, err
	}
	resp := &ExchangeResponse{
		FromTransaction: *fromTransaction,
		ToTransaction:   *toTransaction,
	}
	return resp, nil

}
