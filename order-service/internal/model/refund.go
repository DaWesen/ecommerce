package model

import (
	"time"

	"gorm.io/gorm"
)

type RefundOrder struct {
	RefundNo string  `gorm:"size:32;primaryKey;comment:退款单号"`
	OrderNo  string  `gorm:"size:32;index;not null;comment:订单号"`
	UserID   int64   `gorm:"index;not null;comment:用户ID"`
	Amount   float64 `gorm:"type:decimal(10,2);not null;comment:退款金额"`
	Status   string  `gorm:"size:20;index;not null;default:'pending';comment:状态"`
	Reason   string  `gorm:"size:200;comment:退款原因"`

	// Thrift 中的可选字段
	Processor   string     `gorm:"size:50;comment:处理人"`
	ProcessedAt *time.Time `gorm:"comment:处理时间"`

	// 时间字段
	CreatedAt time.Time      `gorm:"index;autoCreateTime"`
	UpdatedAt time.Time      `gorm:"index;autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`

	// 关联关系
	Order Order `gorm:"foreignKey:OrderNo;references:OrderNo"`
}

func (RefundOrder) TableName() string {
	return "refund_orders"
}

// 添加时间戳方法用于 Thrift 转换
func (r *RefundOrder) GetCreatedAtTimestamp() int64 {
	return r.CreatedAt.Unix()
}

func (r *RefundOrder) GetUpdatedAtTimestamp() int64 {
	return r.UpdatedAt.Unix()
}

// Status 常量
const (
	RefundStatusPending    = "pending"    // PENDING = 0
	RefundStatusApproved   = "approved"   // APPROVED = 1
	RefundStatusRejected   = "rejected"   // REJECTED = 2
	RefundStatusProcessing = "processing" // PROCESSING = 3
	RefundStatusCompleted  = "completed"  // COMPLETED = 4
	RefundStatusFailed     = "failed"     // 额外状态
)
