package stockReservationDao

import (
	"context"
	"ecommerce/order-service/internal/dao/interfaces"
	"ecommerce/order-service/internal/model"
	"time"

	"gorm.io/gorm"
)

type StockReservationRepository struct {
	db *gorm.DB
}

func NewStockReservationRepository(db *gorm.DB) interfaces.IStockReservationRepository {
	return &StockReservationRepository{db: db}
}

// 创建库存预占
func (r *StockReservationRepository) Create(ctx context.Context, reservation *model.StockReservation) error {
	return r.db.WithContext(ctx).Create(reservation).Error
}

// 更新库存预占
func (r *StockReservationRepository) Update(ctx context.Context, reservation *model.StockReservation) error {
	return r.db.WithContext(ctx).Save(reservation).Error
}

// 根据预占ID查询
func (r *StockReservationRepository) FindByReserveID(ctx context.Context, reserveID string) (*model.StockReservation, error) {
	var reservation model.StockReservation
	err := r.db.WithContext(ctx).Where("reserve_id = ?", reserveID).First(&reservation).Error
	return &reservation, err
}

// 根据订单号查询库存预占
func (r *StockReservationRepository) FindByOrderNo(ctx context.Context, orderNo string) ([]*model.StockReservation, error) {
	var reservations []*model.StockReservation
	err := r.db.WithContext(ctx).Where("order_no = ?", orderNo).Find(&reservations).Error
	return reservations, err
}

// 根据预占ID删除
func (r *StockReservationRepository) DeleteByReserveID(ctx context.Context, reserveID string) error {
	return r.db.WithContext(ctx).Where("reserve_id = ?", reserveID).Delete(&model.StockReservation{}).Error
}

// 查询已过期的库存预占
func (r *StockReservationRepository) FindExpiredReservations(ctx context.Context) ([]*model.StockReservation, error) {
	var reservations []*model.StockReservation
	now := time.Now()
	err := r.db.WithContext(ctx).
		Where("status = ? AND expire_time < ?", model.StockStatusReserved, now).
		Find(&reservations).Error
	return reservations, err
}

// 更新库存预占状态
func (r *StockReservationRepository) UpdateStatus(ctx context.Context, reserveID, status string) error {
	return r.db.WithContext(ctx).Model(&model.StockReservation{}).
		Where("reserve_id = ?", reserveID).
		Updates(map[string]interface{}{
			"status":     status,
			"updated_at": time.Now(),
		}).Error
}
