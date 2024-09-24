package service

import (
	"context"
	"github.com/abdelrahman146/digital-wallet/internal/repository"
	repository_mock "github.com/abdelrahman146/digital-wallet/internal/repository/mocks"
	"github.com/abdelrahman146/digital-wallet/pkg/api"
	"github.com/abdelrahman146/digital-wallet/pkg/errs"
	"go.uber.org/mock/gomock"
	"reflect"
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

type TestCase[Service any] struct {
	name          string
	ctx           context.Context
	setupMocks    func(mocks *Mocks, ctx context.Context)
	expectedError string
	expectResult  bool
	testFunc      func(service Service, ctx context.Context) (interface{}, error)
}

type Mocks struct {
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

func NewServiceMocks(ctrl *gomock.Controller) *Mocks {
	auditRepo := repository_mock.NewMockAuditRepo(ctrl)
	accountRepo := repository_mock.NewMockAccountRepo(ctrl)
	transactionRepo := repository_mock.NewMockTransactionRepo(ctrl)
	walletRepo := repository_mock.NewMockWalletRepo(ctrl)
	userRepo := repository_mock.NewMockUserRepo(ctrl)
	tierRepo := repository_mock.NewMockTierRepo(ctrl)
	exchangeRateRepo := repository_mock.NewMockExchangeRateRepo(ctrl)
	programRepo := repository_mock.NewMockProgramRepo(ctrl)
	triggerRepo := repository_mock.NewMockTriggerRepo(ctrl)
	return &Mocks{
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

func SetupTest[Service any](t *testing.T, serviceFactory func(*Mocks) Service) (*gomock.Controller, *Mocks, Service) {
	ctrl := gomock.NewController(t)
	mocks := NewServiceMocks(ctrl)
	service := serviceFactory(mocks)
	return ctrl, mocks, service
}

func RunTestCases[Service any](t *testing.T, serviceFactory func(*Mocks) Service, testCases []TestCase[Service]) {
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl, mocks, service := SetupTest(t, serviceFactory)
			defer ctrl.Finish()
			ctx := context.Background()
			if tc.ctx == nil {
				// Create a basic context, you can customize per test if needed
				ctx = api.CreateAppContext(ctx, api.AppActorUser, test_userId, test_requestId)
			} else {
				ctx = tc.ctx
			}

			// Setup ServiceMocks for each test case
			tc.setupMocks(mocks, ctx)

			// Run the test function for each case
			result, err := tc.testFunc(service, ctx)

			// Validate error
			if tc.expectedError != "" {
				TestExpectError(t, err, tc.expectedError)
			} else if err != nil {
				t.Errorf("expected no error, got %v", err)
			}

			switch result {
			case nil:
				if tc.expectResult {
					t.Errorf("expected result, got %v", result)
				}
			default:
				// Validate result
				if reflect.ValueOf(result).IsZero() && tc.expectResult {
					t.Errorf("expected result, got %v", result)
				}

				// Validate no result
				if !reflect.ValueOf(result).IsZero() && !tc.expectResult {
					t.Errorf("expected no result, got %v", result)
				}
			}
		})
	}
}

func TestExpectError(t *testing.T, err error, expectedCode string) {
	if err == nil || errs.HandleError(err).Code != expectedCode {
		t.Errorf("expected error code %v, got %v", expectedCode, err)
	}
}
