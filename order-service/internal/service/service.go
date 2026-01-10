package service

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"runtime/debug"
	"time"

	"ecommerce/order-service/internal/dao/dao"
	"ecommerce/order-service/internal/dao/interfaces"
	"ecommerce/order-service/internal/model"
	"ecommerce/order-service/kitex_gen/api"

	"github.com/cloudwego/kitex/pkg/klog"
	"gorm.io/gorm"
)

// 错误定义
var (
	ErrUserNotValid                = errors.New("用户状态异常")
	ErrProductNotOnline            = errors.New("商品未上架")
	ErrStockNotEnough              = errors.New("库存不足")
	ErrOrderNotFound               = errors.New("订单不存在")
	ErrOrderItemEmpty              = errors.New("订单商品不能为空")
	ErrOrderStatusWrong            = errors.New("订单状态不正确")
	ErrPermissionDenied            = errors.New("无权操作此订单")
	ErrPaymentFailed               = errors.New("支付失败")
	ErrRefundFailed                = errors.New("退款失败")
	ErrRefundAlreadyExists         = errors.New("该订单已申请退款")
	ErrRefundAmountInvalid         = errors.New("退款金额无效")
	ErrReservationNotFound         = errors.New("库存预占记录不存在")
	ErrReservationExpired          = errors.New("库存预占已过期")
	ErrReservationStatusWrong      = errors.New("库存预占状态不正确")
	ErrTimeoutTaskNotFound         = errors.New("超时任务不存在")
	ErrTimeoutTaskAlreadyProcessed = errors.New("超时任务已处理")
	ErrAddressRequired             = errors.New("收货地址不能为空")
	ErrPhoneRequired               = errors.New("联系电话不能为空")
	ErrOrderAlreadyPaid            = errors.New("订单已支付")
	ErrOrderAlreadyCancelled       = errors.New("订单已取消")
	ErrOrderAlreadyCompleted       = errors.New("订单已完成")
	ErrOrderAlreadyRefunded        = errors.New("订单已退款")
	ErrShippingInfoRequired        = errors.New("物流信息不能为空")
	ErrRefundReasonRequired        = errors.New("退款原因不能为空")
	ErrInvalidOrderStatus          = errors.New("无效的订单状态")
	ErrInvalidRefundStatus         = errors.New("无效的退款状态")
)

// OrderService 订单服务
type OrderService struct {
	db            *gorm.DB
	daoFactory    *dao.DaoFactory
	userClient    interfaces.IUserClient
	productClient interfaces.IProductClient
}

// NewOrderService 创建订单服务实例
func NewOrderService(
	db *gorm.DB,
	daoFactory *dao.DaoFactory,
	userClient interfaces.IUserClient,
	productClient interfaces.IProductClient,
) *OrderService {
	return &OrderService{
		db:            db,
		daoFactory:    daoFactory,
		userClient:    userClient,
		productClient: productClient,
	}
}

// CreateOrder 创建订单
func (s *OrderService) CreateOrder(ctx context.Context, req *api.CreateOrderReq) (*api.CreateOrderResp, error) {
	// 添加 panic 恢复
	defer func() {
		if r := recover(); r != nil {
			klog.Errorf("CreateOrder panic: %v", r)
			debug.PrintStack()
		}
	}()

	klog.Infof("CreateOrder 开始执行，用户ID: %d，商品数量: %d", req.UserId, len(req.Items))

	if req == nil {
		return &api.CreateOrderResp{
			Success: false,
			Code:    400,
			Message: "请求参数为空",
		}, nil
	}

	if len(req.Items) == 0 {
		return &api.CreateOrderResp{
			Success: false,
			Code:    400,
			Message: "订单商品不能为空",
		}, nil
	}

	if req.Address == "" {
		return &api.CreateOrderResp{
			Success: false,
			Code:    400,
			Message: "收货地址不能为空",
		}, nil
	}

	if req.Phone == "" {
		return &api.CreateOrderResp{
			Success: false,
			Code:    400,
			Message: "联系电话不能为空",
		}, nil
	}

	klog.Infof("参数验证通过: 地址=%s, 电话=%s", req.Address, req.Phone)
	var userInfo *interfaces.UserInfo
	if s.userClient != nil {
		klog.Infof("调用用户服务获取用户 %d 信息", req.UserId)
		userInfo, _ = s.userClient.GetUserInfo(ctx, req.UserId)
	}

	// 如果用户信息为空，创建默认用户
	if userInfo == nil {
		klog.Infof("使用默认用户信息")
		userInfo = &interfaces.UserInfo{
			ID:     req.UserId,
			Name:   fmt.Sprintf("用户%d", req.UserId),
			Email:  fmt.Sprintf("user%d@example.com", req.UserId),
			Phone:  "13800138000",
			Status: 1,
		}
	}
	var totalAmount float64
	var orderItems []*model.OrderItem

	klog.Infof("开始处理 %d 个商品", len(req.Items))

	for i, item := range req.Items {
		if item == nil {
			klog.Errorf(" 商品项 %d 为nil", i)
			continue
		}

		var productInfo *interfaces.ProductInfo

		// 获取商品信息
		if s.productClient != nil {
			productInfo, _ = s.productClient.GetProductInfo(ctx, item.ProductId)
		}

		//如果商品信息为空，使用请求数据
		if productInfo == nil {
			productInfo = &interfaces.ProductInfo{
				ID:       item.ProductId,
				Name:     item.ProductName,
				Price:    item.Price,
				Stock:    100,
				Status:   1,
				Category: "默认分类",
				Avatar:   "https://example.com/product.jpg",
			}
		}

		//计算商品总价
		itemTotal := productInfo.Price * float64(item.Quantity)
		totalAmount += itemTotal

		//创建订单项模型
		orderItem := &model.OrderItem{
			ProductID:    item.ProductId,
			ProductName:  productInfo.Name,
			Quantity:     item.Quantity,
			Price:        productInfo.Price,
			ProductImage: productInfo.Avatar,
		}

		if orderItem.ProductImage == "" {
			orderItem.ProductImage = "https://example.com/default.jpg"
		}

		orderItems = append(orderItems, orderItem)
	}

	if len(orderItems) == 0 {
		klog.Error("没有有效的商品项")
		return &api.CreateOrderResp{
			Success: false,
			Code:    400,
			Message: "没有有效的商品项",
		}, nil
	}

	klog.Infof("所有商品处理完成，总金额: %.2f", totalAmount)
	orderNo := s.generateOrderNo()
	klog.Infof("生成订单号: %s", orderNo)
	klog.Info("开始数据库事务")

	//使用一个独立的数据库会话
	db := s.db.WithContext(ctx)

	//开始事务
	tx := db.Begin()
	if tx.Error != nil {
		klog.Errorf("开启事务失败: %v", tx.Error)
		return &api.CreateOrderResp{
			Success: false,
			Code:    500,
			Message: "开启事务失败",
		}, nil
	}

	//确保事务最终被处理
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			klog.Errorf("事务发生panic，已回滚: %v", r)
		}
	}()
	now := time.Now()
	receiver := getReceiver(req.Receiver)

	order := &model.Order{
		OrderNo:     orderNo,
		UserID:      req.UserId,
		TotalAmount: totalAmount,
		Status:      model.OrderStatusPending,
		Address:     req.Address,
		Phone:       req.Phone,
		Receiver:    receiver,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	klog.Infof("创建订单记录")

	// 使用事务保存订单
	if err := tx.Create(order).Error; err != nil {
		tx.Rollback()
		klog.Errorf("创建订单失败: %v", err)
		return &api.CreateOrderResp{
			Success: false,
			Code:    500,
			Message: "创建订单失败",
		}, nil
	}

	klog.Infof("订单创建成功，订单ID: %d", order.ID)
	for _, orderItem := range orderItems {
		orderItem.OrderID = order.ID
		orderItem.OrderNo = orderNo
		orderItem.CreatedAt = now
		orderItem.UpdatedAt = now
	}

	//批量创建订单项
	if err := tx.CreateInBatches(orderItems, len(orderItems)).Error; err != nil {
		tx.Rollback()
		klog.Errorf("创建订单项失败: %v", err)
		return &api.CreateOrderResp{
			Success: false,
			Code:    500,
			Message: "创建订单项失败",
		}, nil
	}

	klog.Infof("订单项创建成功，数量: %d", len(orderItems))
	if err := tx.Commit().Error; err != nil {
		klog.Errorf("提交事务失败: %v", err)
		return &api.CreateOrderResp{
			Success: false,
			Code:    500,
			Message: "提交事务失败",
		}, nil
	}

	klog.Info("事务提交成功")

	//预占库存（异步）
	if s.productClient != nil {
		go func() {
			defer func() {
				if r := recover(); r != nil {
					klog.Errorf("预占库存任务panic: %v", r)
				}
			}()

			for _, orderItem := range orderItems {
				reserveReq := &api.ReserveStockReq{
					OrderNo:       orderNo,
					ProductId:     orderItem.ProductID,
					Quantity:      orderItem.Quantity,
					ExpireSeconds: 900,
				}

				_, err := s.ReserveStock(ctx, reserveReq)
				if err != nil {
					klog.Warnf("预占库存失败: %v", err)
				}
			}
		}()
	}

	//创建超时任务（异步）
	go func() {
		defer func() {
			if r := recover(); r != nil {
				klog.Errorf("创建超时任务panic: %v", r)
			}
		}()

		if s.daoFactory.TimeoutTaskRepo != nil {
			timeoutTask := &model.TimeoutTask{
				TaskID:     s.generateTaskID(),
				OrderNo:    orderNo,
				Type:       model.TimeoutTypeOrderUnpaid,
				Status:     model.TaskStatusPending,
				ExpireTime: time.Now().Add(30 * time.Minute),
				RetryCount: 0,
				CreatedAt:  now,
				UpdatedAt:  now,
			}

			if err := s.daoFactory.TimeoutTaskRepo.Create(context.Background(), timeoutTask); err != nil {
				klog.Warnf("创建超时任务失败: %v", err)
			}
		}
	}()
	paymentUrl := s.generatePaymentUrl(orderNo)
	klog.Infof("订单创建完成! 订单号: %s, 总金额: %.2f", orderNo, totalAmount)

	return &api.CreateOrderResp{
		Success:     true,
		Code:        0,
		Message:     "订单创建成功",
		OrderNo:     orderNo,
		TotalAmount: totalAmount,
		PaymentUrl:  &paymentUrl,
	}, nil
}

// GetOrder 获取订单详情
func (s *OrderService) GetOrder(ctx context.Context, req *api.GetOrderReq) (*api.GetOrderResp, error) {
	//查询订单
	orderRepo := s.daoFactory.OrderRepo
	order, err := orderRepo.FindByOrderNo(ctx, req.OrderNo)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &api.GetOrderResp{
				Success: false,
				Code:    404,
				Message: "订单不存在",
			}, nil
		}
		return &api.GetOrderResp{
			Success: false,
			Code:    500,
			Message: fmt.Sprintf("查询订单失败: %v", err),
		}, nil
	}

	//权限检查
	if req.UserId != nil && *req.UserId > 0 {
		if order.UserID != *req.UserId {
			return &api.GetOrderResp{
				Success: false,
				Code:    403,
				Message: "无权查看此订单",
			}, nil
		}
	}

	//查询订单项
	orderItemRepo := s.daoFactory.OrderItemRepo
	items, err := orderItemRepo.FindByOrderNo(ctx, req.OrderNo)
	if err != nil {
		return &api.GetOrderResp{
			Success: false,
			Code:    500,
			Message: fmt.Sprintf("查询订单项失败: %v", err),
		}, nil
	}

	//转换为API格式
	apiOrder := s.convertToAPIOrder(order, items, nil)

	return &api.GetOrderResp{
		Success: true,
		Code:    0,
		Message: "获取订单成功",
		Order:   apiOrder,
	}, nil
}

// PayOrder 支付订单
func (s *OrderService) PayOrder(ctx context.Context, req *api.PayOrderReq) (*api.PayOrderResp, error) {
	klog.Infof("PayOrder 开始执行，订单号: %s，用户ID: %d", req.OrderNo, req.UserId)

	//查询订单
	klog.Infof("查询订单: %s", req.OrderNo)
	orderRepo := s.daoFactory.OrderRepo
	order, err := orderRepo.FindByOrderNo(ctx, req.OrderNo)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			klog.Errorf("订单不存在: %s", req.OrderNo)
			return &api.PayOrderResp{
				Success: false,
				Code:    404,
				Message: "订单不存在",
			}, nil
		}
		klog.Errorf("查询订单失败: %v", err)
		return &api.PayOrderResp{
			Success: false,
			Code:    500,
			Message: "查询订单失败",
		}, nil
	}

	klog.Infof("找到订单，状态: %s", order.Status)

	//权限检查
	if order.UserID != req.UserId {
		klog.Errorf("权限错误: 订单用户 %d，请求用户 %d", order.UserID, req.UserId)
		return &api.PayOrderResp{
			Success: false,
			Code:    403,
			Message: "无权操作此订单",
		}, nil
	}

	//检查订单状态
	if order.Status != model.OrderStatusPending {
		klog.Warnf("订单状态异常: %s", order.Status)
		return &api.PayOrderResp{
			Success: false,
			Code:    400,
			Message: fmt.Sprintf("订单状态不正确，当前状态: %s", order.Status),
		}, nil
	}

	klog.Infof("订单状态检查通过")

	//处理支付单号
	paymentNo := ""
	if req.PaymentNo != nil && *req.PaymentNo != "" {
		paymentNo = *req.PaymentNo
	} else {
		paymentNo = fmt.Sprintf("PAY%d", time.Now().Unix())
	}

	klog.Infof("支付单号: %s", paymentNo)

	//更新订单状态为已支付
	klog.Info("更新订单状态...")
	now := time.Now()

	err = orderRepo.UpdateStatusAndPayment(ctx, req.OrderNo, model.OrderStatusPaid, paymentNo)
	if err != nil {
		klog.Errorf("更新订单状态失败: %v", err)
		return &api.PayOrderResp{
			Success: false,
			Code:    500,
			Message: "更新订单状态失败",
		}, nil
	}

	klog.Info("订单状态更新成功")

	//删除支付超时任务（异步，不阻塞）
	go func() {
		defer func() {
			if r := recover(); r != nil {
				klog.Errorf("删除超时任务panic: %v", r)
			}
		}()

		if s.daoFactory.TimeoutTaskRepo != nil {
			tasks, _ := s.daoFactory.TimeoutTaskRepo.FindExpiredTasks(context.Background(), model.TimeoutTypeOrderUnpaid, 10)
			for _, task := range tasks {
				if task.OrderNo == req.OrderNo {
					s.daoFactory.TimeoutTaskRepo.Delete(context.Background(), task.TaskID)
					klog.Infof("删除超时任务: %s", task.TaskID)
				}
			}
		}
	}()

	//返回结果
	paidAt := now.Unix()
	klog.Infof("支付成功! 订单号: %s", req.OrderNo)

	return &api.PayOrderResp{
		Success:    true,
		Code:       0,
		Message:    "支付成功",
		NewStatus_: api.OrderStatus_PAID,
		PaidAt:     &paidAt,
	}, nil
}

// CancelOrder 取消订单
func (s *OrderService) CancelOrder(ctx context.Context, req *api.CancelOrderReq) (*api.CancelOrderResp, error) {
	//查询订单
	orderRepo := s.daoFactory.OrderRepo
	order, err := orderRepo.FindByOrderNo(ctx, req.OrderNo)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &api.CancelOrderResp{
				Success: false,
				Code:    404,
				Message: "订单不存在",
			}, nil
		}
		return &api.CancelOrderResp{
			Success: false,
			Code:    500,
			Message: fmt.Sprintf("查询订单失败: %v", err),
		}, nil
	}

	//权限检查
	if order.UserID != req.UserId {
		return &api.CancelOrderResp{
			Success: false,
			Code:    403,
			Message: "无权操作此订单",
		}, nil
	}

	//检查订单状态（只有待支付和已支付的订单可以取消）
	if order.Status != model.OrderStatusPending && order.Status != model.OrderStatusPaid {
		return &api.CancelOrderResp{
			Success: false,
			Code:    400,
			Message: "订单状态不正确，无法取消",
		}, nil
	}

	//开启数据库事务
	tx := s.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return &api.CancelOrderResp{
			Success: false,
			Code:    500,
			Message: fmt.Sprintf("开启事务失败: %v", tx.Error),
		}, nil
	}

	//更新订单状态为已取消
	err = orderRepo.CancelOrder(ctx, req.OrderNo, req.Reason)
	if err != nil {
		tx.Rollback()
		return &api.CancelOrderResp{
			Success: false,
			Code:    500,
			Message: fmt.Sprintf("取消订单失败: %v", err),
		}, nil
	}

	//释放库存
	stockReservationRepo := s.daoFactory.StockReservationRepo
	reservations, err := stockReservationRepo.FindByOrderNo(ctx, req.OrderNo)
	if err == nil && len(reservations) > 0 {
		for _, reservation := range reservations {
			if reservation.Status == model.StockStatusReserved {
				releaseReq := &api.ReleaseStockReq{
					ReserveId: reservation.ReserveID,
					Reason:    "订单取消: " + req.Reason,
				}

				_, err := s.ReleaseStock(ctx, releaseReq)
				if err != nil {
					klog.Errorf("释放库存失败: %v", err)
				}

				// 更新预占状态为已释放
				stockReservationRepo.UpdateStatus(ctx, reservation.ReserveID, model.StockStatusReleased)
			}
		}
	}

	//删除超时任务
	timeoutTaskRepo := s.daoFactory.TimeoutTaskRepo
	tasks, _ := timeoutTaskRepo.FindExpiredTasks(ctx, "", 10)
	for _, task := range tasks {
		if task.OrderNo == req.OrderNo {
			timeoutTaskRepo.Delete(ctx, task.TaskID)
		}
	}

	//提交事务
	if err := tx.Commit().Error; err != nil {
		return &api.CancelOrderResp{
			Success: false,
			Code:    500,
			Message: fmt.Sprintf("提交事务失败: %v", err),
		}, nil
	}

	//返回结果
	now := time.Now().Unix()
	return &api.CancelOrderResp{
		Success:     true,
		Code:        0,
		Message:     "订单取消成功",
		NewStatus_:  api.OrderStatus_CANCELLED,
		CancelledAt: &now,
	}, nil
}

// ListOrders 查询订单列表
func (s *OrderService) ListOrders(ctx context.Context, req *api.ListOrdersReq) (*api.ListOrdersResp, error) {
	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	var orders []*model.Order
	var total int64
	var err error

	orderRepo := s.daoFactory.OrderRepo

	// 构建查询条件
	condition := make(map[string]interface{})
	condition["user_id"] = req.UserId

	if req.Status != nil {
		condition["status"] = s.convertFromAPIOrderStatus(*req.Status)
	}

	if req.StartTime != nil && *req.StartTime > 0 {
		startTime := time.Unix(*req.StartTime, 0)
		condition["created_at"] = []interface{}{">=", startTime}
	}

	if req.EndTime != nil && *req.EndTime > 0 {
		endTime := time.Unix(*req.EndTime, 0)
		condition["created_at"] = []interface{}{"<=", endTime}
	}

	orders, total, err = orderRepo.ListByCondition(ctx, condition, int(req.Page), int(req.PageSize))
	if err != nil {
		return &api.ListOrdersResp{
			Success: false,
			Code:    500,
			Message: fmt.Sprintf("查询订单列表失败: %v", err),
		}, nil
	}

	// 转换为API格式
	var apiOrders []*api.Order
	orderItemRepo := s.daoFactory.OrderItemRepo

	for _, order := range orders {
		items, _ := orderItemRepo.FindByOrderID(ctx, order.ID)
		apiOrders = append(apiOrders, s.convertToAPIOrder(order, items, nil))
	}

	return &api.ListOrdersResp{
		Success:  true,
		Code:     0,
		Message:  "查询成功",
		Total:    int32(total),
		Page:     req.Page,
		PageSize: req.PageSize,
		Orders:   apiOrders,
	}, nil
}

// ApplyRefund 申请退款
func (s *OrderService) ApplyRefund(ctx context.Context, req *api.ApplyRefundReq) (*api.ApplyRefundResp, error) {
	//参数验证
	if req.Reason == "" {
		return &api.ApplyRefundResp{
			Success: false,
			Code:    400,
			Message: "退款原因不能为空",
		}, nil
	}
	if req.Amount != nil && *req.Amount <= 0 {
		return &api.ApplyRefundResp{
			Success: false,
			Code:    400,
			Message: "退款金额无效",
		}, nil
	}

	//查询订单
	orderRepo := s.daoFactory.OrderRepo
	order, err := orderRepo.FindByOrderNo(ctx, req.OrderNo)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &api.ApplyRefundResp{
				Success: false,
				Code:    404,
				Message: "订单不存在",
			}, nil
		}
		return &api.ApplyRefundResp{
			Success: false,
			Code:    500,
			Message: fmt.Sprintf("查询订单失败: %v", err),
		}, nil
	}

	//权限检查
	if order.UserID != req.UserId {
		return &api.ApplyRefundResp{
			Success: false,
			Code:    403,
			Message: "无权操作此订单",
		}, nil
	}

	//检查订单状态（只有已支付、已发货、已完成的订单可以退款）
	if order.Status != model.OrderStatusPaid &&
		order.Status != model.OrderStatusShipped &&
		order.Status != model.OrderStatusCompleted {
		return &api.ApplyRefundResp{
			Success: false,
			Code:    400,
			Message: "订单状态不正确，无法退款",
		}, nil
	}

	//检查是否已申请退款
	refundRepo := s.daoFactory.RefundRepo
	existingRefund, err := refundRepo.FindByOrderNo(ctx, req.OrderNo)
	if err == nil && existingRefund != nil {
		return &api.ApplyRefundResp{
			Success: false,
			Code:    400,
			Message: "该订单已申请退款",
		}, nil
	}

	//计算退款金额
	refundAmount := order.TotalAmount
	if req.Amount != nil && *req.Amount > 0 && *req.Amount < order.TotalAmount {
		refundAmount = *req.Amount
	}

	//创建退款单
	now := time.Now()
	refundNo := s.generateRefundNo()
	refund := &model.RefundOrder{
		RefundNo:  refundNo,
		OrderNo:   req.OrderNo,
		UserID:    req.UserId,
		Amount:    refundAmount,
		Status:    model.RefundStatusPending,
		Reason:    req.Reason,
		CreatedAt: now,
		UpdatedAt: now,
	}

	err = refundRepo.Create(ctx, refund)
	if err != nil {
		return &api.ApplyRefundResp{
			Success: false,
			Code:    500,
			Message: fmt.Sprintf("创建退款单失败: %v", err),
		}, nil
	}

	//更新订单状态为退款中
	if order.Status != model.OrderStatusRefunded {
		orderRepo.UpdateStatus(ctx, req.OrderNo, model.OrderStatusRefunded)
	}

	//返回结果
	return &api.ApplyRefundResp{
		Success:  true,
		Code:     0,
		Message:  "退款申请提交成功",
		RefundNo: refundNo,
		Status:   api.RefundStatus_PENDING,
	}, nil
}

// ProcessRefund 处理退款
func (s *OrderService) ProcessRefund(ctx context.Context, req *api.ProcessRefundReq) (*api.ProcessRefundResp, error) {
	//查询退款单
	refundRepo := s.daoFactory.RefundRepo
	refund, err := refundRepo.FindByRefundNo(ctx, req.RefundNo)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &api.ProcessRefundResp{
				Success: false,
				Code:    404,
				Message: "退款单不存在",
			}, nil
		}
		return &api.ProcessRefundResp{
			Success: false,
			Code:    500,
			Message: fmt.Sprintf("查询退款单失败: %v", err),
		}, nil
	}

	//检查退款单状态
	if refund.Status != model.RefundStatusPending {
		return &api.ProcessRefundResp{
			Success: false,
			Code:    400,
			Message: fmt.Sprintf("退款单状态必须为待处理，当前状态: %s", refund.Status),
		}, nil
	}

	//验证处理动作
	var newStatus string
	var apiStatus api.RefundStatus
	switch req.Action {
	case api.RefundStatus_APPROVED:
		newStatus = model.RefundStatusApproved
		apiStatus = api.RefundStatus_APPROVED
	case api.RefundStatus_REJECTED:
		newStatus = model.RefundStatusRejected
		apiStatus = api.RefundStatus_REJECTED
	default:
		return &api.ProcessRefundResp{
			Success: false,
			Code:    400,
			Message: "无效的退款状态",
		}, nil
	}

	//开启数据库事务
	tx := s.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return &api.ProcessRefundResp{
			Success: false,
			Code:    500,
			Message: fmt.Sprintf("开启事务失败: %v", tx.Error),
		}, nil
	}

	//更新退款单状态
	now := time.Now()
	refund.Status = newStatus
	refund.Processor = fmt.Sprintf("管理员-%d", req.ProcessorId)
	refund.ProcessedAt = &now
	refund.UpdatedAt = now

	if req.Remark != nil && *req.Remark != "" {
		refund.Reason = refund.Reason + " | 处理备注: " + *req.Remark
	}

	err = refundRepo.Update(ctx, refund)
	if err != nil {
		tx.Rollback()
		return &api.ProcessRefundResp{
			Success: false,
			Code:    500,
			Message: fmt.Sprintf("更新退款单失败: %v", err),
		}, nil
	}

	//如果是拒绝退款，恢复订单状态
	if req.Action == api.RefundStatus_REJECTED {
		orderRepo := s.daoFactory.OrderRepo
		orderRepo.UpdateStatus(ctx, refund.OrderNo, model.OrderStatusCompleted)
	}

	//提交事务
	if err := tx.Commit().Error; err != nil {
		return &api.ProcessRefundResp{
			Success: false,
			Code:    500,
			Message: fmt.Sprintf("提交事务失败: %v", err),
		}, nil
	}

	//返回结果
	return &api.ProcessRefundResp{
		Success:    true,
		Code:       0,
		Message:    "退款处理成功",
		NewStatus_: apiStatus,
	}, nil
}

// ReserveStock 预占库存
func (s *OrderService) ReserveStock(ctx context.Context, req *api.ReserveStockReq) (*api.ReserveStockResp, error) {
	//参数验证
	if req.OrderNo == "" {
		return &api.ReserveStockResp{
			Success: false,
			Code:    400,
			Message: "订单号不能为空",
		}, nil
	}
	if req.ProductId <= 0 {
		return &api.ReserveStockResp{
			Success: false,
			Code:    400,
			Message: "商品ID无效",
		}, nil
	}
	if req.Quantity <= 0 {
		return &api.ReserveStockResp{
			Success: false,
			Code:    400,
			Message: "预占数量必须大于0",
		}, nil
	}
	if req.ExpireSeconds <= 0 {
		req.ExpireSeconds = 900 // 默认15分钟
	}

	//检查商品库存
	stockOk, err := s.productClient.CheckStock(ctx, req.ProductId, req.Quantity)
	if err != nil || !stockOk {
		return &api.ReserveStockResp{
			Success: false,
			Code:    400,
			Message: "商品库存不足",
		}, nil
	}

	//创建预占记录
	now := time.Now()
	expireTime := now.Add(time.Duration(req.ExpireSeconds) * time.Second)
	reserveId := s.generateReserveID()

	reservation := &model.StockReservation{
		ReserveID:  reserveId,
		OrderNo:    req.OrderNo,
		ProductID:  req.ProductId,
		Quantity:   req.Quantity,
		Status:     model.StockStatusReserved,
		ExpireTime: expireTime,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	stockReservationRepo := s.daoFactory.StockReservationRepo
	err = stockReservationRepo.Create(ctx, reservation)
	if err != nil {
		return &api.ReserveStockResp{
			Success: false,
			Code:    500,
			Message: fmt.Sprintf("保存预占记录失败: %v", err),
		}, nil
	}

	//创建超时任务
	timeoutTask := &model.TimeoutTask{
		TaskID:     s.generateTaskID(),
		OrderNo:    req.OrderNo,
		Type:       model.TimeoutTypeStockReservation,
		Status:     model.TaskStatusPending,
		ExpireTime: expireTime,
		RetryCount: 0,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	timeoutTaskRepo := s.daoFactory.TimeoutTaskRepo
	err = timeoutTaskRepo.Create(ctx, timeoutTask)
	if err != nil {
		klog.Warnf("创建库存预占超时任务失败: %v", err)
	}

	return &api.ReserveStockResp{
		Success:   true,
		Code:      0,
		Message:   "库存预占成功",
		ReserveId: reserveId,
	}, nil
}

// ReleaseStock 释放库存
func (s *OrderService) ReleaseStock(ctx context.Context, req *api.ReleaseStockReq) (*api.ReleaseStockResp, error) {
	//查询预占记录
	stockReservationRepo := s.daoFactory.StockReservationRepo
	reservation, err := stockReservationRepo.FindByReserveID(ctx, req.ReserveId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &api.ReleaseStockResp{
				Success: false,
				Code:    404,
				Message: "库存预占记录不存在",
			}, nil
		}
		return &api.ReleaseStockResp{
			Success: false,
			Code:    500,
			Message: fmt.Sprintf("查询预占记录失败: %v", err),
		}, nil
	}

	//检查预占状态
	if reservation.Status != model.StockStatusReserved {
		return &api.ReleaseStockResp{
			Success: false,
			Code:    400,
			Message: "库存预占状态不正确",
		}, nil
	}

	//更新预占状态为已释放
	err = stockReservationRepo.UpdateStatus(ctx, req.ReserveId, model.StockStatusReleased)
	if err != nil {
		return &api.ReleaseStockResp{
			Success: false,
			Code:    500,
			Message: fmt.Sprintf("更新预占记录状态失败: %v", err),
		}, nil
	}

	//删除关联的超时任务
	timeoutTaskRepo := s.daoFactory.TimeoutTaskRepo
	tasks, _ := timeoutTaskRepo.FindExpiredTasks(ctx, model.TimeoutTypeStockReservation, 10)
	for _, task := range tasks {
		if task.OrderNo == reservation.OrderNo {
			timeoutTaskRepo.Delete(ctx, task.TaskID)
		}
	}

	return &api.ReleaseStockResp{
		Success: true,
		Code:    0,
		Message: "库存释放成功",
	}, nil
}

// ConfirmStock 确认库存
func (s *OrderService) ConfirmStock(ctx context.Context, req *api.ConfirmStockReq) (*api.ConfirmStockResp, error) {
	//查询预占记录
	stockReservationRepo := s.daoFactory.StockReservationRepo
	reservation, err := stockReservationRepo.FindByReserveID(ctx, req.ReserveId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &api.ConfirmStockResp{
				Success: false,
				Code:    404,
				Message: "库存预占记录不存在",
			}, nil
		}
		return &api.ConfirmStockResp{
			Success: false,
			Code:    500,
			Message: fmt.Sprintf("查询预占记录失败: %v", err),
		}, nil
	}

	//验证订单号
	if reservation.OrderNo != req.OrderNo {
		return &api.ConfirmStockResp{
			Success: false,
			Code:    400,
			Message: "预占记录与订单号不匹配",
		}, nil
	}

	//检查预占状态
	if reservation.Status != model.StockStatusReserved {
		return &api.ConfirmStockResp{
			Success: false,
			Code:    400,
			Message: "库存预占状态不正确",
		}, nil
	}

	//检查是否过期
	if time.Now().After(reservation.ExpireTime) {
		// 自动释放过期的预占
		stockReservationRepo.UpdateStatus(ctx, req.ReserveId, model.StockStatusExpired)
		return &api.ConfirmStockResp{
			Success: false,
			Code:    400,
			Message: "库存预占已过期",
		}, nil
	}

	//更新预占记录状态为已确认
	err = stockReservationRepo.UpdateStatus(ctx, req.ReserveId, model.StockStatusConfirmed)
	if err != nil {
		return &api.ConfirmStockResp{
			Success: false,
			Code:    500,
			Message: fmt.Sprintf("更新预占记录状态失败: %v", err),
		}, nil
	}

	return &api.ConfirmStockResp{
		Success: true,
		Code:    0,
		Message: "库存确认成功",
	}, nil
}

// ProcessTimeout 处理超时
func (s *OrderService) ProcessTimeout(ctx context.Context, req *api.ProcessTimeoutReq) (*api.ProcessTimeoutResp, error) {
	//查询超时任务
	timeoutTaskRepo := s.daoFactory.TimeoutTaskRepo
	task, err := timeoutTaskRepo.FindByTaskID(ctx, req.TaskId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &api.ProcessTimeoutResp{
				Success: false,
				Code:    404,
				Message: "超时任务不存在",
			}, nil
		}
		return &api.ProcessTimeoutResp{
			Success: false,
			Code:    500,
			Message: fmt.Sprintf("查询超时任务失败: %v", err),
		}, nil
	}

	//检查任务类型
	if s.convertToAPITimeoutType(task.Type) != req.Type {
		return &api.ProcessTimeoutResp{
			Success: false,
			Code:    400,
			Message: "任务类型不匹配",
		}, nil
	}

	//检查任务状态
	if task.Status != model.TaskStatusPending {
		return &api.ProcessTimeoutResp{
			Success: false,
			Code:    400,
			Message: "超时任务已处理",
		}, nil
	}

	//处理任务
	results := make(map[string]string)

	switch task.Type {
	case model.TimeoutTypeOrderUnpaid:
		// 处理未支付订单超时
		err = s.processOrderUnpaidTimeout(ctx, task)
		results["action"] = "cancel_order"
		results["order_no"] = task.OrderNo

	case model.TimeoutTypeStockReservation:
		// 处理库存预占超时
		err = s.processStockReservationTimeout(ctx, task)
		results["action"] = "release_stock"
		results["order_no"] = task.OrderNo

	default:
		err = fmt.Errorf("不支持的任务类型: %s", task.Type)
	}

	//更新任务状态
	if err == nil {
		err = timeoutTaskRepo.UpdateStatus(ctx, req.TaskId, model.TaskStatusCompleted)
		results["result"] = "success"
	} else {
		// 处理失败，增加重试次数
		timeoutTaskRepo.IncrementRetryCount(ctx, req.TaskId)
		timeoutTaskRepo.UpdateStatus(ctx, req.TaskId, model.TaskStatusFailed)
		results["result"] = "failed"
		results["error"] = err.Error()
	}

	if err != nil {
		return &api.ProcessTimeoutResp{
			Success: false,
			Code:    500,
			Message: fmt.Sprintf("处理超时任务失败: %v", err),
		}, nil
	}

	return &api.ProcessTimeoutResp{
		Success: true,
		Code:    0,
		Message: "超时任务处理完成",
		Results: results,
	}, nil
}

// GetOrderStats 获取订单统计
func (s *OrderService) GetOrderStats(ctx context.Context, req *api.OrderStatsReq) (*api.OrderStatsResp, error) {
	//构建查询条件
	condition := make(map[string]interface{})
	condition["user_id"] = req.UserId

	if req.StartTime != nil && *req.StartTime > 0 {
		startTime := time.Unix(*req.StartTime, 0)
		condition["created_at"] = []interface{}{">=", startTime}
	}

	if req.EndTime != nil && *req.EndTime > 0 {
		endTime := time.Unix(*req.EndTime, 0)
		condition["created_at"] = []interface{}{"<=", endTime}
	}

	//获取总订单数
	orderRepo := s.daoFactory.OrderRepo
	totalOrders, err := orderRepo.CountByStatus(ctx, req.UserId, "")
	if err != nil {
		return &api.OrderStatsResp{
			Success: false,
			Code:    500,
			Message: fmt.Sprintf("统计订单总数失败: %v", err),
		}, nil
	}

	//获取订单总金额
	totalAmount, err := orderRepo.SumAmountByCondition(ctx, condition)
	if err != nil {
		return &api.OrderStatsResp{
			Success: false,
			Code:    500,
			Message: fmt.Sprintf("统计订单总金额失败: %v", err),
		}, nil
	}

	//获取各状态订单数
	statusCounts := make(map[string]int32)
	statuses := []string{
		model.OrderStatusPending,
		model.OrderStatusPaid,
		model.OrderStatusShipped,
		model.OrderStatusCompleted,
		model.OrderStatusCancelled,
		model.OrderStatusRefunded,
	}

	for _, status := range statuses {
		count, err := orderRepo.CountByStatus(ctx, req.UserId, status)
		if err != nil {
			klog.Errorf("统计%s状态订单数失败: %v", status, err)
			continue
		}
		statusCounts[status] = int32(count)
	}

	//构建响应
	return &api.OrderStatsResp{
		Success:      true,
		Code:         0,
		Message:      "获取统计成功",
		TotalOrders:  int32(totalOrders),
		TotalAmount:  totalAmount,
		StatusCounts: statusCounts,
	}, nil
}

// ShipOrder 发货订单
func (s *OrderService) ShipOrder(ctx context.Context, req *api.PayOrderReq) (*api.PayOrderResp, error) {
	//查询订单
	orderRepo := s.daoFactory.OrderRepo
	order, err := orderRepo.FindByOrderNo(ctx, req.OrderNo)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &api.PayOrderResp{
				Success: false,
				Code:    404,
				Message: "订单不存在",
			}, nil
		}
		return &api.PayOrderResp{
			Success: false,
			Code:    500,
			Message: fmt.Sprintf("查询订单失败: %v", err),
		}, nil
	}

	//检查订单状态
	if order.Status != model.OrderStatusPaid {
		return &api.PayOrderResp{
			Success: false,
			Code:    400,
			Message: fmt.Sprintf("订单状态必须为已支付才能发货，当前状态: %s", order.Status),
		}, nil
	}

	//检查物流信息
	shippingNo := ""
	if req.PaymentNo != nil {
		shippingNo = *req.PaymentNo
	}
	if shippingNo == "" {
		return &api.PayOrderResp{
			Success: false,
			Code:    400,
			Message: "物流单号不能为空",
		}, nil
	}

	//更新订单状态为已发货
	err = orderRepo.UpdateShippingInfo(ctx, req.OrderNo, shippingNo)
	if err != nil {
		return &api.PayOrderResp{
			Success: false,
			Code:    500,
			Message: fmt.Sprintf("更新发货信息失败: %v", err),
		}, nil
	}

	//返回结果
	return &api.PayOrderResp{
		Success:    true,
		Code:       0,
		Message:    "发货成功",
		NewStatus_: api.OrderStatus_SHIPPED,
	}, nil
}

// ConfirmReceipt 确认收货
func (s *OrderService) ConfirmReceipt(ctx context.Context, req *api.PayOrderReq) (*api.PayOrderResp, error) {
	//查询订单
	orderRepo := s.daoFactory.OrderRepo
	order, err := orderRepo.FindByOrderNo(ctx, req.OrderNo)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &api.PayOrderResp{
				Success: false,
				Code:    404,
				Message: "订单不存在",
			}, nil
		}
		return &api.PayOrderResp{
			Success: false,
			Code:    500,
			Message: fmt.Sprintf("查询订单失败: %v", err),
		}, nil
	}

	//权限检查
	if order.UserID != req.UserId {
		return &api.PayOrderResp{
			Success: false,
			Code:    403,
			Message: "无权操作此订单",
		}, nil
	}

	//检查订单状态
	if order.Status != model.OrderStatusShipped {
		return &api.PayOrderResp{
			Success: false,
			Code:    400,
			Message: fmt.Sprintf("订单状态必须为已发货才能确认收货，当前状态: %s", order.Status),
		}, nil
	}

	//更新订单状态为已完成
	err = orderRepo.UpdateStatus(ctx, req.OrderNo, model.OrderStatusCompleted)
	if err != nil {
		return &api.PayOrderResp{
			Success: false,
			Code:    500,
			Message: fmt.Sprintf("确认收货失败: %v", err),
		}, nil
	}

	//返回结果
	return &api.PayOrderResp{
		Success:    true,
		Code:       0,
		Message:    "确认收货成功",
		NewStatus_: api.OrderStatus_COMPLETED,
	}, nil
}

// processOrderUnpaidTimeout 处理订单未支付超时
func (s *OrderService) processOrderUnpaidTimeout(ctx context.Context, task *model.TimeoutTask) error {
	//查询订单
	orderRepo := s.daoFactory.OrderRepo
	order, err := orderRepo.FindByOrderNo(ctx, task.OrderNo)
	if err != nil {
		return fmt.Errorf("查询订单失败: %w", err)
	}

	//检查订单状态（只有待支付订单才需要超时取消）
	if order.Status != model.OrderStatusPending {
		return nil // 订单状态已变更，无需处理
	}

	//取消订单
	_, err = s.CancelOrder(ctx, &api.CancelOrderReq{
		OrderNo: task.OrderNo,
		UserId:  order.UserID,
		Reason:  "支付超时自动取消",
	})
	return err
}

// processStockReservationTimeout 处理库存预占超时
func (s *OrderService) processStockReservationTimeout(ctx context.Context, task *model.TimeoutTask) error {
	//查询订单的库存预占记录
	stockReservationRepo := s.daoFactory.StockReservationRepo
	reservations, err := stockReservationRepo.FindByOrderNo(ctx, task.OrderNo)
	if err != nil {
		return fmt.Errorf("查询库存预占记录失败: %w", err)
	}

	//释放所有未确认的库存预占
	for _, reservation := range reservations {
		if reservation.Status == model.StockStatusReserved && time.Now().After(reservation.ExpireTime) {
			// 更新预占记录状态为已过期
			stockReservationRepo.UpdateStatus(ctx, reservation.ReserveID, model.StockStatusExpired)
		}
	}

	return nil
}

// generateOrderNo 生成订单号
func (s *OrderService) generateOrderNo() string {
	// 格式: ORD + 年月日时分秒 + 4位随机数
	now := time.Now()
	timestamp := now.Format("20060102150405")
	random := rand.Intn(10000)
	return fmt.Sprintf("ORD%s%04d", timestamp, random)
}

// generateRefundNo 生成退款单号
func (s *OrderService) generateRefundNo() string {
	// 格式: REF + 年月日时分秒 + 4位随机数
	now := time.Now()
	timestamp := now.Format("20060102150405")
	random := rand.Intn(10000)
	return fmt.Sprintf("REF%s%04d", timestamp, random)
}

// generateTaskID 生成任务ID
func (s *OrderService) generateTaskID() string {
	// 格式: TASK + 年月日时分秒 + 4位随机数
	now := time.Now()
	timestamp := now.Format("20060102150405")
	random := rand.Intn(10000)
	return fmt.Sprintf("TASK%s%04d", timestamp, random)
}

// generateReserveID 生成预占ID
func (s *OrderService) generateReserveID() string {
	// 格式: RES + 年月日时分秒 + 4位随机数
	now := time.Now()
	timestamp := now.Format("20060102150405")
	random := rand.Intn(10000)
	return fmt.Sprintf("RES%s%04d", timestamp, random)
}

// generatePaymentUrl 生成支付链接
func (s *OrderService) generatePaymentUrl(orderNo string) string {
	return fmt.Sprintf("https://payment.example.com/pay?order_no=%s", orderNo)
}

// getReceiver 获取收货人
func getReceiver(receiver *string) string {
	if receiver != nil && *receiver != "" {
		return *receiver
	}
	return "未填写"
}

// convertToAPIOrder 将model.Order转换为api.Order
func (s *OrderService) convertToAPIOrder(order *model.Order, items []*model.OrderItem, reservations []*model.StockReservation) *api.Order {
	apiOrder := &api.Order{
		Id:          order.ID,
		OrderNo:     order.OrderNo,
		UserId:      order.UserID,
		TotalAmount: order.TotalAmount,
		Status:      s.convertToAPIOrderStatus(order.Status),
		Address:     order.Address,
		Phone:       order.Phone,
		CreatedAt:   order.CreatedAt.Unix(),
		UpdatedAt:   order.UpdatedAt.Unix(),
	}

	// 处理可选字段
	if order.Receiver != "" && order.Receiver != "未填写" {
		receiver := order.Receiver
		apiOrder.Receiver = &receiver
	}
	if order.PaymentNo != "" {
		paymentNo := order.PaymentNo
		apiOrder.PaymentNo = &paymentNo
	}
	if order.ShippingNo != "" {
		shippingNo := order.ShippingNo
		apiOrder.ShippingNo = &shippingNo
	}

	// 转换订单项
	var apiItems []*api.OrderItem
	for _, item := range items {
		apiItem := &api.OrderItem{
			ProductId:   item.ProductID,
			ProductName: item.ProductName,
			Quantity:    item.Quantity,
			Price:       item.Price,
		}
		if item.ProductImage != "" {
			productImage := item.ProductImage
			apiItem.ProductImage = &productImage
		}
		apiItems = append(apiItems, apiItem)
	}
	apiOrder.Items = apiItems

	return apiOrder
}

// convertToAPIRefund 将model.RefundOrder转换为api.RefundOrder
func (s *OrderService) convertToAPIRefund(refund *model.RefundOrder) *api.RefundOrder {
	apiRefund := &api.RefundOrder{
		RefundNo:  refund.RefundNo,
		OrderNo:   refund.OrderNo,
		UserId:    refund.UserID,
		Amount:    refund.Amount,
		Status:    s.convertToAPIRefundStatus(refund.Status),
		Reason:    refund.Reason,
		CreatedAt: refund.CreatedAt.Unix(),
		UpdatedAt: refund.UpdatedAt.Unix(),
	}

	// 处理可选字段
	if refund.Processor != "" {
		processor := refund.Processor
		apiRefund.Processor = &processor
	}
	if refund.ProcessedAt != nil {
		processedAt := refund.ProcessedAt.Unix()
		apiRefund.ProcessedAt = &processedAt
	}

	return apiRefund
}

// convertToAPIOrderStatus 将model状态转换为api状态
func (s *OrderService) convertToAPIOrderStatus(status string) api.OrderStatus {
	switch status {
	case model.OrderStatusPending:
		return api.OrderStatus_PENDING
	case model.OrderStatusPaid:
		return api.OrderStatus_PAID
	case model.OrderStatusShipped:
		return api.OrderStatus_SHIPPED
	case model.OrderStatusCompleted:
		return api.OrderStatus_COMPLETED
	case model.OrderStatusCancelled:
		return api.OrderStatus_CANCELLED
	case model.OrderStatusRefunded:
		return api.OrderStatus_REFUNDED
	default:
		return api.OrderStatus_PENDING
	}
}

// convertFromAPIOrderStatus 将api状态转换为model状态
func (s *OrderService) convertFromAPIOrderStatus(status api.OrderStatus) string {
	switch status {
	case api.OrderStatus_PENDING:
		return model.OrderStatusPending
	case api.OrderStatus_PAID:
		return model.OrderStatusPaid
	case api.OrderStatus_SHIPPED:
		return model.OrderStatusShipped
	case api.OrderStatus_COMPLETED:
		return model.OrderStatusCompleted
	case api.OrderStatus_CANCELLED:
		return model.OrderStatusCancelled
	case api.OrderStatus_REFUNDED:
		return model.OrderStatusRefunded
	default:
		return model.OrderStatusPending
	}
}

// convertToAPIRefundStatus 将model退款状态转换为api状态
func (s *OrderService) convertToAPIRefundStatus(status string) api.RefundStatus {
	switch status {
	case model.RefundStatusPending:
		return api.RefundStatus_PENDING
	case model.RefundStatusApproved:
		return api.RefundStatus_APPROVED
	case model.RefundStatusRejected:
		return api.RefundStatus_REJECTED
	case model.RefundStatusProcessing:
		return api.RefundStatus_PROCESSING
	case model.RefundStatusCompleted:
		return api.RefundStatus_COMPLETED
	default:
		return api.RefundStatus_PENDING
	}
}

// convertToAPITimeoutType 将model超时类型转换为api类型
func (s *OrderService) convertToAPITimeoutType(taskType string) api.TimeoutType {
	switch taskType {
	case model.TimeoutTypeOrderUnpaid:
		return api.TimeoutType_ORDER_UNPAID
	case model.TimeoutTypeStockReservation:
		return api.TimeoutType_STOCK_RESERVATION
	default:
		return api.TimeoutType_ORDER_UNPAID
	}
}
