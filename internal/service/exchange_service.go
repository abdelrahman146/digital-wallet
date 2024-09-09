package service

import (
	"digital-wallet/internal/model"
	"digital-wallet/internal/repository"
	"digital-wallet/pkg/api"
	"digital-wallet/pkg/errs"
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
	exchangeRate, _ := s.repos.ExchangeRate.GetExchangeRateByID(exchangeRateId)
	if exchangeRate == nil {
		return errs.NewNotFoundError("exchangeRate not found", nil)
	}
	return s.repos.ExchangeRate.DeleteExchangeRate(exchangeRateId)
}

func (s *exchangeRateService) Exchange(fromWalletId, toWalletId, userId, actorType, actorId string, amount uint64) (*ExchangeResponse, error) {
	fromWallet, err := s.repos.Wallet.GetWalletByID(fromWalletId)
	if err != nil {
		return nil, err
	}
	if fromWallet == nil {
		return nil, errs.NewNotFoundError("fromWallet not found", nil)
	}
	toWallet, err := s.repos.Wallet.GetWalletByID(toWalletId)
	if err != nil {
		return nil, err
	}
	if toWallet == nil {
		return nil, errs.NewNotFoundError("toWallet not found", nil)
	}
	user, err := s.repos.User.GetUserByID(userId)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errs.NewNotFoundError("user not found", nil)
	}
	exchangeRate, err := s.repos.ExchangeRate.GetExchangeRate(fromWalletId, toWalletId, user.TierID)
	if err != nil {
		return nil, err
	}
	if exchangeRate == nil {
		return nil, errs.NewNotFoundError("exchangeRate not found", nil)
	}
	fromAccount, err := s.repos.Account.GetAccountByUserID(fromWalletId, userId)
	if err != nil {
		return nil, err
	}
	if fromAccount == nil {
		return nil, errs.NewNotFoundError("account not found", nil)
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
		return nil, errs.NewBadRequestError("insufficient balance", nil)
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
