// Code generated by mockery v2.50.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	osmosis "github.com/skip-mev/connect/v2/providers/apis/defi/osmosis"
)

// Client is an autogenerated mock type for the Client type
type Client struct {
	mock.Mock
}

type Client_Expecter struct {
	mock *mock.Mock
}

func (_m *Client) EXPECT() *Client_Expecter {
	return &Client_Expecter{mock: &_m.Mock}
}

// SpotPrice provides a mock function with given fields: ctx, poolID, baseAsset, quoteAsset
func (_m *Client) SpotPrice(ctx context.Context, poolID uint64, baseAsset string, quoteAsset string) (osmosis.WrappedSpotPriceResponse, error) {
	ret := _m.Called(ctx, poolID, baseAsset, quoteAsset)

	if len(ret) == 0 {
		panic("no return value specified for SpotPrice")
	}

	var r0 osmosis.WrappedSpotPriceResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint64, string, string) (osmosis.WrappedSpotPriceResponse, error)); ok {
		return rf(ctx, poolID, baseAsset, quoteAsset)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint64, string, string) osmosis.WrappedSpotPriceResponse); ok {
		r0 = rf(ctx, poolID, baseAsset, quoteAsset)
	} else {
		r0 = ret.Get(0).(osmosis.WrappedSpotPriceResponse)
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint64, string, string) error); ok {
		r1 = rf(ctx, poolID, baseAsset, quoteAsset)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Client_SpotPrice_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SpotPrice'
type Client_SpotPrice_Call struct {
	*mock.Call
}

// SpotPrice is a helper method to define mock.On call
//   - ctx context.Context
//   - poolID uint64
//   - baseAsset string
//   - quoteAsset string
func (_e *Client_Expecter) SpotPrice(ctx interface{}, poolID interface{}, baseAsset interface{}, quoteAsset interface{}) *Client_SpotPrice_Call {
	return &Client_SpotPrice_Call{Call: _e.mock.On("SpotPrice", ctx, poolID, baseAsset, quoteAsset)}
}

func (_c *Client_SpotPrice_Call) Run(run func(ctx context.Context, poolID uint64, baseAsset string, quoteAsset string)) *Client_SpotPrice_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uint64), args[2].(string), args[3].(string))
	})
	return _c
}

func (_c *Client_SpotPrice_Call) Return(_a0 osmosis.WrappedSpotPriceResponse, _a1 error) *Client_SpotPrice_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Client_SpotPrice_Call) RunAndReturn(run func(context.Context, uint64, string, string) (osmosis.WrappedSpotPriceResponse, error)) *Client_SpotPrice_Call {
	_c.Call.Return(run)
	return _c
}

// NewClient creates a new instance of Client. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *Client {
	mock := &Client{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
