// Code generated by mockery v2.50.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// NodeClient is an autogenerated mock type for the NodeClient type
type NodeClient struct {
	mock.Mock
}

type NodeClient_Expecter struct {
	mock *mock.Mock
}

func (_m *NodeClient) EXPECT() *NodeClient_Expecter {
	return &NodeClient_Expecter{mock: &_m.Mock}
}

// DeriveNodeIdentifier provides a mock function with no fields
func (_m *NodeClient) DeriveNodeIdentifier() (string, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for DeriveNodeIdentifier")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func() (string, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NodeClient_DeriveNodeIdentifier_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeriveNodeIdentifier'
type NodeClient_DeriveNodeIdentifier_Call struct {
	*mock.Call
}

// DeriveNodeIdentifier is a helper method to define mock.On call
func (_e *NodeClient_Expecter) DeriveNodeIdentifier() *NodeClient_DeriveNodeIdentifier_Call {
	return &NodeClient_DeriveNodeIdentifier_Call{Call: _e.mock.On("DeriveNodeIdentifier")}
}

func (_c *NodeClient_DeriveNodeIdentifier_Call) Run(run func()) *NodeClient_DeriveNodeIdentifier_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *NodeClient_DeriveNodeIdentifier_Call) Return(_a0 string, _a1 error) *NodeClient_DeriveNodeIdentifier_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *NodeClient_DeriveNodeIdentifier_Call) RunAndReturn(run func() (string, error)) *NodeClient_DeriveNodeIdentifier_Call {
	_c.Call.Return(run)
	return _c
}

// NewNodeClient creates a new instance of NodeClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewNodeClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *NodeClient {
	mock := &NodeClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}