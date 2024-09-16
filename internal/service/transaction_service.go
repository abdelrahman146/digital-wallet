package service

import (
	"context"
	"digital-wallet/internal/model"
	"digital-wallet/internal/repository"
	"digital-wallet/pkg/api"
	"digital-wallet/pkg/errs"
	"digital-wallet/pkg/logger"
	"digital-wallet/pkg/types"
	"digital-wallet/pkg/validator"
	"github.com/shopspring/decimal"
	"time"
)

type TransactionService interface {
	// CreateTransaction creates a transaction
	CreateTransaction(ctx context.Context, walletId, accountId, actorType, actorId string, req *TransactionRequest) (*model.Transaction, error)
	// Exchange exchanges between two accounts for the same user
	Exchange(ctx context.Context, fromWalletId, toWalletId, userId string, amount uint64) (*ExchangeResponse, error)
	// GetAccountTransactions returns a list of transactions for an account
	GetAccountTransactions(ctx context.Context, walletId, accountId string, page int, limit int) (*api.List[model.Transaction], error)
	// GetAccountTransactionSum returns the sum of transactions for an account
	GetAccountTransactionSum(ctx context.Context, walletId, accountId string) (uint64, error)
	// GetWalletTransactions returns a list of transactions for a wallet
	GetWalletTransactions(ctx context.Context, walletId string, page int, limit int) (*api.List[model.Transaction], error)
	// GetWalletTransactionSum returns the sum of transactions for a wallet
	GetWalletTransactionSum(ctx context.Context, walletId string) (uint64, error)
	// GetAccountExpiringTransactionsSum returns the sum of expiring transactions for an account
	GetAccountExpiringTransactionsSum(ctx context.Context, accountId string) (uint64, error)
	// GetExpiredWalletTransactions returns a list of expired transactions for a wallet
	GetExpiredWalletTransactions(ctx context.Context, walletId string) ([]model.Transaction, error)
}

type transactionService struct {
	repos *repository.Repos
}

func NewTransactionService(repos *repository.Repos) TransactionService {
	return &transactionService{repos: repos}
}

func (s *transactionService) CreateTransaction(ctx context.Context, walletId, accountId, actorType, actorId string, req *TransactionRequest) (*model.Transaction, error) {
	if err := validator.GetValidator().ValidateStruct(req); err != nil {
		fields := validator.GetValidator().GetValidationErrors(err)
		api.GetLogger(ctx).Error("Invalid transaction request", logger.Field("fields", fields), logger.Field("request", req))
		return nil, errs.NewValidationError("Invalid transaction request", "", fields)
	}
	wallet, err := s.repos.Wallet.FetchWalletByID(ctx, walletId)
	if wallet == nil {
		return nil, errs.NewNotFoundError("wallet not found", "WALLET_NOT_FOUND", err)
	}
	account, _ := s.repos.Account.FetchAccountByID(ctx, accountId)
	if account == nil {
		return nil, errs.NewNotFoundError("Account not found", "ACCOUNT_NOT_FOUND", err)
	}
	if err := api.IsAuthorizedUser(ctx, account.UserID); err != nil {
		api.GetLogger(ctx).Error("Unauthorized", logger.Field("userId", account.UserID))
		return nil, err
	}
	transaction := &model.Transaction{
		AccountID: accountId,
		WalletID:  walletId,
		Amount:    req.Amount,
		Reason:    req.Reason,
		Type:      req.Type,
		Metadata:  req.Metadata,
		ProgramID: req.ProgramID,
	}
	if req.Type == model.TransactionTypeCredit && wallet.PointsExpireAfter != nil {
		expireAt := time.Now().Add(wallet.PointsExpireAfter.Duration())
		transaction.ExpireAt = &expireAt
	}
	if err := s.repos.Transaction.CreateTransaction(ctx, transaction, account.Version); err != nil {
		return nil, err
	}
	return transaction, nil
}

func (s *transactionService) Exchange(ctx context.Context, fromWalletId, toWalletId, userId string, amount uint64) (*ExchangeResponse, error) {
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

	// Get Exchange Rate
	exchangeRate, err := s.repos.ExchangeRate.FetchExchangeRate(ctx, fromWalletId, toWalletId, user.TierID)
	if exchangeRate == nil {
		api.GetLogger(ctx).Error("Exchange Rate not found", logger.Field("fromWalletId", fromWalletId), logger.Field("toWalletId", toWalletId), logger.Field("tierId", user.TierID))
		return nil, errs.NewNotFoundError("Exchange Rate not found", "EXCHANGE_RATE_NOT_FOUND", err)
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
		Reason:    model.TransactionReasonExchange,
		Type:      model.TransactionTypeDebit,
	}
	fromTransaction.Metadata = make(types.JSONB)
	fromTransaction.Metadata["toWalletId"] = toWalletId
	fromTransaction.Metadata["toAccountId"] = toAccount.ID

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
		Reason:    model.TransactionReasonExchange,
	}
	toTransaction.Metadata = make(types.JSONB)
	toTransaction.Metadata["fromWalletId"] = fromWalletId
	toTransaction.Metadata["fromAccountId"] = fromAccount.ID
	toTransaction.Metadata["ExchangedAmount"] = amount
	toTransaction.Metadata["exchangeRate"] = exchangeRate.ExchangeRate.String()

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

	if err := s.repos.Transaction.PerformExchange(ctx, from, to); err != nil {
		return nil, err
	}
	resp := &ExchangeResponse{
		FromTransaction: *fromTransaction,
		ToTransaction:   *toTransaction,
	}
	return resp, nil
}

func (s *transactionService) GetAccountTransactions(ctx context.Context, walletId, accountId string, page int, limit int) (*api.List[model.Transaction], error) {
	wallet, err := s.repos.Wallet.FetchWalletByID(ctx, walletId)
	if wallet == nil {
		api.GetLogger(ctx).Error("Wallet not found", logger.Field("walletId", walletId))
		return nil, errs.NewNotFoundError("Wallet not found", "WALLET_NOT_FOUND", err)
	}
	account, err := s.repos.Account.FetchAccountByID(ctx, accountId)
	if account == nil {
		api.GetLogger(ctx).Error("Account not found", logger.Field("walletId", walletId), logger.Field("accountId", accountId))
		return nil, errs.NewNotFoundError("Account not found", "ACCOUNT_NOT_FOUND", err)
	}
	if err := api.IsAuthorizedUser(ctx, account.UserID); err != nil {
		api.GetLogger(ctx).Error("Unauthorized", logger.Field("userId", account.UserID))
		return nil, err
	}
	transactions, err := s.repos.Transaction.FetchAccountTransactions(ctx, accountId, page, limit)
	if err != nil {
		return nil, err
	}
	total, err := s.repos.Transaction.CountAccountTransactions(ctx, accountId)
	if err != nil {
		return nil, err
	}
	return &api.List[model.Transaction]{Items: transactions, Page: page, Limit: limit, Total: total}, nil
}

func (s *transactionService) GetAccountTransactionSum(ctx context.Context, walletId, accountId string) (uint64, error) {
	wallet, err := s.repos.Wallet.FetchWalletByID(ctx, walletId)
	if wallet == nil {
		api.GetLogger(ctx).Error("Wallet not found", logger.Field("walletId", walletId))
		return 0, errs.NewNotFoundError("Wallet not found", "WALLET_NOT_FOUND", err)
	}
	account, err := s.repos.Account.FetchAccountByID(ctx, accountId)
	if account == nil {
		api.GetLogger(ctx).Error("Account not found", logger.Field("walletId", walletId), logger.Field("accountId", accountId))
		return 0, errs.NewNotFoundError("Account not found", "ACCOUNT_NOT_FOUND", err)
	}
	if err := api.IsAuthorizedUser(ctx, account.UserID); err != nil {
		api.GetLogger(ctx).Error("Unauthorized", logger.Field("userId", account.UserID))
		return 0, err
	}
	sum, err := s.repos.Transaction.SumAccountTransactions(ctx, accountId)
	return sum, err
}

func (s *transactionService) GetWalletTransactions(ctx context.Context, walletId string, page int, limit int) (*api.List[model.Transaction], error) {
	if err := api.IsAdmin(ctx); err != nil {
		api.GetLogger(ctx).Error("Unauthorized", logger.Field("actor", api.GetActor(ctx)), logger.Field("userId", api.GetUserID(ctx)))
		return nil, err
	}
	wallet, err := s.repos.Wallet.FetchWalletByID(ctx, walletId)
	if wallet == nil {
		api.GetLogger(ctx).Error("Wallet not found", logger.Field("walletId", walletId))
		return nil, errs.NewNotFoundError("Wallet not found", "WALLET_NOT_FOUND", err)
	}
	transactions, err := s.repos.Transaction.FetchWalletTransactions(ctx, walletId, page, limit)
	if err != nil {
		return nil, err
	}
	total, err := s.repos.Transaction.CountWalletTransactions(ctx, walletId)
	if err != nil {
		return nil, err
	}
	return &api.List[model.Transaction]{Items: transactions, Page: page, Limit: limit, Total: total}, nil
}

func (s *transactionService) GetWalletTransactionSum(ctx context.Context, walletId string) (uint64, error) {
	if err := api.IsAdmin(ctx); err != nil {
		api.GetLogger(ctx).Error("Unauthorized", logger.Field("actor", api.GetActor(ctx)), logger.Field("userId", api.GetUserID(ctx)))
		return 0, err
	}
	wallet, err := s.repos.Wallet.FetchWalletByID(ctx, walletId)
	if wallet == nil {
		api.GetLogger(ctx).Error("Wallet not found", logger.Field("walletId", walletId))
		return 0, errs.NewNotFoundError("Wallet not found", "WALLET_NOT_FOUND", err)
	}
	sum, err := s.repos.Transaction.SumAccountTransactions(ctx, walletId)
	if err != nil {
		return 0, err
	}
	return sum, nil
}

func (s *transactionService) GetAccountExpiringTransactionsSum(ctx context.Context, accountId string) (uint64, error) {
	account, err := s.repos.Account.FetchAccountByID(ctx, accountId)
	if account == nil {
		api.GetLogger(ctx).Error("Account not found", logger.Field("accountId", accountId))
		return 0, errs.NewNotFoundError("Account not found", "ACCOUNT_NOT_FOUND", err)
	}
	if err := api.IsAuthorizedUser(ctx, account.UserID); err != nil {
		api.GetLogger(ctx).Error("Unauthorized", logger.Field("userId", account.UserID))
		return 0, err
	}
	wallet, err := s.repos.Wallet.FetchWalletByID(ctx, account.WalletID)
	if wallet == nil {
		api.GetLogger(ctx).Error("Wallet not found", logger.Field("walletId", account.WalletID))
		return 0, errs.NewNotFoundError("Wallet not found", "WALLET_NOT_FOUND", err)
	}
	if wallet.PointsExpireAfter == nil {
		api.GetLogger(ctx).Error("Wallet points expire after not set", logger.Field("walletId", account.WalletID))
		return 0, errs.NewBadRequestError("Wallet points expire after not set", "WALLET_POINTS_EXPIRE_AFTER_NOT_SET", nil)
	}
	sum, err := s.repos.Transaction.SumExpiringAccountTransactions(ctx, accountId, *wallet.PointsExpireAfter)
	if err != nil {
		return 0, err
	}
	return sum, nil
}

func (s *transactionService) GetExpiredWalletTransactions(ctx context.Context, walletId string) ([]model.Transaction, error) {
	if err := api.IsAdmin(ctx); err != nil {
		api.GetLogger(ctx).Error("Unauthorized", logger.Field("actor", api.GetActor(ctx)), logger.Field("userId", api.GetUserID(ctx)))
		return nil, err
	}
	wallet, err := s.repos.Wallet.FetchWalletByID(ctx, walletId)
	if wallet == nil {
		api.GetLogger(ctx).Error("Wallet not found", logger.Field("walletId", walletId))
		return nil, errs.NewNotFoundError("Wallet not found", "WALLET_NOT_FOUND", err)
	}
	transactions, err := s.repos.Transaction.FetchExpiredWalletTransactions(ctx, walletId)
	if err != nil {
		return nil, err
	}
	return transactions, nil
}
