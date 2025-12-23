package main

import (
	"context"
	api "ecommerce/product-service/kitex_gen/api"
)

// ProductServiceImpl implements the last service interface defined in the IDL.
type ProductServiceImpl struct{}

// CreateProduct implements the ProductServiceImpl interface.
func (s *ProductServiceImpl) CreateProduct(ctx context.Context, req *api.CreateProductReq) (resp *api.CreateProductResp, err error) {
	// TODO: Your code here...
	return
}

// GetProduct implements the ProductServiceImpl interface.
func (s *ProductServiceImpl) GetProduct(ctx context.Context, req *api.GetProductReq) (resp *api.GetProductResp, err error) {
	// TODO: Your code here...
	return
}

// UpdateProduct implements the ProductServiceImpl interface.
func (s *ProductServiceImpl) UpdateProduct(ctx context.Context, req *api.UpdateProductReq) (resp *api.UpdateProductResp, err error) {
	// TODO: Your code here...
	return
}

// OnlineProduct implements the ProductServiceImpl interface.
func (s *ProductServiceImpl) OnlineProduct(ctx context.Context, req *api.OnlineProductReq) (resp *api.OnlineProductResp, err error) {
	// TODO: Your code here...
	return
}

// OfflineProduct implements the ProductServiceImpl interface.
func (s *ProductServiceImpl) OfflineProduct(ctx context.Context, req *api.OfflineProductReq) (resp *api.OfflineProductResp, err error) {
	// TODO: Your code here...
	return
}

// DeleteProduct implements the ProductServiceImpl interface.
func (s *ProductServiceImpl) DeleteProduct(ctx context.Context, req *api.DeleteProductReq) (resp *api.DeleteProductResp, err error) {
	// TODO: Your code here...
	return
}

// UserSearchProducts implements the ProductServiceImpl interface.
func (s *ProductServiceImpl) UserSearchProducts(ctx context.Context, req *api.UserSearchProductsReq) (resp *api.UserSearchProductsResp, err error) {
	// TODO: Your code here...
	return
}

// AdminSearchProducts implements the ProductServiceImpl interface.
func (s *ProductServiceImpl) AdminSearchProducts(ctx context.Context, req *api.AdminSearchProductsReq) (resp *api.AdminSearchProductsResp, err error) {
	// TODO: Your code here...
	return
}
