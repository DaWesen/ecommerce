package model

import (
	"time"

	"gorm.io/gorm"
)

type StockReservation struct {
	ReserveID  string    `gorm:"size:32;primaryKey;comment:预占ID"`
	OrderNo    string    `gorm:"size:32;index;not null;comment:订单号"`
	ProductID  int64     `gorm:"index;not null;comment:商品ID"`
	Quantity   int32     `gorm:"not null;comment:预占数量"`
	Status     string    `gorm:"size:20;index;not null;default:'reserved';comment:状态"`
	ExpireTime time.Time `gorm:"index;comment:过期时间"`

	// 时间字段
	CreatedAt time.Time      `gorm:"index;autoCreateTime"`
	UpdatedAt time.Time      `gorm:"index;autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`

	// 关联关系
	Order Order `gorm:"foreignKey:OrderNo;references:OrderNo"`
}

func (StockReservation) TableName() string {
	return "stock_reservations"
}

// 添加时间戳方法用于 Thrift 转换
func (s *StockReservation) GetCreatedAtTimestamp() int64 {
	return s.CreatedAt.Unix()
}

func (s *StockReservation) GetUpdatedAtTimestamp() int64 {
	return s.UpdatedAt.Unix()
}

func (s *StockReservation) GetExpireTimeTimestamp() int64 {
	return s.ExpireTime.Unix()
}

// Status 常量
const (
	StockStatusReserved  = "reserved"  // 已预占
	StockStatusConfirmed = "confirmed" // 已确认（扣减）
	StockStatusReleased  = "released"  // 已释放
	StockStatusExpired   = "expired"   // 已过期
)
