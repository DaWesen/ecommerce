package client

import (
	"context"
	"ecommerce/order-service/internal/dao/interfaces"
	"ecommerce/order-service/kitex_gen/api"
	"ecommerce/order-service/kitex_gen/api/productservice"

	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/klog"
)

type ProductClient struct {
	client productservice.Client
}

func NewProductClient(addr string) (interfaces.IProductClient, error) {
	c, err := productservice.NewClient("product-service",
		client.WithHostPorts(addr),
	)
	if err != nil {
		return nil, err
	}
	return &ProductClient{client: c}, nil
}

func (pc *ProductClient) GetProductInfo(ctx context.Context, productID int64) (*interfaces.ProductInfo, error) {
	req := &api.GetProductReq{
		Id: productID,
	}

	resp, err := pc.client.GetProduct(ctx, req)
	if err != nil {
		klog.Errorf("GetProductInfo failed: %v", err)
		return nil, err
	}

	if !resp.Success || resp.Product == nil {
		return nil, nil
	}

	productInfo := &interfaces.ProductInfo{
		ID:       resp.Product.Id,
		Name:     resp.Product.Name,
		Price:    resp.Product.Price,
		Stock:    resp.Product.Stock,
		Status:   int32(resp.Product.Status),
		Category: resp.Product.Category,
		Avatar:   resp.Product.Avatar,
	}

	if resp.Product.Brand != nil {
		productInfo.Brand = *resp.Product.Brand
	}

	return productInfo, nil
}

func (pc *ProductClient) BatchGetProducts(ctx context.Context, productIDs []int64) (map[int64]*interfaces.ProductInfo, error) {
	result := make(map[int64]*interfaces.ProductInfo)
	for _, id := range productIDs {
		productInfo, err := pc.GetProductInfo(ctx, id)
		if err == nil && productInfo != nil {
			result[id] = productInfo
		}
	}
	return result, nil
}

func (pc *ProductClient) CheckStock(ctx context.Context, productID int64, quantity int32) (bool, error) {
	productInfo, err := pc.GetProductInfo(ctx, productID)
	if err != nil {
		return false, err
	}
	return productInfo != nil && productInfo.Stock >= quantity, nil
}
