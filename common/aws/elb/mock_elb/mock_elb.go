// Automatically generated by MockGen. DO NOT EDIT!
// Source: github.com/quintilesims/layer0/common/aws/elb (interfaces: Provider)

package mock_elb

import (
	gomock "github.com/golang/mock/gomock"
	elb "github.com/quintilesims/layer0/common/aws/elb"
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

func (_m *MockProvider) ConfigureHealthCheck(_param0 string, _param1 *elb.HealthCheck) error {
	ret := _m.ctrl.Call(_m, "ConfigureHealthCheck", _param0, _param1)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockProviderRecorder) ConfigureHealthCheck(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "ConfigureHealthCheck", arg0, arg1)
}

func (_m *MockProvider) CreateLoadBalancer(_param0 string, _param1 string, _param2 []*string, _param3 []*string, _param4 []*elb.Listener) (*string, error) {
	ret := _m.ctrl.Call(_m, "CreateLoadBalancer", _param0, _param1, _param2, _param3, _param4)
	ret0, _ := ret[0].(*string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockProviderRecorder) CreateLoadBalancer(arg0, arg1, arg2, arg3, arg4 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "CreateLoadBalancer", arg0, arg1, arg2, arg3, arg4)
}

func (_m *MockProvider) CreateLoadBalancerListeners(_param0 string, _param1 []*elb.Listener) error {
	ret := _m.ctrl.Call(_m, "CreateLoadBalancerListeners", _param0, _param1)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockProviderRecorder) CreateLoadBalancerListeners(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "CreateLoadBalancerListeners", arg0, arg1)
}

func (_m *MockProvider) DeleteLoadBalancer(_param0 string) error {
	ret := _m.ctrl.Call(_m, "DeleteLoadBalancer", _param0)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockProviderRecorder) DeleteLoadBalancer(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "DeleteLoadBalancer", arg0)
}

func (_m *MockProvider) DeleteLoadBalancerListeners(_param0 string, _param1 []*elb.Listener) error {
	ret := _m.ctrl.Call(_m, "DeleteLoadBalancerListeners", _param0, _param1)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockProviderRecorder) DeleteLoadBalancerListeners(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "DeleteLoadBalancerListeners", arg0, arg1)
}

func (_m *MockProvider) DeregisterInstancesFromLoadBalancer(_param0 string, _param1 []string) error {
	ret := _m.ctrl.Call(_m, "DeregisterInstancesFromLoadBalancer", _param0, _param1)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockProviderRecorder) DeregisterInstancesFromLoadBalancer(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "DeregisterInstancesFromLoadBalancer", arg0, arg1)
}

func (_m *MockProvider) DescribeInstanceHealth(_param0 string) ([]*elb.InstanceState, error) {
	ret := _m.ctrl.Call(_m, "DescribeInstanceHealth", _param0)
	ret0, _ := ret[0].([]*elb.InstanceState)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockProviderRecorder) DescribeInstanceHealth(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "DescribeInstanceHealth", arg0)
}

func (_m *MockProvider) DescribeLoadBalancer(_param0 string) (*elb.LoadBalancerDescription, error) {
	ret := _m.ctrl.Call(_m, "DescribeLoadBalancer", _param0)
	ret0, _ := ret[0].(*elb.LoadBalancerDescription)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockProviderRecorder) DescribeLoadBalancer(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "DescribeLoadBalancer", arg0)
}

func (_m *MockProvider) DescribeLoadBalancers() ([]*elb.LoadBalancerDescription, error) {
	ret := _m.ctrl.Call(_m, "DescribeLoadBalancers")
	ret0, _ := ret[0].([]*elb.LoadBalancerDescription)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockProviderRecorder) DescribeLoadBalancers() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "DescribeLoadBalancers")
}

func (_m *MockProvider) RegisterInstancesWithLoadBalancer(_param0 string, _param1 []string) error {
	ret := _m.ctrl.Call(_m, "RegisterInstancesWithLoadBalancer", _param0, _param1)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockProviderRecorder) RegisterInstancesWithLoadBalancer(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "RegisterInstancesWithLoadBalancer", arg0, arg1)
}
