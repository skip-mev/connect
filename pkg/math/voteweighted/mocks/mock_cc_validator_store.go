// Code generated by mockery v2.49.0. DO NOT EDIT.

package mocks

import (
	consumertypes "github.com/cosmos/interchain-security/v6/x/ccv/consumer/types"
	mock "github.com/stretchr/testify/mock"

	types "github.com/cosmos/cosmos-sdk/types"
)

// CCValidatorStore is an autogenerated mock type for the CCValidatorStore type
type CCValidatorStore struct {
	mock.Mock
}

type CCValidatorStore_Expecter struct {
	mock *mock.Mock
}

func (_m *CCValidatorStore) EXPECT() *CCValidatorStore_Expecter {
	return &CCValidatorStore_Expecter{mock: &_m.Mock}
}

// GetAllCCValidator provides a mock function with given fields: ctx
func (_m *CCValidatorStore) GetAllCCValidator(ctx types.Context) []consumertypes.CrossChainValidator {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for GetAllCCValidator")
	}

	var r0 []consumertypes.CrossChainValidator
	if rf, ok := ret.Get(0).(func(types.Context) []consumertypes.CrossChainValidator); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]consumertypes.CrossChainValidator)
		}
	}

	return r0
}

// CCValidatorStore_GetAllCCValidator_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetAllCCValidator'
type CCValidatorStore_GetAllCCValidator_Call struct {
	*mock.Call
}

// GetAllCCValidator is a helper method to define mock.On call
//   - ctx types.Context
func (_e *CCValidatorStore_Expecter) GetAllCCValidator(ctx interface{}) *CCValidatorStore_GetAllCCValidator_Call {
	return &CCValidatorStore_GetAllCCValidator_Call{Call: _e.mock.On("GetAllCCValidator", ctx)}
}

func (_c *CCValidatorStore_GetAllCCValidator_Call) Run(run func(ctx types.Context)) *CCValidatorStore_GetAllCCValidator_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(types.Context))
	})
	return _c
}

func (_c *CCValidatorStore_GetAllCCValidator_Call) Return(_a0 []consumertypes.CrossChainValidator) *CCValidatorStore_GetAllCCValidator_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *CCValidatorStore_GetAllCCValidator_Call) RunAndReturn(run func(types.Context) []consumertypes.CrossChainValidator) *CCValidatorStore_GetAllCCValidator_Call {
	_c.Call.Return(run)
	return _c
}

// GetCCValidator provides a mock function with given fields: ctx, addr
func (_m *CCValidatorStore) GetCCValidator(ctx types.Context, addr []byte) (consumertypes.CrossChainValidator, bool) {
	ret := _m.Called(ctx, addr)

	if len(ret) == 0 {
		panic("no return value specified for GetCCValidator")
	}

	var r0 consumertypes.CrossChainValidator
	var r1 bool
	if rf, ok := ret.Get(0).(func(types.Context, []byte) (consumertypes.CrossChainValidator, bool)); ok {
		return rf(ctx, addr)
	}
	if rf, ok := ret.Get(0).(func(types.Context, []byte) consumertypes.CrossChainValidator); ok {
		r0 = rf(ctx, addr)
	} else {
		r0 = ret.Get(0).(consumertypes.CrossChainValidator)
	}

	if rf, ok := ret.Get(1).(func(types.Context, []byte) bool); ok {
		r1 = rf(ctx, addr)
	} else {
		r1 = ret.Get(1).(bool)
	}

	return r0, r1
}

// CCValidatorStore_GetCCValidator_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetCCValidator'
type CCValidatorStore_GetCCValidator_Call struct {
	*mock.Call
}

// GetCCValidator is a helper method to define mock.On call
//   - ctx types.Context
//   - addr []byte
func (_e *CCValidatorStore_Expecter) GetCCValidator(ctx interface{}, addr interface{}) *CCValidatorStore_GetCCValidator_Call {
	return &CCValidatorStore_GetCCValidator_Call{Call: _e.mock.On("GetCCValidator", ctx, addr)}
}

func (_c *CCValidatorStore_GetCCValidator_Call) Run(run func(ctx types.Context, addr []byte)) *CCValidatorStore_GetCCValidator_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(types.Context), args[1].([]byte))
	})
	return _c
}

func (_c *CCValidatorStore_GetCCValidator_Call) Return(_a0 consumertypes.CrossChainValidator, _a1 bool) *CCValidatorStore_GetCCValidator_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *CCValidatorStore_GetCCValidator_Call) RunAndReturn(run func(types.Context, []byte) (consumertypes.CrossChainValidator, bool)) *CCValidatorStore_GetCCValidator_Call {
	_c.Call.Return(run)
	return _c
}

// NewCCValidatorStore creates a new instance of CCValidatorStore. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewCCValidatorStore(t interface {
	mock.TestingT
	Cleanup(func())
}) *CCValidatorStore {
	mock := &CCValidatorStore{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
