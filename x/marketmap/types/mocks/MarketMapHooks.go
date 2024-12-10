// Code generated by mockery v2.46.0. DO NOT EDIT.

package mocks

import (
	marketmaptypes "github.com/skip-mev/slinky/x/marketmap/types"
	mock "github.com/stretchr/testify/mock"

	types "github.com/cosmos/cosmos-sdk/types"
)

// MarketMapHooks is an autogenerated mock type for the MarketMapHooks type
type MarketMapHooks struct {
	mock.Mock
}

// AfterMarketCreated provides a mock function with given fields: ctx, market
func (_m *MarketMapHooks) AfterMarketCreated(ctx types.Context, market marketmaptypes.Market) error {
	ret := _m.Called(ctx, market)

	if len(ret) == 0 {
		panic("no return value specified for AfterMarketCreated")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(types.Context, marketmaptypes.Market) error); ok {
		r0 = rf(ctx, market)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// AfterMarketGenesis provides a mock function with given fields: ctx, tickers
func (_m *MarketMapHooks) AfterMarketGenesis(ctx types.Context, tickers map[string]marketmaptypes.Market) error {
	ret := _m.Called(ctx, tickers)

	if len(ret) == 0 {
		panic("no return value specified for AfterMarketGenesis")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(types.Context, map[string]marketmaptypes.Market) error); ok {
		r0 = rf(ctx, tickers)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// AfterMarketRemoved provides a mock function with given fields: ctx, market
func (_m *MarketMapHooks) AfterMarketRemoved(ctx types.Context, market marketmaptypes.Market) error {
	ret := _m.Called(ctx, market)

	if len(ret) == 0 {
		panic("no return value specified for AfterMarketRemoved")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(types.Context, marketmaptypes.Market) error); ok {
		r0 = rf(ctx, market)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// AfterMarketUpdated provides a mock function with given fields: ctx, market
func (_m *MarketMapHooks) AfterMarketUpdated(ctx types.Context, market marketmaptypes.Market) error {
	ret := _m.Called(ctx, market)

	if len(ret) == 0 {
		panic("no return value specified for AfterMarketUpdated")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(types.Context, marketmaptypes.Market) error); ok {
		r0 = rf(ctx, market)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewMarketMapHooks creates a new instance of MarketMapHooks. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMarketMapHooks(t interface {
	mock.TestingT
	Cleanup(func())
}) *MarketMapHooks {
	mock := &MarketMapHooks{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
