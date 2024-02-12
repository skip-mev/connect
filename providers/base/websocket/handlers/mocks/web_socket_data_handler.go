// Code generated by mockery v2.40.3. DO NOT EDIT.

package mocks

import (
	handlers "github.com/skip-mev/slinky/providers/base/websocket/handlers"
	mock "github.com/stretchr/testify/mock"

	types "github.com/skip-mev/slinky/providers/types"
)

// WebSocketDataHandler is an autogenerated mock type for the WebSocketDataHandler type
type WebSocketDataHandler[K types.ResponseKey, V types.ResponseValue] struct {
	mock.Mock
}

// CreateMessages provides a mock function with given fields: ids
func (_m *WebSocketDataHandler[K, V]) CreateMessages(ids []K) ([]handlers.WebsocketEncodedMessage, error) {
	ret := _m.Called(ids)

	if len(ret) == 0 {
		panic("no return value specified for CreateMessages")
	}

	var r0 []handlers.WebsocketEncodedMessage
	var r1 error
	if rf, ok := ret.Get(0).(func([]K) ([]handlers.WebsocketEncodedMessage, error)); ok {
		return rf(ids)
	}
	if rf, ok := ret.Get(0).(func([]K) []handlers.WebsocketEncodedMessage); ok {
		r0 = rf(ids)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]handlers.WebsocketEncodedMessage)
		}
	}

	if rf, ok := ret.Get(1).(func([]K) error); ok {
		r1 = rf(ids)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// HandleMessage provides a mock function with given fields: message
func (_m *WebSocketDataHandler[K, V]) HandleMessage(message []byte) (types.GetResponse[K, V], []handlers.WebsocketEncodedMessage, error) {
	ret := _m.Called(message)

	if len(ret) == 0 {
		panic("no return value specified for HandleMessage")
	}

	var r0 types.GetResponse[K, V]
	var r1 []handlers.WebsocketEncodedMessage
	var r2 error
	if rf, ok := ret.Get(0).(func([]byte) (types.GetResponse[K, V], []handlers.WebsocketEncodedMessage, error)); ok {
		return rf(message)
	}
	if rf, ok := ret.Get(0).(func([]byte) types.GetResponse[K, V]); ok {
		r0 = rf(message)
	} else {
		r0 = ret.Get(0).(types.GetResponse[K, V])
	}

	if rf, ok := ret.Get(1).(func([]byte) []handlers.WebsocketEncodedMessage); ok {
		r1 = rf(message)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).([]handlers.WebsocketEncodedMessage)
		}
	}

	if rf, ok := ret.Get(2).(func([]byte) error); ok {
		r2 = rf(message)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// HeartBeatMessages provides a mock function with given fields:
func (_m *WebSocketDataHandler[K, V]) HeartBeatMessages() ([]handlers.WebsocketEncodedMessage, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for HeartBeatMessages")
	}

	var r0 []handlers.WebsocketEncodedMessage
	var r1 error
	if rf, ok := ret.Get(0).(func() ([]handlers.WebsocketEncodedMessage, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() []handlers.WebsocketEncodedMessage); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]handlers.WebsocketEncodedMessage)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewWebSocketDataHandler creates a new instance of WebSocketDataHandler. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewWebSocketDataHandler[K types.ResponseKey, V types.ResponseValue](t interface {
	mock.TestingT
	Cleanup(func())
}) *WebSocketDataHandler[K, V] {
	mock := &WebSocketDataHandler[K, V]{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
