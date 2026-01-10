package refundDao

import (
	"context"
	"ecommerce/order-service/internal/dao/interfaces"
	"ecommerce/order-service/internal/model"
	"time"

	"gorm.io/gorm"
)

type RefundRepository struct {
	db *gorm.DB
}

func NewRefundRepository(db *gorm.DB) interfaces.IRefundRepository {
	return &RefundRepository{db: db}
}
func (r *RefundRepository) Create(ctx context.Context, refund *model.RefundOrder) error {
	return r.db.WithContext(ctx).Create(refund).Error
}

// 更新退款单
func (r *RefundRepository) Update(ctx context.Context, refund *model.RefundOrder) error {
	return r.db.WithContext(ctx).Save(refund).Error
}

// 根据退款单号查询
func (r *RefundRepository) FindByRefundNo(ctx context.Context, refundNo string) (*model.RefundOrder, error) {
	var refund model.RefundOrder
	err := r.db.WithContext(ctx).Where("refund_no = ?", refundNo).First(&refund).Error
	return &refund, err
}

// 根据订单号查询退款单
func (r *RefundRepository) FindByOrderNo(ctx context.Context, orderNo string) (*model.RefundOrder, error) {
	var refund model.RefundOrder
	err := r.db.WithContext(ctx).Where("order_no = ?", orderNo).First(&refund).Error
	return &refund, err
}

// 根据用户ID查询退款单列表
func (r *RefundRepository) ListByUserID(ctx context.Context, userID int64, page, pageSize int) ([]*model.RefundOrder, int64, error) {
	var refunds []*model.RefundOrder
	var total int64

	err := r.db.WithContext(ctx).Model(&model.RefundOrder{}).
		Where("user_id = ?", userID).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err = r.db.WithContext(ctx).Where("user_id = ?", userID).
		Offset(offset).Limit(pageSize).
		Order("created_at DESC").
		Find(&refunds).Error

	return refunds, total, err
}

// 根据条件查询退款单列表
func (r *RefundRepository) ListByCondition(ctx context.Context, condition map[string]interface{}, page, pageSize int) ([]*model.RefundOrder, int64, error) {
	var refunds []*model.RefundOrder
	var total int64

	db := r.db.WithContext(ctx).Model(&model.RefundOrder{})

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
	err = db.Offset(offset).Limit(pageSize).
		Order("created_at DESC").
		Find(&refunds).Error

	return refunds, total, err
}

// 更新退款单状态
func (r *RefundRepository) UpdateStatus(ctx context.Context, refundNo string, status string) error {
	updates := map[string]interface{}{
		"status":     status,
		"updated_at": time.Now(),
	}

	if status == model.RefundStatusApproved || status == model.RefundStatusRejected {
		now := time.Now()
		updates["processed_at"] = &now
	}

	return r.db.WithContext(ctx).Model(&model.RefundOrder{}).
		Where("refund_no = ?", refundNo).
		Updates(updates).Error
}
