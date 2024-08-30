// Code generated by mockery v2.45.0. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"

	handlers "github.com/skip-mev/connect/v2/providers/base/websocket/handlers"
)

// WebSocketConnHandler is an autogenerated mock type for the WebSocketConnHandler type
type WebSocketConnHandler struct {
	mock.Mock
}

type WebSocketConnHandler_Expecter struct {
	mock *mock.Mock
}

func (_m *WebSocketConnHandler) EXPECT() *WebSocketConnHandler_Expecter {
	return &WebSocketConnHandler_Expecter{mock: &_m.Mock}
}

// Close provides a mock function with given fields:
func (_m *WebSocketConnHandler) Close() error {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Close")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// WebSocketConnHandler_Close_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Close'
type WebSocketConnHandler_Close_Call struct {
	*mock.Call
}

// Close is a helper method to define mock.On call
func (_e *WebSocketConnHandler_Expecter) Close() *WebSocketConnHandler_Close_Call {
	return &WebSocketConnHandler_Close_Call{Call: _e.mock.On("Close")}
}

func (_c *WebSocketConnHandler_Close_Call) Run(run func()) *WebSocketConnHandler_Close_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *WebSocketConnHandler_Close_Call) Return(_a0 error) *WebSocketConnHandler_Close_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *WebSocketConnHandler_Close_Call) RunAndReturn(run func() error) *WebSocketConnHandler_Close_Call {
	_c.Call.Return(run)
	return _c
}

// Copy provides a mock function with given fields:
func (_m *WebSocketConnHandler) Copy() handlers.WebSocketConnHandler {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Copy")
	}

	var r0 handlers.WebSocketConnHandler
	if rf, ok := ret.Get(0).(func() handlers.WebSocketConnHandler); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(handlers.WebSocketConnHandler)
		}
	}

	return r0
}

// WebSocketConnHandler_Copy_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Copy'
type WebSocketConnHandler_Copy_Call struct {
	*mock.Call
}

// Copy is a helper method to define mock.On call
func (_e *WebSocketConnHandler_Expecter) Copy() *WebSocketConnHandler_Copy_Call {
	return &WebSocketConnHandler_Copy_Call{Call: _e.mock.On("Copy")}
}

func (_c *WebSocketConnHandler_Copy_Call) Run(run func()) *WebSocketConnHandler_Copy_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *WebSocketConnHandler_Copy_Call) Return(_a0 handlers.WebSocketConnHandler) *WebSocketConnHandler_Copy_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *WebSocketConnHandler_Copy_Call) RunAndReturn(run func() handlers.WebSocketConnHandler) *WebSocketConnHandler_Copy_Call {
	_c.Call.Return(run)
	return _c
}

// Dial provides a mock function with given fields:
func (_m *WebSocketConnHandler) Dial() error {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Dial")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// WebSocketConnHandler_Dial_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Dial'
type WebSocketConnHandler_Dial_Call struct {
	*mock.Call
}

// Dial is a helper method to define mock.On call
func (_e *WebSocketConnHandler_Expecter) Dial() *WebSocketConnHandler_Dial_Call {
	return &WebSocketConnHandler_Dial_Call{Call: _e.mock.On("Dial")}
}

func (_c *WebSocketConnHandler_Dial_Call) Run(run func()) *WebSocketConnHandler_Dial_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *WebSocketConnHandler_Dial_Call) Return(_a0 error) *WebSocketConnHandler_Dial_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *WebSocketConnHandler_Dial_Call) RunAndReturn(run func() error) *WebSocketConnHandler_Dial_Call {
	_c.Call.Return(run)
	return _c
}

// Read provides a mock function with given fields:
func (_m *WebSocketConnHandler) Read() ([]byte, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Read")
	}

	var r0 []byte
	var r1 error
	if rf, ok := ret.Get(0).(func() ([]byte, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() []byte); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// WebSocketConnHandler_Read_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Read'
type WebSocketConnHandler_Read_Call struct {
	*mock.Call
}

// Read is a helper method to define mock.On call
func (_e *WebSocketConnHandler_Expecter) Read() *WebSocketConnHandler_Read_Call {
	return &WebSocketConnHandler_Read_Call{Call: _e.mock.On("Read")}
}

func (_c *WebSocketConnHandler_Read_Call) Run(run func()) *WebSocketConnHandler_Read_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *WebSocketConnHandler_Read_Call) Return(_a0 []byte, _a1 error) *WebSocketConnHandler_Read_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *WebSocketConnHandler_Read_Call) RunAndReturn(run func() ([]byte, error)) *WebSocketConnHandler_Read_Call {
	_c.Call.Return(run)
	return _c
}

// Write provides a mock function with given fields: message
func (_m *WebSocketConnHandler) Write(message []byte) error {
	ret := _m.Called(message)

	if len(ret) == 0 {
		panic("no return value specified for Write")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func([]byte) error); ok {
		r0 = rf(message)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// WebSocketConnHandler_Write_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Write'
type WebSocketConnHandler_Write_Call struct {
	*mock.Call
}

// Write is a helper method to define mock.On call
//   - message []byte
func (_e *WebSocketConnHandler_Expecter) Write(message interface{}) *WebSocketConnHandler_Write_Call {
	return &WebSocketConnHandler_Write_Call{Call: _e.mock.On("Write", message)}
}

func (_c *WebSocketConnHandler_Write_Call) Run(run func(message []byte)) *WebSocketConnHandler_Write_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].([]byte))
	})
	return _c
}

func (_c *WebSocketConnHandler_Write_Call) Return(_a0 error) *WebSocketConnHandler_Write_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *WebSocketConnHandler_Write_Call) RunAndReturn(run func([]byte) error) *WebSocketConnHandler_Write_Call {
	_c.Call.Return(run)
	return _c
}

// NewWebSocketConnHandler creates a new instance of WebSocketConnHandler. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewWebSocketConnHandler(t interface {
	mock.TestingT
	Cleanup(func())
}) *WebSocketConnHandler {
	mock := &WebSocketConnHandler{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
