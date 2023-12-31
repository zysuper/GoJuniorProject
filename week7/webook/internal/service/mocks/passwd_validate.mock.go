// Code generated by MockGen. DO NOT EDIT.
// Source: ./webook/internal/service/passwd_validate.go
//
// Generated by this command:
//
//	mockgen -source=./webook/internal/service/passwd_validate.go -package=svcmocks -destination=./webook/internal/service/mocks/passwd_validate.mock.go
//
// Package svcmocks is a generated GoMock package.
package svcmocks

import (
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockPasswordValidateService is a mock of PasswordValidateService interface.
type MockPasswordValidateService struct {
	ctrl     *gomock.Controller
	recorder *MockPasswordValidateServiceMockRecorder
}

// MockPasswordValidateServiceMockRecorder is the mock recorder for MockPasswordValidateService.
type MockPasswordValidateServiceMockRecorder struct {
	mock *MockPasswordValidateService
}

// NewMockPasswordValidateService creates a new mock instance.
func NewMockPasswordValidateService(ctrl *gomock.Controller) *MockPasswordValidateService {
	mock := &MockPasswordValidateService{ctrl: ctrl}
	mock.recorder = &MockPasswordValidateServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPasswordValidateService) EXPECT() *MockPasswordValidateServiceMockRecorder {
	return m.recorder
}

// ComparePassword mocks base method.
func (m *MockPasswordValidateService) ComparePassword(hashedPasswd, passwd string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ComparePassword", hashedPasswd, passwd)
	ret0, _ := ret[0].(error)
	return ret0
}

// ComparePassword indicates an expected call of ComparePassword.
func (mr *MockPasswordValidateServiceMockRecorder) ComparePassword(hashedPasswd, passwd any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ComparePassword", reflect.TypeOf((*MockPasswordValidateService)(nil).ComparePassword), hashedPasswd, passwd)
}

// Hash mocks base method.
func (m *MockPasswordValidateService) Hash(password string) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Hash", password)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Hash indicates an expected call of Hash.
func (mr *MockPasswordValidateServiceMockRecorder) Hash(password any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Hash", reflect.TypeOf((*MockPasswordValidateService)(nil).Hash), password)
}
