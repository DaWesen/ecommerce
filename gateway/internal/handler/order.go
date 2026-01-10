package handler

import (
	"context"
	"strconv"

	"ecommerce/gateway/internal/client"
	"ecommerce/gateway/pkg/response"
	"ecommerce/order-service/kitex_gen/api"

	"github.com/cloudwego/hertz/pkg/app"
)

func safeInt64Ptr(v *int64) int64 {
	if v == nil {
		return 0
	}
	return *v
}

// CreateOrder 创建订单
func CreateOrder(clientManager *client.ClientManager) app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		userID, err := getUserIDFromContext(ctx)
		if err != nil || userID == 0 {
			response.Error(ctx, 401, "用户未登录")
			return
		}

		var req api.CreateOrderReq
		if err := ctx.BindAndValidate(&req); err != nil {
			response.Error(ctx, 400, "参数错误: "+err.Error())
			return
		}

		req.UserId = userID

		resp, err := clientManager.OrderClient.CreateOrder(c, &req)
		if err != nil {
			response.Error(ctx, 500, "创建订单失败: "+err.Error())
			return
		}

		if !resp.Success {
			response.Error(ctx, int(resp.Code), resp.Message)
			return
		}

		response.Success(ctx, map[string]interface{}{
			"order_no":     resp.OrderNo,
			"total_amount": resp.TotalAmount,
			"message":      resp.Message,
			"payment_url":  safeString(resp.PaymentUrl),
		})
	}
}

// GetOrder 获取订单详情
func GetOrder(clientManager *client.ClientManager) app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		orderNo := ctx.Param("order_no")
		if orderNo == "" {
			response.Error(ctx, 400, "订单号不能为空")
			return
		}

		userID, _ := getUserIDFromContext(ctx)

		req := &api.GetOrderReq{
			OrderNo: orderNo,
		}

		if userID > 0 {
			req.UserId = &userID
		}

		resp, err := clientManager.OrderClient.GetOrder(c, req)
		if err != nil {
			response.Error(ctx, 500, "获取订单失败: "+err.Error())
			return
		}

		if !resp.Success {
			response.Error(ctx, int(resp.Code), resp.Message)
			return
		}

		response.Success(ctx, resp.Order)
	}
}

// ListOrders 查询订单列表
func ListOrders(clientManager *client.ClientManager) app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		userID, err := getUserIDFromContext(ctx)
		if userID == 0 {
			response.Error(ctx, 401, "用户未登录")
			return
		}

		page, _ := strconv.Atoi(ctx.Query("page"))
		pageSize, _ := strconv.Atoi(ctx.Query("page_size"))

		if page <= 0 {
			page = 1
		}
		if pageSize <= 0 {
			pageSize = 10
		}
		if pageSize > 100 {
			pageSize = 100
		}

		req := &api.ListOrdersReq{
			UserId:   userID,
			Page:     int32(page),
			PageSize: int32(pageSize),
		}

		resp, err := clientManager.OrderClient.ListOrders(c, req)
		if err != nil {
			response.Error(ctx, 500, "查询订单列表失败: "+err.Error())
			return
		}

		response.SuccessWithPagination(ctx, resp.Orders, int64(resp.Total), page, pageSize)
	}
}

// PayOrder 支付订单
func PayOrder(clientManager *client.ClientManager) app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		orderNo := ctx.Param("order_no")
		if orderNo == "" {
			response.Error(ctx, 400, "订单号不能为空")
			return
		}

		userID, err := getUserIDFromContext(ctx)
		if userID == 0 {
			response.Error(ctx, 401, "用户未登录")
			return
		}

		var req api.PayOrderReq
		if err := ctx.BindAndValidate(&req); err != nil {
			response.Error(ctx, 400, "参数错误: "+err.Error())
			return
		}

		req.OrderNo = orderNo
		req.UserId = userID

		resp, err := clientManager.OrderClient.PayOrder(c, &req)
		if err != nil {
			response.Error(ctx, 500, "支付订单失败: "+err.Error())
			return
		}

		if !resp.Success {
			response.Error(ctx, int(resp.Code), resp.Message)
			return
		}

		response.Success(ctx, map[string]interface{}{
			"order_no":   orderNo,
			"message":    resp.Message,
			"new_status": resp.NewStatus_,
			"paid_at":    safeInt64Ptr(resp.PaidAt),
		})
	}
}

// CancelOrder 取消订单
func CancelOrder(clientManager *client.ClientManager) app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		orderNo := ctx.Param("order_no")
		if orderNo == "" {
			response.Error(ctx, 400, "订单号不能为空")
			return
		}

		userID, err := getUserIDFromContext(ctx)
		if userID == 0 {
			response.Error(ctx, 401, "用户未登录")
			return
		}

		var req api.CancelOrderReq
		if err := ctx.BindAndValidate(&req); err != nil {
			response.Error(ctx, 400, "参数错误: "+err.Error())
			return
		}

		req.OrderNo = orderNo
		req.UserId = userID

		resp, err := clientManager.OrderClient.CancelOrder(c, &req)
		if err != nil {
			response.Error(ctx, 500, "取消订单失败: "+err.Error())
			return
		}

		if !resp.Success {
			response.Error(ctx, int(resp.Code), resp.Message)
			return
		}

		response.Success(ctx, map[string]interface{}{
			"order_no":     orderNo,
			"message":      resp.Message,
			"new_status":   resp.NewStatus_,
			"cancelled_at": safeInt64Ptr(resp.CancelledAt),
		})
	}
}

// ApplyRefund 申请退款
func ApplyRefund(clientManager *client.ClientManager) app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		orderNo := ctx.Param("order_no")
		if orderNo == "" {
			response.Error(ctx, 400, "订单号不能为空")
			return
		}

		userID, err := getUserIDFromContext(ctx)
		if userID == 0 {
			response.Error(ctx, 401, "用户未登录")
			return
		}

		var req api.ApplyRefundReq
		if err := ctx.BindAndValidate(&req); err != nil {
			response.Error(ctx, 400, "参数错误: "+err.Error())
			return
		}

		req.OrderNo = orderNo
		req.UserId = userID

		resp, err := clientManager.OrderClient.ApplyRefund(c, &req)
		if err != nil {
			response.Error(ctx, 500, "申请退款失败: "+err.Error())
			return
		}

		if !resp.Success {
			response.Error(ctx, int(resp.Code), resp.Message)
			return
		}

		response.Success(ctx, map[string]interface{}{
			"refund_no": resp.RefundNo,
			"order_no":  orderNo,
			"message":   resp.Message,
			"status":    resp.Status,
		})
	}
}

// ProcessRefund 处理退款（管理员）
func ProcessRefund(clientManager *client.ClientManager) app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		refundNo := ctx.Param("refund_no")
		if refundNo == "" {
			response.Error(ctx, 400, "退款单号不能为空")
			return
		}

		var req api.ProcessRefundReq
		if err := ctx.BindAndValidate(&req); err != nil {
			response.Error(ctx, 400, "参数错误: "+err.Error())
			return
		}

		req.RefundNo = refundNo

		resp, err := clientManager.OrderClient.ProcessRefund(c, &req)
		if err != nil {
			response.Error(ctx, 500, "处理退款失败: "+err.Error())
			return
		}

		if !resp.Success {
			response.Error(ctx, int(resp.Code), resp.Message)
			return
		}

		response.Success(ctx, map[string]interface{}{
			"refund_no":  refundNo,
			"message":    resp.Message,
			"new_status": resp.NewStatus_,
		})
	}
}

// ShipOrder 发货订单（管理员）
func ShipOrder(clientManager *client.ClientManager) app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		orderNo := ctx.Param("order_no")
		if orderNo == "" {
			response.Error(ctx, 400, "订单号不能为空")
			return
		}

		req := &api.PayOrderReq{
			OrderNo: orderNo,
		}

		resp, err := clientManager.OrderClient.ShipOrder(c, req)
		if err != nil {
			response.Error(ctx, 500, "发货失败: "+err.Error())
			return
		}

		if !resp.Success {
			response.Error(ctx, int(resp.Code), resp.Message)
			return
		}

		response.Success(ctx, map[string]interface{}{
			"order_no":   orderNo,
			"message":    resp.Message,
			"new_status": resp.NewStatus_,
		})
	}
}

// ConfirmReceipt 确认收货
func ConfirmReceipt(clientManager *client.ClientManager) app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		orderNo := ctx.Param("order_no")
		if orderNo == "" {
			response.Error(ctx, 400, "订单号不能为空")
			return
		}

		userID, err := getUserIDFromContext(ctx)
		if userID == 0 {
			response.Error(ctx, 401, "用户未登录")
			return
		}

		req := &api.PayOrderReq{
			OrderNo: orderNo,
			UserId:  userID,
		}

		resp, err := clientManager.OrderClient.ConfirmReceipt(c, req)
		if err != nil {
			response.Error(ctx, 500, "确认收货失败: "+err.Error())
			return
		}

		if !resp.Success {
			response.Error(ctx, int(resp.Code), resp.Message)
			return
		}

		response.Success(ctx, map[string]interface{}{
			"order_no":   orderNo,
			"message":    resp.Message,
			"new_status": resp.NewStatus_,
		})
	}
}

// ListAllOrders 获取所有订单（管理员）
func ListAllOrders(clientManager *client.ClientManager) app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		page, _ := strconv.Atoi(ctx.Query("page"))
		pageSize, _ := strconv.Atoi(ctx.Query("page_size"))

		if page <= 0 {
			page = 1
		}
		if pageSize <= 0 {
			pageSize = 10
		}
		if pageSize > 100 {
			pageSize = 100
		}

		req := &api.ListOrdersReq{
			Page:     int32(page),
			PageSize: int32(pageSize),
		}

		resp, err := clientManager.OrderClient.ListOrders(c, req)
		if err != nil {
			response.Error(ctx, 500, "获取订单列表失败: "+err.Error())
			return
		}

		response.SuccessWithPagination(ctx, resp.Orders, int64(resp.Total), page, pageSize)
	}
}

// GetOrderStats 获取订单统计（管理员）
func GetOrderStats(clientManager *client.ClientManager) app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		req := &api.OrderStatsReq{}

		resp, err := clientManager.OrderClient.GetOrderStats(c, req)
		if err != nil {
			response.Error(ctx, 500, "获取订单统计失败: "+err.Error())
			return
		}

		if !resp.Success {
			response.Error(ctx, int(resp.Code), resp.Message)
			return
		}

		response.Success(ctx, map[string]interface{}{
			"total_orders":  resp.TotalOrders,
			"total_amount":  resp.TotalAmount,
			"status_counts": resp.StatusCounts,
		})
	}
}
