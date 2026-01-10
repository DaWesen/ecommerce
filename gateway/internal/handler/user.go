package handler

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"ecommerce/gateway/internal/client"
	"ecommerce/gateway/pkg/response"
	"ecommerce/user-service/kitex_gen/api"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

// 辅助函数：获取当前用户ID（从上下文中）
func getUserIDFromContext(ctx *app.RequestContext) (int64, error) {
	userIDStr := ctx.GetString("user_id")
	if userIDStr == "" {
		return 0, errors.New("用户未登录")
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		return 0, errors.New("用户ID格式错误")
	}
	return userID, nil
}

// 辅助函数：从头部获取Token
func getTokenFromHeader(ctx *app.RequestContext) string {
	authHeader := string(ctx.GetHeader("Authorization"))
	if authHeader == "" {
		return ""
	}

	// 去掉 "Bearer " 前缀
	return strings.TrimPrefix(authHeader, "Bearer ")
}

// 安全转换函数
func safeStringPtr(v *string) string {
	if v == nil {
		return ""
	}
	return *v
}

// CreateUser 用户注册
func CreateUser(clientManager *client.ClientManager) app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		var req api.RegisterReq
		if err := ctx.BindAndValidate(&req); err != nil {
			response.Error(ctx, 400, "参数错误: "+err.Error())
			return
		}

		resp, err := clientManager.UserClient.Register(c, &req)
		if err != nil {
			hlog.Errorf("注册失败: %v", err)
			response.Error(ctx, 500, "注册失败: "+err.Error())
			return
		}

		if !resp.Success {
			response.Error(ctx, int(resp.Code), safeStringPtr(resp.Message))
			return
		}

		response.Success(ctx, map[string]interface{}{
			"user_id": resp.Id,
			"token":   resp.Token,
			"message": "注册成功",
		})
	}
}

// UserLogin 用户登录
func UserLogin(clientManager *client.ClientManager) app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		var req api.LoginReq
		if err := ctx.BindAndValidate(&req); err != nil {
			response.Error(ctx, 400, "参数错误: "+err.Error())
			return
		}

		resp, err := clientManager.UserClient.Login(c, &req)
		if err != nil {
			hlog.Errorf("登录失败: %v", err)
			response.Error(ctx, 500, "登录失败: "+err.Error())
			return
		}

		if !resp.Success {
			response.Error(ctx, int(resp.Code), safeStringPtr(resp.Message))
			return
		}

		response.Success(ctx, map[string]interface{}{
			"user_id": resp.Id,
			"token":   resp.Token,
			"message": "登录成功",
		})
	}
}

// GetUserProfile 获取用户信息
func GetUserProfile(clientManager *client.ClientManager) app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		idStr := ctx.Param("id")
		userID, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			response.Error(ctx, 400, "用户ID格式错误")
			return
		}

		token := getTokenFromHeader(ctx)
		if token == "" {
			response.Error(ctx, 401, "缺少认证token")
			return
		}

		req := &api.GetUserProfileReq{
			Id:    userID,
			Token: token,
		}

		resp, err := clientManager.UserClient.GetUserProfile(c, req)
		if err != nil {
			hlog.Errorf("获取用户信息失败: %v", err)
			response.Error(ctx, 500, "获取用户信息失败: "+err.Error())
			return
		}

		if !resp.Success {
			response.Error(ctx, int(resp.Code), safeStringPtr(resp.Message))
			return
		}

		response.Success(ctx, resp.User)
	}
}

// UpdateUserProfile 更新用户信息
func UpdateUserProfile(clientManager *client.ClientManager) app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		userID, err := getUserIDFromContext(ctx)
		if err != nil {
			response.Error(ctx, 401, "用户未登录")
			return
		}

		var req api.UpdateUserReq
		if err := ctx.BindAndValidate(&req); err != nil {
			response.Error(ctx, 400, "参数错误: "+err.Error())
			return
		}

		// 使用管理员的更新接口，因为普通用户接口不支持通过ID更新
		adminReq := &api.UpdateUserProfileReq{
			UserId: userID,
			Name:   req.Name,
			Avatar: req.Avatar,
			Bio:    req.Bio,
			Gender: req.Gender,
			Token:  req.Token,
		}

		resp, err := clientManager.UserClient.AdminUpdateUserProfile(c, adminReq)
		if err != nil {
			hlog.Errorf("更新用户信息失败: %v", err)
			response.Error(ctx, 500, "更新用户信息失败: "+err.Error())
			return
		}

		if !resp.Success {
			response.Error(ctx, int(resp.Code), safeStringPtr(resp.Message))
			return
		}

		response.Success(ctx, map[string]interface{}{
			"user":    resp.User,
			"message": "更新成功",
		})
	}
}

// ChangePassword 修改密码
func ChangePassword(clientManager *client.ClientManager) app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		userID, err := getUserIDFromContext(ctx)
		if err != nil {
			response.Error(ctx, 401, "用户未登录")
			return
		}

		var req api.ChangePasswordReq
		if err := ctx.BindAndValidate(&req); err != nil {
			response.Error(ctx, 400, "参数错误: "+err.Error())
			return
		}

		// 使用管理员的更新接口
		adminReq := &api.UpdatePasswordReq{
			UserId:       userID,
			NewPassword_: req.NewPassword_,
			Token:        req.Token,
		}

		resp, err := clientManager.UserClient.AdminUpdatePassword(c, adminReq)
		if err != nil {
			hlog.Errorf("修改密码失败: %v", err)
			response.Error(ctx, 500, "修改密码失败: "+err.Error())
			return
		}

		if !resp.Success {
			response.Error(ctx, int(resp.Code), safeStringPtr(resp.Message))
			return
		}

		response.Success(ctx, map[string]interface{}{
			"message": "密码修改成功",
		})
	}
}

// GetCurrentUser 获取当前登录用户信息
func GetCurrentUser(clientManager *client.ClientManager) app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		userID, err := getUserIDFromContext(ctx)
		if err != nil {
			response.Error(ctx, 401, "用户未登录")
			return
		}

		token := getTokenFromHeader(ctx)
		if token == "" {
			response.Error(ctx, 401, "缺少认证token")
			return
		}

		req := &api.GetUserProfileReq{
			Id:    userID,
			Token: token,
		}

		resp, err := clientManager.UserClient.GetUserProfile(c, req)
		if err != nil {
			hlog.Errorf("获取用户信息失败: %v", err)
			response.Error(ctx, 500, "获取用户信息失败: "+err.Error())
			return
		}

		if !resp.Success {
			response.Error(ctx, int(resp.Code), safeStringPtr(resp.Message))
			return
		}

		response.Success(ctx, resp.User)
	}
}

// UserLogout 用户登出
func UserLogout(clientManager *client.ClientManager) app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		token := getTokenFromHeader(ctx)
		if token == "" {
			response.Error(ctx, 401, "缺少认证token")
			return
		}

		req := &api.LogoutReq{
			Token: token,
		}

		resp, err := clientManager.UserClient.Logout(c, req)
		if err != nil {
			hlog.Errorf("登出失败: %v", err)
			response.Error(ctx, 500, "登出失败: "+err.Error())
			return
		}

		if !resp.Success {
			response.Error(ctx, int(resp.Code), safeStringPtr(resp.Message))
			return
		}

		response.Success(ctx, map[string]interface{}{
			"message": "登出成功",
		})
	}
}

// ListUsers 用户列表（管理员）
func ListUsers(clientManager *client.ClientManager) app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		token := getTokenFromHeader(ctx)
		if token == "" {
			response.Error(ctx, 401, "缺少认证token")
			return
		}

		page, _ := strconv.Atoi(string(ctx.Query("page")))
		pageSize, _ := strconv.Atoi(string(ctx.Query("page_size")))

		if page <= 0 {
			page = 1
		}
		if pageSize <= 0 {
			pageSize = 10
		}
		if pageSize > 100 {
			pageSize = 100
		}

		// 解析可选参数
		var statusPtr *api.UserStatus
		statusStr := string(ctx.Query("status"))
		if statusStr != "" {
			if status, err := strconv.ParseInt(statusStr, 10, 64); err == nil {
				statusVal := api.UserStatus(status)
				statusPtr = &statusVal
			}
		}

		keyword := string(ctx.Query("keyword"))
		var keywordPtr *string
		if keyword != "" {
			keywordPtr = &keyword
		}

		req := &api.ListUsersReq{
			Page:     int32(page),
			PageSize: int32(pageSize),
			Status:   statusPtr,
			Keyword:  keywordPtr,
			Token:    token,
		}

		resp, err := clientManager.UserClient.ListUsers(c, req)
		if err != nil {
			hlog.Errorf("获取用户列表失败: %v", err)
			response.Error(ctx, 500, "获取用户列表失败: "+err.Error())
			return
		}

		response.SuccessWithPagination(ctx, resp.Users, int64(resp.Total), page, pageSize)
	}
}

// SearchUsers 搜索用户（管理员）
func SearchUsers(clientManager *client.ClientManager) app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		token := getTokenFromHeader(ctx)
		if token == "" {
			response.Error(ctx, 401, "缺少认证token")
			return
		}

		keyword := string(ctx.Query("keyword"))
		if keyword == "" {
			response.Error(ctx, 400, "搜索关键词不能为空")
			return
		}

		page, _ := strconv.Atoi(string(ctx.Query("page")))
		pageSize, _ := strconv.Atoi(string(ctx.Query("page_size")))

		if page <= 0 {
			page = 1
		}
		if pageSize <= 0 {
			pageSize = 10
		}
		if pageSize > 100 {
			pageSize = 100
		}

		req := &api.SearchUsersReq{
			Keyword:  keyword,
			Page:     int32(page),
			PageSize: int32(pageSize),
			Token:    token,
		}

		resp, err := clientManager.UserClient.SearchUsers(c, req)
		if err != nil {
			hlog.Errorf("搜索用户失败: %v", err)
			response.Error(ctx, 500, "搜索用户失败: "+err.Error())
			return
		}

		response.SuccessWithPagination(ctx, resp.Users, int64(resp.Total), page, pageSize)
	}
}

// UpdateUserStatus 更新用户状态（管理员）
func UpdateUserStatus(clientManager *client.ClientManager) app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		token := getTokenFromHeader(ctx)
		if token == "" {
			response.Error(ctx, 401, "缺少认证token")
			return
		}

		userID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
		if err != nil {
			response.Error(ctx, 400, "用户ID格式错误")
			return
		}

		var req api.UpdateUserStatusReq
		if err := ctx.BindAndValidate(&req); err != nil {
			response.Error(ctx, 400, "参数错误: "+err.Error())
			return
		}

		req.UserId = userID
		req.Token = token

		resp, err := clientManager.UserClient.UpdateUserStatus(c, &req)
		if err != nil {
			hlog.Errorf("更新用户状态失败: %v", err)
			response.Error(ctx, 500, "更新用户状态失败: "+err.Error())
			return
		}

		if !resp.Success {
			response.Error(ctx, int(resp.Code), safeStringPtr(resp.Message))
			return
		}

		response.Success(ctx, map[string]interface{}{
			"message": "用户状态更新成功",
		})
	}
}

// DeleteUser 删除用户（管理员）
func DeleteUser(clientManager *client.ClientManager) app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		token := getTokenFromHeader(ctx)
		if token == "" {
			response.Error(ctx, 401, "缺少认证token")
			return
		}

		userID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
		if err != nil {
			response.Error(ctx, 400, "用户ID格式错误")
			return
		}

		req := &api.DeleteUserReq{
			UserId: userID,
			Token:  token,
		}

		resp, err := clientManager.UserClient.DeleteUser(c, req)
		if err != nil {
			hlog.Errorf("删除用户失败: %v", err)
			response.Error(ctx, 500, "删除用户失败: "+err.Error())
			return
		}

		if !resp.Success {
			response.Error(ctx, int(resp.Code), safeStringPtr(resp.Message))
			return
		}

		response.Success(ctx, map[string]interface{}{
			"message": "用户删除成功",
		})
	}
}

// BanUser 封禁用户（管理员）
func BanUser(clientManager *client.ClientManager) app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		token := getTokenFromHeader(ctx)
		if token == "" {
			response.Error(ctx, 401, "缺少认证token")
			return
		}

		userID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
		if err != nil {
			response.Error(ctx, 400, "用户ID格式错误")
			return
		}

		var req api.BanUserReq
		if err := ctx.BindAndValidate(&req); err != nil {
			response.Error(ctx, 400, "参数错误: "+err.Error())
			return
		}

		req.UserId = userID
		req.Token = token

		resp, err := clientManager.UserClient.BanUser(c, &req)
		if err != nil {
			hlog.Errorf("封禁用户失败: %v", err)
			response.Error(ctx, 500, "封禁用户失败: "+err.Error())
			return
		}

		if !resp.Success {
			response.Error(ctx, int(resp.Code), safeStringPtr(resp.Message))
			return
		}

		response.Success(ctx, map[string]interface{}{
			"message": "用户封禁成功",
		})
	}
}

// UnbanUser 解封用户（管理员）
func UnbanUser(clientManager *client.ClientManager) app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		token := getTokenFromHeader(ctx)
		if token == "" {
			response.Error(ctx, 401, "缺少认证token")
			return
		}

		userID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
		if err != nil {
			response.Error(ctx, 400, "用户ID格式错误")
			return
		}

		req := &api.UnbanUserReq{
			UserId: userID,
			Token:  token,
		}

		resp, err := clientManager.UserClient.UnbanUser(c, req)
		if err != nil {
			hlog.Errorf("解封用户失败: %v", err)
			response.Error(ctx, 500, "解封用户失败: "+err.Error())
			return
		}

		if !resp.Success {
			response.Error(ctx, int(resp.Code), safeStringPtr(resp.Message))
			return
		}

		response.Success(ctx, map[string]interface{}{
			"message": "用户解封成功",
		})
	}
}
