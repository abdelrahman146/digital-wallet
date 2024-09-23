// Code generated by MockGen. DO NOT EDIT.
// Source: internal/service/trigger_service.go
//
// Generated by this command:
//
//	mockgen -source=internal/service/trigger_service.go -destination=internal/service/mocks/trigger_service_mock.go -package=service_mock
//

// Package service_mock is a generated GoMock package.
package service_mock

import (
	context "context"
	reflect "reflect"

	model "github.com/abdelrahman146/digital-wallet/internal/model"
	service "github.com/abdelrahman146/digital-wallet/internal/service"
	api "github.com/abdelrahman146/digital-wallet/pkg/api"
	gomock "go.uber.org/mock/gomock"
)

// MockTriggerService is a mock of TriggerService interface.
type MockTriggerService struct {
	ctrl     *gomock.Controller
	recorder *MockTriggerServiceMockRecorder
}

// MockTriggerServiceMockRecorder is the mock recorder for MockTriggerService.
type MockTriggerServiceMockRecorder struct {
	mock *MockTriggerService
}

// NewMockTriggerService creates a new mock instance.
func NewMockTriggerService(ctrl *gomock.Controller) *MockTriggerService {
	mock := &MockTriggerService{ctrl: ctrl}
	mock.recorder = &MockTriggerServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTriggerService) EXPECT() *MockTriggerServiceMockRecorder {
	return m.recorder
}

// CreateTrigger mocks base method.
func (m *MockTriggerService) CreateTrigger(ctx context.Context, req service.CreateTriggerRequest) (*model.Trigger, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateTrigger", ctx, req)
	ret0, _ := ret[0].(*model.Trigger)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateTrigger indicates an expected call of CreateTrigger.
func (mr *MockTriggerServiceMockRecorder) CreateTrigger(ctx, req any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTrigger", reflect.TypeOf((*MockTriggerService)(nil).CreateTrigger), ctx, req)
}

// DeleteTrigger mocks base method.
func (m *MockTriggerService) DeleteTrigger(ctx context.Context, triggerId uint64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteTrigger", ctx, triggerId)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteTrigger indicates an expected call of DeleteTrigger.
func (mr *MockTriggerServiceMockRecorder) DeleteTrigger(ctx, triggerId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteTrigger", reflect.TypeOf((*MockTriggerService)(nil).DeleteTrigger), ctx, triggerId)
}

// GetTrigger mocks base method.
func (m *MockTriggerService) GetTrigger(ctx context.Context, triggerId uint64) (*model.Trigger, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTrigger", ctx, triggerId)
	ret0, _ := ret[0].(*model.Trigger)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTrigger indicates an expected call of GetTrigger.
func (mr *MockTriggerServiceMockRecorder) GetTrigger(ctx, triggerId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTrigger", reflect.TypeOf((*MockTriggerService)(nil).GetTrigger), ctx, triggerId)
}

// ListTriggers mocks base method.
func (m *MockTriggerService) ListTriggers(ctx context.Context, page, limit int) (*api.List[model.Trigger], error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListTriggers", ctx, page, limit)
	ret0, _ := ret[0].(*api.List[model.Trigger])
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListTriggers indicates an expected call of ListTriggers.
func (mr *MockTriggerServiceMockRecorder) ListTriggers(ctx, page, limit any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListTriggers", reflect.TypeOf((*MockTriggerService)(nil).ListTriggers), ctx, page, limit)
}

// UpdateTrigger mocks base method.
func (m *MockTriggerService) UpdateTrigger(ctx context.Context, triggerId uint64, req service.UpdateTriggerRequest) (*model.Trigger, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateTrigger", ctx, triggerId, req)
	ret0, _ := ret[0].(*model.Trigger)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateTrigger indicates an expected call of UpdateTrigger.
func (mr *MockTriggerServiceMockRecorder) UpdateTrigger(ctx, triggerId, req any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateTrigger", reflect.TypeOf((*MockTriggerService)(nil).UpdateTrigger), ctx, triggerId, req)
}
