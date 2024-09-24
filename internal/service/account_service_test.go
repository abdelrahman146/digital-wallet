package service

import (
	"context"
	"github.com/abdelrahman146/digital-wallet/internal/model"
	"github.com/abdelrahman146/digital-wallet/pkg/api"
	"github.com/abdelrahman146/digital-wallet/pkg/errs"
	"go.uber.org/mock/gomock"
	"testing"
)

func setupAccountServiceTest(t *testing.T) (*gomock.Controller, *serviceMocks, AccountService) {
	ctrl := gomock.NewController(t)
	mocks := newServiceMocks(ctrl)
	accountService := NewAccountService(mocks.repos)
	return ctrl, mocks, accountService
}

func TestAccountService_CreateAccount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mocks := newServiceMocks(ctrl)
	accountService := NewAccountService(mocks.repos)

	t.Run("success case", func(t *testing.T) {
		userId := "user123"
		walletId := "wallet123"
		ctx := api.CreateAppContext(context.Background(), api.AppActorUser, userId, "mock-123")
		mocks.userRepo.EXPECT().FetchUserByID(ctx, userId).Return(&model.User{ID: userId}, nil)
		mocks.walletRepo.EXPECT().FetchWalletByID(ctx, walletId).Return(&model.Wallet{ID: walletId}, nil)
		mocks.accountRepo.EXPECT().CreateAccount(ctx, gomock.Any()).Return(nil)

		account, err := accountService.CreateAccount(ctx, walletId, userId)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if account == nil {
			t.Error("expected valid account, got nil")
		}
	})

	t.Run("user not found", func(t *testing.T) {
		userId := "user123"
		walletId := "wallet123"
		ctx := api.CreateAppContext(context.Background(), api.AppActorUser, userId, "mock-123")
		mocks.userRepo.EXPECT().FetchUserByID(ctx, userId).Return(nil, nil)

		account, err := accountService.CreateAccount(ctx, walletId, userId)
		testExpectError(t, err, "USER_NOT_FOUND")
		if account != nil {
			t.Error("expected nil account, got valid")
		}
	})

	t.Run("wallet not found", func(t *testing.T) {
		userId := "user123"
		walletId := "wallet123"
		ctx := api.CreateAppContext(context.Background(), api.AppActorUser, userId, "mock-123")
		mocks.userRepo.EXPECT().FetchUserByID(ctx, userId).Return(&model.User{ID: userId}, nil)
		mocks.walletRepo.EXPECT().FetchWalletByID(ctx, walletId).Return(nil, nil)

		account, err := accountService.CreateAccount(ctx, walletId, userId)
		testExpectError(t, err, "WALLET_NOT_FOUND")
		if account != nil {
			t.Error("expected nil account, got valid")
		}
	})

	t.Run("unauthorized user", func(t *testing.T) {
		userId := "user123"
		walletId := "wallet123"
		ctx := api.CreateAppContext(context.Background(), api.AppActorUser, "another-user", "mock-123")
		mocks.userRepo.EXPECT().FetchUserByID(ctx, userId).Return(&model.User{ID: userId}, nil)
		account, err := accountService.CreateAccount(ctx, walletId, userId)
		testExpectError(t, err, "UNAUTHORIZED")
		if account != nil {
			t.Error("expected nil account, got valid")
		}
	})

	t.Run("repository create account error", func(t *testing.T) {
		userId := "user123"
		walletId := "wallet123"
		ctx := api.CreateAppContext(context.Background(), api.AppActorUser, userId, "mock-123")
		mocks.userRepo.EXPECT().FetchUserByID(ctx, userId).Return(&model.User{ID: userId}, nil)
		mocks.walletRepo.EXPECT().FetchWalletByID(ctx, walletId).Return(&model.Wallet{ID: walletId}, nil)
		mocks.accountRepo.EXPECT().CreateAccount(ctx, gomock.Any()).Return(errs.NewInternalError("db error", "", nil))

		account, err := accountService.CreateAccount(ctx, walletId, userId)
		testExpectError(t, err, "INTERNAL_ERROR")
		if account != nil {
			t.Error("expected nil account, got valid")
		}
	})
}

func TestAccountService_GetAccount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mocks := newServiceMocks(ctrl)
	accountService := NewAccountService(mocks.repos)

	t.Run("success case", func(t *testing.T) {
		accountId := "account123"
		userId := "user123"
		ctx := api.CreateAppContext(context.Background(), api.AppActorUser, userId, "mock-123")
		mocks.accountRepo.EXPECT().FetchAccountByID(ctx, accountId).Return(&model.Account{ID: accountId, UserID: userId}, nil)

		account, err := accountService.GetAccount(ctx, accountId)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if account == nil {
			t.Error("expected valid account, got nil")
		}
	})

	t.Run("account not found", func(t *testing.T) {
		accountId := "account123"
		ctx := api.CreateAppContext(context.Background(), api.AppActorUser, "user123", "mock-123")
		mocks.accountRepo.EXPECT().FetchAccountByID(ctx, accountId).Return(nil, nil)

		account, err := accountService.GetAccount(ctx, accountId)
		testExpectError(t, err, "ACCOUNT_NOT_FOUND")
		if account != nil {
			t.Error("expected nil account, got valid")
		}
	})

	t.Run("unauthorized user", func(t *testing.T) {
		accountId := "account123"
		userId := "user123"
		ctx := api.CreateAppContext(context.Background(), api.AppActorUser, "differentUser", "mock-123")
		mocks.accountRepo.EXPECT().FetchAccountByID(ctx, accountId).Return(&model.Account{ID: accountId, UserID: userId}, nil)

		account, err := accountService.GetAccount(ctx, accountId)
		testExpectError(t, err, "UNAUTHORIZED")
		if account != nil {
			t.Error("expected nil account, got valid")
		}
	})
}

func TestAccountService_GetWalletAccounts(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mocks := newServiceMocks(ctrl)
	accountService := NewAccountService(mocks.repos)

	t.Run("success case", func(t *testing.T) {
		walletId := "wallet123"
		ctx := api.CreateAppContext(context.Background(), api.AppActorAdmin, "admin123", "mock-123")
		mocks.walletRepo.EXPECT().FetchWalletByID(ctx, walletId).Return(&model.Wallet{ID: walletId}, nil)
		mocks.accountRepo.EXPECT().FetchWalletAccounts(ctx, walletId, 1, 10).Return([]model.Account{
			{ID: "account1", WalletID: walletId},
			{ID: "account2", WalletID: walletId},
		}, nil)
		mocks.accountRepo.EXPECT().CountWalletAccounts(ctx, walletId).Return(int64(2), nil)

		accounts, err := accountService.GetWalletAccounts(ctx, walletId, 1, 10)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if len(accounts.Items) != 2 {
			t.Errorf("expected 2 accounts, got %v", len(accounts.Items))
		}
	})

	t.Run("wallet not found", func(t *testing.T) {
		walletId := "wallet123"
		ctx := api.CreateAppContext(context.Background(), api.AppActorAdmin, "admin123", "mock-123")
		mocks.walletRepo.EXPECT().FetchWalletByID(ctx, walletId).Return(nil, nil)

		accounts, err := accountService.GetWalletAccounts(ctx, walletId, 1, 10)
		testExpectError(t, err, "WALLET_NOT_FOUND")
		if accounts != nil {
			t.Error("expected nil accounts, got valid")
		}
	})

	t.Run("unauthorized user", func(t *testing.T) {
		walletId := "wallet123"
		ctx := api.CreateAppContext(context.Background(), api.AppActorUser, "user123", "mock-123")

		accounts, err := accountService.GetWalletAccounts(ctx, walletId, 1, 10)
		testExpectError(t, err, "UNAUTHORIZED")
		if accounts != nil {
			t.Error("expected nil accounts, got valid")
		}
	})
}

func TestAccountService_GetWalletAccountsSum(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mocks := newServiceMocks(ctrl)
	accountService := NewAccountService(mocks.repos)

	t.Run("success case", func(t *testing.T) {
		walletId := "wallet123"
		ctx := api.CreateAppContext(context.Background(), api.AppActorAdmin, "admin123", "mock-123")
		mocks.walletRepo.EXPECT().FetchWalletByID(ctx, walletId).Return(&model.Wallet{ID: walletId}, nil)
		mocks.accountRepo.EXPECT().SumWalletAccounts(ctx, walletId).Return(uint64(1000), nil)

		sum, err := accountService.GetWalletAccountsSum(ctx, walletId)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if sum != 1000 {
			t.Errorf("expected sum to be 1000, got %v", sum)
		}
	})

	t.Run("wallet not found", func(t *testing.T) {
		walletId := "wallet123"
		ctx := api.CreateAppContext(context.Background(), api.AppActorAdmin, "admin123", "mock-123")
		mocks.walletRepo.EXPECT().FetchWalletByID(ctx, walletId).Return(nil, nil)

		sum, err := accountService.GetWalletAccountsSum(ctx, walletId)
		testExpectError(t, err, "WALLET_NOT_FOUND")
		if sum != 0 {
			t.Errorf("expected sum to be 0, got %v", sum)
		}
	})

	t.Run("unauthorized user", func(t *testing.T) {
		walletId := "wallet123"
		ctx := api.CreateAppContext(context.Background(), api.AppActorUser, "user123", "mock-123")

		sum, err := accountService.GetWalletAccountsSum(ctx, walletId)
		testExpectError(t, err, "UNAUTHORIZED")
		if sum != 0 {
			t.Errorf("expected sum to be 0, got %v", sum)
		}
	})
}

func TestAccountService_DeleteAccount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mocks := newServiceMocks(ctrl)
	accountService := NewAccountService(mocks.repos)

	t.Run("success case", func(t *testing.T) {
		accountId := "account123"
		ctx := api.CreateAppContext(context.Background(), api.AppActorAdmin, "admin123", "mock-123")
		mocks.accountRepo.EXPECT().FetchAccountByID(ctx, accountId).Return(&model.Account{ID: accountId}, nil)
		mocks.accountRepo.EXPECT().DeleteAccount(ctx, gomock.Any()).Return(nil)

		err := accountService.DeleteAccount(ctx, accountId)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("account not found", func(t *testing.T) {
		accountId := "account123"
		ctx := api.CreateAppContext(context.Background(), api.AppActorAdmin, "admin123", "mock-123")
		mocks.accountRepo.EXPECT().FetchAccountByID(ctx, accountId).Return(nil, nil)

		err := accountService.DeleteAccount(ctx, accountId)
		testExpectError(t, err, "ACCOUNT_NOT_FOUND")
	})

	t.Run("unauthorized user", func(t *testing.T) {
		accountId := "account123"
		ctx := api.CreateAppContext(context.Background(), api.AppActorUser, "user123", "mock-123")

		err := accountService.DeleteAccount(ctx, accountId)
		testExpectError(t, err, "UNAUTHORIZED")
	})
}
