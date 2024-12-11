// Code generated by mockery v2.50.0. DO NOT EDIT.

package mocks

import (
	context "context"

	math "cosmossdk.io/math"

	mock "github.com/stretchr/testify/mock"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	types "github.com/cosmos/cosmos-sdk/types"
)

// ValidatorStore is an autogenerated mock type for the ValidatorStore type
type ValidatorStore struct {
	mock.Mock
}

type ValidatorStore_Expecter struct {
	mock *mock.Mock
}

func (_m *ValidatorStore) EXPECT() *ValidatorStore_Expecter {
	return &ValidatorStore_Expecter{mock: &_m.Mock}
}

// TotalBondedTokens provides a mock function with given fields: ctx
func (_m *ValidatorStore) TotalBondedTokens(ctx context.Context) (math.Int, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for TotalBondedTokens")
	}

	var r0 math.Int
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (math.Int, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) math.Int); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Get(0).(math.Int)
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ValidatorStore_TotalBondedTokens_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'TotalBondedTokens'
type ValidatorStore_TotalBondedTokens_Call struct {
	*mock.Call
}

// TotalBondedTokens is a helper method to define mock.On call
//   - ctx context.Context
func (_e *ValidatorStore_Expecter) TotalBondedTokens(ctx interface{}) *ValidatorStore_TotalBondedTokens_Call {
	return &ValidatorStore_TotalBondedTokens_Call{Call: _e.mock.On("TotalBondedTokens", ctx)}
}

func (_c *ValidatorStore_TotalBondedTokens_Call) Run(run func(ctx context.Context)) *ValidatorStore_TotalBondedTokens_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *ValidatorStore_TotalBondedTokens_Call) Return(_a0 math.Int, _a1 error) *ValidatorStore_TotalBondedTokens_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *ValidatorStore_TotalBondedTokens_Call) RunAndReturn(run func(context.Context) (math.Int, error)) *ValidatorStore_TotalBondedTokens_Call {
	_c.Call.Return(run)
	return _c
}

// ValidatorByConsAddr provides a mock function with given fields: ctx, addr
func (_m *ValidatorStore) ValidatorByConsAddr(ctx context.Context, addr types.ConsAddress) (stakingtypes.ValidatorI, error) {
	ret := _m.Called(ctx, addr)

	if len(ret) == 0 {
		panic("no return value specified for ValidatorByConsAddr")
	}

	var r0 stakingtypes.ValidatorI
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, types.ConsAddress) (stakingtypes.ValidatorI, error)); ok {
		return rf(ctx, addr)
	}
	if rf, ok := ret.Get(0).(func(context.Context, types.ConsAddress) stakingtypes.ValidatorI); ok {
		r0 = rf(ctx, addr)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(stakingtypes.ValidatorI)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, types.ConsAddress) error); ok {
		r1 = rf(ctx, addr)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ValidatorStore_ValidatorByConsAddr_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ValidatorByConsAddr'
type ValidatorStore_ValidatorByConsAddr_Call struct {
	*mock.Call
}

// ValidatorByConsAddr is a helper method to define mock.On call
//   - ctx context.Context
//   - addr types.ConsAddress
func (_e *ValidatorStore_Expecter) ValidatorByConsAddr(ctx interface{}, addr interface{}) *ValidatorStore_ValidatorByConsAddr_Call {
	return &ValidatorStore_ValidatorByConsAddr_Call{Call: _e.mock.On("ValidatorByConsAddr", ctx, addr)}
}

func (_c *ValidatorStore_ValidatorByConsAddr_Call) Run(run func(ctx context.Context, addr types.ConsAddress)) *ValidatorStore_ValidatorByConsAddr_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(types.ConsAddress))
	})
	return _c
}

func (_c *ValidatorStore_ValidatorByConsAddr_Call) Return(_a0 stakingtypes.ValidatorI, _a1 error) *ValidatorStore_ValidatorByConsAddr_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *ValidatorStore_ValidatorByConsAddr_Call) RunAndReturn(run func(context.Context, types.ConsAddress) (stakingtypes.ValidatorI, error)) *ValidatorStore_ValidatorByConsAddr_Call {
	_c.Call.Return(run)
	return _c
}

// NewValidatorStore creates a new instance of ValidatorStore. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewValidatorStore(t interface {
	mock.TestingT
	Cleanup(func())
}) *ValidatorStore {
	mock := &ValidatorStore{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
