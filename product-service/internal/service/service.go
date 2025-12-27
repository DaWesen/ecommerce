package service

import (
	"context"
	"ecommerce/product-service/internal/model"
	"ecommerce/product-service/internal/repository"
	"ecommerce/product-service/kitex_gen/api"
	"errors"
	"fmt"
	"time"
)

type ProductService interface {
	CreateProduct(ctx context.Context, req *api.CreateProductReq) (*api.CreateProductResp, error)
	GetProduct(ctx context.Context, id int64) (*api.GetProductResp, error)
	UpdateProduct(ctx context.Context, req *api.UpdateProductReq) (*api.UpdateProductResp, error)
	DeleteProduct(ctx context.Context, id int64) (*api.DeleteProductResp, error)

	OnlineProduct(ctx context.Context, id int64) (*api.OnlineProductResp, error)
	OfflineProduct(ctx context.Context, id int64) (*api.OfflineProductResp, error)

	UserSearchProducts(ctx context.Context, req *api.UserSearchProductsReq) (*api.UserSearchProductsResp, error)
	AdminSearchProducts(ctx context.Context, req *api.AdminSearchProductsReq) (*api.AdminSearchProductsResp, error)

	UpdateStock(ctx context.Context, id int64, delta int32) error
	CheckStock(cttx context.Context, id int64, quantity int32) (bool, error)
}

type productServiceImpl struct {
	productRepo repository.ProductRepository
}

func NewProductService(productRepo repository.ProductRepository) ProductService {
	return &productServiceImpl{
		productRepo: productRepo,
	}
}

// 创建商品
func (s *productServiceImpl) CreateProduct(ctx context.Context, req *api.CreateProductReq) (*api.CreateProductResp, error) {
	if req.Name == "" {
		return &api.CreateProductResp{
			Success: false,
			Code:    400,
			Message: stringPtr("商品名称不能为空"),
		}, nil
	}
	if req.Price <= 0 {
		return &api.CreateProductResp{
			Success: false,
			Code:    400,
			Message: stringPtr("商品价格必须大于0"),
		}, nil
	}
	if req.Stock < 0 {
		return &api.CreateProductResp{
			Success: false,
			Code:    400,
			Message: stringPtr("库存不能为负数"),
		}, nil
	}
	now := time.Now().Unix()
	product := &model.Product{
		Name:      req.Name,
		Avatar:    req.Avatar,
		Category:  req.Category,
		Price:     req.Price,
		Stock:     req.Stock,
		Status:    model.ProductStatus(req.Status),
		CreatedAt: now,
		UpdatedAt: now,
	}
	if req.Brand != nil && *req.Brand != "" {
		product.Brand = *req.Brand
	}
	err := s.productRepo.Create(ctx, product)
	if err != nil {
		fmt.Printf("创建商品失败: %v\n", err)
		return &api.CreateProductResp{
			Success: false,
			Code:    500,
			Message: stringPtr("创建商品失败，请稍后重试"),
		}, nil
	}
	return &api.CreateProductResp{
		Success: true,
		Code:    0,
		Message: stringPtr("创建成功"),
		Product: s.convertToAPIProduct(product),
	}, nil
}

// 获取商品
func (s *productServiceImpl) GetProduct(ctx context.Context, id int64) (*api.GetProductResp, error) {
	product, err := s.productRepo.FindByID(ctx, id)
	if err != nil {
		return &api.GetProductResp{
			Success: false,
			Code:    500,
			Message: stringPtr("查询商品失败"),
			Product: s.convertToAPIProduct(product),
		}, nil
	}
	if product == nil {
		return &api.GetProductResp{
			Success: false,
			Code:    404,
			Message: stringPtr("商品不存在"),
		}, nil
	}

	return &api.GetProductResp{
		Success: true,
		Code:    0,
		Message: stringPtr("查询成功"),
		Product: s.convertToAPIProduct(product),
	}, nil
}

// 更新商品
func (s *productServiceImpl) UpdateProduct(ctx context.Context, req *api.UpdateProductReq) (*api.UpdateProductResp, error) {
	product, err := s.productRepo.FindByID(ctx, req.GetId())
	if err != nil {
		return &api.UpdateProductResp{
			Success: false,
			Code:    500,
			Message: stringPtr("查询商品失败"),
		}, nil
	}
	if product == nil {
		return &api.UpdateProductResp{
			Success: false,
			Code:    404,
			Message: stringPtr("商品不存在"),
		}, nil
	}
	if req.Name != nil {
		product.Name = *req.Name
	}
	if req.Avatar != nil {
		product.Avatar = *req.Avatar
	}
	if req.Category != nil {
		product.Category = *req.Category
	}
	if req.Price != nil {
		product.Price = *req.Price
	}
	if req.Stock != nil {
		product.Stock = *req.Stock
	}
	if req.Status != nil {
		product.Status = model.ProductStatus(*req.Status)
	}
	if req.Brand != nil {
		product.Brand = *req.Brand
	}
	product.UpdatedAt = time.Now().Unix()
	err = s.productRepo.Update(ctx, product)
	if err != nil {
		return &api.UpdateProductResp{
			Success: false,
			Code:    500,
			Message: stringPtr("更新商品失败"),
		}, nil
	}
	return &api.UpdateProductResp{
		Success: true,
		Code:    0,
		Message: stringPtr("更新成功"),
		Product: s.convertToAPIProduct(product),
	}, nil
}

// 删除商品
func (s *productServiceImpl) DeleteProduct(ctx context.Context, id int64) (*api.DeleteProductResp, error) {
	product, err := s.productRepo.FindByID(ctx, id)
	if err != nil {
		return &api.DeleteProductResp{
			Success: false,
			Code:    500,
			Message: stringPtr("查询商品失败"),
		}, nil
	}
	if product == nil {
		return &api.DeleteProductResp{
			Success: false,
			Code:    404,
			Message: stringPtr("商品不存在"),
		}, nil
	}
	err = s.productRepo.UpdateStatus(ctx, id, model.ProductStatusDELETED)
	if err != nil {
		return &api.DeleteProductResp{
			Success: false,
			Code:    500,
			Message: stringPtr("删除商品失败"),
		}, nil
	}
	return &api.DeleteProductResp{
		Success: true,
		Code:    0,
		Message: stringPtr("删除成功"),
	}, nil
}

// 上架商品
func (s *productServiceImpl) OnlineProduct(ctx context.Context, id int64) (*api.OnlineProductResp, error) {
	product, err := s.productRepo.FindByID(ctx, id)
	if err != nil {
		return &api.OnlineProductResp{
			Success: false,
			Code:    500,
			Message: stringPtr("查询商品失败"),
		}, nil
	}
	if product == nil {
		return &api.OnlineProductResp{
			Success: false,
			Code:    404,
			Message: stringPtr("商品不存在"),
		}, nil
	}
	oldStatus := product.Status
	err = s.productRepo.Online(ctx, id)
	if err != nil {
		return &api.OnlineProductResp{
			Success: false,
			Code:    500,
			Message: stringPtr("上架商品失败"),
		}, nil
	}
	return &api.OnlineProductResp{
		Success:    true,
		Code:       0,
		Message:    stringPtr("上架成功"),
		OldStatus:  api.ProductStatus(oldStatus),
		NewStatus_: api.ProductStatus(model.ProductStatusONLINE),
		OperatedAt: time.Now().Unix(),
	}, nil
}

// 下架
func (s *productServiceImpl) OfflineProduct(ctx context.Context, id int64) (*api.OfflineProductResp, error) {
	product, err := s.productRepo.FindByID(ctx, id)
	if err != nil {
		return &api.OfflineProductResp{
			Success: false,
			Code:    500,
			Message: stringPtr("查询商品失败"),
		}, nil
	}

	if product == nil {
		return &api.OfflineProductResp{
			Success: false,
			Code:    404,
			Message: stringPtr("商品不存在"),
		}, nil
	}
	oldStatus := product.Status
	err = s.productRepo.Offline(ctx, id)
	if err != nil {
		return &api.OfflineProductResp{
			Success: false,
			Code:    500,
			Message: stringPtr("下架商品失败"),
		}, nil
	}
	return &api.OfflineProductResp{
		Success:    true,
		Code:       0,
		Message:    stringPtr("下架成功"),
		OldStatus:  api.ProductStatus(oldStatus),
		NewStatus_: api.ProductStatus(model.ProductStatusOFFLINE),
		OperatedAt: time.Now().Unix(),
	}, nil
}

// 用户搜索
func (s *productServiceImpl) UserSearchProducts(ctx context.Context, req *api.UserSearchProductsReq) (*api.UserSearchProductsResp, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 || req.PageSize > 100 {
		req.PageSize = 20
	}
	products, total, err := s.productRepo.SearchForUser(ctx,
		req.Category,
		req.MinPrice,
		req.MaxPrice,
		req.Keyword,
		req.Page,
		req.PageSize,
	)
	if err != nil {
		return &api.UserSearchProductsResp{
			Success: false,
			Code:    500,
			Message: stringPtr("搜索商品失败"),
		}, nil
	}
	apiProducts := make([]*api.SimpleProduct, 0, len(products))
	for _, p := range products {
		apiProducts = append(apiProducts, s.convertToAPISimpleProduct(p))
	}
	return &api.UserSearchProductsResp{
		Success:  true,
		Code:     0,
		Message:  stringPtr("搜索成功"),
		Total:    int32(total),
		Page:     req.Page,
		PageSize: req.PageSize,
		Products: apiProducts,
	}, nil
}

// 管理员搜索
func (s *productServiceImpl) AdminSearchProducts(ctx context.Context, req *api.AdminSearchProductsReq) (*api.AdminSearchProductsResp, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 || req.PageSize > 100 {
		req.PageSize = 20
	}
	products, total, err := s.productRepo.SearchForAdmin(ctx,
		req.Id,
		req.Category,
		req.MinPrice,
		req.MaxPrice,
		req.Keyword,
		req.Page,
		req.PageSize,
	)

	if err != nil {
		return &api.AdminSearchProductsResp{
			Success: false,
			Code:    500,
			Message: stringPtr("搜索商品失败"),
		}, nil
	}
	apiProducts := make([]*api.Product, 0, len(products))
	for _, p := range products {
		apiProducts = append(apiProducts, s.convertToAPIProduct(p))
	}
	return &api.AdminSearchProductsResp{
		Success:  true,
		Code:     0,
		Message:  stringPtr("搜索成功"),
		Total:    int32(total),
		Page:     req.Page,
		PageSize: req.PageSize,
		Products: apiProducts,
	}, nil
}

// 更新库存
func (s *productServiceImpl) UpdateStock(ctx context.Context, id int64, delta int32) error {

	success, err := s.productRepo.UpdateStock(ctx, id, delta)
	if err != nil {
		return err
	}

	if !success {
		return errors.New("更新库存失败")
	}

	return nil
}

// 检查库存
func (s *productServiceImpl) CheckStock(ctx context.Context, id int64, quantity int32) (bool, error) {
	return s.productRepo.CheckStock(ctx, id, quantity)
}

// 模型转化
func (s *productServiceImpl) convertToAPIProduct(p *model.Product) *api.Product {
	product := &api.Product{
		Id:        p.ID,
		Name:      p.Name,
		Avatar:    p.Avatar,
		Category:  p.Category,
		Price:     p.Price,
		Stock:     p.Stock,
		Status:    api.ProductStatus(p.Status),
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}

	if p.Brand != "" {
		product.Brand = &p.Brand
	}

	return product
}

// 模型转化
func (s *productServiceImpl) convertToAPISimpleProduct(p *model.Product) *api.SimpleProduct {
	simpleProduct := &api.SimpleProduct{
		Id:       p.ID,
		Category: p.Category,
		Price:    p.Price,
		Stock:    p.Stock,
		Status:   api.ProductStatus(p.Status),
		Name:     p.Name,
		Avatar:   p.Avatar,
	}

	if p.Brand != "" {
		simpleProduct.Brand = &p.Brand
	}

	return simpleProduct
}

// 字符串转化
func stringPtr(s string) *string {
	return &s
}
