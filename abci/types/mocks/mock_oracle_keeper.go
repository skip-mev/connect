// Code generated by mockery v2.50.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	oracletypes "github.com/skip-mev/connect/v2/x/oracle/types"

	types "github.com/skip-mev/connect/v2/pkg/types"
)

// OracleKeeper is an autogenerated mock type for the OracleKeeper type
type OracleKeeper struct {
	mock.Mock
}

type OracleKeeper_Expecter struct {
	mock *mock.Mock
}

func (_m *OracleKeeper) EXPECT() *OracleKeeper_Expecter {
	return &OracleKeeper_Expecter{mock: &_m.Mock}
}

// GetAllCurrencyPairs provides a mock function with given fields: ctx
func (_m *OracleKeeper) GetAllCurrencyPairs(ctx context.Context) []types.CurrencyPair {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for GetAllCurrencyPairs")
	}

	var r0 []types.CurrencyPair
	if rf, ok := ret.Get(0).(func(context.Context) []types.CurrencyPair); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]types.CurrencyPair)
		}
	}

	return r0
}

// OracleKeeper_GetAllCurrencyPairs_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetAllCurrencyPairs'
type OracleKeeper_GetAllCurrencyPairs_Call struct {
	*mock.Call
}

// GetAllCurrencyPairs is a helper method to define mock.On call
//   - ctx context.Context
func (_e *OracleKeeper_Expecter) GetAllCurrencyPairs(ctx interface{}) *OracleKeeper_GetAllCurrencyPairs_Call {
	return &OracleKeeper_GetAllCurrencyPairs_Call{Call: _e.mock.On("GetAllCurrencyPairs", ctx)}
}

func (_c *OracleKeeper_GetAllCurrencyPairs_Call) Run(run func(ctx context.Context)) *OracleKeeper_GetAllCurrencyPairs_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *OracleKeeper_GetAllCurrencyPairs_Call) Return(_a0 []types.CurrencyPair) *OracleKeeper_GetAllCurrencyPairs_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *OracleKeeper_GetAllCurrencyPairs_Call) RunAndReturn(run func(context.Context) []types.CurrencyPair) *OracleKeeper_GetAllCurrencyPairs_Call {
	_c.Call.Return(run)
	return _c
}

// SetPriceForCurrencyPair provides a mock function with given fields: ctx, cp, qp
func (_m *OracleKeeper) SetPriceForCurrencyPair(ctx context.Context, cp types.CurrencyPair, qp oracletypes.QuotePrice) error {
	ret := _m.Called(ctx, cp, qp)

	if len(ret) == 0 {
		panic("no return value specified for SetPriceForCurrencyPair")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, types.CurrencyPair, oracletypes.QuotePrice) error); ok {
		r0 = rf(ctx, cp, qp)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// OracleKeeper_SetPriceForCurrencyPair_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SetPriceForCurrencyPair'
type OracleKeeper_SetPriceForCurrencyPair_Call struct {
	*mock.Call
}

// SetPriceForCurrencyPair is a helper method to define mock.On call
//   - ctx context.Context
//   - cp types.CurrencyPair
//   - qp oracletypes.QuotePrice
func (_e *OracleKeeper_Expecter) SetPriceForCurrencyPair(ctx interface{}, cp interface{}, qp interface{}) *OracleKeeper_SetPriceForCurrencyPair_Call {
	return &OracleKeeper_SetPriceForCurrencyPair_Call{Call: _e.mock.On("SetPriceForCurrencyPair", ctx, cp, qp)}
}

func (_c *OracleKeeper_SetPriceForCurrencyPair_Call) Run(run func(ctx context.Context, cp types.CurrencyPair, qp oracletypes.QuotePrice)) *OracleKeeper_SetPriceForCurrencyPair_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(types.CurrencyPair), args[2].(oracletypes.QuotePrice))
	})
	return _c
}

func (_c *OracleKeeper_SetPriceForCurrencyPair_Call) Return(_a0 error) *OracleKeeper_SetPriceForCurrencyPair_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *OracleKeeper_SetPriceForCurrencyPair_Call) RunAndReturn(run func(context.Context, types.CurrencyPair, oracletypes.QuotePrice) error) *OracleKeeper_SetPriceForCurrencyPair_Call {
	_c.Call.Return(run)
	return _c
}

// NewOracleKeeper creates a new instance of OracleKeeper. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewOracleKeeper(t interface {
	mock.TestingT
	Cleanup(func())
}) *OracleKeeper {
	mock := &OracleKeeper{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
