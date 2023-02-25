// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// ISQLRelational is an autogenerated mock type for the ISQLRelational type
type ISQLRelational struct {
	mock.Mock
}

// Execute provides a mock function with given fields: query, dataModel
func (_m *ISQLRelational) Execute(query string, dataModel interface{}) (interface{}, error) {
	ret := _m.Called(query, dataModel)

	var r0 interface{}
	if rf, ok := ret.Get(0).(func(string, interface{}) interface{}); ok {
		r0 = rf(query, dataModel)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, interface{}) error); ok {
		r1 = rf(query, dataModel)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}