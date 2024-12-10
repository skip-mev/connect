// Code generated by mockery v2.44.1. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	types "github.com/skip-mev/slinky/providers/types"
)

// Provider is an autogenerated mock type for the Provider type
type Provider[K types.ResponseKey, V types.ResponseValue] struct {
	mock.Mock
}

// GetData provides a mock function with given fields:
func (_m *Provider[K, V]) GetData() map[K]types.ResolvedResult[V] {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetData")
	}

	var r0 map[K]types.ResolvedResult[V]
	if rf, ok := ret.Get(0).(func() map[K]types.ResolvedResult[V]); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[K]types.ResolvedResult[V])
		}
	}

	return r0
}

// IsRunning provides a mock function with given fields:
func (_m *Provider[K, V]) IsRunning() bool {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for IsRunning")
	}

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// Name provides a mock function with given fields:
func (_m *Provider[K, V]) Name() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Name")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// Start provides a mock function with given fields: _a0
func (_m *Provider[K, V]) Start(_a0 context.Context) error {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for Start")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Type provides a mock function with given fields:
func (_m *Provider[K, V]) Type() types.ProviderType {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Type")
	}

	var r0 types.ProviderType
	if rf, ok := ret.Get(0).(func() types.ProviderType); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(types.ProviderType)
	}

	return r0
}

// NewProvider creates a new instance of Provider. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewProvider[K types.ResponseKey, V types.ResponseValue](t interface {
	mock.TestingT
	Cleanup(func())
}) *Provider[K, V] {
	mock := &Provider[K, V]{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}