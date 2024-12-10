// Code generated by mockery v2.46.0. DO NOT EDIT.

package mocks

import (
	math "cosmossdk.io/math"
	types "github.com/cosmos/cosmos-sdk/types"
	mock "github.com/stretchr/testify/mock"
)

// StakingKeeper is an autogenerated mock type for the StakingKeeper type
type StakingKeeper struct {
	mock.Mock
}

// GetValidatorStake provides a mock function with given fields: ctx, val
func (_m *StakingKeeper) GetValidatorStake(ctx types.Context, val types.ValAddress) (math.Int, bool) {
	ret := _m.Called(ctx, val)

	if len(ret) == 0 {
		panic("no return value specified for GetValidatorStake")
	}

	var r0 math.Int
	var r1 bool
	if rf, ok := ret.Get(0).(func(types.Context, types.ValAddress) (math.Int, bool)); ok {
		return rf(ctx, val)
	}
	if rf, ok := ret.Get(0).(func(types.Context, types.ValAddress) math.Int); ok {
		r0 = rf(ctx, val)
	} else {
		r0 = ret.Get(0).(math.Int)
	}

	if rf, ok := ret.Get(1).(func(types.Context, types.ValAddress) bool); ok {
		r1 = rf(ctx, val)
	} else {
		r1 = ret.Get(1).(bool)
	}

	return r0, r1
}

// Slash provides a mock function with given fields: ctx, val, amount
func (_m *StakingKeeper) Slash(ctx types.Context, val types.ValAddress, amount math.Int) error {
	ret := _m.Called(ctx, val, amount)

	if len(ret) == 0 {
		panic("no return value specified for Slash")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(types.Context, types.ValAddress, math.Int) error); ok {
		r0 = rf(ctx, val, amount)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewStakingKeeper creates a new instance of StakingKeeper. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewStakingKeeper(t interface {
	mock.TestingT
	Cleanup(func())
}) *StakingKeeper {
	mock := &StakingKeeper{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
