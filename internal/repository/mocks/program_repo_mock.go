// Code generated by MockGen. DO NOT EDIT.
// Source: internal/repository/program_repo.go
//
// Generated by this command:
//
//	mockgen -source=internal/repository/program_repo.go -destination=internal/repository/mocks/program_repo_mock.go -package=repository_mock
//

// Package repository_mock is a generated GoMock package.
package repository_mock

import (
	context "context"
	reflect "reflect"

	model "github.com/abdelrahman146/digital-wallet/internal/model"
	gomock "go.uber.org/mock/gomock"
)

// MockProgramRepo is a mock of ProgramRepo interface.
type MockProgramRepo struct {
	ctrl     *gomock.Controller
	recorder *MockProgramRepoMockRecorder
}

// MockProgramRepoMockRecorder is the mock recorder for MockProgramRepo.
type MockProgramRepoMockRecorder struct {
	mock *MockProgramRepo
}

// NewMockProgramRepo creates a new mock instance.
func NewMockProgramRepo(ctrl *gomock.Controller) *MockProgramRepo {
	mock := &MockProgramRepo{ctrl: ctrl}
	mock.recorder = &MockProgramRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockProgramRepo) EXPECT() *MockProgramRepoMockRecorder {
	return m.recorder
}

// CountPrograms mocks base method.
func (m *MockProgramRepo) CountPrograms(ctx context.Context) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CountPrograms", ctx)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CountPrograms indicates an expected call of CountPrograms.
func (mr *MockProgramRepoMockRecorder) CountPrograms(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CountPrograms", reflect.TypeOf((*MockProgramRepo)(nil).CountPrograms), ctx)
}

// CountProgramsByWalletID mocks base method.
func (m *MockProgramRepo) CountProgramsByWalletID(ctx context.Context, walletID uint64) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CountProgramsByWalletID", ctx, walletID)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CountProgramsByWalletID indicates an expected call of CountProgramsByWalletID.
func (mr *MockProgramRepoMockRecorder) CountProgramsByWalletID(ctx, walletID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CountProgramsByWalletID", reflect.TypeOf((*MockProgramRepo)(nil).CountProgramsByWalletID), ctx, walletID)
}

// CreateProgram mocks base method.
func (m *MockProgramRepo) CreateProgram(ctx context.Context, program *model.Program) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateProgram", ctx, program)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateProgram indicates an expected call of CreateProgram.
func (mr *MockProgramRepoMockRecorder) CreateProgram(ctx, program any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateProgram", reflect.TypeOf((*MockProgramRepo)(nil).CreateProgram), ctx, program)
}

// DeleteProgram mocks base method.
func (m *MockProgramRepo) DeleteProgram(ctx context.Context, program *model.Program) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteProgram", ctx, program)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteProgram indicates an expected call of DeleteProgram.
func (mr *MockProgramRepoMockRecorder) DeleteProgram(ctx, program any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteProgram", reflect.TypeOf((*MockProgramRepo)(nil).DeleteProgram), ctx, program)
}

// FetchProgramByID mocks base method.
func (m *MockProgramRepo) FetchProgramByID(ctx context.Context, id uint64) (*model.Program, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FetchProgramByID", ctx, id)
	ret0, _ := ret[0].(*model.Program)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FetchProgramByID indicates an expected call of FetchProgramByID.
func (mr *MockProgramRepoMockRecorder) FetchProgramByID(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FetchProgramByID", reflect.TypeOf((*MockProgramRepo)(nil).FetchProgramByID), ctx, id)
}

// FetchPrograms mocks base method.
func (m *MockProgramRepo) FetchPrograms(ctx context.Context, page, limit int) ([]model.Program, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FetchPrograms", ctx, page, limit)
	ret0, _ := ret[0].([]model.Program)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FetchPrograms indicates an expected call of FetchPrograms.
func (mr *MockProgramRepoMockRecorder) FetchPrograms(ctx, page, limit any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FetchPrograms", reflect.TypeOf((*MockProgramRepo)(nil).FetchPrograms), ctx, page, limit)
}

// FetchProgramsByWalletID mocks base method.
func (m *MockProgramRepo) FetchProgramsByWalletID(ctx context.Context, walletID uint64) ([]model.Program, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FetchProgramsByWalletID", ctx, walletID)
	ret0, _ := ret[0].([]model.Program)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FetchProgramsByWalletID indicates an expected call of FetchProgramsByWalletID.
func (mr *MockProgramRepoMockRecorder) FetchProgramsByWalletID(ctx, walletID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FetchProgramsByWalletID", reflect.TypeOf((*MockProgramRepo)(nil).FetchProgramsByWalletID), ctx, walletID)
}

// FetchTriggerPrograms mocks base method.
func (m *MockProgramRepo) FetchTriggerPrograms(ctx context.Context, triggerSlug string) ([]*model.Program, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FetchTriggerPrograms", ctx, triggerSlug)
	ret0, _ := ret[0].([]*model.Program)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FetchTriggerPrograms indicates an expected call of FetchTriggerPrograms.
func (mr *MockProgramRepoMockRecorder) FetchTriggerPrograms(ctx, triggerSlug any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FetchTriggerPrograms", reflect.TypeOf((*MockProgramRepo)(nil).FetchTriggerPrograms), ctx, triggerSlug)
}

// UpdateProgram mocks base method.
func (m *MockProgramRepo) UpdateProgram(ctx context.Context, program *model.Program) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateProgram", ctx, program)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateProgram indicates an expected call of UpdateProgram.
func (mr *MockProgramRepoMockRecorder) UpdateProgram(ctx, program any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateProgram", reflect.TypeOf((*MockProgramRepo)(nil).UpdateProgram), ctx, program)
}