package service

import (
	"context"
	"github.com/abdelrahman146/digital-wallet/internal/model"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestAccountService_CreateAccount(t *testing.T) {
	testcases := []TestCase[AccountService]{
		{
			name: "CreateAccount: Success case",
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
	}
	serviceFactory := func(mocks *Mocks) AccountService {
		return NewAccountService(mocks.repos)
	}
	RunTestCases[AccountService](t, serviceFactory, testcases)
}
