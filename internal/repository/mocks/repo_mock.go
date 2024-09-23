package repository_mock

import (
	"github.com/abdelrahman146/digital-wallet/internal/repository"
	"go.uber.org/mock/gomock"
	"testing"
)

func NewRepoMock(t *testing.T) *repository.Repos {
	ctrl := gomock.NewController(t)
	return &repository.Repos{
		Audit:        NewMockAuditRepo(ctrl),
		Transaction:  NewMockTransactionRepo(ctrl),
		Account:      NewMockAccountRepo(ctrl),
		Wallet:       NewMockWalletRepo(ctrl),
		Tier:         NewMockTierRepo(ctrl),
		User:         NewMockUserRepo(ctrl),
		ExchangeRate: NewMockExchangeRateRepo(ctrl),
		Program:      NewMockProgramRepo(ctrl),
		Trigger:      NewMockTriggerRepo(ctrl),
	}
}
