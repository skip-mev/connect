// Code generated by mockery v2.45.0. DO NOT EDIT.

package mocks

import (
	types "github.com/cometbft/cometbft/abci/types"
	mock "github.com/stretchr/testify/mock"
)

// ExtendedCommitCodec is an autogenerated mock type for the ExtendedCommitCodec type
type ExtendedCommitCodec struct {
	mock.Mock
}

type ExtendedCommitCodec_Expecter struct {
	mock *mock.Mock
}

func (_m *ExtendedCommitCodec) EXPECT() *ExtendedCommitCodec_Expecter {
	return &ExtendedCommitCodec_Expecter{mock: &_m.Mock}
}

// Decode provides a mock function with given fields: _a0
func (_m *ExtendedCommitCodec) Decode(_a0 []byte) (types.ExtendedCommitInfo, error) {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for Decode")
	}

	var r0 types.ExtendedCommitInfo
	var r1 error
	if rf, ok := ret.Get(0).(func([]byte) (types.ExtendedCommitInfo, error)); ok {
		return rf(_a0)
	}
	if rf, ok := ret.Get(0).(func([]byte) types.ExtendedCommitInfo); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(types.ExtendedCommitInfo)
	}

	if rf, ok := ret.Get(1).(func([]byte) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ExtendedCommitCodec_Decode_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Decode'
type ExtendedCommitCodec_Decode_Call struct {
	*mock.Call
}

// Decode is a helper method to define mock.On call
//   - _a0 []byte
func (_e *ExtendedCommitCodec_Expecter) Decode(_a0 interface{}) *ExtendedCommitCodec_Decode_Call {
	return &ExtendedCommitCodec_Decode_Call{Call: _e.mock.On("Decode", _a0)}
}

func (_c *ExtendedCommitCodec_Decode_Call) Run(run func(_a0 []byte)) *ExtendedCommitCodec_Decode_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].([]byte))
	})
	return _c
}

func (_c *ExtendedCommitCodec_Decode_Call) Return(_a0 types.ExtendedCommitInfo, _a1 error) *ExtendedCommitCodec_Decode_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *ExtendedCommitCodec_Decode_Call) RunAndReturn(run func([]byte) (types.ExtendedCommitInfo, error)) *ExtendedCommitCodec_Decode_Call {
	_c.Call.Return(run)
	return _c
}

// Encode provides a mock function with given fields: _a0
func (_m *ExtendedCommitCodec) Encode(_a0 types.ExtendedCommitInfo) ([]byte, error) {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for Encode")
	}

	var r0 []byte
	var r1 error
	if rf, ok := ret.Get(0).(func(types.ExtendedCommitInfo) ([]byte, error)); ok {
		return rf(_a0)
	}
	if rf, ok := ret.Get(0).(func(types.ExtendedCommitInfo) []byte); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	if rf, ok := ret.Get(1).(func(types.ExtendedCommitInfo) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ExtendedCommitCodec_Encode_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Encode'
type ExtendedCommitCodec_Encode_Call struct {
	*mock.Call
}

// Encode is a helper method to define mock.On call
//   - _a0 types.ExtendedCommitInfo
func (_e *ExtendedCommitCodec_Expecter) Encode(_a0 interface{}) *ExtendedCommitCodec_Encode_Call {
	return &ExtendedCommitCodec_Encode_Call{Call: _e.mock.On("Encode", _a0)}
}

func (_c *ExtendedCommitCodec_Encode_Call) Run(run func(_a0 types.ExtendedCommitInfo)) *ExtendedCommitCodec_Encode_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(types.ExtendedCommitInfo))
	})
	return _c
}

func (_c *ExtendedCommitCodec_Encode_Call) Return(_a0 []byte, _a1 error) *ExtendedCommitCodec_Encode_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *ExtendedCommitCodec_Encode_Call) RunAndReturn(run func(types.ExtendedCommitInfo) ([]byte, error)) *ExtendedCommitCodec_Encode_Call {
	_c.Call.Return(run)
	return _c
}

// NewExtendedCommitCodec creates a new instance of ExtendedCommitCodec. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewExtendedCommitCodec(t interface {
	mock.TestingT
	Cleanup(func())
}) *ExtendedCommitCodec {
	mock := &ExtendedCommitCodec{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
