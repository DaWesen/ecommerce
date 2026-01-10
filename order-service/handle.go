package main

import (
	"context"
	"ecommerce/order-service/internal/service"
	"ecommerce/order-service/kitex_gen/api"

	"github.com/cloudwego/kitex/pkg/klog"
)

// OrderServiceImpl 实现 Thrift 定义的接口
type OrderServiceImpl struct {
	orderService *service.OrderService
}

// NewOrderServiceImpl 创建处理器
func NewOrderServiceImpl(orderService *service.OrderService) *OrderServiceImpl {
	return &OrderServiceImpl{
		orderService: orderService,
	}
}

// CreateOrder 创建订单
func (h *OrderServiceImpl) CreateOrder(ctx context.Context, req *api.CreateOrderReq) (resp *api.CreateOrderResp, err error) {
	klog.Infof("CreateOrder called with userId: %d, items: %d", req.UserId, len(req.Items))
	return h.orderService.CreateOrder(ctx, req)
}

// GetOrder 获取订单详情
func (h *OrderServiceImpl) GetOrder(ctx context.Context, req *api.GetOrderReq) (resp *api.GetOrderResp, err error) {
	klog.Infof("GetOrder called with orderNo: %s", req.OrderNo)
	return h.orderService.GetOrder(ctx, req)
}

// ListOrders 查询订单列表
func (h *OrderServiceImpl) ListOrders(ctx context.Context, req *api.ListOrdersReq) (resp *api.ListOrdersResp, err error) {
	klog.Infof("ListOrders called with userId: %d, page: %d", req.UserId, req.Page)
	return h.orderService.ListOrders(ctx, req)
}

// PayOrder 支付订单
func (h *OrderServiceImpl) PayOrder(ctx context.Context, req *api.PayOrderReq) (resp *api.PayOrderResp, err error) {
	klog.Infof("PayOrder called with orderNo: %s, userId: %d", req.OrderNo, req.UserId)
	return h.orderService.PayOrder(ctx, req)
}

// CancelOrder 取消订单
func (h *OrderServiceImpl) CancelOrder(ctx context.Context, req *api.CancelOrderReq) (resp *api.CancelOrderResp, err error) {
	klog.Infof("CancelOrder called with orderNo: %s, userId: %d", req.OrderNo, req.UserId)
	return h.orderService.CancelOrder(ctx, req)
}

// ApplyRefund 申请退款
func (h *OrderServiceImpl) ApplyRefund(ctx context.Context, req *api.ApplyRefundReq) (resp *api.ApplyRefundResp, err error) {
	klog.Infof("ApplyRefund called with orderNo: %s, userId: %d", req.OrderNo, req.UserId)
	return h.orderService.ApplyRefund(ctx, req)
}

// ProcessRefund 处理退款
func (h *OrderServiceImpl) ProcessRefund(ctx context.Context, req *api.ProcessRefundReq) (resp *api.ProcessRefundResp, err error) {
	klog.Infof("ProcessRefund called with refundNo: %s, processorId: %d", req.RefundNo, req.ProcessorId)
	return h.orderService.ProcessRefund(ctx, req)
}

// ReserveStock 库存预占
func (h *OrderServiceImpl) ReserveStock(ctx context.Context, req *api.ReserveStockReq) (resp *api.ReserveStockResp, err error) {
	klog.Infof("ReserveStock called with orderNo: %s, productId: %d", req.OrderNo, req.ProductId)
	return h.orderService.ReserveStock(ctx, req)
}

// ReleaseStock 释放库存
func (h *OrderServiceImpl) ReleaseStock(ctx context.Context, req *api.ReleaseStockReq) (resp *api.ReleaseStockResp, err error) {
	klog.Infof("ReleaseStock called with reserveId: %s", req.ReserveId)
	return h.orderService.ReleaseStock(ctx, req)
}

// ConfirmStock 确认库存
func (h *OrderServiceImpl) ConfirmStock(ctx context.Context, req *api.ConfirmStockReq) (resp *api.ConfirmStockResp, err error) {
	klog.Infof("ConfirmStock called with reserveId: %s", req.ReserveId)
	return h.orderService.ConfirmStock(ctx, req)
}

// ProcessTimeout 处理超时
func (h *OrderServiceImpl) ProcessTimeout(ctx context.Context, req *api.ProcessTimeoutReq) (resp *api.ProcessTimeoutResp, err error) {
	klog.Infof("ProcessTimeout called with taskId: %s", req.TaskId)
	return h.orderService.ProcessTimeout(ctx, req)
}

// GetOrderStats 获取订单统计
func (h *OrderServiceImpl) GetOrderStats(ctx context.Context, req *api.OrderStatsReq) (resp *api.OrderStatsResp, err error) {
	klog.Infof("GetOrderStats called with userId: %d", req.UserId)
	return h.orderService.GetOrderStats(ctx, req)
}

// UpdateOrderStatus 更新订单状态
func (h *OrderServiceImpl) UpdateOrderStatus(ctx context.Context, req *api.CancelOrderReq) (resp *api.CancelOrderResp, err error) {
	klog.Infof("UpdateOrderStatus called with orderNo: %s", req.OrderNo)
	return h.orderService.CancelOrder(ctx, req)
}

// ShipOrder 发货
func (h *OrderServiceImpl) ShipOrder(ctx context.Context, req *api.PayOrderReq) (resp *api.PayOrderResp, err error) {
	klog.Infof("ShipOrder called with orderNo: %s", req.OrderNo)
	return h.orderService.ShipOrder(ctx, req)
}

// ConfirmReceipt 确认收货
func (h *OrderServiceImpl) ConfirmReceipt(ctx context.Context, req *api.PayOrderReq) (resp *api.PayOrderResp, err error) {
	klog.Infof("ConfirmReceipt called with orderNo: %s", req.OrderNo)
	return h.orderService.ConfirmReceipt(ctx, req)
}
