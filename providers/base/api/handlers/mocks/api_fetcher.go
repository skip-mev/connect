// Code generated by mockery v2.45.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	types "github.com/skip-mev/connect/v2/providers/types"
)

// APIFetcher is an autogenerated mock type for the APIFetcher type
type APIFetcher[K types.ResponseKey, V types.ResponseValue] struct {
	mock.Mock
}

type APIFetcher_Expecter[K types.ResponseKey, V types.ResponseValue] struct {
	mock *mock.Mock
}

func (_m *APIFetcher[K, V]) EXPECT() *APIFetcher_Expecter[K, V] {
	return &APIFetcher_Expecter[K, V]{mock: &_m.Mock}
}

// Fetch provides a mock function with given fields: ctx, ids
func (_m *APIFetcher[K, V]) Fetch(ctx context.Context, ids []K) types.GetResponse[K, V] {
	ret := _m.Called(ctx, ids)

	if len(ret) == 0 {
		panic("no return value specified for Fetch")
	}

	var r0 types.GetResponse[K, V]
	if rf, ok := ret.Get(0).(func(context.Context, []K) types.GetResponse[K, V]); ok {
		r0 = rf(ctx, ids)
	} else {
		r0 = ret.Get(0).(types.GetResponse[K, V])
	}

	return r0
}

// APIFetcher_Fetch_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Fetch'
type APIFetcher_Fetch_Call[K types.ResponseKey, V types.ResponseValue] struct {
	*mock.Call
}

// Fetch is a helper method to define mock.On call
//   - ctx context.Context
//   - ids []K
func (_e *APIFetcher_Expecter[K, V]) Fetch(ctx interface{}, ids interface{}) *APIFetcher_Fetch_Call[K, V] {
	return &APIFetcher_Fetch_Call[K, V]{Call: _e.mock.On("Fetch", ctx, ids)}
}

func (_c *APIFetcher_Fetch_Call[K, V]) Run(run func(ctx context.Context, ids []K)) *APIFetcher_Fetch_Call[K, V] {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].([]K))
	})
	return _c
}

func (_c *APIFetcher_Fetch_Call[K, V]) Return(_a0 types.GetResponse[K, V]) *APIFetcher_Fetch_Call[K, V] {
	_c.Call.Return(_a0)
	return _c
}

func (_c *APIFetcher_Fetch_Call[K, V]) RunAndReturn(run func(context.Context, []K) types.GetResponse[K, V]) *APIFetcher_Fetch_Call[K, V] {
	_c.Call.Return(run)
	return _c
}

// NewAPIFetcher creates a new instance of APIFetcher. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewAPIFetcher[K types.ResponseKey, V types.ResponseValue](t interface {
	mock.TestingT
	Cleanup(func())
}) *APIFetcher[K, V] {
	mock := &APIFetcher[K, V]{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
