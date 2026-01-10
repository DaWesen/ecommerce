package client

import (
	"context"
	"fmt"
	"time"

	"ecommerce/gateway/config"
	"ecommerce/order-service/kitex_gen/api/orderservice"
	"ecommerce/product-service/kitex_gen/api/productservice"
	"ecommerce/user-service/kitex_gen/api/userservice"

	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/endpoint"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/loadbalance"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/transport"
)

// ClientManager 管理所有RPC客户端
type ClientManager struct {
	UserClient    *UserClient
	ProductClient *ProductClient
	OrderClient   *OrderClient
}

// NewClientManager 创建客户端管理器
func NewClientManager(cfg *config.Config) (*ClientManager, error) {
	// 创建用户服务客户端
	userClient, err := newUserClient(cfg.Services.UserService)
	if err != nil {
		return nil, fmt.Errorf("创建用户服务客户端失败: %w", err)
	}

	// 创建商品服务客户端
	productClient, err := newProductClient(cfg.Services.ProductService)
	if err != nil {
		return nil, fmt.Errorf("创建商品服务客户端失败: %w", err)
	}

	// 创建订单服务客户端
	orderClient, err := newOrderClient(cfg.Services.OrderService)
	if err != nil {
		return nil, fmt.Errorf("创建订单服务客户端失败: %w", err)
	}

	klog.Info("✅ 所有RPC客户端创建成功")

	return &ClientManager{
		UserClient:    userClient,
		ProductClient: productClient,
		OrderClient:   orderClient,
	}, nil
}

// newUserClient 创建用户服务客户端
func newUserClient(cfg config.ServiceConfig) (*UserClient, error) {
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	opts := []client.Option{
		client.WithHostPorts(addr),
		client.WithTransportProtocol(transport.TTHeader),
		client.WithLoadBalancer(loadbalance.NewWeightedBalancer()),
		client.WithRPCTimeout(time.Duration(cfg.Timeout) * time.Millisecond),
		client.WithConnectTimeout(3 * time.Second),
		client.WithMiddleware(clientMiddleware),
	}

	// 创建客户端
	c, err := userservice.NewClient(cfg.Name, opts...)
	if err != nil {
		return nil, err
	}

	return &UserClient{client: c}, nil
}

// newProductClient 创建商品服务客户端
func newProductClient(cfg config.ServiceConfig) (*ProductClient, error) {
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	opts := []client.Option{
		client.WithHostPorts(addr),
		client.WithTransportProtocol(transport.TTHeader),
		client.WithLoadBalancer(loadbalance.NewWeightedBalancer()),
		client.WithRPCTimeout(time.Duration(cfg.Timeout) * time.Millisecond),
		client.WithConnectTimeout(3 * time.Second),
		client.WithMiddleware(clientMiddleware),
	}

	c, err := productservice.NewClient(cfg.Name, opts...)
	if err != nil {
		return nil, err
	}

	return &ProductClient{client: c}, nil
}

// newOrderClient 创建订单服务客户端
func newOrderClient(cfg config.ServiceConfig) (*OrderClient, error) {
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	opts := []client.Option{
		client.WithHostPorts(addr),
		client.WithTransportProtocol(transport.TTHeader),
		client.WithLoadBalancer(loadbalance.NewWeightedBalancer()),
		client.WithRPCTimeout(time.Duration(cfg.Timeout) * time.Millisecond),
		client.WithConnectTimeout(3 * time.Second),
		client.WithMiddleware(clientMiddleware),
	}

	c, err := orderservice.NewClient(cfg.Name, opts...)
	if err != nil {
		return nil, err
	}

	return &OrderClient{client: c}, nil
}

// clientMiddleware 客户端中间件
func clientMiddleware(next endpoint.Endpoint) endpoint.Endpoint {
	return func(ctx context.Context, req, resp interface{}) (err error) {
		start := time.Now()

		// 调用下一个中间件或实际的服务
		err = next(ctx, req, resp)

		duration := time.Since(start)

		// 记录RPC调用日志
		ri := rpcinfo.GetRPCInfo(ctx)
		if ri != nil {
			klog.Infof("RPC调用: 服务=%s, 方法=%s, 耗时=%v, 错误=%v",
				ri.To().ServiceName(),
				ri.To().Method(),
				duration,
				err,
			)
		}

		return err
	}
}

// Close 关闭所有客户端连接
func (cm *ClientManager) Close() {
	klog.Info("正在关闭RPC客户端连接...")
}
