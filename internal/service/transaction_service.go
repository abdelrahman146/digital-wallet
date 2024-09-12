package service

import (
	"digital-wallet/internal/model"
	"digital-wallet/internal/repository"
	"digital-wallet/pkg/api"
	"digital-wallet/pkg/errs"
	"digital-wallet/pkg/logger"
	"digital-wallet/pkg/validator"
)

type TransactionService interface {
	CreateTransaction(walletId, accountId, actorType, actorId string, req *TransactionRequest) (*model.Transaction, error)
	GetTransactionsByAccountID(walletId, accountId string, page int, limit int) (*api.List[model.Transaction], error)
	GetTransactionsSumByAccountID(walletId, accountId string) (uint64, error)
	GetTransactions(walletId string, page int, limit int) (*api.List[model.Transaction], error)
	GetTransactionsSum(walletId string) (uint64, error)
}

type transactionService struct {
	repos *repository.Repos
}

func NewTransactionService(repos *repository.Repos) TransactionService {
	return &transactionService{repos: repos}
}

func (s *transactionService) CreateTransaction(walletId, accountId, actorType, actorId string, req *TransactionRequest) (*model.Transaction, error) {
	if err := validator.GetValidator().ValidateStruct(req); err != nil {
		fields := validator.GetValidator().GetValidationErrors(err)
		logger.GetLogger().Error("Invalid transaction request", logger.Field("fields", fields), logger.Field("request", req))
		return nil, errs.NewValidationError("Invalid transaction request", "", fields)
	}
	wallet, err := s.repos.Wallet.GetWalletByID(walletId)
	if wallet == nil {
		return nil, errs.NewNotFoundError("wallet not found", "WALLET_NOT_FOUND", err)
	}
	account, _ := s.repos.Account.GetAccountByID(walletId, accountId)
	if account == nil {
		return nil, errs.NewNotFoundError("Account not found", "ACCOUNT_NOT_FOUND", err)
	}
	transaction := &model.Transaction{
		AccountID: accountId,
		Amount:    req.Amount,
		ActorType: actorType,
		ActorID:   actorId,
		Metadata:  req.Metadata,
		ProgramID: req.ProgramID,
	}
	switch req.Type {
	case "credit":
		transaction.Type = model.TransactionTypeCredit
	case "debit":
		transaction.Type = model.TransactionTypeDebit
	}
	if err := s.repos.Transaction.Create(walletId, transaction, account.Version); err != nil {
		return nil, err
	}
	return transaction, nil
}

func (s *transactionService) GetTransactionsByAccountID(walletId, accountId string, page int, limit int) (*api.List[model.Transaction], error) {
	wallet, err := s.repos.Wallet.GetWalletByID(walletId)
	if wallet == nil {
		logger.GetLogger().Error("Wallet not found", logger.Field("walletId", walletId))
		return nil, errs.NewNotFoundError("Wallet not found", "WALLET_NOT_FOUND", err)
	}
	account, err := s.repos.Account.GetAccountByID(walletId, accountId)
	if account == nil {
		logger.GetLogger().Error("Account not found", logger.Field("walletId", walletId), logger.Field("accountId", accountId))
		return nil, errs.NewNotFoundError("Account not found", "ACCOUNT_NOT_FOUND", err)
	}
	transactions, err := s.repos.Transaction.GetTransactionsByAccountID(walletId, accountId, page, limit)
	if err != nil {
		return nil, err
	}
	total, err := s.repos.Transaction.GetTotalTransactionsByAccountID(walletId, accountId)
	if err != nil {
		return nil, err
	}
	return &api.List[model.Transaction]{Items: transactions, Page: page, Limit: limit, Total: total}, nil
}

func (s *transactionService) GetTransactionsSumByAccountID(walletId, accountId string) (uint64, error) {
	wallet, err := s.repos.Wallet.GetWalletByID(walletId)
	if wallet == nil {
		logger.GetLogger().Error("Wallet not found", logger.Field("walletId", walletId))
		return 0, errs.NewNotFoundError("Wallet not found", "WALLET_NOT_FOUND", err)
	}
	account, err := s.repos.Account.GetAccountByID(walletId, accountId)
	if account == nil {
		logger.GetLogger().Error("Account not found", logger.Field("walletId", walletId), logger.Field("accountId", accountId))
		return 0, errs.NewNotFoundError("Account not found", "ACCOUNT_NOT_FOUND", err)
	}
	sum, err := s.repos.Transaction.GetTransactionsSumByAccountID(walletId, accountId)
	return sum, err
}

func (s *transactionService) GetTransactions(walletId string, page int, limit int) (*api.List[model.Transaction], error) {
	wallet, err := s.repos.Wallet.GetWalletByID(walletId)
	if wallet == nil {
		logger.GetLogger().Error("Wallet not found", logger.Field("walletId", walletId))
		return nil, errs.NewNotFoundError("Wallet not found", "WALLET_NOT_FOUND", err)
	}
	transactions, err := s.repos.Transaction.GetTransactions(walletId, page, limit)
	if err != nil {
		return nil, err
	}
	total, err := s.repos.Transaction.GetTotalTransactions(walletId)
	if err != nil {
		return nil, err
	}
	return &api.List[model.Transaction]{Items: transactions, Page: page, Limit: limit, Total: total}, nil
}

func (s *transactionService) GetTransactionsSum(walletId string) (uint64, error) {
	wallet, err := s.repos.Wallet.GetWalletByID(walletId)
	if wallet == nil {
		logger.GetLogger().Error("Wallet not found", logger.Field("walletId", walletId))
		return 0, errs.NewNotFoundError("Wallet not found", "WALLET_NOT_FOUND", err)
	}
	sum, err := s.repos.Transaction.GetTransactionsSum(walletId)
	if err != nil {
		return 0, err
	}
	return sum, nil
}
