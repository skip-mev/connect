// Code generated by mockery v2.44.2. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"

	marketmaptypes "github.com/skip-mev/connect/v2/x/marketmap/types"

	types "github.com/cosmos/cosmos-sdk/types"
)

// MarketMapKeeper is an autogenerated mock type for the MarketMapKeeper type
type MarketMapKeeper struct {
	mock.Mock
}

// GetMarket provides a mock function with given fields: ctx, tickerStr
func (_m *MarketMapKeeper) GetMarket(ctx types.Context, tickerStr string) (marketmaptypes.Market, error) {
	ret := _m.Called(ctx, tickerStr)

	if len(ret) == 0 {
		panic("no return value specified for GetMarket")
	}

	var r0 marketmaptypes.Market
	var r1 error
	if rf, ok := ret.Get(0).(func(types.Context, string) (marketmaptypes.Market, error)); ok {
		return rf(ctx, tickerStr)
	}
	if rf, ok := ret.Get(0).(func(types.Context, string) marketmaptypes.Market); ok {
		r0 = rf(ctx, tickerStr)
	} else {
		r0 = ret.Get(0).(marketmaptypes.Market)
	}

	if rf, ok := ret.Get(1).(func(types.Context, string) error); ok {
		r1 = rf(ctx, tickerStr)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewMarketMapKeeper creates a new instance of MarketMapKeeper. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMarketMapKeeper(t interface {
	mock.TestingT
	Cleanup(func())
}) *MarketMapKeeper {
	mock := &MarketMapKeeper{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
