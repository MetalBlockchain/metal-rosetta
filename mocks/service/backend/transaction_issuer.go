// Code generated by mockery v2.20.2. DO NOT EDIT.

package chain

import (
	context "context"

	ids "github.com/MetalBlockchain/metalgo/ids"
	mock "github.com/stretchr/testify/mock"
)

// TransactionIssuer is an autogenerated mock type for the TransactionIssuer type
type TransactionIssuer struct {
	mock.Mock
}

// IssueTx provides a mock function with given fields: ctx, txByte
func (_m *TransactionIssuer) IssueTx(ctx context.Context, txByte []byte) (ids.ID, error) {
	ret := _m.Called(ctx, txByte)

	var r0 ids.ID
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, []byte) (ids.ID, error)); ok {
		return rf(ctx, txByte)
	}
	if rf, ok := ret.Get(0).(func(context.Context, []byte) ids.ID); ok {
		r0 = rf(ctx, txByte)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(ids.ID)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, []byte) error); ok {
		r1 = rf(ctx, txByte)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewTransactionIssuer interface {
	mock.TestingT
	Cleanup(func())
}

// NewTransactionIssuer creates a new instance of TransactionIssuer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewTransactionIssuer(t mockConstructorTestingTNewTransactionIssuer) *TransactionIssuer {
	mock := &TransactionIssuer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
