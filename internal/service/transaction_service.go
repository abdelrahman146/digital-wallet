package service

import (
	"digital-wallet/internal/model"
	"digital-wallet/internal/repository"
	"digital-wallet/pkg/api"
	"digital-wallet/pkg/errs"
	"digital-wallet/pkg/validator"
)

type TransactionService interface {
	Deposit(req *DepositRequest) (*model.Transaction, error)
	Withdraw(req *WithdrawRequest) (*model.Transaction, error)
	Refund(req *RefundRequest) (*model.Transaction, error)
	Purchase(req *PurchaseRequest) (*model.Transaction, error)
	Transfer(req *TransferRequest) (*model.Transaction, error)
	GetTransactionsByWalletID(walletId string, page int, limit int) (*api.List[model.Transaction], error)
	GetTransactionsSumByWalletID(walletId string) (float64, error)
	GetTransactionsSum() (float64, error)
}

type transactionService struct {
	repos *repository.Repos
}

func NewTransactionService(repos *repository.Repos) TransactionService {
	return &transactionService{repos: repos}
}

func (s *transactionService) Deposit(req *DepositRequest) (*model.Transaction, error) {
	if err := validator.GetValidator().ValidateStruct(req); err != nil {
		fields := validator.GetValidator().GetValidationErrors(err)
		return nil, errs.NewValidationError("invalid deposit request", fields)
	}
	wallet, err := s.repos.Wallet.GetWalletByUserID(req.UserID)
	if err != nil {
		return nil, errs.NewNotFoundError("wallet not found", err)
	}
	transaction := &model.Transaction{
		WalletID:      wallet.ID,
		Amount:        req.Amount,
		Type:          model.TransactionTypeDeposit,
		ReferenceID:   &req.BankTransactionId,
		ReferenceType: &model.TransactionReferenceTypeBankTransaction,
		InitiatedBy:   req.InitiatedBy,
	}
	if err = s.repos.Transaction.Create(transaction, wallet.Version); err != nil {
		return nil, errs.NewBadRequestError("failed to deposit", err)
	}
	return transaction, nil
}

func (s *transactionService) Withdraw(req *WithdrawRequest) (*model.Transaction, error) {
	if err := validator.GetValidator().ValidateStruct(req); err != nil {
		fields := validator.GetValidator().GetValidationErrors(err)
		return nil, errs.NewValidationError("invalid withdraw request", fields)
	}
	wallet, err := s.repos.Wallet.GetWalletByUserID(req.UserID)
	if err != nil {
		return nil, errs.NewNotFoundError("wallet not found", err)
	}
	transaction := &model.Transaction{
		WalletID:      wallet.ID,
		Amount:        req.Amount,
		Type:          model.TransactionTypeWithdraw,
		ReferenceID:   &req.BankTransactionId,
		ReferenceType: &model.TransactionReferenceTypeBankTransaction,
		InitiatedBy:   req.InitiatedBy,
	}
	if err = s.repos.Transaction.Create(transaction, wallet.Version); err != nil {
		return nil, errs.NewBadRequestError("failed to withdraw", err)
	}
	return transaction, nil
}

func (s *transactionService) Refund(req *RefundRequest) (*model.Transaction, error) {
	if err := validator.GetValidator().ValidateStruct(req); err != nil {
		fields := validator.GetValidator().GetValidationErrors(err)
		return nil, errs.NewValidationError("invalid refund request", fields)
	}
	wallet, err := s.repos.Wallet.GetWalletByUserID(req.UserID)
	if err != nil {
		return nil, errs.NewNotFoundError("wallet not found", err)
	}
	transaction := &model.Transaction{
		WalletID:      wallet.ID,
		Amount:        req.Amount,
		Type:          model.TransactionTypeRefund,
		ReferenceID:   &req.OrderId,
		ReferenceType: &model.TransactionReferenceTypeOrder,
		InitiatedBy:   req.InitiatedBy,
	}
	if err = s.repos.Transaction.Create(transaction, wallet.Version); err != nil {
		return nil, errs.NewBadRequestError("failed to refund", err)
	}
	return transaction, nil
}

func (s *transactionService) Purchase(req *PurchaseRequest) (*model.Transaction, error) {
	if err := validator.GetValidator().ValidateStruct(req); err != nil {
		fields := validator.GetValidator().GetValidationErrors(err)
		return nil, errs.NewValidationError("invalid purchase request", fields)
	}
	wallet, err := s.repos.Wallet.GetWalletByUserID(req.UserID)
	if err != nil {
		return nil, errs.NewNotFoundError("wallet not found", err)
	}
	transaction := &model.Transaction{
		WalletID:      wallet.ID,
		Amount:        req.Amount,
		Type:          model.TransactionTypePurchase,
		ReferenceID:   &req.OrderId,
		ReferenceType: &model.TransactionReferenceTypeOrder,
		InitiatedBy:   req.InitiatedBy,
	}
	if err = s.repos.Transaction.Create(transaction, wallet.Version); err != nil {
		return nil, errs.NewBadRequestError("failed to purchase", err)
	}
	return transaction, nil
}

func (s *transactionService) Transfer(req *TransferRequest) (*model.Transaction, error) {
	if err := validator.GetValidator().ValidateStruct(req); err != nil {
		fields := validator.GetValidator().GetValidationErrors(err)
		return nil, errs.NewValidationError("invalid transfer request", fields)
	}
	fromWallet, err := s.repos.Wallet.GetWalletByUserID(req.FromUserID)
	if err != nil {
		return nil, errs.NewNotFoundError("from wallet not found", err)
	}
	toWallet, err := s.repos.Wallet.GetWalletByUserID(req.ToUserID)
	if err != nil {
		return nil, errs.NewNotFoundError("to wallet not found", err)
	}
	transactionOut := &model.Transaction{
		WalletID:    fromWallet.ID,
		Amount:      -req.Amount,
		Type:        model.TransactionTypeTransferOut,
		InitiatedBy: req.InitiatedBy,
	}
	transactionIn := &model.Transaction{
		WalletID:      toWallet.ID,
		Amount:        req.Amount,
		Type:          model.TransactionTypeTransferIn,
		ReferenceType: &model.TransactionReferenceTypeTransfer,
		InitiatedBy:   req.InitiatedBy,
	}
	if err = s.repos.Transaction.Transfer(transactionOut, fromWallet.Version, transactionIn, toWallet.Version); err != nil {
		return nil, errs.NewBadRequestError("failed to transfer", err)
	}
	return transactionOut, nil
}

func (s *transactionService) GetTransactionsByWalletID(walletId string, page int, limit int) (*api.List[model.Transaction], error) {
	transactions, err := s.repos.Transaction.GetTransactionsByWalletID(walletId, page, limit)
	if err != nil {
		return nil, errs.NewInternalError("failed to get transactions", err)
	}
	total, err := s.repos.Transaction.GetTotalTransactionsByWalletID(walletId)
	if err != nil {
		return nil, errs.NewInternalError("failed to get total transactions", err)
	}
	return &api.List[model.Transaction]{Items: transactions, Total: total}, nil
}

func (s *transactionService) GetTransactionsSumByWalletID(walletId string) (float64, error) {
	sum, err := s.repos.Transaction.GetTransactionsSumByWalletID(walletId)
	if err != nil {
		return 0, errs.NewInternalError("failed to get transactions sum", err)
	}
	return sum, nil
}

func (s *transactionService) GetTransactionsSum() (float64, error) {
	sum, err := s.repos.Transaction.GetTransactionsSum()
	if err != nil {
		return 0, errs.NewInternalError("failed to get transactions sum", err)
	}
	return sum, nil
}
