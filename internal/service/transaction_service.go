package service

import (
	"digital-wallet/internal/model"
	"digital-wallet/internal/repository"
	"digital-wallet/pkg/api"
	"digital-wallet/pkg/errs"
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
		return nil, errs.NewValidationError("invalid transaction request", fields)
	}
	wallet, _ := s.repos.Wallet.GetWalletByID(walletId)
	if wallet == nil {
		return nil, errs.NewNotFoundError("wallet not found", nil)
	}
	account, _ := s.repos.Account.GetAccountByID(walletId, accountId)
	if account == nil {
		return nil, errs.NewNotFoundError("account not found", nil)
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
		return nil, errs.NewInternalError("failed to create transaction", err)
	}
	return transaction, nil
}

func (s *transactionService) GetTransactionsByAccountID(walletId, accountId string, page int, limit int) (*api.List[model.Transaction], error) {
	wallet, _ := s.repos.Wallet.GetWalletByID(walletId)
	if wallet == nil {
		return nil, errs.NewNotFoundError("wallet not found", nil)
	}
	account, _ := s.repos.Account.GetAccountByID(walletId, accountId)
	if account == nil {
		return nil, errs.NewNotFoundError("account not found", nil)
	}
	transactions, err := s.repos.Transaction.GetTransactionsByAccountID(walletId, accountId, page, limit)
	if err != nil {
		return nil, errs.NewInternalError("failed to get transactions", err)
	}
	total, err := s.repos.Transaction.GetTotalTransactionsByAccountID(walletId, accountId)
	if err != nil {
		return nil, errs.NewInternalError("failed to get total transactions", err)
	}
	return &api.List[model.Transaction]{Items: transactions, Page: page, Limit: limit, Total: total}, nil
}

func (s *transactionService) GetTransactionsSumByAccountID(walletId, accountId string) (uint64, error) {
	wallet, _ := s.repos.Wallet.GetWalletByID(walletId)
	if wallet == nil {
		return 0, errs.NewNotFoundError("wallet not found", nil)
	}
	account, _ := s.repos.Account.GetAccountByID(walletId, accountId)
	if account == nil {
		return 0, errs.NewNotFoundError("account not found", nil)
	}
	sum, err := s.repos.Transaction.GetTransactionsSumByAccountID(walletId, accountId)
	if err != nil {
		return 0, errs.NewInternalError("failed to get transactions sum", err)
	}
	return sum, nil
}

func (s *transactionService) GetTransactions(walletId string, page int, limit int) (*api.List[model.Transaction], error) {
	wallet, _ := s.repos.Wallet.GetWalletByID(walletId)
	if wallet == nil {
		return nil, errs.NewNotFoundError("wallet not found", nil)
	}
	transactions, err := s.repos.Transaction.GetTransactions(walletId, page, limit)
	if err != nil {
		return nil, errs.NewInternalError("failed to get transactions", err)
	}
	total, err := s.repos.Transaction.GetTotalTransactions(walletId)
	if err != nil {
		return nil, errs.NewInternalError("failed to get total transactions", err)
	}
	return &api.List[model.Transaction]{Items: transactions, Page: page, Limit: limit, Total: total}, nil
}

func (s *transactionService) GetTransactionsSum(walletId string) (uint64, error) {
	wallet, _ := s.repos.Wallet.GetWalletByID(walletId)
	if wallet == nil {
		return 0, errs.NewNotFoundError("wallet not found", nil)
	}
	sum, err := s.repos.Transaction.GetTransactionsSum(walletId)
	if err != nil {
		return 0, errs.NewInternalError("failed to get transactions sum", err)
	}
	return sum, nil
}
