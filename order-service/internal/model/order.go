package model

import (
	"time"

	"gorm.io/gorm"
)

type Order struct {
	ID          int64   `gorm:"primaryKey;autoIncrement"`
	OrderNo     string  `gorm:"size:32;uniqueIndex;not null;comment:订单号"`
	UserID      int64   `gorm:"index;not null;comment:用户ID"`
	TotalAmount float64 `gorm:"type:decimal(10,2);not null;comment:总金额"`
	Status      string  `gorm:"size:20;index;not null;default:'pending';comment:状态"`

	Address  string `gorm:"size:200;comment:收货地址"`
	Phone    string `gorm:"size:20;comment:联系电话"`
	Receiver string `gorm:"size:50;comment:收货人姓名"`

	// 业务单号
	PaymentNo  string `gorm:"size:100;index;comment:支付单号"`
	ShippingNo string `gorm:"size:100;index;comment:物流单号"`

	// 时间字段
	PaidAt      *time.Time     `gorm:"comment:支付时间"`
	ShippedAt   *time.Time     `gorm:"comment:发货时间"`
	DeliveredAt *time.Time     `gorm:"comment:送达时间"`
	CancelledAt *time.Time     `gorm:"comment:取消时间"`
	CreatedAt   time.Time      `gorm:"index;autoCreateTime"`
	UpdatedAt   time.Time      `gorm:"index;autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`

	// 关联关系
	Items []OrderItem `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE"`
}

func (Order) TableName() string {
	return "orders"
}

// 添加时间戳方法用于 Thrift 转换
func (o *Order) GetCreatedAtTimestamp() int64 {
	return o.CreatedAt.Unix()
}

func (o *Order) GetUpdatedAtTimestamp() int64 {
	return o.UpdatedAt.Unix()
}

// Status 常量
const (
	OrderStatusPending   = "pending"   // PENDING = 0
	OrderStatusPaid      = "paid"      // PAID = 1
	OrderStatusShipped   = "shipped"   // SHIPPED = 2
	OrderStatusCompleted = "completed" // COMPLETED = 3
	OrderStatusCancelled = "cancelled" // CANCELLED = 4
	OrderStatusRefunded  = "refunded"  // REFUNDED = 5
)
