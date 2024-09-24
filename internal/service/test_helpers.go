package service

import (
	"context"
	"github.com/abdelrahman146/digital-wallet/internal/repository"
	repository_mock "github.com/abdelrahman146/digital-wallet/internal/repository/mocks"
	"github.com/abdelrahman146/digital-wallet/pkg/api"
	"github.com/abdelrahman146/digital-wallet/pkg/errs"
	"go.uber.org/mock/gomock"
	"testing"
)

// Constants for test data
const (
	test_userId    = "user-123"
	test_walletId  = "wallet-123"
	test_accountId = "account-123"
	test_adminId   = "admin-123"
	test_requestId = "mock-123"
)

type TestCase struct {
	name          string
	setupMocks    func(mocks *serviceMocks, ctx context.Context)
	expectedError string
	expectedValid bool
}

type serviceMocks struct {
	auditRepo        *repository_mock.MockAuditRepo
	accountRepo      *repository_mock.MockAccountRepo
	transactionRepo  *repository_mock.MockTransactionRepo
	walletRepo       *repository_mock.MockWalletRepo
	userRepo         *repository_mock.MockUserRepo
	tierRepo         *repository_mock.MockTierRepo
	exchangeRateRepo *repository_mock.MockExchangeRateRepo
	programRepo      *repository_mock.MockProgramRepo
	triggerRepo      *repository_mock.MockTriggerRepo
	repos            *repository.Repos
}

func newServiceMocks(ctrl *gomock.Controller) *serviceMocks {
	auditRepo := repository_mock.NewMockAuditRepo(ctrl)
	accountRepo := repository_mock.NewMockAccountRepo(ctrl)
	transactionRepo := repository_mock.NewMockTransactionRepo(ctrl)
	walletRepo := repository_mock.NewMockWalletRepo(ctrl)
	userRepo := repository_mock.NewMockUserRepo(ctrl)
	tierRepo := repository_mock.NewMockTierRepo(ctrl)
	exchangeRateRepo := repository_mock.NewMockExchangeRateRepo(ctrl)
	programRepo := repository_mock.NewMockProgramRepo(ctrl)
	triggerRepo := repository_mock.NewMockTriggerRepo(ctrl)
	return &serviceMocks{
		auditRepo:        auditRepo,
		accountRepo:      accountRepo,
		transactionRepo:  transactionRepo,
		walletRepo:       walletRepo,
		userRepo:         userRepo,
		tierRepo:         tierRepo,
		exchangeRateRepo: exchangeRateRepo,
		programRepo:      programRepo,
		triggerRepo:      triggerRepo,
		repos: &repository.Repos{
			Audit:        auditRepo,
			Account:      accountRepo,
			Transaction:  transactionRepo,
			Wallet:       walletRepo,
			User:         userRepo,
			Tier:         tierRepo,
			ExchangeRate: exchangeRateRepo,
			Program:      programRepo,
			Trigger:      triggerRepo,
		},
	}
}

func testExpectError(t *testing.T, err error, expectedCode string) {
	if err == nil || errs.HandleError(err).Code != expectedCode {
		t.Errorf("expected error code %v, got %v", expectedCode, err)
	}
}

func runTestCases(t *testing.T, testCases []TestCase, testFunc func() (interface{}, error)) {
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl, mocks, accountService := setupTest(t)
			defer ctrl.Finish()
			ctx := api.CreateAppContext(context.Background(), api.AppActorUser, "user123", "mock-123")

			// Setup Mocks
			tc.setupMocks(mocks, ctx)

			// Run test function
			result, err := testFunc()

			// Validate error
			if tc.expectedError != "" {
				testExpectError(t, err, tc.expectedError)
			} else if err != nil {
				t.Errorf("expected no error, got %v", err)
			}

			// Validate result
			if tc.expectedValid && result == nil {
				t.Error("expected valid result, got nil")
			}
		})
	}
}
