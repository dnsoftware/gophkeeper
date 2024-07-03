// Code generated by MockGen. DO NOT EDIT.
// Source: internal/server/domain/entity/entity.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	entity "github.com/dnsoftware/gophkeeper/internal/server/domain/entity"
	gomock "github.com/golang/mock/gomock"
)

// MockEntityRepo is a mock of EntityRepo interface.
type MockEntityRepo struct {
	ctrl     *gomock.Controller
	recorder *MockEntityRepoMockRecorder
}

// MockEntityRepoMockRecorder is the mock recorder for MockEntityRepo.
type MockEntityRepoMockRecorder struct {
	mock *MockEntityRepo
}

// NewMockEntityRepo creates a new mock instance.
func NewMockEntityRepo(ctrl *gomock.Controller) *MockEntityRepo {
	mock := &MockEntityRepo{ctrl: ctrl}
	mock.recorder = &MockEntityRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockEntityRepo) EXPECT() *MockEntityRepoMockRecorder {
	return m.recorder
}

// CreateEntity mocks base method.
func (m *MockEntityRepo) CreateEntity(ctx context.Context, entity entity.EntityModel) (int32, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateEntity", ctx, entity)
	ret0, _ := ret[0].(int32)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateEntity indicates an expected call of CreateEntity.
func (mr *MockEntityRepoMockRecorder) CreateEntity(ctx, entity interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateEntity", reflect.TypeOf((*MockEntityRepo)(nil).CreateEntity), ctx, entity)
}

// DeleteEntity mocks base method.
func (m *MockEntityRepo) DeleteEntity(ctx context.Context, id, userID int32) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteEntity", ctx, id, userID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteEntity indicates an expected call of DeleteEntity.
func (mr *MockEntityRepoMockRecorder) DeleteEntity(ctx, id, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteEntity", reflect.TypeOf((*MockEntityRepo)(nil).DeleteEntity), ctx, id, userID)
}

// GetBinaryFilenameByEntityID mocks base method.
func (m *MockEntityRepo) GetBinaryFilenameByEntityID(ctx context.Context, entityID int32) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBinaryFilenameByEntityID", ctx, entityID)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBinaryFilenameByEntityID indicates an expected call of GetBinaryFilenameByEntityID.
func (mr *MockEntityRepoMockRecorder) GetBinaryFilenameByEntityID(ctx, entityID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBinaryFilenameByEntityID", reflect.TypeOf((*MockEntityRepo)(nil).GetBinaryFilenameByEntityID), ctx, entityID)
}

// GetEntity mocks base method.
func (m *MockEntityRepo) GetEntity(ctx context.Context, id int32) (entity.EntityModel, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetEntity", ctx, id)
	ret0, _ := ret[0].(entity.EntityModel)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetEntity indicates an expected call of GetEntity.
func (mr *MockEntityRepoMockRecorder) GetEntity(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetEntity", reflect.TypeOf((*MockEntityRepo)(nil).GetEntity), ctx, id)
}

// GetEntityListByType mocks base method.
func (m *MockEntityRepo) GetEntityListByType(ctx context.Context, etype string, userID int32) (map[int32][]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetEntityListByType", ctx, etype, userID)
	ret0, _ := ret[0].(map[int32][]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetEntityListByType indicates an expected call of GetEntityListByType.
func (mr *MockEntityRepoMockRecorder) GetEntityListByType(ctx, etype, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetEntityListByType", reflect.TypeOf((*MockEntityRepo)(nil).GetEntityListByType), ctx, etype, userID)
}

// SetChunkCountForCryptoBinary mocks base method.
func (m *MockEntityRepo) SetChunkCountForCryptoBinary(ctx context.Context, entityID, chunkCount int32) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetChunkCountForCryptoBinary", ctx, entityID, chunkCount)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetChunkCountForCryptoBinary indicates an expected call of SetChunkCountForCryptoBinary.
func (mr *MockEntityRepoMockRecorder) SetChunkCountForCryptoBinary(ctx, entityID, chunkCount interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetChunkCountForCryptoBinary", reflect.TypeOf((*MockEntityRepo)(nil).SetChunkCountForCryptoBinary), ctx, entityID, chunkCount)
}

// UpdateEntity mocks base method.
func (m *MockEntityRepo) UpdateEntity(ctx context.Context, entity entity.EntityModel) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateEntity", ctx, entity)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateEntity indicates an expected call of UpdateEntity.
func (mr *MockEntityRepoMockRecorder) UpdateEntity(ctx, entity interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateEntity", reflect.TypeOf((*MockEntityRepo)(nil).UpdateEntity), ctx, entity)
}

// MockFieldRepo is a mock of FieldRepo interface.
type MockFieldRepo struct {
	ctrl     *gomock.Controller
	recorder *MockFieldRepoMockRecorder
}

// MockFieldRepoMockRecorder is the mock recorder for MockFieldRepo.
type MockFieldRepoMockRecorder struct {
	mock *MockFieldRepo
}

// NewMockFieldRepo creates a new mock instance.
func NewMockFieldRepo(ctrl *gomock.Controller) *MockFieldRepo {
	mock := &MockFieldRepo{ctrl: ctrl}
	mock.recorder = &MockFieldRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockFieldRepo) EXPECT() *MockFieldRepoMockRecorder {
	return m.recorder
}

// IsFieldType mocks base method.
func (m *MockFieldRepo) IsFieldType(ctx context.Context, id int32, ftype string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsFieldType", ctx, id, ftype)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsFieldType indicates an expected call of IsFieldType.
func (mr *MockFieldRepoMockRecorder) IsFieldType(ctx, id, ftype interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsFieldType", reflect.TypeOf((*MockFieldRepo)(nil).IsFieldType), ctx, id, ftype)
}
