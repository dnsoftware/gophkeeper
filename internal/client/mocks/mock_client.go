// Code generated by MockGen. DO NOT EDIT.
// Source: internal/client/domain/client.go

// Package mocksclient is a generated GoMock package.
package mocksclient

import (
	reflect "reflect"

	domain "github.com/dnsoftware/gophkeeper/internal/client/domain"
	gomock "github.com/golang/mock/gomock"
)

// MockSender is a mock of Sender interface.
type MockSender struct {
	ctrl     *gomock.Controller
	recorder *MockSenderMockRecorder
}

// MockSenderMockRecorder is the mock recorder for MockSender.
type MockSenderMockRecorder struct {
	mock *MockSender
}

// NewMockSender creates a new mock instance.
func NewMockSender(ctrl *gomock.Controller) *MockSender {
	mock := &MockSender{ctrl: ctrl}
	mock.recorder = &MockSenderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSender) EXPECT() *MockSenderMockRecorder {
	return m.recorder
}

// AddEntity mocks base method.
func (m *MockSender) AddEntity(ae domain.Entity) (int32, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddEntity", ae)
	ret0, _ := ret[0].(int32)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddEntity indicates an expected call of AddEntity.
func (mr *MockSenderMockRecorder) AddEntity(ae interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddEntity", reflect.TypeOf((*MockSender)(nil).AddEntity), ae)
}

// DownloadBinary mocks base method.
func (m *MockSender) DownloadBinary(entityId int32, fileName string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DownloadBinary", entityId, fileName)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DownloadBinary indicates an expected call of DownloadBinary.
func (mr *MockSenderMockRecorder) DownloadBinary(entityId, fileName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DownloadBinary", reflect.TypeOf((*MockSender)(nil).DownloadBinary), entityId, fileName)
}

// DownloadCryptoBinary mocks base method.
func (m *MockSender) DownloadCryptoBinary(entityId int32, fileName string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DownloadCryptoBinary", entityId, fileName)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DownloadCryptoBinary indicates an expected call of DownloadCryptoBinary.
func (mr *MockSenderMockRecorder) DownloadCryptoBinary(entityId, fileName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DownloadCryptoBinary", reflect.TypeOf((*MockSender)(nil).DownloadCryptoBinary), entityId, fileName)
}

// EntityCodes mocks base method.
func (m *MockSender) EntityCodes() ([]*domain.EntityCode, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EntityCodes")
	ret0, _ := ret[0].([]*domain.EntityCode)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// EntityCodes indicates an expected call of EntityCodes.
func (mr *MockSenderMockRecorder) EntityCodes() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EntityCodes", reflect.TypeOf((*MockSender)(nil).EntityCodes))
}

// Fields mocks base method.
func (m *MockSender) Fields(etype string) ([]*domain.Field, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Fields", etype)
	ret0, _ := ret[0].([]*domain.Field)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Fields indicates an expected call of Fields.
func (mr *MockSenderMockRecorder) Fields(etype interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Fields", reflect.TypeOf((*MockSender)(nil).Fields), etype)
}

// Login mocks base method.
func (m *MockSender) Login(login, password string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Login", login, password)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Login indicates an expected call of Login.
func (mr *MockSenderMockRecorder) Login(login, password interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Login", reflect.TypeOf((*MockSender)(nil).Login), login, password)
}

// Registration mocks base method.
func (m *MockSender) Registration(login, password, password2 string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Registration", login, password, password2)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Registration indicates an expected call of Registration.
func (mr *MockSenderMockRecorder) Registration(login, password, password2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Registration", reflect.TypeOf((*MockSender)(nil).Registration), login, password, password2)
}

// UploadBinary mocks base method.
func (m *MockSender) UploadBinary(entityId int32, file string) (int32, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UploadBinary", entityId, file)
	ret0, _ := ret[0].(int32)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UploadBinary indicates an expected call of UploadBinary.
func (mr *MockSenderMockRecorder) UploadBinary(entityId, file interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UploadBinary", reflect.TypeOf((*MockSender)(nil).UploadBinary), entityId, file)
}

// UploadCryptoBinary mocks base method.
func (m *MockSender) UploadCryptoBinary(entityId int32, file string) (int32, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UploadCryptoBinary", entityId, file)
	ret0, _ := ret[0].(int32)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UploadCryptoBinary indicates an expected call of UploadCryptoBinary.
func (mr *MockSenderMockRecorder) UploadCryptoBinary(entityId, file interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UploadCryptoBinary", reflect.TypeOf((*MockSender)(nil).UploadCryptoBinary), entityId, file)
}