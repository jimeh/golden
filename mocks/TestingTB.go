// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// TestingTB is an autogenerated mock type for the TestingTB type
type TestingTB struct {
	mock.Mock
}

// Fatalf provides a mock function with given fields: format, a
func (_m *TestingTB) Fatalf(format string, a ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, format)
	_ca = append(_ca, a...)
	_m.Called(_ca...)
}

// Logf provides a mock function with given fields: format, a
func (_m *TestingTB) Logf(format string, a ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, format)
	_ca = append(_ca, a...)
	_m.Called(_ca...)
}

// Name provides a mock function with given fields:
func (_m *TestingTB) Name() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}
