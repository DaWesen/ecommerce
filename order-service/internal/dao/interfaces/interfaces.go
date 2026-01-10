package interfaces

import (
	"context"
	"ecommerce/order-service/internal/model"
)

// UserInfo 用户服务返回的信息
type UserInfo struct {
	ID        int64
	Name      string
	Email     string
	Phone     string
	Status    int32
	Avatar    string
	LastLogin int64
}

// ProductInfo 商品服务返回的信息
type ProductInfo struct {
	ID       int64
	Name     string
	Price    float64
	Stock    int32
	Status   int32
	Category string
	Avatar   string
	Brand    string
}

type IOrderRepository interface {
	// 基础CRUD
	Create(ctx context.Context, order *model.Order) error
	Update(ctx context.Context, order *model.Order) error
	Delete(ctx context.Context, id int64) error
	FindByID(ctx context.Context, id int64) (*model.Order, error)
	SoftDelete(ctx context.Context, id int64) error

	// 查询方法
	FindByOrderNo(ctx context.Context, orderNo string) (*model.Order, error)
	FindByUserID(ctx context.Context, userID int64, status string, page, pageSize int) ([]*model.Order, int64, error)
	ListByCondition(ctx context.Context, condition map[string]interface{}, page, pageSize int) ([]*model.Order, int64, error)
	CountByStatus(ctx context.Context, userID int64, status string) (int64, error)
	SumAmountByCondition(ctx context.Context, condition map[string]interface{}) (float64, error)

	// 业务方法
	UpdateStatus(ctx context.Context, orderNo string, status string) error
	UpdatePaymentInfo(ctx context.Context, orderNo, paymentNo string) error
	UpdateShippingInfo(ctx context.Context, orderNo, shippingNo string) error
	UpdateStatusAndPayment(ctx context.Context, orderNo, status, paymentNo string) error
	CancelOrder(ctx context.Context, orderNo string, reason string) error
}

// 订单项接口
type IOrderItemRepository interface {
	CreateBatch(ctx context.Context, items []*model.OrderItem) error
	FindByOrderID(ctx context.Context, orderID int64) ([]*model.OrderItem, error)
	FindByOrderNo(ctx context.Context, orderNo string) ([]*model.OrderItem, error)
	DeleteByOrderID(ctx context.Context, orderID int64) error
}

// 退款单接口
type IRefundRepository interface {
	Create(ctx context.Context, refund *model.RefundOrder) error
	Update(ctx context.Context, refund *model.RefundOrder) error
	FindByRefundNo(ctx context.Context, refundNo string) (*model.RefundOrder, error)
	FindByOrderNo(ctx context.Context, orderNo string) (*model.RefundOrder, error)
	ListByUserID(ctx context.Context, userID int64, page, pageSize int) ([]*model.RefundOrder, int64, error)
	ListByCondition(ctx context.Context, condition map[string]interface{}, page, pageSize int) ([]*model.RefundOrder, int64, error)
	UpdateStatus(ctx context.Context, refundNo string, status string) error
}

// 库存预占接口
type IStockReservationRepository interface {
	Create(ctx context.Context, reservation *model.StockReservation) error
	Update(ctx context.Context, reservation *model.StockReservation) error
	FindByReserveID(ctx context.Context, reserveID string) (*model.StockReservation, error)
	FindByOrderNo(ctx context.Context, orderNo string) ([]*model.StockReservation, error)
	DeleteByReserveID(ctx context.Context, reserveID string) error
	FindExpiredReservations(ctx context.Context) ([]*model.StockReservation, error)
	UpdateStatus(ctx context.Context, reserveID, status string) error
}

// 超时任务接口
type ITimeoutTaskRepository interface {
	Create(ctx context.Context, task *model.TimeoutTask) error
	Update(ctx context.Context, task *model.TimeoutTask) error
	FindByTaskID(ctx context.Context, taskID string) (*model.TimeoutTask, error)
	FindExpiredTasks(ctx context.Context, taskType string, limit int) ([]*model.TimeoutTask, error)
	Delete(ctx context.Context, taskID string) error
	UpdateStatus(ctx context.Context, taskID, status string) error
	IncrementRetryCount(ctx context.Context, taskID string) error
}

// 外部服务客户端接口
type IUserClient interface {
	GetUserInfo(ctx context.Context, userID int64) (*UserInfo, error)
	BatchGetUsers(ctx context.Context, userIDs []int64) (map[int64]*UserInfo, error)
	ValidateUser(ctx context.Context, userID int64) (bool, error)
}

type IProductClient interface {
	GetProductInfo(ctx context.Context, productID int64) (*ProductInfo, error)
	BatchGetProducts(ctx context.Context, productIDs []int64) (map[int64]*ProductInfo, error)
	CheckStock(ctx context.Context, productID int64, quantity int32) (bool, error)
}
