package client

import (
	"ecommerce/order-service/kitex_gen/api/orderservice"
	"ecommerce/product-service/kitex_gen/api/productservice"
	"ecommerce/user-service/kitex_gen/api/userservice"
)

// 客户端类型定义
type (
	UserClient struct {
		client userservice.Client
	}

	ProductClient struct {
		client productservice.Client
	}

	OrderClient struct {
		client orderservice.Client
	}
)
