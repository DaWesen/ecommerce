package dao

import (
	"ecommerce/order-service/internal/dao/interfaces"
	"ecommerce/order-service/internal/dao/orderDao"
	"ecommerce/order-service/internal/dao/orderItemDao"
	"ecommerce/order-service/internal/dao/refundDao"
	"ecommerce/order-service/internal/dao/stockReservationDao"
	"ecommerce/order-service/internal/dao/timeOutTaskDao"

	"gorm.io/gorm"
)

type DaoFactory struct {
	OrderRepo            interfaces.IOrderRepository
	OrderItemRepo        interfaces.IOrderItemRepository
	RefundRepo           interfaces.IRefundRepository
	StockReservationRepo interfaces.IStockReservationRepository
	TimeoutTaskRepo      interfaces.ITimeoutTaskRepository
}

func NewDaoFactory(db *gorm.DB) *DaoFactory {
	return &DaoFactory{
		OrderRepo:            orderDao.NewOrderRepository(db),
		OrderItemRepo:        orderItemDao.NewOrderItemRepository(db),
		RefundRepo:           refundDao.NewRefundRepository(db),
		StockReservationRepo: stockReservationDao.NewStockReservationRepository(db),
		TimeoutTaskRepo:      timeOutTaskDao.NewTimeOutTaskRepository(db),
	}
}
