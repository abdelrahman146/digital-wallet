// Code generated by MockGen. DO NOT EDIT.
// Source: internal/service/program_service.go
//
// Generated by this command:
//
//	mockgen -source=internal/service/program_service.go -destination=internal/service/mocks/program_service_mock.go -package=service_mock
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

// MockProgramService is a mock of ProgramService interface.
type MockProgramService struct {
	ctrl     *gomock.Controller
	recorder *MockProgramServiceMockRecorder
}

// MockProgramServiceMockRecorder is the mock recorder for MockProgramService.
type MockProgramServiceMockRecorder struct {
	mock *MockProgramService
}

// NewMockProgramService creates a new mock instance.
func NewMockProgramService(ctrl *gomock.Controller) *MockProgramService {
	mock := &MockProgramService{ctrl: ctrl}
	mock.recorder = &MockProgramServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockProgramService) EXPECT() *MockProgramServiceMockRecorder {
	return m.recorder
}

// CreateProgram mocks base method.
func (m *MockProgramService) CreateProgram(ctx context.Context, req service.CreateProgramRequest) (*model.Program, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateProgram", ctx, req)
	ret0, _ := ret[0].(*model.Program)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateProgram indicates an expected call of CreateProgram.
func (mr *MockProgramServiceMockRecorder) CreateProgram(ctx, req any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateProgram", reflect.TypeOf((*MockProgramService)(nil).CreateProgram), ctx, req)
}

// DeleteProgram mocks base method.
func (m *MockProgramService) DeleteProgram(ctx context.Context, id uint64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteProgram", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteProgram indicates an expected call of DeleteProgram.
func (mr *MockProgramServiceMockRecorder) DeleteProgram(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteProgram", reflect.TypeOf((*MockProgramService)(nil).DeleteProgram), ctx, id)
}

// GetProgram mocks base method.
func (m *MockProgramService) GetProgram(ctx context.Context, id uint64) (*model.Program, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProgram", ctx, id)
	ret0, _ := ret[0].(*model.Program)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProgram indicates an expected call of GetProgram.
func (mr *MockProgramServiceMockRecorder) GetProgram(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProgram", reflect.TypeOf((*MockProgramService)(nil).GetProgram), ctx, id)
}

// ListPrograms mocks base method.
func (m *MockProgramService) ListPrograms(ctx context.Context, page, limit int) (*api.List[model.Program], error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListPrograms", ctx, page, limit)
	ret0, _ := ret[0].(*api.List[model.Program])
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListPrograms indicates an expected call of ListPrograms.
func (mr *MockProgramServiceMockRecorder) ListPrograms(ctx, page, limit any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListPrograms", reflect.TypeOf((*MockProgramService)(nil).ListPrograms), ctx, page, limit)
}

// UpdateProgram mocks base method.
func (m *MockProgramService) UpdateProgram(ctx context.Context, id uint64, req service.UpdateProgramRequest) (*model.Program, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateProgram", ctx, id, req)
	ret0, _ := ret[0].(*model.Program)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateProgram indicates an expected call of UpdateProgram.
func (mr *MockProgramServiceMockRecorder) UpdateProgram(ctx, id, req any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateProgram", reflect.TypeOf((*MockProgramService)(nil).UpdateProgram), ctx, id, req)
}
