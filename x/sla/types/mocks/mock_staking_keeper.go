// Code generated by mockery v2.43.2. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	types "github.com/cosmos/cosmos-sdk/types"
)

// StakingKeeper is an autogenerated mock type for the StakingKeeper type
type StakingKeeper struct {
	mock.Mock
}

// GetLastValidatorPower provides a mock function with given fields: ctx, operator
func (_m *StakingKeeper) GetLastValidatorPower(ctx context.Context, operator types.ValAddress) (int64, error) {
	ret := _m.Called(ctx, operator)

	if len(ret) == 0 {
		panic("no return value specified for GetLastValidatorPower")
	}

	var r0 int64
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, types.ValAddress) (int64, error)); ok {
		return rf(ctx, operator)
	}
	if rf, ok := ret.Get(0).(func(context.Context, types.ValAddress) int64); ok {
		r0 = rf(ctx, operator)
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func(context.Context, types.ValAddress) error); ok {
		r1 = rf(ctx, operator)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
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
