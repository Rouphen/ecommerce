// Code generated by mockery v2.16.0. DO NOT EDIT.

package mocks

import (
	context "context"
	domain "ecommerce/order-service/pkg/domain"

	mock "github.com/stretchr/testify/mock"
)

// OrderUsecase is an autogenerated mock type for the OrderUsecase type
type OrderUsecase struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, order
func (_m *OrderUsecase) Create(ctx context.Context, order *domain.Order) *domain.OrderResponse {
	ret := _m.Called(ctx, order)

	var r0 *domain.OrderResponse
	if rf, ok := ret.Get(0).(func(context.Context, *domain.Order) *domain.OrderResponse); ok {
		r0 = rf(ctx, order)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.OrderResponse)
		}
	}

	return r0
}

// Delete provides a mock function with given fields: ctx, orderID
func (_m *OrderUsecase) Delete(ctx context.Context, orderID int64) *domain.OrderResponse {
	ret := _m.Called(ctx, orderID)

	var r0 *domain.OrderResponse
	if rf, ok := ret.Get(0).(func(context.Context, int64) *domain.OrderResponse); ok {
		r0 = rf(ctx, orderID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.OrderResponse)
		}
	}

	return r0
}

type mockConstructorTestingTNewOrderUsecase interface {
	mock.TestingT
	Cleanup(func())
}

// NewOrderUsecase creates a new instance of OrderUsecase. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewOrderUsecase(t mockConstructorTestingTNewOrderUsecase) *OrderUsecase {
	mock := &OrderUsecase{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
