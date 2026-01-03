package main

import (
	"context"
	"fmt"
	"log"

	"ecommerce/user-service/kitex_gen/api"
	"ecommerce/user-service/kitex_gen/api/userservice"

	"github.com/bytedance/gopkg/cloud/metainfo"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/endpoint"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/transmeta"
	"github.com/cloudwego/kitex/transport"
)

func main() {
	fmt.Println("🧐 用户服务 RPC 客户端测试 🧐")

	//创建客户端 - 使用 TTHeader 传输协议
	c, err := userservice.NewClient(
		"user.service",
		client.WithHostPorts("localhost:50052"),
		client.WithTransportProtocol(transport.TTHeader),        // 使用 TTHeader 协议
		client.WithMetaHandler(transmeta.ClientTTHeaderHandler), // 添加 TTHeader 处理器
		client.WithMiddleware(clientMiddleware),                 // 添加客户端中间件
	)
	if err != nil {
		log.Fatal("创建客户端失败:", err)
	}

	ctx := context.Background()

	//注册用户（不需要token）
	fmt.Println("\n1. 测试用户注册")
	registerReq := &api.RegisterReq{
		Name:     "测试用户",
		Email:    "test@example.com",
		Password: "test123456",
		Phone:    "13800138000",
	}

	registerResp, err := c.Register(ctx, registerReq)
	if err != nil {
		log.Printf("用户注册失败: %v", err)
	} else {
		if registerResp.Success {
			fmt.Printf("✅ 注册成功: ID=%d\n", registerResp.Id)
		} else {
			fmt.Printf("❌ 注册失败: %s (代码: %d)\n",
				getStringValue(registerResp.Message), registerResp.Code)
		}
	}

	//登录（不需要token）
	fmt.Println("\n2. 测试用户登录")
	loginReq := &api.LoginReq{
		Phone:    "13800138000",
		Password: "test123456",
	}

	loginResp, err := c.Login(ctx, loginReq)
	if err != nil {
		log.Printf("用户登录失败: %v", err)
	} else {
		if loginResp.Success {
			fmt.Printf("✅ 登录成功: ID=%d, Token=%s...\n",
				loginResp.Id, safeSubstring(loginResp.Token, 0, 20))

			//保存token用于后续需要认证的请求
			token := loginResp.Token
			userID := loginResp.Id

			//使用新的方式传递认证信息
			testWithNewMethod(c, ctx, userID, token)
		} else {
			fmt.Printf("❌ 登录失败: %s (代码: %d)\n",
				getStringValue(loginResp.Message), loginResp.Code)
		}
	}

	fmt.Println("\n=== 测试完成 ===")
}

// 客户端中间件：在发送请求前打印 metainfo
func clientMiddleware(next endpoint.Endpoint) endpoint.Endpoint {
	return func(ctx context.Context, req, resp interface{}) (err error) {
		//打印请求的 metainfo
		if val, ok := metainfo.GetValue(ctx, "authorization"); ok {
			klog.Infof("客户端中间件: authorization = %s...", safeSubstring(val, 0, 30))
		}
		if val, ok := metainfo.GetValue(ctx, "Authorization"); ok {
			klog.Infof("客户端中间件: Authorization = %s...", safeSubstring(val, 0, 30))
		}

		return next(ctx, req, resp)
	}
}

// 使用新的方式传递认证信息
func testWithNewMethod(c userservice.Client, ctx context.Context, userID int64, token string) {
	fmt.Println("\n--- 测试需要认证的接口 ---")

	//使用 metainfo 传递，并在请求结构体中同时设置 token
	fmt.Println("\n方法1: 使用 metainfo 传递 + 请求体中的 token")

	//同时设置 metainfo 和请求结构体中的 token
	authCtx := metainfo.WithPersistentValue(ctx, "Authorization", "Bearer "+token)

	//测试获取用户资料 - 同时在 metainfo 和请求结构体中设置 token
	fmt.Println("\n3. 测试获取用户资料")
	getProfileReq := &api.GetUserProfileReq{
		Id:    userID,
		Token: token, //在请求结构体中设置 token
	}

	getProfileResp, err := c.GetUserProfile(authCtx, getProfileReq)
	if err != nil {
		fmt.Printf("❌ 获取用户资料失败: %v\n", err)
	} else {
		if getProfileResp.Success && getProfileResp.User != nil {
			fmt.Printf("✅ 获取用户资料成功:\n")
			fmt.Printf("   ID: %d\n", getProfileResp.User.Id)
			fmt.Printf("   姓名: %s\n", getProfileResp.User.Name)
			if getProfileResp.User.Email != "" {
				fmt.Printf("   邮箱: %s\n", getProfileResp.User.Email)
			}
			if getProfileResp.User.Phone != "" {
				fmt.Printf("   手机号: %s\n", getProfileResp.User.Phone)
			}
			fmt.Printf("   状态: %s\n", userStatusToString(getProfileResp.User.Status))
		} else {
			fmt.Printf("❌ 获取用户资料失败: %s (代码: %d)\n",
				getStringValue(getProfileResp.Message), getProfileResp.Code)
		}
	}

	//测试其他需要认证的方法
	fmt.Println("\n4. 测试更新用户信息")
	updateUserReq := &api.UpdateUserReq{
		Name:  stringPtr("更新后的名字"),
		Token: token,
	}
	updateResp, err := c.UpdateUser(authCtx, updateUserReq)
	if err != nil {
		fmt.Printf("❌ 更新用户信息失败: %v\n", err)
	} else {
		if updateResp.Success {
			fmt.Printf("✅ 更新用户信息成功\n")
		} else {
			fmt.Printf("❌ 更新用户信息失败: %s (代码: %d)\n",
				getStringValue(updateResp.Message), updateResp.Code)
		}
	}

	//测试获取用户状态
	fmt.Println("\n5. 测试获取用户状态")
	getStatusReq := &api.GetUserStatusReq{
		UserId: userID,
		Token:  token,
	}
	getStatusResp, err := c.GetUserStatus(authCtx, getStatusReq)
	if err != nil {
		fmt.Printf("❌ 获取用户状态失败: %v\n", err)
	} else {
		if getStatusResp.Success {
			fmt.Printf("✅ 获取用户状态成功:\n")
			fmt.Printf("   状态: %s\n", userStatusToString(getStatusResp.Status))
			fmt.Printf("   是否被封禁: %v\n", getStatusResp.IsBanned)
			fmt.Printf("   是否被删除: %v\n", getStatusResp.IsDeleted)
		} else {
			fmt.Printf("❌ 获取用户状态失败: %s (代码: %d)\n",
				getStringValue(getStatusResp.Message), getStatusResp.Code)
		}
	}
}

// 辅助函数
func getStringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func stringPtr(s string) *string {
	return &s
}

func safeSubstring(s string, start, end int) string {
	if s == "" {
		return ""
	}
	if start < 0 {
		start = 0
	}
	if end > len(s) {
		end = len(s)
	}
	if start >= end {
		return ""
	}
	return s[start:end]
}

func int32Ptr(i int32) *int32 {
	return &i
}

func userStatusToString(status api.UserStatus) string {
	switch status {
	case api.UserStatus_BANNED:
		return "封禁"
	case api.UserStatus_ACTIVE:
		return "活跃"
	case api.UserStatus_POWER:
		return "管理员"
	case api.UserStatus_Deleted:
		return "已删除"
	default:
		return fmt.Sprintf("未知(%d)", status)
	}
}
