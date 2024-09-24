package service_mock

import (
	"github.com/abdelrahman146/digital-wallet/internal/service"
	"go.uber.org/mock/gomock"
	"testing"
)

func NewServiceMock(t *testing.T) (*service.Services, *gomock.Controller) {
	ctrl := gomock.NewController(t)
	return &service.Services{
		Audit:        NewMockAuditService(ctrl),
		Transaction:  NewMockTransactionService(ctrl),
		Account:      NewMockAccountService(ctrl),
		Wallet:       NewMockWalletService(ctrl),
		Tier:         NewMockTierService(ctrl),
		User:         NewMockUserService(ctrl),
		ExchangeRate: NewMockExchangeRateService(ctrl),
		Trigger:      NewMockTriggerService(ctrl),
		Program:      NewMockProgramService(ctrl),
	}, ctrl
}

func CloseServiceMock(ctrl *gomock.Controller) {
	ctrl.Finish()
}
