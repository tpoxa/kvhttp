// Code generated by mockery v2.20.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// IStore is an autogenerated mock type for the IStore type
type IStore struct {
	mock.Mock
}

// Delete provides a mock function with given fields: ctx, key
func (_m *IStore) Delete(ctx context.Context, key string) {
	_m.Called(ctx, key)
}

// Get provides a mock function with given fields: ctx, key
func (_m *IStore) Get(ctx context.Context, key string) ([]byte, bool) {
	ret := _m.Called(ctx, key)

	var r0 []byte
	var r1 bool
	if rf, ok := ret.Get(0).(func(context.Context, string) ([]byte, bool)); ok {
		return rf(ctx, key)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) []byte); ok {
		r0 = rf(ctx, key)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) bool); ok {
		r1 = rf(ctx, key)
	} else {
		r1 = ret.Get(1).(bool)
	}

	return r0, r1
}

// Set provides a mock function with given fields: ctx, key, value
func (_m *IStore) Set(ctx context.Context, key string, value []byte) {
	_m.Called(ctx, key, value)
}

// SetIfCheckSumMatch provides a mock function with given fields: ctx, key, value, checksum
func (_m *IStore) SetIfCheckSumMatch(ctx context.Context, key string, value []byte, checksum string) error {
	ret := _m.Called(ctx, key, value, checksum)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, []byte, string) error); ok {
		r0 = rf(ctx, key, value, checksum)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewIStore interface {
	mock.TestingT
	Cleanup(func())
}

// NewIStore creates a new instance of IStore. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewIStore(t mockConstructorTestingTNewIStore) *IStore {
	mock := &IStore{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}