// Automatically generated by MockGen. DO NOT EDIT!
// Source: github.com/quintilesims/layer0/common/aws/s3 (interfaces: Provider)

package mock_s3

import (
	gomock "github.com/golang/mock/gomock"
	os "os"
)

// Mock of Provider interface
type MockProvider struct {
	ctrl     *gomock.Controller
	recorder *_MockProviderRecorder
}

// Recorder for MockProvider (not exported)
type _MockProviderRecorder struct {
	mock *MockProvider
}

func NewMockProvider(ctrl *gomock.Controller) *MockProvider {
	mock := &MockProvider{ctrl: ctrl}
	mock.recorder = &_MockProviderRecorder{mock}
	return mock
}

func (_m *MockProvider) EXPECT() *_MockProviderRecorder {
	return _m.recorder
}

func (_m *MockProvider) DeleteObject(_param0 string, _param1 string) error {
	ret := _m.ctrl.Call(_m, "DeleteObject", _param0, _param1)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockProviderRecorder) DeleteObject(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "DeleteObject", arg0, arg1)
}

func (_m *MockProvider) GetObject(_param0 string, _param1 string) ([]byte, error) {
	ret := _m.ctrl.Call(_m, "GetObject", _param0, _param1)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockProviderRecorder) GetObject(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetObject", arg0, arg1)
}

func (_m *MockProvider) GetObjectToFile(_param0 string, _param1 string, _param2 string, _param3 os.FileMode) error {
	ret := _m.ctrl.Call(_m, "GetObjectToFile", _param0, _param1, _param2, _param3)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockProviderRecorder) GetObjectToFile(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetObjectToFile", arg0, arg1, arg2, arg3)
}

func (_m *MockProvider) ListObjects(_param0 string, _param1 string) ([]string, error) {
	ret := _m.ctrl.Call(_m, "ListObjects", _param0, _param1)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockProviderRecorder) ListObjects(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "ListObjects", arg0, arg1)
}

func (_m *MockProvider) PutObject(_param0 string, _param1 string, _param2 []byte) error {
	ret := _m.ctrl.Call(_m, "PutObject", _param0, _param1, _param2)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockProviderRecorder) PutObject(arg0, arg1, arg2 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "PutObject", arg0, arg1, arg2)
}

func (_m *MockProvider) PutObjectFromFile(_param0 string, _param1 string, _param2 string) error {
	ret := _m.ctrl.Call(_m, "PutObjectFromFile", _param0, _param1, _param2)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockProviderRecorder) PutObjectFromFile(arg0, arg1, arg2 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "PutObjectFromFile", arg0, arg1, arg2)
}
