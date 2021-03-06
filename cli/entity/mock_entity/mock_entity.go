// Automatically generated by MockGen. DO NOT EDIT!
// Source: github.com/quintilesims/layer0/cli/entity (interfaces: Entity)

package mock_entity

import (
	gomock "github.com/golang/mock/gomock"
	table "github.com/quintilesims/layer0/cli/printer/table"
)

// Mock of Entity interface
type MockEntity struct {
	ctrl     *gomock.Controller
	recorder *_MockEntityRecorder
}

// Recorder for MockEntity (not exported)
type _MockEntityRecorder struct {
	mock *MockEntity
}

func NewMockEntity(ctrl *gomock.Controller) *MockEntity {
	mock := &MockEntity{ctrl: ctrl}
	mock.recorder = &_MockEntityRecorder{mock}
	return mock
}

func (_m *MockEntity) EXPECT() *_MockEntityRecorder {
	return _m.recorder
}

func (_m *MockEntity) Table() table.Table {
	ret := _m.ctrl.Call(_m, "Table")
	ret0, _ := ret[0].(table.Table)
	return ret0
}

func (_mr *_MockEntityRecorder) Table() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Table")
}
