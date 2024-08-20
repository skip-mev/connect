// Code generated by mockery v2.44.2. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	handlers "github.com/skip-mev/connect/v2/providers/base/websocket/handlers"

	types "github.com/skip-mev/connect/v2/providers/types"
)

// WebSocketQueryHandler is an autogenerated mock type for the WebSocketQueryHandler type
type WebSocketQueryHandler[K types.ResponseKey, V types.ResponseValue] struct {
	mock.Mock
}

// Copy provides a mock function with given fields:
func (_m *WebSocketQueryHandler[K, V]) Copy() handlers.WebSocketQueryHandler[K, V] {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Copy")
	}

	var r0 handlers.WebSocketQueryHandler[K, V]
	if rf, ok := ret.Get(0).(func() handlers.WebSocketQueryHandler[K, V]); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(handlers.WebSocketQueryHandler[K, V])
		}
	}

	return r0
}

// Start provides a mock function with given fields: ctx, ids, responseCh
func (_m *WebSocketQueryHandler[K, V]) Start(ctx context.Context, ids []K, responseCh chan<- types.GetResponse[K, V]) error {
	ret := _m.Called(ctx, ids, responseCh)

	if len(ret) == 0 {
		panic("no return value specified for Start")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, []K, chan<- types.GetResponse[K, V]) error); ok {
		r0 = rf(ctx, ids, responseCh)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewWebSocketQueryHandler creates a new instance of WebSocketQueryHandler. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewWebSocketQueryHandler[K types.ResponseKey, V types.ResponseValue](t interface {
	mock.TestingT
	Cleanup(func())
}) *WebSocketQueryHandler[K, V] {
	mock := &WebSocketQueryHandler[K, V]{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
