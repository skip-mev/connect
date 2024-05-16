// Code generated by mockery v2.30.1. DO NOT EDIT.

package mocks

import (
	http "net/http"

	metrics "github.com/skip-mev/slinky/providers/base/api/metrics"
	mock "github.com/stretchr/testify/mock"

	time "time"

	types "github.com/skip-mev/slinky/providers/types"
)

// APIMetrics is an autogenerated mock type for the APIMetrics type
type APIMetrics struct {
	mock.Mock
}

// AddHTTPStatusCode provides a mock function with given fields: providerName, resp
func (_m *APIMetrics) AddHTTPStatusCode(providerName string, resp *http.Response) {
	_m.Called(providerName, resp)
}

// AddProviderResponse provides a mock function with given fields: providerName, id, errorCode
func (_m *APIMetrics) AddProviderResponse(providerName string, id string, errorCode types.ErrorCode) {
	_m.Called(providerName, id, errorCode)
}

// AddRPCStatusCode provides a mock function with given fields: providerName, endpoint, code
func (_m *APIMetrics) AddRPCStatusCode(providerName string, endpoint string, code metrics.RPCCode) {
	_m.Called(providerName, endpoint, code)
}

// ObserveProviderResponseLatency provides a mock function with given fields: providerName, endpoint, duration
func (_m *APIMetrics) ObserveProviderResponseLatency(providerName string, endpoint string, duration time.Duration) {
	_m.Called(providerName, endpoint, duration)
}

// NewAPIMetrics creates a new instance of APIMetrics. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewAPIMetrics(t interface {
	mock.TestingT
	Cleanup(func())
}) *APIMetrics {
	mock := &APIMetrics{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
