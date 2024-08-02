// Code generated by mockery v2.44.1. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"

	incentivestypes "github.com/skip-mev/slinky/x/incentives/types"

	types "github.com/cosmos/cosmos-sdk/types"
)

// IncentiveKeeper is an autogenerated mock type for the IncentiveKeeper type
type IncentiveKeeper struct {
	mock.Mock
}

// AddIncentives provides a mock function with given fields: ctx, incentives
func (_m *IncentiveKeeper) AddIncentives(ctx types.Context, incentives []incentivestypes.Incentive) error {
	ret := _m.Called(ctx, incentives)

	if len(ret) == 0 {
		panic("no return value specified for AddIncentives")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(types.Context, []incentivestypes.Incentive) error); ok {
		r0 = rf(ctx, incentives)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewIncentiveKeeper creates a new instance of IncentiveKeeper. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewIncentiveKeeper(t interface {
	mock.TestingT
	Cleanup(func())
}) *IncentiveKeeper {
	mock := &IncentiveKeeper{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
