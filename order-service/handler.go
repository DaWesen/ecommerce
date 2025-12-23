package main

import (
	"context"
	api "ecommerce/order-service/kitex_gen/api"
)

// OrderServiceImpl implements the last service interface defined in the IDL.
type OrderServiceImpl struct{}

// CreateOrder implements the OrderServiceImpl interface.
func (s *OrderServiceImpl) CreateOrder(ctx context.Context, req *api.CreateOrderReq) (resp *api.CreateOrderResp, err error) {
	// TODO: Your code here...
	return
}

// PayOrder implements the OrderServiceImpl interface.
func (s *OrderServiceImpl) PayOrder(ctx context.Context, req *api.PayOrderReq) (resp *api.PayOrderResp, err error) {
	// TODO: Your code here...
	return
}

// CancelOrder implements the OrderServiceImpl interface.
func (s *OrderServiceImpl) CancelOrder(ctx context.Context, req *api.CancelOrderReq) (resp *api.CancelOrderResp, err error) {
	// TODO: Your code here...
	return
}

// GetOrderDetail implements the OrderServiceImpl interface.
func (s *OrderServiceImpl) GetOrderDetail(ctx context.Context, req *api.GetOrderDetailReq) (resp *api.GetOrderDetailResp, err error) {
	// TODO: Your code here...
	return
}

// QueryOrders implements the OrderServiceImpl interface.
func (s *OrderServiceImpl) QueryOrders(ctx context.Context, req *api.QueryOrderReq) (resp *api.QueryOrderResp, err error) {
	// TODO: Your code here...
	return
}

// UpdateOrderStatus implements the OrderServiceImpl interface.
func (s *OrderServiceImpl) UpdateOrderStatus(ctx context.Context, req *api.UpdateOrderStatusReq) (resp *api.UpdateOrderStatusResp, err error) {
	// TODO: Your code here...
	return
}
