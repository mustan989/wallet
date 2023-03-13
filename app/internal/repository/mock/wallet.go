// Code generated by MockGen. DO NOT EDIT.
// Source: ./repository/wallet.go

// Package mock_repository is a generated GoMock package.
package mock_repository

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	model "github.com/mustan989/wallet/model"
)

// MockWallet is a mock of Wallet interface.
type MockWallet struct {
	ctrl     *gomock.Controller
	recorder *MockWalletMockRecorder
}

// MockWalletMockRecorder is the mock recorder for MockWallet.
type MockWalletMockRecorder struct {
	mock *MockWallet
}

// NewMockWallet creates a new mock instance.
func NewMockWallet(ctrl *gomock.Controller) *MockWallet {
	mock := &MockWallet{ctrl: ctrl}
	mock.recorder = &MockWalletMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockWallet) EXPECT() *MockWalletMockRecorder {
	return m.recorder
}

// CountAll mocks base method.
func (m *MockWallet) CountAll(ctx context.Context, filter *model.WalletFilter) (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CountAll", ctx, filter)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CountAll indicates an expected call of CountAll.
func (mr *MockWalletMockRecorder) CountAll(ctx, filter interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CountAll", reflect.TypeOf((*MockWallet)(nil).CountAll), ctx, filter)
}

// Create mocks base method.
func (m *MockWallet) Create(ctx context.Context, data *model.Wallet) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, data)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockWalletMockRecorder) Create(ctx, data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockWallet)(nil).Create), ctx, data)
}

// DeleteByID mocks base method.
func (m *MockWallet) DeleteByID(ctx context.Context, id uint64) (*model.Wallet, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteByID", ctx, id)
	ret0, _ := ret[0].(*model.Wallet)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteByID indicates an expected call of DeleteByID.
func (mr *MockWalletMockRecorder) DeleteByID(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteByID", reflect.TypeOf((*MockWallet)(nil).DeleteByID), ctx, id)
}

// FindAll mocks base method.
func (m *MockWallet) FindAll(ctx context.Context, filter *model.WalletFilter) ([]*model.Wallet, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindAll", ctx, filter)
	ret0, _ := ret[0].([]*model.Wallet)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindAll indicates an expected call of FindAll.
func (mr *MockWalletMockRecorder) FindAll(ctx, filter interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindAll", reflect.TypeOf((*MockWallet)(nil).FindAll), ctx, filter)
}

// FindByID mocks base method.
func (m *MockWallet) FindByID(ctx context.Context, id uint64) (*model.Wallet, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByID", ctx, id)
	ret0, _ := ret[0].(*model.Wallet)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByID indicates an expected call of FindByID.
func (mr *MockWalletMockRecorder) FindByID(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByID", reflect.TypeOf((*MockWallet)(nil).FindByID), ctx, id)
}

// Update mocks base method.
func (m *MockWallet) Update(ctx context.Context, data *model.Wallet) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, data)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockWalletMockRecorder) Update(ctx, data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockWallet)(nil).Update), ctx, data)
}