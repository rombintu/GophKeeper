// Code generated by MockGen. DO NOT EDIT.
// Source: internal/proto/keeper/keeper_grpc.pb.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	keeper "github.com/rombintu/GophKeeper/internal/proto/keeper"
	grpc "google.golang.org/grpc"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// MockKeeperClient is a mock of KeeperClient interface.
type MockKeeperClient struct {
	ctrl     *gomock.Controller
	recorder *MockKeeperClientMockRecorder
}

// MockKeeperClientMockRecorder is the mock recorder for MockKeeperClient.
type MockKeeperClientMockRecorder struct {
	mock *MockKeeperClient
}

// NewMockKeeperClient creates a new mock instance.
func NewMockKeeperClient(ctrl *gomock.Controller) *MockKeeperClient {
	mock := &MockKeeperClient{ctrl: ctrl}
	mock.recorder = &MockKeeperClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockKeeperClient) EXPECT() *MockKeeperClientMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockKeeperClient) Create(ctx context.Context, in *keeper.CreateRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Create", varargs...)
	ret0, _ := ret[0].(*emptypb.Empty)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockKeeperClientMockRecorder) Create(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockKeeperClient)(nil).Create), varargs...)
}

// CreateMany mocks base method.
func (m *MockKeeperClient) CreateMany(ctx context.Context, in *keeper.CreateBatchRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "CreateMany", varargs...)
	ret0, _ := ret[0].(*emptypb.Empty)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateMany indicates an expected call of CreateMany.
func (mr *MockKeeperClientMockRecorder) CreateMany(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateMany", reflect.TypeOf((*MockKeeperClient)(nil).CreateMany), varargs...)
}

// Delete mocks base method.
func (m *MockKeeperClient) Delete(ctx context.Context, in *keeper.DeleteRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Delete", varargs...)
	ret0, _ := ret[0].(*emptypb.Empty)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Delete indicates an expected call of Delete.
func (mr *MockKeeperClientMockRecorder) Delete(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockKeeperClient)(nil).Delete), varargs...)
}

// Fetch mocks base method.
func (m *MockKeeperClient) Fetch(ctx context.Context, in *keeper.FetchRequest, opts ...grpc.CallOption) (*keeper.FetchResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Fetch", varargs...)
	ret0, _ := ret[0].(*keeper.FetchResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Fetch indicates an expected call of Fetch.
func (mr *MockKeeperClientMockRecorder) Fetch(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Fetch", reflect.TypeOf((*MockKeeperClient)(nil).Fetch), varargs...)
}

// MockKeeperServer is a mock of KeeperServer interface.
type MockKeeperServer struct {
	ctrl     *gomock.Controller
	recorder *MockKeeperServerMockRecorder
}

// MockKeeperServerMockRecorder is the mock recorder for MockKeeperServer.
type MockKeeperServerMockRecorder struct {
	mock *MockKeeperServer
}

// NewMockKeeperServer creates a new mock instance.
func NewMockKeeperServer(ctrl *gomock.Controller) *MockKeeperServer {
	mock := &MockKeeperServer{ctrl: ctrl}
	mock.recorder = &MockKeeperServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockKeeperServer) EXPECT() *MockKeeperServerMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockKeeperServer) Create(arg0 context.Context, arg1 *keeper.CreateRequest) (*emptypb.Empty, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0, arg1)
	ret0, _ := ret[0].(*emptypb.Empty)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockKeeperServerMockRecorder) Create(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockKeeperServer)(nil).Create), arg0, arg1)
}

// CreateMany mocks base method.
func (m *MockKeeperServer) CreateMany(arg0 context.Context, arg1 *keeper.CreateBatchRequest) (*emptypb.Empty, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateMany", arg0, arg1)
	ret0, _ := ret[0].(*emptypb.Empty)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateMany indicates an expected call of CreateMany.
func (mr *MockKeeperServerMockRecorder) CreateMany(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateMany", reflect.TypeOf((*MockKeeperServer)(nil).CreateMany), arg0, arg1)
}

// Delete mocks base method.
func (m *MockKeeperServer) Delete(arg0 context.Context, arg1 *keeper.DeleteRequest) (*emptypb.Empty, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0, arg1)
	ret0, _ := ret[0].(*emptypb.Empty)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Delete indicates an expected call of Delete.
func (mr *MockKeeperServerMockRecorder) Delete(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockKeeperServer)(nil).Delete), arg0, arg1)
}

// Fetch mocks base method.
func (m *MockKeeperServer) Fetch(arg0 context.Context, arg1 *keeper.FetchRequest) (*keeper.FetchResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Fetch", arg0, arg1)
	ret0, _ := ret[0].(*keeper.FetchResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Fetch indicates an expected call of Fetch.
func (mr *MockKeeperServerMockRecorder) Fetch(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Fetch", reflect.TypeOf((*MockKeeperServer)(nil).Fetch), arg0, arg1)
}

// mustEmbedUnimplementedKeeperServer mocks base method.
func (m *MockKeeperServer) mustEmbedUnimplementedKeeperServer() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "mustEmbedUnimplementedKeeperServer")
}

// mustEmbedUnimplementedKeeperServer indicates an expected call of mustEmbedUnimplementedKeeperServer.
func (mr *MockKeeperServerMockRecorder) mustEmbedUnimplementedKeeperServer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "mustEmbedUnimplementedKeeperServer", reflect.TypeOf((*MockKeeperServer)(nil).mustEmbedUnimplementedKeeperServer))
}

// MockUnsafeKeeperServer is a mock of UnsafeKeeperServer interface.
type MockUnsafeKeeperServer struct {
	ctrl     *gomock.Controller
	recorder *MockUnsafeKeeperServerMockRecorder
}

// MockUnsafeKeeperServerMockRecorder is the mock recorder for MockUnsafeKeeperServer.
type MockUnsafeKeeperServerMockRecorder struct {
	mock *MockUnsafeKeeperServer
}

// NewMockUnsafeKeeperServer creates a new mock instance.
func NewMockUnsafeKeeperServer(ctrl *gomock.Controller) *MockUnsafeKeeperServer {
	mock := &MockUnsafeKeeperServer{ctrl: ctrl}
	mock.recorder = &MockUnsafeKeeperServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUnsafeKeeperServer) EXPECT() *MockUnsafeKeeperServerMockRecorder {
	return m.recorder
}

// mustEmbedUnimplementedKeeperServer mocks base method.
func (m *MockUnsafeKeeperServer) mustEmbedUnimplementedKeeperServer() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "mustEmbedUnimplementedKeeperServer")
}

// mustEmbedUnimplementedKeeperServer indicates an expected call of mustEmbedUnimplementedKeeperServer.
func (mr *MockUnsafeKeeperServerMockRecorder) mustEmbedUnimplementedKeeperServer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "mustEmbedUnimplementedKeeperServer", reflect.TypeOf((*MockUnsafeKeeperServer)(nil).mustEmbedUnimplementedKeeperServer))
}
