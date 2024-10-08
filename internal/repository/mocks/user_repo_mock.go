// Code generated by MockGen. DO NOT EDIT.
// Source: internal/repository/user_repo.go
//
// Generated by this command:
//
//	mockgen -source=internal/repository/user_repo.go -destination=internal/repository/mocks/user_repo_mock.go -package=repository_mock
//

// Package repository_mock is a generated GoMock package.
package repository_mock

import (
	context "context"
	reflect "reflect"

	model "github.com/abdelrahman146/digital-wallet/internal/model"
	gomock "go.uber.org/mock/gomock"
)

// MockUserRepo is a mock of UserRepo interface.
type MockUserRepo struct {
	ctrl     *gomock.Controller
	recorder *MockUserRepoMockRecorder
}

// MockUserRepoMockRecorder is the mock recorder for MockUserRepo.
type MockUserRepoMockRecorder struct {
	mock *MockUserRepo
}

// NewMockUserRepo creates a new mock instance.
func NewMockUserRepo(ctrl *gomock.Controller) *MockUserRepo {
	mock := &MockUserRepo{ctrl: ctrl}
	mock.recorder = &MockUserRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserRepo) EXPECT() *MockUserRepoMockRecorder {
	return m.recorder
}

// CountTotalUsers mocks base method.
func (m *MockUserRepo) CountTotalUsers(ctx context.Context) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CountTotalUsers", ctx)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CountTotalUsers indicates an expected call of CountTotalUsers.
func (mr *MockUserRepoMockRecorder) CountTotalUsers(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CountTotalUsers", reflect.TypeOf((*MockUserRepo)(nil).CountTotalUsers), ctx)
}

// CountUsersByTierID mocks base method.
func (m *MockUserRepo) CountUsersByTierID(ctx context.Context, tierId string) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CountUsersByTierID", ctx, tierId)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CountUsersByTierID indicates an expected call of CountUsersByTierID.
func (mr *MockUserRepoMockRecorder) CountUsersByTierID(ctx, tierId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CountUsersByTierID", reflect.TypeOf((*MockUserRepo)(nil).CountUsersByTierID), ctx, tierId)
}

// CreateUser mocks base method.
func (m *MockUserRepo) CreateUser(ctx context.Context, user *model.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", ctx, user)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockUserRepoMockRecorder) CreateUser(ctx, user any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockUserRepo)(nil).CreateUser), ctx, user)
}

// DeleteUser mocks base method.
func (m *MockUserRepo) DeleteUser(ctx context.Context, user *model.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteUser", ctx, user)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteUser indicates an expected call of DeleteUser.
func (mr *MockUserRepoMockRecorder) DeleteUser(ctx, user any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteUser", reflect.TypeOf((*MockUserRepo)(nil).DeleteUser), ctx, user)
}

// FetchUserByID mocks base method.
func (m *MockUserRepo) FetchUserByID(ctx context.Context, userId string) (*model.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FetchUserByID", ctx, userId)
	ret0, _ := ret[0].(*model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FetchUserByID indicates an expected call of FetchUserByID.
func (mr *MockUserRepoMockRecorder) FetchUserByID(ctx, userId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FetchUserByID", reflect.TypeOf((*MockUserRepo)(nil).FetchUserByID), ctx, userId)
}

// FetchUsers mocks base method.
func (m *MockUserRepo) FetchUsers(ctx context.Context, page, limit int) ([]model.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FetchUsers", ctx, page, limit)
	ret0, _ := ret[0].([]model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FetchUsers indicates an expected call of FetchUsers.
func (mr *MockUserRepoMockRecorder) FetchUsers(ctx, page, limit any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FetchUsers", reflect.TypeOf((*MockUserRepo)(nil).FetchUsers), ctx, page, limit)
}

// FetchUsersByTierID mocks base method.
func (m *MockUserRepo) FetchUsersByTierID(ctx context.Context, tierId string, page, limit int) ([]model.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FetchUsersByTierID", ctx, tierId, page, limit)
	ret0, _ := ret[0].([]model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FetchUsersByTierID indicates an expected call of FetchUsersByTierID.
func (mr *MockUserRepoMockRecorder) FetchUsersByTierID(ctx, tierId, page, limit any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FetchUsersByTierID", reflect.TypeOf((*MockUserRepo)(nil).FetchUsersByTierID), ctx, tierId, page, limit)
}

// UpdateUserTier mocks base method.
func (m *MockUserRepo) UpdateUserTier(ctx context.Context, userId, tierId string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUserTier", ctx, userId, tierId)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateUserTier indicates an expected call of UpdateUserTier.
func (mr *MockUserRepoMockRecorder) UpdateUserTier(ctx, userId, tierId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUserTier", reflect.TypeOf((*MockUserRepo)(nil).UpdateUserTier), ctx, userId, tierId)
}
