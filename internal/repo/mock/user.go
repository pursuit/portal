// Code generated by MockGen. DO NOT EDIT.
// Source: user.go

// Package mock_repo is a generated GoMock package.
package mock_repo

import (
	context "context"
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"
	internal "github.com/pursuit/portal/internal"
	repo "github.com/pursuit/portal/internal/repo"
)

// MockUser is a mock of User interface.
type MockUser struct {
	ctrl     *gomock.Controller
	recorder *MockUserMockRecorder
}

// MockUserMockRecorder is the mock recorder for MockUser.
type MockUserMockRecorder struct {
	mock *MockUser
}

// NewMockUser creates a new mock instance.
func NewMockUser(ctrl *gomock.Controller) *MockUser {
	mock := &MockUser{ctrl: ctrl}
	mock.recorder = &MockUserMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUser) EXPECT() *MockUserMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockUser) Create(ctx context.Context, db repo.DB, username string, hashedPassword []byte, now time.Time) (int, *internal.E) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, db, username, hashedPassword, now)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(*internal.E)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockUserMockRecorder) Create(ctx, db, username, hashedPassword, now interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockUser)(nil).Create), ctx, db, username, hashedPassword, now)
}
