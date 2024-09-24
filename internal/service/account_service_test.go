package service

import (
	"context"
	"github.com/abdelrahman146/digital-wallet/internal/model"
	"github.com/abdelrahman146/digital-wallet/pkg/api"
	"github.com/abdelrahman146/digital-wallet/pkg/errs"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestAccountService_CreateAccount(t *testing.T) {
	testcases := []TestCase[AccountService]{
		{
			name: "Success case",
			setupMocks: func(mocks *Mocks, ctx context.Context) {
				mocks.userRepo.EXPECT().FetchUserByID(ctx, test_userId).Return(&model.User{ID: test_userId}, nil)
				mocks.walletRepo.EXPECT().FetchWalletByID(ctx, test_walletId).Return(&model.Wallet{ID: test_walletId}, nil)
				mocks.accountRepo.EXPECT().CreateAccount(ctx, gomock.Any()).Return(nil)
			},
			testFunc: func(service AccountService, ctx context.Context) (interface{}, error) {
				return service.CreateAccount(ctx, test_walletId, test_userId)
			},
			expectResult: true,
		},
		{
			name: "User not found",
			setupMocks: func(mocks *Mocks, ctx context.Context) {
				mocks.userRepo.EXPECT().FetchUserByID(ctx, test_userId).Return(nil, nil)
			},
			expectedError: "USER_NOT_FOUND",
			testFunc: func(service AccountService, ctx context.Context) (interface{}, error) {
				return service.CreateAccount(ctx, test_walletId, test_userId)
			},
			expectResult: false,
		},
		{
			name: "User not authorized",
			ctx:  api.CreateAppContext(context.Background(), api.AppActorUser, "unauthorized-user", test_requestId),
			setupMocks: func(mocks *Mocks, ctx context.Context) {
				mocks.userRepo.EXPECT().FetchUserByID(ctx, test_userId).Return(&model.User{ID: test_userId}, nil)
			},
			expectedError: "UNAUTHORIZED",
			testFunc: func(service AccountService, ctx context.Context) (interface{}, error) {
				return service.CreateAccount(ctx, test_walletId, test_userId)
			},
			expectResult: false,
		},
		{
			name: "Wallet not found",
			setupMocks: func(mocks *Mocks, ctx context.Context) {
				mocks.userRepo.EXPECT().FetchUserByID(ctx, test_userId).Return(&model.User{ID: test_userId}, nil)
				mocks.walletRepo.EXPECT().FetchWalletByID(ctx, test_walletId).Return(nil, nil)
			},
			expectedError: "WALLET_NOT_FOUND",
			testFunc: func(service AccountService, ctx context.Context) (interface{}, error) {
				return service.CreateAccount(ctx, test_walletId, test_userId)
			},
			expectResult: false,
		},
		{
			name: "Repo Error Creating Account",
			setupMocks: func(mocks *Mocks, ctx context.Context) {
				mocks.userRepo.EXPECT().FetchUserByID(ctx, test_userId).Return(&model.User{ID: test_userId}, nil)
				mocks.walletRepo.EXPECT().FetchWalletByID(ctx, test_walletId).Return(&model.Wallet{ID: test_walletId}, nil)
				mocks.accountRepo.EXPECT().CreateAccount(ctx, gomock.Any()).Return(errs.NewInternalError("DB error", "", nil))
			},
			expectedError: "INTERNAL_ERROR",
			testFunc: func(service AccountService, ctx context.Context) (interface{}, error) {
				return service.CreateAccount(ctx, test_walletId, test_userId)
			},
			expectResult: false,
		},
	}
	serviceFactory := func(mocks *Mocks) AccountService {
		return NewAccountService(mocks.repos)
	}
	RunTestCases[AccountService](t, serviceFactory, testcases)
}

func TestAccountService_GetAccount(t *testing.T) {
	testcases := []TestCase[AccountService]{
		{
			name: "Success case",
			setupMocks: func(mocks *Mocks, ctx context.Context) {
				mocks.accountRepo.EXPECT().FetchAccountByID(ctx, test_accountId).Return(&model.Account{ID: test_accountId, UserID: test_userId}, nil)
			},
			testFunc: func(service AccountService, ctx context.Context) (interface{}, error) {
				return service.GetAccount(ctx, test_accountId)
			},
			expectResult: true,
		},
		{
			name: "Account not found",
			setupMocks: func(mocks *Mocks, ctx context.Context) {
				mocks.accountRepo.EXPECT().FetchAccountByID(ctx, test_accountId).Return(nil, nil)
			},
			expectedError: "ACCOUNT_NOT_FOUND",
			testFunc: func(service AccountService, ctx context.Context) (interface{}, error) {
				return service.GetAccount(ctx, test_accountId)
			},
			expectResult: false,
		},
		{
			name: "User not authorized",
			ctx:  api.CreateAppContext(context.Background(), api.AppActorUser, "unauthorized-user", test_requestId),
			setupMocks: func(mocks *Mocks, ctx context.Context) {
				mocks.accountRepo.EXPECT().FetchAccountByID(ctx, test_accountId).Return(&model.Account{ID: test_accountId, UserID: test_userId}, nil)
			},
			expectedError: "UNAUTHORIZED",
			testFunc: func(service AccountService, ctx context.Context) (interface{}, error) {
				return service.GetAccount(ctx, test_accountId)
			},
			expectResult: false,
		},
	}
	serviceFactory := func(mocks *Mocks) AccountService {
		return NewAccountService(mocks.repos)
	}
	RunTestCases[AccountService](t, serviceFactory, testcases)
}

func TestAccountService_DeleteAccount(t *testing.T) {
	testcases := []TestCase[AccountService]{
		{
			name: "Success case",
			ctx:  api.CreateAppContext(context.Background(), api.AppActorAdmin, test_userId, test_requestId),
			setupMocks: func(mocks *Mocks, ctx context.Context) {
				mocks.accountRepo.EXPECT().FetchAccountByID(ctx, test_accountId).Return(&model.Account{ID: test_accountId, UserID: test_userId}, nil)
				mocks.accountRepo.EXPECT().DeleteAccount(ctx, gomock.Any()).Return(nil)
			},
			testFunc: func(service AccountService, ctx context.Context) (interface{}, error) {
				return nil, service.DeleteAccount(ctx, test_accountId)
			},
			expectResult: false,
		},
		{
			name: "Account not found",
			ctx:  api.CreateAppContext(context.Background(), api.AppActorAdmin, test_userId, test_requestId),
			setupMocks: func(mocks *Mocks, ctx context.Context) {
				mocks.accountRepo.EXPECT().FetchAccountByID(ctx, test_accountId).Return(nil, nil)
			},
			expectedError: "ACCOUNT_NOT_FOUND",
			testFunc: func(service AccountService, ctx context.Context) (interface{}, error) {
				return nil, service.DeleteAccount(ctx, test_accountId)
			},
			expectResult: false,
		},
		{
			name:          "User must be admin",
			setupMocks:    func(mocks *Mocks, ctx context.Context) {},
			expectedError: "UNAUTHORIZED",
			testFunc: func(service AccountService, ctx context.Context) (interface{}, error) {
				return nil, service.DeleteAccount(ctx, test_accountId)
			},
			expectResult: false,
		},
		{
			name: "Repo Error Deleting Account",
			ctx:  api.CreateAppContext(context.Background(), api.AppActorAdmin, test_userId, test_requestId),
			setupMocks: func(mocks *Mocks, ctx context.Context) {
				mocks.accountRepo.EXPECT().FetchAccountByID(ctx, test_accountId).Return(&model.Account{ID: test_accountId, UserID: test_userId}, nil)
				mocks.accountRepo.EXPECT().DeleteAccount(ctx, gomock.Any()).Return(errs.NewInternalError("DB error", "", nil))
			},
			expectedError: "INTERNAL_ERROR",
			testFunc: func(service AccountService, ctx context.Context) (interface{}, error) {
				return nil, service.DeleteAccount(ctx, test_accountId)
			},
			expectResult: false,
		},
	}
	serviceFactory := func(mocks *Mocks) AccountService {
		return NewAccountService(mocks.repos)
	}
	RunTestCases[AccountService](t, serviceFactory, testcases)
}

func TestAccountService_GetWalletAccountsSum(t *testing.T) {
	testcases := []TestCase[AccountService]{
		{
			name: "Success case",
			ctx:  api.CreateAppContext(context.Background(), api.AppActorAdmin, test_userId, test_requestId),
			setupMocks: func(mocks *Mocks, ctx context.Context) {
				mocks.walletRepo.EXPECT().FetchWalletByID(ctx, test_walletId).Return(&model.Wallet{ID: test_walletId}, nil)
				mocks.accountRepo.EXPECT().SumWalletAccounts(ctx, test_walletId).Return(uint64(100), nil)
			},
			testFunc: func(service AccountService, ctx context.Context) (interface{}, error) {
				return service.GetWalletAccountsSum(ctx, test_walletId)
			},
			expectResult: true,
		},
		{
			name: "Wallet not found",
			ctx:  api.CreateAppContext(context.Background(), api.AppActorAdmin, test_userId, test_requestId),
			setupMocks: func(mocks *Mocks, ctx context.Context) {
				mocks.walletRepo.EXPECT().FetchWalletByID(ctx, test_walletId).Return(nil, nil)
			},
			expectedError: "WALLET_NOT_FOUND",
			testFunc: func(service AccountService, ctx context.Context) (interface{}, error) {
				return service.GetWalletAccountsSum(ctx, test_walletId)
			},
			expectResult: false,
		},
		{
			name:          "User must be admin",
			setupMocks:    func(mocks *Mocks, ctx context.Context) {},
			expectedError: "UNAUTHORIZED",
			testFunc: func(service AccountService, ctx context.Context) (interface{}, error) {
				return service.GetWalletAccountsSum(ctx, test_walletId)
			},
			expectResult: false,
		},
		{
			name: "Repo Error Summing Accounts",
			ctx:  api.CreateAppContext(context.Background(), api.AppActorAdmin, test_userId, test_requestId),
			setupMocks: func(mocks *Mocks, ctx context.Context) {
				mocks.walletRepo.EXPECT().FetchWalletByID(ctx, test_walletId).Return(&model.Wallet{ID: test_walletId}, nil)
				mocks.accountRepo.EXPECT().SumWalletAccounts(ctx, test_walletId).Return(uint64(0), errs.NewInternalError("DB error", "", nil))
			},
			expectedError: "INTERNAL_ERROR",
			testFunc: func(service AccountService, ctx context.Context) (interface{}, error) {
				return service.GetWalletAccountsSum(ctx, test_walletId)
			},
			expectResult: false,
		},
	}
	serviceFactory := func(mocks *Mocks) AccountService {
		return NewAccountService(mocks.repos)
	}
	RunTestCases[AccountService](t, serviceFactory, testcases)
}
