package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"ecommerce/order-service/kitex_gen/api"
	"ecommerce/order-service/kitex_gen/api/orderservice"

	"github.com/cloudwego/kitex/client"
)

func main() {
	fmt.Println("🚀 订单服务完整测试")

	// 创建客户端
	c, err := orderservice.NewClient(
		"order.service",
		client.WithHostPorts("localhost:50053"),
		client.WithRPCTimeout(30*time.Second),
	)
	if err != nil {
		log.Fatalf("❌ 创建客户端失败: %v", err)
	}

	fmt.Println("✅ 客户端创建成功")

	ctx := context.Background()

	// 测试用例
	fmt.Println("\n📋 执行测试用例:")

	// 测试1: 创建订单
	fmt.Println("\n1. 测试创建订单")
	orderNo := createSimpleOrder(ctx, c)
	if orderNo == "" {
		log.Fatal("❌ 创建订单失败，停止测试")
	}

	// 测试2: 查询刚创建的订单
	fmt.Println("\n2. 查询刚创建的订单")
	getOrder(ctx, c, orderNo)

	// 测试3: 查询订单列表
	fmt.Println("\n3. 查询订单列表")
	listOrders(ctx, c, 100001) // 使用创建订单的用户ID

	// 测试4: 支付订单
	fmt.Println("\n4. 测试支付订单")
	payOrder(ctx, c, orderNo, 100001)

	// 测试5: 再次查询订单
	fmt.Println("\n5. 支付后查询订单")
	getOrder(ctx, c, orderNo)

	// 测试6: 订单统计
	fmt.Println("\n6. 订单统计")
	orderStats(ctx, c, 100001)

	fmt.Println("\n🎉 所有测试完成!")
}

// createSimpleOrder 创建简单订单
func createSimpleOrder(ctx context.Context, c orderservice.Client) string {
	fmt.Println("   发送创建订单请求...")

	req := &api.CreateOrderReq{
		UserId:  100001, // 固定用户ID，方便后续查询
		Address: "测试地址123号",
		Phone:   "13800138000",
		Items: []*api.OrderItem{
			{
				ProductId:   1001,
				ProductName: "测试商品",
				Quantity:    1,
				Price:       99.99,
			},
		},
	}

	start := time.Now()
	resp, err := c.CreateOrder(ctx, req)
	duration := time.Since(start)

	if err != nil {
		fmt.Printf("   ❌ RPC错误: %v\n", err)
		return ""
	}

	fmt.Printf("   ⏱️  响应时间: %v\n", duration)

	if resp.Success {
		fmt.Printf("   ✅ 创建成功!\n")
		fmt.Printf("      订单号: %s\n", resp.OrderNo)
		fmt.Printf("      总金额: ¥%.2f\n", resp.TotalAmount)
		if resp.PaymentUrl != nil {
			fmt.Printf("      支付链接: %s\n", *resp.PaymentUrl)
		}
		return resp.OrderNo
	} else {
		fmt.Printf("   ❌ 创建失败: %s (代码: %d)\n", resp.Message, resp.Code)
		return ""
	}
}

// getOrder 查询订单详情
func getOrder(ctx context.Context, c orderservice.Client, orderNo string) {
	if orderNo == "" {
		fmt.Println("   ⚠️  没有订单号，跳过查询")
		return
	}

	fmt.Printf("   查询订单: %s\n", orderNo)

	userID := int64(100001) // 使用创建订单的用户ID
	req := &api.GetOrderReq{
		OrderNo: orderNo,
		UserId:  &userID,
	}

	resp, err := c.GetOrder(ctx, req)
	if err != nil {
		fmt.Printf("   ❌ 查询失败: %v\n", err)
		return
	}

	if resp.Success {
		fmt.Println("   ✅ 查询成功!")
		if resp.Order != nil {
			order := resp.Order
			fmt.Printf("      订单号: %s\n", order.OrderNo)
			fmt.Printf("      用户ID: %d\n", order.UserId)
			fmt.Printf("      总金额: ¥%.2f\n", order.TotalAmount)
			fmt.Printf("      状态: %v\n", order.Status)
			fmt.Printf("      地址: %s\n", order.Address)
			fmt.Printf("      电话: %s\n", order.Phone)

			if len(order.Items) > 0 {
				fmt.Printf("      商品数量: %d\n", len(order.Items))
				for i, item := range order.Items {
					fmt.Printf("        %d. %s x%d ¥%.2f\n",
						i+1, item.ProductName, item.Quantity, item.Price)
				}
			}
		}
	} else {
		fmt.Printf("   ❌ 查询失败: %s (代码: %d)\n", resp.Message, resp.Code)
	}
}

// listOrders 查询订单列表
func listOrders(ctx context.Context, c orderservice.Client, userID int64) {
	fmt.Printf("   查询用户 %d 的订单列表\n", userID)

	req := &api.ListOrdersReq{
		UserId:   userID,
		Page:     1,
		PageSize: 10,
	}

	resp, err := c.ListOrders(ctx, req)
	if err != nil {
		fmt.Printf("   ❌ 查询失败: %v\n", err)
		return
	}

	if resp.Success {
		fmt.Printf("   ✅ 查询成功: %s\n", resp.Message)
		fmt.Printf("      总订单数: %d\n", resp.Total)
		fmt.Printf("      当前页: %d\n", resp.Page)
		fmt.Printf("      返回订单数: %d\n", len(resp.Orders))

		if len(resp.Orders) > 0 {
			fmt.Println("      订单列表:")
			for i, order := range resp.Orders {
				fmt.Printf("        %d. %s - ¥%.2f - %v\n",
					i+1, order.OrderNo, order.TotalAmount, order.Status)
			}
		} else {
			fmt.Println("      没有订单")
		}
	} else {
		fmt.Printf("   ❌ 查询失败: %s (代码: %d)\n", resp.Message, resp.Code)
	}
}

// payOrder 支付订单
func payOrder(ctx context.Context, c orderservice.Client, orderNo string, userID int64) {
	if orderNo == "" {
		fmt.Println("   ⚠️  没有订单号，跳过支付")
		return
	}

	fmt.Printf("   支付订单: %s\n", orderNo)

	// 先查询订单状态
	getReq := &api.GetOrderReq{
		OrderNo: orderNo,
		UserId:  &userID,
	}

	getResp, err := c.GetOrder(ctx, getReq)
	if err != nil || !getResp.Success || getResp.Order == nil {
		fmt.Println("   ⚠️  无法获取订单信息，跳过支付")
		return
	}

	order := getResp.Order

	// 检查订单状态
	if order.Status == api.OrderStatus_PAID {
		fmt.Println("   ⚠️  订单已支付")
		return
	}

	if order.Status != api.OrderStatus_PENDING {
		fmt.Printf("   ⚠️  订单状态为 %v，无法支付\n", order.Status)
		return
	}

	// 执行支付
	payReq := &api.PayOrderReq{
		OrderNo: orderNo,
		UserId:  userID,
	}

	paymentNo := fmt.Sprintf("PAY%d", time.Now().Unix())
	payReq.PaymentNo = &paymentNo

	payResp, err := c.PayOrder(ctx, payReq)
	if err != nil {
		fmt.Printf("   ❌ 支付失败: %v\n", err)
		return
	}

	if payResp.Success {
		fmt.Println("   ✅ 支付成功!")
		fmt.Printf("      新状态: %v\n", payResp.NewStatus_)
	} else {
		fmt.Printf("   ❌ 支付失败: %s (代码: %d)\n", payResp.Message, payResp.Code)
	}
}

// orderStats 订单统计
func orderStats(ctx context.Context, c orderservice.Client, userID int64) {
	fmt.Printf("   用户 %d 的订单统计\n", userID)

	req := &api.OrderStatsReq{
		UserId: userID,
	}

	// 设置时间范围（最近30天）
	endTime := time.Now().Unix()
	startTime := time.Now().AddDate(0, 0, -30).Unix()
	req.StartTime = &startTime
	req.EndTime = &endTime

	resp, err := c.GetOrderStats(ctx, req)
	if err != nil {
		fmt.Printf("   ❌ 统计失败: %v\n", err)
		return
	}

	if resp.Success {
		fmt.Println("   ✅ 统计成功!")
		fmt.Printf("      总订单数: %d\n", resp.TotalOrders)
		fmt.Printf("      总金额: ¥%.2f\n", resp.TotalAmount)

		if len(resp.StatusCounts) > 0 {
			fmt.Println("      各状态订单数:")
			for status, count := range resp.StatusCounts {
				fmt.Printf("        %s: %d\n", status, count)
			}
		}
	} else {
		fmt.Printf("   ❌ 统计失败: %s (代码: %d)\n", resp.Message, resp.Code)
	}
}
