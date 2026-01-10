package model

import (
	"time"

	"gorm.io/gorm"
)

type TimeoutTask struct {
	TaskID     string    `gorm:"size:32;primaryKey;comment:任务ID"`
	OrderNo    string    `gorm:"size:32;index;not null;comment:订单号"`
	Type       string    `gorm:"size:50;index;not null;comment:任务类型"`
	Status     string    `gorm:"size:20;index;not null;default:'pending';comment:状态"`
	ExpireTime time.Time `gorm:"index;not null;comment:过期时间"`
	RetryCount int32     `gorm:"not null;default:0;comment:重试次数"`

	//时间字段
	CreatedAt time.Time      `gorm:"index;autoCreateTime"`
	UpdatedAt time.Time      `gorm:"index;autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`

	//关联关系
	Order Order `gorm:"foreignKey:OrderNo;references:OrderNo"`
}

func (TimeoutTask) TableName() string {
	return "timeout_tasks"
}

// 添加时间戳方法用于 Thrift 转换
func (t *TimeoutTask) GetCreatedAtTimestamp() int64 {
	return t.CreatedAt.Unix()
}

func (t *TimeoutTask) GetUpdatedAtTimestamp() int64 {
	return t.UpdatedAt.Unix()
}

func (t *TimeoutTask) GetExpireTimeTimestamp() int64 {
	return t.ExpireTime.Unix()
}

// Type 常量
const (
	TimeoutTypeOrderUnpaid      = "order_unpaid"
	TimeoutTypeStockReservation = "stock_reservation"
)

// Status 常量
const (
	TaskStatusPending    = "pending"    // 待处理
	TaskStatusProcessing = "processing" // 处理中
	TaskStatusCompleted  = "completed"  // 已完成
	TaskStatusFailed     = "failed"     // 已失败
)
