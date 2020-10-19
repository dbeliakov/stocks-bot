// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/dbeliakov/stocks-bot/storage (interfaces: Storage)

// Package mock is a generated GoMock package.
package mock

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockStorage is a mock of Storage interface
type MockStorage struct {
	ctrl     *gomock.Controller
	recorder *MockStorageMockRecorder
}

// MockStorageMockRecorder is the mock recorder for MockStorage
type MockStorageMockRecorder struct {
	mock *MockStorage
}

// NewMockStorage creates a new mock instance
func NewMockStorage(ctrl *gomock.Controller) *MockStorage {
	mock := &MockStorage{ctrl: ctrl}
	mock.recorder = &MockStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockStorage) EXPECT() *MockStorageMockRecorder {
	return m.recorder
}

// AddSymbol mocks base method
func (m *MockStorage) AddSymbol(arg0 int64, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddSymbol", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddSymbol indicates an expected call of AddSymbol
func (mr *MockStorageMockRecorder) AddSymbol(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddSymbol", reflect.TypeOf((*MockStorage)(nil).AddSymbol), arg0, arg1)
}

// GetState mocks base method
func (m *MockStorage) GetState(arg0 int64) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetState", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetState indicates an expected call of GetState
func (mr *MockStorageMockRecorder) GetState(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetState", reflect.TypeOf((*MockStorage)(nil).GetState), arg0)
}

// Init mocks base method
func (m *MockStorage) Init() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Init")
	ret0, _ := ret[0].(error)
	return ret0
}

// Init indicates an expected call of Init
func (mr *MockStorageMockRecorder) Init() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Init", reflect.TypeOf((*MockStorage)(nil).Init))
}

// RemoveSymbol mocks base method
func (m *MockStorage) RemoveSymbol(arg0 int64, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveSymbol", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveSymbol indicates an expected call of RemoveSymbol
func (mr *MockStorageMockRecorder) RemoveSymbol(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveSymbol", reflect.TypeOf((*MockStorage)(nil).RemoveSymbol), arg0, arg1)
}

// SetState mocks base method
func (m *MockStorage) SetState(arg0 int64, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetState", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetState indicates an expected call of SetState
func (mr *MockStorageMockRecorder) SetState(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetState", reflect.TypeOf((*MockStorage)(nil).SetState), arg0, arg1)
}

// Symbols mocks base method
func (m *MockStorage) Symbols(arg0 int64) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Symbols", arg0)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Symbols indicates an expected call of Symbols
func (mr *MockStorageMockRecorder) Symbols(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Symbols", reflect.TypeOf((*MockStorage)(nil).Symbols), arg0)
}
