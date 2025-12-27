package main

import (
	"context"
	"fmt"
	"log"

	"ecommerce/product-service/kitex_gen/api"
	"ecommerce/product-service/kitex_gen/api/productservice"

	"github.com/cloudwego/kitex/client"
)

func main() {
	fmt.Println("🧐 Product Service RPC 客户端测试🧐")
	//创建客户端
	c, err := productservice.NewClient(
		"product.service",
		client.WithHostPorts("localhost:50051"),
	)
	if err != nil {
		log.Fatal("创建客户端失败:", err)
	}
	ctx := context.Background()
	//测试创建产品
	fmt.Println("\n1. 测试创建产品")
	createReq := &api.CreateProductReq{
		Name:     "圣园未花玩偶",
		Avatar:   "https://example.com/iphone.jpg",
		Category: "玩偶",
		Price:    300,
		Stock:    100,
	}
	brand := "崔尼蒂"
	createReq.Brand = &brand

	createResp, err := c.CreateProduct(ctx, createReq)
	if err != nil {
		log.Printf("创建产品失败: %v", err)
	} else {
		printResponse("创建产品", createResp)
	}
	//测试获取产品
	if createResp != nil && createResp.Success {
		fmt.Println("\n2. 测试获取产品")
		getReq := &api.GetProductReq{Id: createResp.Product.Id}
		getResp, err := c.GetProduct(ctx, getReq)
		if err != nil {
			log.Printf("获取产品失败: %v", err)
		} else {
			printResponse("获取产品", getResp)
		}
	}
	//测试搜索产品
	fmt.Println("\n3. 测试搜索产品")
	searchReq := &api.UserSearchProductsReq{
		Category: &createReq.Category,
		Page:     1,
		PageSize: 10,
	}
	searchResp, err := c.UserSearchProducts(ctx, searchReq)
	if err != nil {
		log.Printf("搜索产品失败: %v", err)
	} else {
		printResponse("搜索产品", searchResp)
		if searchResp.Success {
			fmt.Printf("   找到 %d 个产品:\n", len(searchResp.Products))
			for i, p := range searchResp.Products {
				fmt.Printf("   %d. %s - ¥%.2f (库存: %d)\n", i+1, p.Name, p.Price, p.Stock)
			}
		}
	}
	//测试上架产品
	if createResp != nil && createResp.Success {
		fmt.Println("\n4. 测试上架产品")
		onlineReq := &api.OnlineProductReq{Id: createResp.Product.Id}
		onlineResp, err := c.OnlineProduct(ctx, onlineReq)
		if err != nil {
			log.Printf("上架产品失败: %v", err)
		} else {
			printResponse("上架产品", onlineResp)
		}
	}
	//测试管理员搜索
	fmt.Println("\n5. 测试管理员搜索")
	adminSearchReq := &api.AdminSearchProductsReq{
		Page:     1,
		PageSize: 10,
	}
	adminSearchResp, err := c.AdminSearchProducts(ctx, adminSearchReq)
	if err != nil {
		log.Printf("管理员搜索失败: %v", err)
	} else {
		printResponse("管理员搜索", adminSearchResp)
	}

	fmt.Println("\n=== 测试完成 ===")
}

// 打印响应结果
func printResponse(operation string, resp interface{}) {
	switch r := resp.(type) {
	case *api.CreateProductResp:
		if r.Success {
			fmt.Printf("✅ %s成功: %s (ID: %d)\n", operation, r.Message, r.Product.Id)
		} else {
			fmt.Printf("❌ %s失败: %s (代码: %d)\n", operation, r.Message, r.Code)
		}
	case *api.GetProductResp:
		if r.Success {
			fmt.Printf("✅ %s成功: %s\n", operation, r.Message)
		} else {
			fmt.Printf("❌ %s失败: %s (代码: %d)\n", operation, r.Message, r.Code)
		}
	case *api.UserSearchProductsResp:
		if r.Success {
			fmt.Printf("✅ %s成功: 找到 %d 个产品\n", operation, r.Total)
		} else {
			fmt.Printf("❌ %s失败: %s (代码: %d)\n", operation, r.Message, r.Code)
		}
	case *api.OnlineProductResp:
		if r.Success {
			fmt.Printf("✅ %s成功: 状态从 %d 改为 %d\n", operation, r.OldStatus, r.NewStatus_)
		} else {
			fmt.Printf("❌ %s失败: %s (代码: %d)\n", operation, r.Message, r.Code)
		}
	case *api.AdminSearchProductsResp:
		if r.Success {
			fmt.Printf("✅ %s成功: 找到 %d 个产品\n", operation, r.Total)
		} else {
			fmt.Printf("❌ %s失败: %s (代码: %d)\n", operation, r.Message, r.Code)
		}
	default:
		fmt.Printf("未知响应类型%T\n", resp)
	}
}
