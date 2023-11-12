// Code generated by MockGen. DO NOT EDIT.
// Source: ./webook/internal/service/sms/types.go
//
// Generated by this command:
//
//	mockgen -source=./webook/internal/service/sms/types.go -package=smsmocks -destination=./webook/internal/service/sms/mocks/sms.mock.go
//
// Package smsmocks is a generated GoMock package.
package smsmocks

import (
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockSmsService is a mock of SmsService interface.
type MockSmsService struct {
	ctrl     *gomock.Controller
	recorder *MockSmsServiceMockRecorder
}

// MockSmsServiceMockRecorder is the mock recorder for MockSmsService.
type MockSmsServiceMockRecorder struct {
	mock *MockSmsService
}

// NewMockSmsService creates a new mock instance.
func NewMockSmsService(ctrl *gomock.Controller) *MockSmsService {
	mock := &MockSmsService{ctrl: ctrl}
	mock.recorder = &MockSmsServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSmsService) EXPECT() *MockSmsServiceMockRecorder {
	return m.recorder
}

// Send mocks base method.
func (m *MockSmsService) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	m.ctrl.T.Helper()
	varargs := []any{ctx, tplId, args}
	for _, a := range numbers {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Send", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Send indicates an expected call of Send.
func (mr *MockSmsServiceMockRecorder) Send(ctx, tplId, args any, numbers ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, tplId, args}, numbers...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Send", reflect.TypeOf((*MockSmsService)(nil).Send), varargs...)
}