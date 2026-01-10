package client

import (
	"context"
	"ecommerce/order-service/kitex_gen/api"
)

// CreateOrder 创建订单
func (oc *OrderClient) CreateOrder(ctx context.Context, req *api.CreateOrderReq) (*api.CreateOrderResp, error) {
	return oc.client.CreateOrder(ctx, req)
}

// PayOrder 支付订单
func (oc *OrderClient) PayOrder(ctx context.Context, req *api.PayOrderReq) (*api.PayOrderResp, error) {
	return oc.client.PayOrder(ctx, req)
}

// CancelOrder 取消订单
func (oc *OrderClient) CancelOrder(ctx context.Context, req *api.CancelOrderReq) (*api.CancelOrderResp, error) {
	return oc.client.CancelOrder(ctx, req)
}

// GetOrder 获取订单详情
func (oc *OrderClient) GetOrder(ctx context.Context, req *api.GetOrderReq) (*api.GetOrderResp, error) {
	return oc.client.GetOrder(ctx, req)
}

// ListOrders 获取订单列表
func (oc *OrderClient) ListOrders(ctx context.Context, req *api.ListOrdersReq) (*api.ListOrdersResp, error) {
	return oc.client.ListOrders(ctx, req)
}

// ApplyRefund 申请退款
func (oc *OrderClient) ApplyRefund(ctx context.Context, req *api.ApplyRefundReq) (*api.ApplyRefundResp, error) {
	return oc.client.ApplyRefund(ctx, req)
}

// ProcessRefund 处理退款（管理员）
func (oc *OrderClient) ProcessRefund(ctx context.Context, req *api.ProcessRefundReq) (*api.ProcessRefundResp, error) {
	return oc.client.ProcessRefund(ctx, req)
}

// ReserveStock 预留库存
func (oc *OrderClient) ReserveStock(ctx context.Context, req *api.ReserveStockReq) (*api.ReserveStockResp, error) {
	return oc.client.ReserveStock(ctx, req)
}

// ReleaseStock 释放库存
func (oc *OrderClient) ReleaseStock(ctx context.Context, req *api.ReleaseStockReq) (*api.ReleaseStockResp, error) {
	return oc.client.ReleaseStock(ctx, req)
}

// ConfirmStock 确认库存
func (oc *OrderClient) ConfirmStock(ctx context.Context, req *api.ConfirmStockReq) (*api.ConfirmStockResp, error) {
	return oc.client.ConfirmStock(ctx, req)
}

// ProcessTimeout 处理超时
func (oc *OrderClient) ProcessTimeout(ctx context.Context, req *api.ProcessTimeoutReq) (*api.ProcessTimeoutResp, error) {
	return oc.client.ProcessTimeout(ctx, req)
}

// GetOrderStats 获取订单统计
func (oc *OrderClient) GetOrderStats(ctx context.Context, req *api.OrderStatsReq) (*api.OrderStatsResp, error) {
	return oc.client.GetOrderStats(ctx, req)
}

// UpdateOrderStatus 更新订单状态
func (oc *OrderClient) UpdateOrderStatus(ctx context.Context, req *api.CancelOrderReq) (*api.CancelOrderResp, error) {
	return oc.client.UpdateOrderStatus(ctx, req)
}

// ShipOrder 发货订单
func (oc *OrderClient) ShipOrder(ctx context.Context, req *api.PayOrderReq) (*api.PayOrderResp, error) {
	return oc.client.ShipOrder(ctx, req)
}

// ConfirmReceipt 确认收货
func (oc *OrderClient) ConfirmReceipt(ctx context.Context, req *api.PayOrderReq) (*api.PayOrderResp, error) {
	return oc.client.ConfirmReceipt(ctx, req)
}
