// Code generated by MockGen. DO NOT EDIT.
// Source: base.go

// Package mock_user is a generated GoMock package.
package mock_user

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	internal "github.com/pursuit/portal/internal"
)

// MockService is a mock of Service interface.
type MockService struct {
	ctrl     *gomock.Controller
	recorder *MockServiceMockRecorder
}

// MockServiceMockRecorder is the mock recorder for MockService.
type MockServiceMockRecorder struct {
	mock *MockService
}

// NewMockService creates a new mock instance.
func NewMockService(ctrl *gomock.Controller) *MockService {
	mock := &MockService{ctrl: ctrl}
	mock.recorder = &MockServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockService) EXPECT() *MockServiceMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockService) Create(ctx context.Context, username string, password []byte) (int, *internal.E) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, username, password)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(*internal.E)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockServiceMockRecorder) Create(ctx, username, password interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockService)(nil).Create), ctx, username, password)
}

// Login mocks base method.
func (m *MockService) Login(ctx context.Context, username string, password []byte) (string, *internal.E) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Login", ctx, username, password)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(*internal.E)
	return ret0, ret1
}

// Login indicates an expected call of Login.
func (mr *MockServiceMockRecorder) Login(ctx, username, password interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Login", reflect.TypeOf((*MockService)(nil).Login), ctx, username, password)
}
