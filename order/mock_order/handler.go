// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/qdm12/goshutdown/order (interfaces: Handler)

// Package mock_order is a generated GoMock package.
package mock_order

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	handler "github.com/qdm12/goshutdown/handler"
)

// MockHandler is a mock of Handler interface.
type MockHandler struct {
	ctrl     *gomock.Controller
	recorder *MockHandlerMockRecorder
}

// MockHandlerMockRecorder is the mock recorder for MockHandler.
type MockHandlerMockRecorder struct {
	mock *MockHandler
}

// NewMockHandler creates a new mock instance.
func NewMockHandler(ctrl *gomock.Controller) *MockHandler {
	mock := &MockHandler{ctrl: ctrl}
	mock.recorder = &MockHandlerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockHandler) EXPECT() *MockHandlerMockRecorder {
	return m.recorder
}

// Append mocks base method.
func (m *MockHandler) Append(arg0 ...handler.Handler) {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range arg0 {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Append", varargs...)
}

// Append indicates an expected call of Append.
func (mr *MockHandlerMockRecorder) Append(arg0 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Append", reflect.TypeOf((*MockHandler)(nil).Append), arg0...)
}

// IsCritical mocks base method.
func (m *MockHandler) IsCritical() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsCritical")
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsCritical indicates an expected call of IsCritical.
func (mr *MockHandlerMockRecorder) IsCritical() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsCritical", reflect.TypeOf((*MockHandler)(nil).IsCritical))
}

// Name mocks base method.
func (m *MockHandler) Name() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Name")
	ret0, _ := ret[0].(string)
	return ret0
}

// Name indicates an expected call of Name.
func (mr *MockHandlerMockRecorder) Name() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Name", reflect.TypeOf((*MockHandler)(nil).Name))
}

// Shutdown mocks base method.
func (m *MockHandler) Shutdown(arg0 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Shutdown", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Shutdown indicates an expected call of Shutdown.
func (mr *MockHandlerMockRecorder) Shutdown(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Shutdown", reflect.TypeOf((*MockHandler)(nil).Shutdown), arg0)
}
