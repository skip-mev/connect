// Code generated by mockery v2.40.2. DO NOT EDIT.

package mocks

import (
	pkgtypes "github.com/skip-mev/slinky/pkg/types"
	mock "github.com/stretchr/testify/mock"

	types "github.com/cosmos/cosmos-sdk/types"
)

// OracleKeeper is an autogenerated mock type for the OracleKeeper type
type OracleKeeper struct {
	mock.Mock
}

// GetAllCurrencyPairs provides a mock function with given fields: ctx
func (_m *OracleKeeper) GetAllCurrencyPairs(ctx types.Context) []pkgtypes.CurrencyPair {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for GetAllCurrencyPairs")
	}

	var r0 []pkgtypes.CurrencyPair
	if rf, ok := ret.Get(0).(func(types.Context) []pkgtypes.CurrencyPair); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]pkgtypes.CurrencyPair)
		}
	}

	return r0
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
