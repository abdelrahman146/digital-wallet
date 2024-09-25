package service

import (
	"context"
	"github.com/abdelrahman146/digital-wallet/internal/model"
	"github.com/abdelrahman146/digital-wallet/pkg/api"
	"testing"
)

func TestAuditService_GetActorAuditLogs(t *testing.T) {
	testcases := []TestCase[AuditService]{
		{
			name: "success case",
			ctx:  api.CreateAppContext(context.Background(), api.AppActorAdmin, test_userId, test_requestId),
			setupMocks: func(mocks *Mocks, ctx context.Context) {
				mocks.auditRepo.EXPECT().FetchActorAuditLogs(ctx, api.AppActorAdmin, test_userId, 1, 10).Return([]model.Audit{
					{
						ID:        "123",
						Operation: model.AuditOperationCreate,
						Actor:     api.AppActorUser,
						ActorID:   test_userId,
						Table:     "programs",
						RecordID:  "12",
					},
				}, nil)
				mocks.auditRepo.EXPECT().CountActorAuditLogs(ctx, api.AppActorAdmin, test_userId).Return(int64(1), nil)
			},
			testFunc: func(service AuditService, ctx context.Context) (interface{}, error) {
				return service.GetActorAuditLogs(ctx, api.AppActorAdmin, test_userId, 1, 10)
			},
			expectResult: true,
		},
	}
	auditServiceFactory := func(mocks *Mocks) AuditService {
		return NewAuditService(mocks.repos)
	}
	RunTestCases[AuditService](t, auditServiceFactory, testcases)
}
