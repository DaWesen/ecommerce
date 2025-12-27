package main

import (
	"context"
	"ecommerce/product-service/internal/repository"
	"ecommerce/product-service/internal/service"
	api "ecommerce/product-service/kitex_gen/api"
	"ecommerce/product-service/pkg/config"
	"ecommerce/product-service/pkg/database"
	"fmt"
	"log"
)

// ProductServiceImpl implements the last service interface defined in the IDL.
type ProductServiceImpl struct {
	productService service.ProductService
}

// 创建处理器
func NewProductServiceImpl() (*ProductServiceImpl, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("加载配置失败: %v", err)
	}
	db, _, err := database.NewDatabase(&cfg.Database)
	if err != nil {
		return nil, fmt.Errorf("数据库连接失败: %v", err)
	}
	productRepo := repository.NewProductRepository(db)
	productService := service.NewProductService(productRepo)
	return &ProductServiceImpl{
		productService: productService,
	}, nil
}

// CreateProduct implements the ProductServiceImpl interface.
func (s *ProductServiceImpl) CreateProduct(ctx context.Context, req *api.CreateProductReq) (resp *api.CreateProductResp, err error) {
	// TODO: Your code here...
	log.Printf("接收到创建商品请求:name=%s,price=%.2f", req.Name, req.Price)
	return s.productService.CreateProduct(ctx, req)
}

// GetProduct implements the ProductServiceImpl interface.
func (s *ProductServiceImpl) GetProduct(ctx context.Context, req *api.GetProductReq) (resp *api.GetProductResp, err error) {
	// TODO: Your code here...
	log.Printf("接收到获取商品请求: id=%d", req.GetId())
	return s.productService.GetProduct(ctx, req.GetId())
}

// UpdateProduct implements the ProductServiceImpl interface.
func (s *ProductServiceImpl) UpdateProduct(ctx context.Context, req *api.UpdateProductReq) (resp *api.UpdateProductResp, err error) {
	// TODO: Your code here...
	log.Printf("接收到更新商品请求: id=%d", req.GetId())
	return s.productService.UpdateProduct(ctx, req)
}

// OnlineProduct implements the ProductServiceImpl interface.
func (s *ProductServiceImpl) OnlineProduct(ctx context.Context, req *api.OnlineProductReq) (resp *api.OnlineProductResp, err error) {
	// TODO: Your code here...
	log.Printf("接收到上架商品请求: id=%d", req.GetId())
	return s.productService.OnlineProduct(ctx, req.GetId())
}

// OfflineProduct implements the ProductServiceImpl interface.
func (s *ProductServiceImpl) OfflineProduct(ctx context.Context, req *api.OfflineProductReq) (resp *api.OfflineProductResp, err error) {
	// TODO: Your code here...
	log.Printf("接收到下架商品请求:id=%d", req.GetId())
	return s.productService.OfflineProduct(ctx, req.GetId())
}

// DeleteProduct implements the ProductServiceImpl interface.
func (s *ProductServiceImpl) DeleteProduct(ctx context.Context, req *api.DeleteProductReq) (resp *api.DeleteProductResp, err error) {
	// TODO: Your code here...
	log.Printf("接收到删除商品请求:id=%d", req.GetId())
	return s.productService.DeleteProduct(ctx, req.GetId())
}

// UserSearchProducts implements the ProductServiceImpl interface.
func (s *ProductServiceImpl) UserSearchProducts(ctx context.Context, req *api.UserSearchProductsReq) (resp *api.UserSearchProductsResp, err error) {
	// TODO: Your code here...
	log.Printf("接收到用户搜索商品的请求")
	return s.productService.UserSearchProducts(ctx, req)
}

// AdminSearchProducts implements the ProductServiceImpl interface.
func (s *ProductServiceImpl) AdminSearchProducts(ctx context.Context, req *api.AdminSearchProductsReq) (resp *api.AdminSearchProductsResp, err error) {
	// TODO: Your code here...
	log.Printf("接收到管理员搜索商品的请求")
	return s.productService.AdminSearchProducts(ctx, req)
}
