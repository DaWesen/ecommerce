package orderDao

import (
	"context"
	"ecommerce/order-service/internal/dao/interfaces"
	"ecommerce/order-service/internal/model"
	"time"

	"gorm.io/gorm"
)

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) interfaces.IOrderRepository {
	return &OrderRepository{db: db}
}

// 创建订单
func (r *OrderRepository) Create(ctx context.Context, order *model.Order) error {
	return r.db.WithContext(ctx).Create(order).Error
}

// 更新订单
func (r *OrderRepository) Update(ctx context.Context, order *model.Order) error {
	return r.db.WithContext(ctx).Save(order).Error
}

// 硬删除订单
func (r *OrderRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Unscoped().Delete(&model.Order{}, "id = ?", id).Error
}

// 软删除订单
func (r *OrderRepository) SoftDelete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&model.Order{}, "id = ?", id).Error
}

// 根据ID查询订单
func (r *OrderRepository) FindByID(ctx context.Context, id int64) (*model.Order, error) {
	var order model.Order
	err := r.db.WithContext(ctx).Preload("Items").Where("id = ?", id).First(&order).Error
	return &order, err
}

// 根据订单号查询订单
func (r *OrderRepository) FindByOrderNo(ctx context.Context, orderNo string) (*model.Order, error) {
	var order model.Order
	err := r.db.WithContext(ctx).Preload("Items").Where("order_no = ?", orderNo).First(&order).Error
	return &order, err
}

// 根据用户ID查询订单列表
func (r *OrderRepository) FindByUserID(ctx context.Context, userID int64, status string, page, pageSize int) ([]*model.Order, int64, error) {
	var orders []*model.Order
	var total int64

	db := r.db.WithContext(ctx).Model(&model.Order{}).Where("user_id = ?", userID)

	if status != "" {
		db = db.Where("status = ?", status)
	}

	err := db.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err = db.Preload("Items").Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&orders).Error

	return orders, total, err
}

// 根据条件查询订单列表
func (r *OrderRepository) ListByCondition(ctx context.Context, condition map[string]interface{}, page, pageSize int) ([]*model.Order, int64, error) {
	var orders []*model.Order
	var total int64

	db := r.db.WithContext(ctx).Model(&model.Order{})

	for key, value := range condition {
		switch v := value.(type) {
		case []interface{}:
			if len(v) == 2 {
				db = db.Where(key+" "+v[0].(string)+" ?", v[1])
			}
		default:
			db = db.Where(key+" = ?", value)
		}
	}

	err := db.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err = db.Preload("Items").Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&orders).Error

	return orders, total, err
}

// 统计用户各状态订单数量
func (r *OrderRepository) CountByStatus(ctx context.Context, userID int64, status string) (int64, error) {
	var count int64
	db := r.db.WithContext(ctx).Model(&model.Order{}).Where("user_id = ?", userID)

	if status != "" {
		db = db.Where("status = ?", status)
	}

	err := db.Count(&count).Error
	return count, err
}

// SumAmountByCondition 根据条件统计订单金额
func (r *OrderRepository) SumAmountByCondition(ctx context.Context, condition map[string]interface{}) (float64, error) {
	var total float64

	db := r.db.WithContext(ctx).Model(&model.Order{})

	// 应用条件
	for key, value := range condition {
		switch v := value.(type) {
		case []interface{}:
			if len(v) == 2 {
				operator := v[0].(string)
				val := v[1]
				// 对于时间字段，需要特殊处理
				if t, ok := val.(time.Time); ok {
					db = db.Where(key+" "+operator+" ?", t.Format("2006-01-02 15:04:05"))
				} else {
					db = db.Where(key+" "+operator+" ?", val)
				}
			}
		default:
			db = db.Where(key+" = ?", value)
		}
	}

	// 使用Raw SQL来避免复杂的查询构建问题
	// 首先构建查询字符串
	query := db.Select("SUM(total_amount) as total")

	// 执行查询
	rows, err := query.Rows()
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	// 读取结果
	if rows.Next() {
		var totalPtr *float64
		err = rows.Scan(&totalPtr)
		if err != nil {
			return 0, err
		}
		if totalPtr != nil {
			total = *totalPtr
		}
	}

	return total, nil
}

// 更新订单状态
func (r *OrderRepository) UpdateStatus(ctx context.Context, orderNo string, status string) error {
	updates := map[string]interface{}{
		"status":     status,
		"updated_at": time.Now(),
	}

	switch status {
	case model.OrderStatusPaid:
		now := time.Now()
		updates["paid_at"] = &now
	case model.OrderStatusShipped:
		now := time.Now()
		updates["shipped_at"] = &now
	case model.OrderStatusCancelled:
		now := time.Now()
		updates["cancelled_at"] = &now
	case model.OrderStatusCompleted:
		now := time.Now()
		updates["delivered_at"] = &now
	}

	return r.db.WithContext(ctx).Model(&model.Order{}).
		Where("order_no = ?", orderNo).
		Updates(updates).Error
}

// 更新支付信息
func (r *OrderRepository) UpdatePaymentInfo(ctx context.Context, orderNo, paymentNo string) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&model.Order{}).
		Where("order_no = ?", orderNo).
		Updates(map[string]interface{}{
			"payment_no": paymentNo,
			"paid_at":    &now,
			"status":     model.OrderStatusPaid,
			"updated_at": time.Now(),
		}).Error
}

// 更新物流信息
func (r *OrderRepository) UpdateShippingInfo(ctx context.Context, orderNo, shippingNo string) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&model.Order{}).
		Where("order_no = ?", orderNo).
		Updates(map[string]interface{}{
			"shipping_no": shippingNo,
			"shipped_at":  &now,
			"status":      model.OrderStatusShipped,
			"updated_at":  time.Now(),
		}).Error
}

// 更新状态和支付信息
func (r *OrderRepository) UpdateStatusAndPayment(ctx context.Context, orderNo, status, paymentNo string) error {
	updates := map[string]interface{}{
		"status":     status,
		"payment_no": paymentNo,
		"updated_at": time.Now(),
	}

	if status == model.OrderStatusPaid {
		now := time.Now()
		updates["paid_at"] = &now
	}

	return r.db.WithContext(ctx).Model(&model.Order{}).
		Where("order_no = ?", orderNo).
		Updates(updates).Error
}

// 取消订单
func (r *OrderRepository) CancelOrder(ctx context.Context, orderNo string, reason string) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&model.Order{}).
		Where("order_no = ?", orderNo).
		Updates(map[string]interface{}{
			"status":       model.OrderStatusCancelled,
			"cancelled_at": &now,
			"updated_at":   time.Now(),
		}).Error
}
