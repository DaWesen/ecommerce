package repository

import (
	"context"
	"ecommerce/product-service/internal/model"
	"errors"

	"gorm.io/gorm"
)

// 货物接口
type ProductRepository interface {
	//基础CRUD
	Create(ctx context.Context, product *model.Product) error
	FindByID(ctx context.Context, id int64) (*model.Product, error)
	Update(ctx context.Context, product *model.Product) error
	Delete(ctx context.Context, id int64) error

	//状态
	UpdateStatus(ctx context.Context, id int64, status model.ProductStatus) error
	Online(ctx context.Context, id int64) error
	Offline(ctx context.Context, id int64) error

	//搜索
	SearchForUser(ctx context.Context,
		category *string,
		minPrice, maxPrice *float64,
		keyword *string,
		page, pageSize int32) ([]*model.Product, int64, error)

	SearchForAdmin(ctx context.Context,
		id *int64,
		category *string,
		minPrice, maxPrice *float64,
		keyword *string,
		page, pageSize int32) ([]*model.Product, int64, error)

	//库存管理
	UpdateStock(ctx context.Context, id int64, delta int32) (bool, error)
	CheckStock(ctx context.Context, id int64, quantity int32) (bool, error)
}

type productRepositoryImpl struct {
	db *gorm.DB
}

// 创建商品存储实例
func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepositoryImpl{db: db}
}

// 创建商品
func (r *productRepositoryImpl) Create(ctx context.Context, product *model.Product) error {
	return r.db.WithContext(ctx).Create(product).Error
}

// 根据ID查找商品
func (r *productRepositoryImpl) FindByID(ctx context.Context, id int64) (*model.Product, error) {
	var product model.Product
	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&product).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &product, nil
}

// 更新商品
func (r *productRepositoryImpl) Update(ctx context.Context, product *model.Product) error {
	return r.db.WithContext(ctx).Save(product).Error
}

// 删除商品
func (r *productRepositoryImpl) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).
		Model(&model.Product{}).
		Where("id = ?", id).
		Update("status", model.ProductStatusDELETED).Error
}

// 更新商品状态
func (r *productRepositoryImpl) UpdateStatus(ctx context.Context, id int64, status model.ProductStatus) error {
	return r.db.WithContext(ctx).
		Model(&model.Product{}).
		Where("id = ?", id).
		Update("status", status).Error
}

// 上架商品
func (r *productRepositoryImpl) Online(ctx context.Context, id int64) error {
	return r.UpdateStatus(ctx, id, model.ProductStatusONLINE)
}

// 下架商品
func (r *productRepositoryImpl) Offline(ctx context.Context, id int64) error {
	return r.UpdateStatus(ctx, id, model.ProductStatusOFFLINE)
}

// 用户搜索商品
func (r *productRepositoryImpl) SearchForUser(ctx context.Context, category *string,
	minPrice, maxPrice *float64, keyword *string,
	page, pageSize int32) ([]*model.Product, int64, error) {
	query := r.db.WithContext(ctx).Model(&model.Product{}).
		Where("status = ?", model.ProductStatusONLINE)
	if category != nil && *category != "" {
		query = query.Where("category = ?", *category)
	}
	if minPrice != nil {
		query = query.Where("price >= ?", *minPrice)
	}
	if maxPrice != nil {
		query = query.Where("price <= ?", *maxPrice)
	}
	if keyword != nil && *keyword != "" {
		query = query.Where("name Like ?", "%"+*keyword+"%")
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var products []*model.Product
	offset := (page - 1) * pageSize
	err := query.Offset(int(offset)).
		Limit(int(pageSize)).
		Order("created_at DESC").
		Find(&products).Error
	return products, total, err
}

// 管理员搜素商品
func (r *productRepositoryImpl) SearchForAdmin(ctx context.Context, id *int64,
	category *string, minPrice, maxPrice *float64,
	keyword *string, page, pageSize int32) ([]*model.Product, int64, error) {
	query := r.db.WithContext(ctx).Model(&model.Product{}).
		Where("status != ?", model.ProductStatusDELETED)
	if id != nil && *id > 0 {
		query = query.Where("id = ?", *id)
	}
	if category != nil && *category != "" {
		query = query.Where("category = ?", *category)
	}
	if minPrice != nil {
		query = query.Where("price >= ?", *minPrice)
	}
	if maxPrice != nil {
		query = query.Where("price <= ?", *maxPrice)
	}
	if keyword != nil && *keyword != "" {
		query = query.Where("name Like ?", "%"+*keyword+"%")
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var products []*model.Product
	offset := (page - 1) * pageSize
	err := query.Offset(int(offset)).
		Limit(int(pageSize)).
		Order("created_at DESC").
		Find(&products).Error
	return products, total, err
}

// 更新库存
func (r *productRepositoryImpl) UpdateStock(ctx context.Context,
	id int64, delta int32) (bool, error) {
	if delta < 0 {
		var currentStock int32
		err := r.db.WithContext(ctx).
			Model(&model.Product{}).
			Select("stock").
			Where("id = ?", id).
			Scan(&currentStock).Error

		if err != nil {
			return false, err
		}

		if currentStock+delta < 0 {
			return false, errors.New("库存不足")
		}
	}
	result := r.db.WithContext(ctx).
		Model(&model.Product{}).
		Where("id = ? AND stock + ? >= 0", id, delta).
		Update("stock", gorm.Expr("stock + ?", delta))

	return result.RowsAffected > 0, result.Error
}

// 检查库存
func (r *productRepositoryImpl) CheckStock(ctx context.Context,
	id int64, quantity int32) (bool, error) {
	var stock int32
	err := r.db.WithContext(ctx).
		Model(&model.Product{}).
		Select("stock").
		Where("id = ? AND status = ?", id, model.ProductStatusONLINE).
		Scan(&stock).Error

	if err != nil {
		return false, err
	}

	return stock >= quantity, nil
}
