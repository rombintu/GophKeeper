// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/rombintu/GophKeeper/internal/storage (interfaces: ClientManager)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	keeper "github.com/rombintu/GophKeeper/internal/proto/keeper"
)

// MockClientManager is a mock of ClientManager interface.
type MockClientManager struct {
	ctrl     *gomock.Controller
	recorder *MockClientManagerMockRecorder
}

// MockClientManagerMockRecorder is the mock recorder for MockClientManager.
type MockClientManagerMockRecorder struct {
	mock *MockClientManager
}

// NewMockClientManager creates a new mock instance.
func NewMockClientManager(ctrl *gomock.Controller) *MockClientManager {
	mock := &MockClientManager{ctrl: ctrl}
	mock.recorder = &MockClientManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockClientManager) EXPECT() *MockClientManagerMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockClientManager) Close(arg0 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockClientManagerMockRecorder) Close(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockClientManager)(nil).Close), arg0)
}

// Configure mocks base method.
func (m *MockClientManager) Configure(arg0 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Configure", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Configure indicates an expected call of Configure.
func (mr *MockClientManagerMockRecorder) Configure(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Configure", reflect.TypeOf((*MockClientManager)(nil).Configure), arg0)
}

// Get mocks base method.
func (m *MockClientManager) Get(arg0 context.Context, arg1 []byte) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0, arg1)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockClientManagerMockRecorder) Get(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockClientManager)(nil).Get), arg0, arg1)
}

// GetMap mocks base method.
func (m *MockClientManager) GetMap(arg0 context.Context) (map[string]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMap", arg0)
	ret0, _ := ret[0].(map[string]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMap indicates an expected call of GetMap.
func (mr *MockClientManagerMockRecorder) GetMap(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMap", reflect.TypeOf((*MockClientManager)(nil).GetMap), arg0)
}

// Open mocks base method.
func (m *MockClientManager) Open(arg0 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Open", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Open indicates an expected call of Open.
func (mr *MockClientManagerMockRecorder) Open(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Open", reflect.TypeOf((*MockClientManager)(nil).Open), arg0)
}

// Ping mocks base method.
func (m *MockClientManager) Ping(arg0 context.Context, arg1 bool) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Ping", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Ping indicates an expected call of Ping.
func (mr *MockClientManagerMockRecorder) Ping(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Ping", reflect.TypeOf((*MockClientManager)(nil).Ping), arg0, arg1)
}

// SecretCreate mocks base method.
func (m *MockClientManager) SecretCreate(arg0 context.Context, arg1 *keeper.Secret) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SecretCreate", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// SecretCreate indicates an expected call of SecretCreate.
func (mr *MockClientManagerMockRecorder) SecretCreate(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SecretCreate", reflect.TypeOf((*MockClientManager)(nil).SecretCreate), arg0, arg1)
}

// SecretCreateBatch mocks base method.
func (m *MockClientManager) SecretCreateBatch(arg0 context.Context, arg1 []*keeper.Secret) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SecretCreateBatch", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// SecretCreateBatch indicates an expected call of SecretCreateBatch.
func (mr *MockClientManagerMockRecorder) SecretCreateBatch(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SecretCreateBatch", reflect.TypeOf((*MockClientManager)(nil).SecretCreateBatch), arg0, arg1)
}

// SecretGetBatch mocks base method.
func (m *MockClientManager) SecretGetBatch(arg0 context.Context) ([]*keeper.Secret, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SecretGetBatch", arg0)
	ret0, _ := ret[0].([]*keeper.Secret)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SecretGetBatch indicates an expected call of SecretGetBatch.
func (mr *MockClientManagerMockRecorder) SecretGetBatch(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SecretGetBatch", reflect.TypeOf((*MockClientManager)(nil).SecretGetBatch), arg0)
}

// SecretList mocks base method.
func (m *MockClientManager) SecretList(arg0 context.Context, arg1 string) ([]*keeper.Secret, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SecretList", arg0, arg1)
	ret0, _ := ret[0].([]*keeper.Secret)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SecretList indicates an expected call of SecretList.
func (mr *MockClientManagerMockRecorder) SecretList(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SecretList", reflect.TypeOf((*MockClientManager)(nil).SecretList), arg0, arg1)
}

// SecretPurge mocks base method.
func (m *MockClientManager) SecretPurge(arg0 context.Context, arg1 *keeper.Secret) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SecretPurge", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// SecretPurge indicates an expected call of SecretPurge.
func (mr *MockClientManagerMockRecorder) SecretPurge(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SecretPurge", reflect.TypeOf((*MockClientManager)(nil).SecretPurge), arg0, arg1)
}

// Set mocks base method.
func (m *MockClientManager) Set(arg0 context.Context, arg1, arg2 []byte) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Set", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// Set indicates an expected call of Set.
func (mr *MockClientManagerMockRecorder) Set(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Set", reflect.TypeOf((*MockClientManager)(nil).Set), arg0, arg1, arg2)
}
