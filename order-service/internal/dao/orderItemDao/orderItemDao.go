package orderItemDao

import (
	"context"
	"ecommerce/order-service/internal/dao/interfaces"
	"ecommerce/order-service/internal/model"

	"gorm.io/gorm"
)

type OrderItemRepository struct {
	db *gorm.DB
}

func NewOrderItemRepository(db *gorm.DB) interfaces.IOrderItemRepository {
	return &OrderItemRepository{db: db}
}

// 批量创建订单项
func (r *OrderItemRepository) CreateBatch(ctx context.Context, items []*model.OrderItem) error {
	if len(items) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).CreateInBatches(items, 100).Error
}

// 根据订单ID查询订单项
func (r *OrderItemRepository) FindByOrderID(ctx context.Context, orderID int64) ([]*model.OrderItem, error) {
	var items []*model.OrderItem
	err := r.db.WithContext(ctx).Where("order_id = ?", orderID).Find(&items).Error
	return items, err
}

// 根据订单号查询订单项
func (r *OrderItemRepository) FindByOrderNo(ctx context.Context, orderNo string) ([]*model.OrderItem, error) {
	var items []*model.OrderItem
	err := r.db.WithContext(ctx).Where("order_no = ?", orderNo).Find(&items).Error
	return items, err
}

// 根据订单ID删除订单项
func (r *OrderItemRepository) DeleteByOrderID(ctx context.Context, orderID int64) error {
	return r.db.WithContext(ctx).Where("order_id = ?", orderID).Delete(&model.OrderItem{}).Error
}
