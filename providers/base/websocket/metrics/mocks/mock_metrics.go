// Code generated by mockery v2.46.0. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"

	metrics "github.com/skip-mev/connect/v2/providers/base/websocket/metrics"

	time "time"
)

// WebSocketMetrics is an autogenerated mock type for the WebSocketMetrics type
type WebSocketMetrics struct {
	mock.Mock
}

type WebSocketMetrics_Expecter struct {
	mock *mock.Mock
}

func (_m *WebSocketMetrics) EXPECT() *WebSocketMetrics_Expecter {
	return &WebSocketMetrics_Expecter{mock: &_m.Mock}
}

// AddWebSocketConnectionStatus provides a mock function with given fields: provider, status
func (_m *WebSocketMetrics) AddWebSocketConnectionStatus(provider string, status metrics.ConnectionStatus) {
	_m.Called(provider, status)
}

// WebSocketMetrics_AddWebSocketConnectionStatus_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'AddWebSocketConnectionStatus'
type WebSocketMetrics_AddWebSocketConnectionStatus_Call struct {
	*mock.Call
}

// AddWebSocketConnectionStatus is a helper method to define mock.On call
//   - provider string
//   - status metrics.ConnectionStatus
func (_e *WebSocketMetrics_Expecter) AddWebSocketConnectionStatus(provider interface{}, status interface{}) *WebSocketMetrics_AddWebSocketConnectionStatus_Call {
	return &WebSocketMetrics_AddWebSocketConnectionStatus_Call{Call: _e.mock.On("AddWebSocketConnectionStatus", provider, status)}
}

func (_c *WebSocketMetrics_AddWebSocketConnectionStatus_Call) Run(run func(provider string, status metrics.ConnectionStatus)) *WebSocketMetrics_AddWebSocketConnectionStatus_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(metrics.ConnectionStatus))
	})
	return _c
}

func (_c *WebSocketMetrics_AddWebSocketConnectionStatus_Call) Return() *WebSocketMetrics_AddWebSocketConnectionStatus_Call {
	_c.Call.Return()
	return _c
}

func (_c *WebSocketMetrics_AddWebSocketConnectionStatus_Call) RunAndReturn(run func(string, metrics.ConnectionStatus)) *WebSocketMetrics_AddWebSocketConnectionStatus_Call {
	_c.Call.Return(run)
	return _c
}

// AddWebSocketDataHandlerStatus provides a mock function with given fields: provider, status
func (_m *WebSocketMetrics) AddWebSocketDataHandlerStatus(provider string, status metrics.HandlerStatus) {
	_m.Called(provider, status)
}

// WebSocketMetrics_AddWebSocketDataHandlerStatus_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'AddWebSocketDataHandlerStatus'
type WebSocketMetrics_AddWebSocketDataHandlerStatus_Call struct {
	*mock.Call
}

// AddWebSocketDataHandlerStatus is a helper method to define mock.On call
//   - provider string
//   - status metrics.HandlerStatus
func (_e *WebSocketMetrics_Expecter) AddWebSocketDataHandlerStatus(provider interface{}, status interface{}) *WebSocketMetrics_AddWebSocketDataHandlerStatus_Call {
	return &WebSocketMetrics_AddWebSocketDataHandlerStatus_Call{Call: _e.mock.On("AddWebSocketDataHandlerStatus", provider, status)}
}

func (_c *WebSocketMetrics_AddWebSocketDataHandlerStatus_Call) Run(run func(provider string, status metrics.HandlerStatus)) *WebSocketMetrics_AddWebSocketDataHandlerStatus_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(metrics.HandlerStatus))
	})
	return _c
}

func (_c *WebSocketMetrics_AddWebSocketDataHandlerStatus_Call) Return() *WebSocketMetrics_AddWebSocketDataHandlerStatus_Call {
	_c.Call.Return()
	return _c
}

func (_c *WebSocketMetrics_AddWebSocketDataHandlerStatus_Call) RunAndReturn(run func(string, metrics.HandlerStatus)) *WebSocketMetrics_AddWebSocketDataHandlerStatus_Call {
	_c.Call.Return(run)
	return _c
}

// ObserveWebSocketLatency provides a mock function with given fields: provider, duration
func (_m *WebSocketMetrics) ObserveWebSocketLatency(provider string, duration time.Duration) {
	_m.Called(provider, duration)
}

// WebSocketMetrics_ObserveWebSocketLatency_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ObserveWebSocketLatency'
type WebSocketMetrics_ObserveWebSocketLatency_Call struct {
	*mock.Call
}

// ObserveWebSocketLatency is a helper method to define mock.On call
//   - provider string
//   - duration time.Duration
func (_e *WebSocketMetrics_Expecter) ObserveWebSocketLatency(provider interface{}, duration interface{}) *WebSocketMetrics_ObserveWebSocketLatency_Call {
	return &WebSocketMetrics_ObserveWebSocketLatency_Call{Call: _e.mock.On("ObserveWebSocketLatency", provider, duration)}
}

func (_c *WebSocketMetrics_ObserveWebSocketLatency_Call) Run(run func(provider string, duration time.Duration)) *WebSocketMetrics_ObserveWebSocketLatency_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(time.Duration))
	})
	return _c
}

func (_c *WebSocketMetrics_ObserveWebSocketLatency_Call) Return() *WebSocketMetrics_ObserveWebSocketLatency_Call {
	_c.Call.Return()
	return _c
}

func (_c *WebSocketMetrics_ObserveWebSocketLatency_Call) RunAndReturn(run func(string, time.Duration)) *WebSocketMetrics_ObserveWebSocketLatency_Call {
	_c.Call.Return(run)
	return _c
}

// NewWebSocketMetrics creates a new instance of WebSocketMetrics. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewWebSocketMetrics(t interface {
	mock.TestingT
	Cleanup(func())
}) *WebSocketMetrics {
	mock := &WebSocketMetrics{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
