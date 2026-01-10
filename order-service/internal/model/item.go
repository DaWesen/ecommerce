package model

import (
	"time"

	"gorm.io/gorm"
)

type OrderItem struct {
	ID           int64   `gorm:"primaryKey;autoIncrement"`
	OrderID      int64   `gorm:"index;not null;comment:订单ID"`
	OrderNo      string  `gorm:"size:32;index;not null;comment:订单号"`
	ProductID    int64   `gorm:"index;not null;comment:商品ID"`
	ProductName  string  `gorm:"size:100;not null;comment:商品名称"`
	Quantity     int32   `gorm:"not null;default:1;comment:数量"`
	Price        float64 `gorm:"type:decimal(10,2);not null;comment:单价"`
	ProductImage string  `gorm:"size:500;comment:商品图片"`

	// 时间字段
	CreatedAt time.Time      `gorm:"index;autoCreateTime"`
	UpdatedAt time.Time      `gorm:"index;autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`

	// 关联关系
	Order Order `gorm:"foreignKey:OrderID"`
}

func (OrderItem) TableName() string {
	return "order_items"
}
