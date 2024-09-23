// Code generated by MockGen. DO NOT EDIT.
// Source: internal/repository/wallet_repo.go
//
// Generated by this command:
//
//	mockgen -source=internal/repository/wallet_repo.go -destination=internal/repository/mocks/wallet_repo_mock.go -package=repository_mock
//

// Package repository_mock is a generated GoMock package.
package repository_mock

import (
	context "context"
	reflect "reflect"

	model "github.com/abdelrahman146/digital-wallet/internal/model"
	gomock "go.uber.org/mock/gomock"
)

// MockWalletRepo is a mock of WalletRepo interface.
type MockWalletRepo struct {
	ctrl     *gomock.Controller
	recorder *MockWalletRepoMockRecorder
}

// MockWalletRepoMockRecorder is the mock recorder for MockWalletRepo.
type MockWalletRepoMockRecorder struct {
	mock *MockWalletRepo
}

// NewMockWalletRepo creates a new mock instance.
func NewMockWalletRepo(ctrl *gomock.Controller) *MockWalletRepo {
	mock := &MockWalletRepo{ctrl: ctrl}
	mock.recorder = &MockWalletRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockWalletRepo) EXPECT() *MockWalletRepoMockRecorder {
	return m.recorder
}

// CountTotalWallets mocks base method.
func (m *MockWalletRepo) CountTotalWallets(ctx context.Context) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CountTotalWallets", ctx)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CountTotalWallets indicates an expected call of CountTotalWallets.
func (mr *MockWalletRepoMockRecorder) CountTotalWallets(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CountTotalWallets", reflect.TypeOf((*MockWalletRepo)(nil).CountTotalWallets), ctx)
}

// CreateWallet mocks base method.
func (m *MockWalletRepo) CreateWallet(ctx context.Context, wallet *model.Wallet) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateWallet", ctx, wallet)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateWallet indicates an expected call of CreateWallet.
func (mr *MockWalletRepoMockRecorder) CreateWallet(ctx, wallet any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateWallet", reflect.TypeOf((*MockWalletRepo)(nil).CreateWallet), ctx, wallet)
}

// DeleteWallet mocks base method.
func (m *MockWalletRepo) DeleteWallet(ctx context.Context, wallet *model.Wallet) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteWallet", ctx, wallet)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteWallet indicates an expected call of DeleteWallet.
func (mr *MockWalletRepoMockRecorder) DeleteWallet(ctx, wallet any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteWallet", reflect.TypeOf((*MockWalletRepo)(nil).DeleteWallet), ctx, wallet)
}

// FetchWalletByID mocks base method.
func (m *MockWalletRepo) FetchWalletByID(ctx context.Context, walletId string) (*model.Wallet, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FetchWalletByID", ctx, walletId)
	ret0, _ := ret[0].(*model.Wallet)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FetchWalletByID indicates an expected call of FetchWalletByID.
func (mr *MockWalletRepoMockRecorder) FetchWalletByID(ctx, walletId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FetchWalletByID", reflect.TypeOf((*MockWalletRepo)(nil).FetchWalletByID), ctx, walletId)
}

// FetchWallets mocks base method.
func (m *MockWalletRepo) FetchWallets(ctx context.Context, page, limit int) ([]model.Wallet, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FetchWallets", ctx, page, limit)
	ret0, _ := ret[0].([]model.Wallet)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FetchWallets indicates an expected call of FetchWallets.
func (mr *MockWalletRepoMockRecorder) FetchWallets(ctx, page, limit any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FetchWallets", reflect.TypeOf((*MockWalletRepo)(nil).FetchWallets), ctx, page, limit)
}

// UpdateWallet mocks base method.
func (m *MockWalletRepo) UpdateWallet(ctx context.Context, wallet *model.Wallet) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateWallet", ctx, wallet)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateWallet indicates an expected call of UpdateWallet.
func (mr *MockWalletRepoMockRecorder) UpdateWallet(ctx, wallet any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateWallet", reflect.TypeOf((*MockWalletRepo)(nil).UpdateWallet), ctx, wallet)
}
