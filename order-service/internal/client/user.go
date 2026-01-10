package client

import (
	"context"
	"ecommerce/order-service/internal/dao/interfaces"
	"ecommerce/order-service/kitex_gen/api"
	"ecommerce/order-service/kitex_gen/api/userservice"

	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/klog"
)

type UserClient struct {
	client userservice.Client
}

func NewUserClient(addr string) (*UserClient, error) {
	c, err := userservice.NewClient("user-service",
		client.WithHostPorts(addr),
	)
	if err != nil {
		return nil, err
	}
	return &UserClient{client: c}, nil
}

func (uc *UserClient) GetUserInfo(ctx context.Context, userID int64) (*interfaces.UserInfo, error) {
	req := &api.GetUserProfileReq{
		Id: userID,
	}

	resp, err := uc.client.GetUserProfile(ctx, req)
	if err != nil {
		klog.Errorf("GetUserInfo failed: %v", err)
		return nil, err
	}

	if !resp.Success || resp.User == nil {
		return nil, nil
	}

	userInfo := &interfaces.UserInfo{
		ID:     resp.User.Id,
		Name:   resp.User.Name,
		Email:  resp.User.Email,
		Phone:  resp.User.Phone,
		Status: int32(resp.User.Status),
	}

	// 处理可选字段
	if resp.User.Avatar != nil {
		userInfo.Avatar = *resp.User.Avatar
	}
	if resp.User.LastLogin != nil {
		userInfo.LastLogin = *resp.User.LastLogin
	}

	return userInfo, nil
}

func (uc *UserClient) BatchGetUsers(ctx context.Context, userIDs []int64) (map[int64]*interfaces.UserInfo, error) {
	// 由于Thrift定义没有批量查询方法，这里循环调用单个查询
	result := make(map[int64]*interfaces.UserInfo)
	for _, id := range userIDs {
		userInfo, err := uc.GetUserInfo(ctx, id)
		if err == nil && userInfo != nil {
			result[id] = userInfo
		}
	}
	return result, nil
}

func (uc *UserClient) ValidateUser(ctx context.Context, userID int64) (bool, error) {
	userInfo, err := uc.GetUserInfo(ctx, userID)
	if err != nil {
		return false, err
	}
	return userInfo != nil && userInfo.Status == 1, nil // ACTIVE状态
}
