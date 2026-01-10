package client

import (
	"context"
	"ecommerce/product-service/kitex_gen/api"
)

// CreateProduct 创建商品（管理员）
func (pc *ProductClient) CreateProduct(ctx context.Context, req *api.CreateProductReq) (*api.CreateProductResp, error) {
	return pc.client.CreateProduct(ctx, req)
}

// GetProduct 获取商品详情
func (pc *ProductClient) GetProduct(ctx context.Context, req *api.GetProductReq) (*api.GetProductResp, error) {
	return pc.client.GetProduct(ctx, req)
}

// UpdateProduct 更新商品（管理员）
func (pc *ProductClient) UpdateProduct(ctx context.Context, req *api.UpdateProductReq) (*api.UpdateProductResp, error) {
	return pc.client.UpdateProduct(ctx, req)
}

// OnlineProduct 上架商品（管理员）
func (pc *ProductClient) OnlineProduct(ctx context.Context, req *api.OnlineProductReq) (*api.OnlineProductResp, error) {
	return pc.client.OnlineProduct(ctx, req)
}

// OfflineProduct 下架商品（管理员）
func (pc *ProductClient) OfflineProduct(ctx context.Context, req *api.OfflineProductReq) (*api.OfflineProductResp, error) {
	return pc.client.OfflineProduct(ctx, req)
}

// DeleteProduct 删除商品（管理员）
func (pc *ProductClient) DeleteProduct(ctx context.Context, req *api.DeleteProductReq) (*api.DeleteProductResp, error) {
	return pc.client.DeleteProduct(ctx, req)
}

// UserSearchProducts 用户搜索商品
func (pc *ProductClient) UserSearchProducts(ctx context.Context, req *api.UserSearchProductsReq) (*api.UserSearchProductsResp, error) {
	return pc.client.UserSearchProducts(ctx, req)
}

// AdminSearchProducts 管理员搜索商品
func (pc *ProductClient) AdminSearchProducts(ctx context.Context, req *api.AdminSearchProductsReq) (*api.AdminSearchProductsResp, error) {
	return pc.client.AdminSearchProducts(ctx, req)
}
