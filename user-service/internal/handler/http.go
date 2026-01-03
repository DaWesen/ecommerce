package handler

import (
	"context"
	"strconv"

	"ecommerce/user-service/internal/service"
	api "ecommerce/user-service/kitex_gen/api"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

type UserHTTPHandler struct {
	userService service.UserService
}

func NewUserHTTPHandler(userService service.UserService) *UserHTTPHandler {
	return &UserHTTPHandler{
		userService: userService,
	}
}

// HTTP响应包装器
type HTTPResponse struct {
	Code    int32       `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Success bool        `json:"success"`
}

func stringPtr(s string) *string {
	return &s
}

func getStringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// 注册用户
func (h *UserHTTPHandler) Register(ctx context.Context, c *app.RequestContext) {
	var req api.RegisterReq
	if err := c.BindAndValidate(&req); err != nil {
		hlog.CtxErrorf(ctx, "参数绑定失败: %v", err)
		c.JSON(consts.StatusBadRequest, HTTPResponse{
			Success: false,
			Code:    400,
			Message: "参数错误",
			Error:   err.Error(),
		})
		return
	}

	resp, err := h.userService.Register(ctx, &req)
	if err != nil {
		hlog.CtxErrorf(ctx, "创建用户失败: %v", err)
		c.JSON(consts.StatusInternalServerError, HTTPResponse{
			Success: false,
			Code:    500,
			Message: "创建用户失败",
			Error:   err.Error(),
		})
		return
	}

	if !resp.Success {
		c.JSON(consts.StatusBadRequest, HTTPResponse{
			Success: false,
			Code:    resp.Code,
			Message: getStringValue(resp.Message),
		})
		return
	}

	c.JSON(consts.StatusOK, HTTPResponse{
		Success: true,
		Code:    0,
		Message: getStringValue(resp.Message),
		Data: map[string]interface{}{
			"id":    resp.Id,
			"token": resp.Token,
		},
	})
}

// 用户登录
func (h *UserHTTPHandler) Login(ctx context.Context, c *app.RequestContext) {
	var req api.LoginReq
	if err := c.BindAndValidate(&req); err != nil {
		hlog.CtxErrorf(ctx, "参数绑定失败: %v", err)
		c.JSON(consts.StatusBadRequest, HTTPResponse{
			Success: false,
			Code:    400,
			Message: "参数错误",
			Error:   err.Error(),
		})
		return
	}

	resp, err := h.userService.Login(ctx, &req)
	if err != nil {
		hlog.CtxErrorf(ctx, "用户登录失败: %v", err)
		c.JSON(consts.StatusInternalServerError, HTTPResponse{
			Success: false,
			Code:    500,
			Message: "登录失败",
			Error:   err.Error(),
		})
		return
	}

	if !resp.Success {
		status := consts.StatusBadRequest
		if resp.Code == 404 {
			status = consts.StatusNotFound
		} else if resp.Code == 401 {
			status = consts.StatusUnauthorized
		} else if resp.Code == 403 {
			status = consts.StatusForbidden
		}
		c.JSON(status, HTTPResponse{
			Success: false,
			Code:    resp.Code,
			Message: getStringValue(resp.Message),
		})
		return
	}

	c.JSON(consts.StatusOK, HTTPResponse{
		Success: true,
		Code:    0,
		Message: getStringValue(resp.Message),
		Data: map[string]interface{}{
			"id":    resp.Id,
			"token": resp.Token,
		},
	})
}

// 更新用户信息
func (h *UserHTTPHandler) UpdateUser(ctx context.Context, c *app.RequestContext) {
	var req api.UpdateUserReq
	if err := c.BindAndValidate(&req); err != nil {
		hlog.CtxErrorf(ctx, "参数绑定失败: %v", err)
		c.JSON(consts.StatusBadRequest, HTTPResponse{
			Success: false,
			Code:    400,
			Message: "参数错误",
			Error:   err.Error(),
		})
		return
	}

	resp, err := h.userService.UpdateUser(ctx, &req)
	if err != nil {
		hlog.CtxErrorf(ctx, "更新用户信息失败: %v", err)
		c.JSON(consts.StatusInternalServerError, HTTPResponse{
			Success: false,
			Code:    500,
			Message: "更新失败",
			Error:   err.Error(),
		})
		return
	}

	if !resp.Success {
		c.JSON(consts.StatusBadRequest, HTTPResponse{
			Success: false,
			Code:    resp.Code,
			Message: getStringValue(resp.Message),
		})
		return
	}

	c.JSON(consts.StatusOK, HTTPResponse{
		Success: true,
		Code:    0,
		Message: getStringValue(resp.Message),
	})
}

// 修改密码
func (h *UserHTTPHandler) ChangePassword(ctx context.Context, c *app.RequestContext) {
	var req api.ChangePasswordReq
	if err := c.BindAndValidate(&req); err != nil {
		hlog.CtxErrorf(ctx, "参数绑定失败: %v", err)
		c.JSON(consts.StatusBadRequest, HTTPResponse{
			Success: false,
			Code:    400,
			Message: "参数错误",
			Error:   err.Error(),
		})
		return
	}

	resp, err := h.userService.ChangePassword(ctx, &req)
	if err != nil {
		hlog.CtxErrorf(ctx, "修改密码失败: %v", err)
		c.JSON(consts.StatusInternalServerError, HTTPResponse{
			Success: false,
			Code:    500,
			Message: "修改密码失败",
			Error:   err.Error(),
		})
		return
	}

	if !resp.Success {
		c.JSON(consts.StatusBadRequest, HTTPResponse{
			Success: false,
			Code:    resp.Code,
			Message: getStringValue(resp.Message),
		})
		return
	}

	c.JSON(consts.StatusOK, HTTPResponse{
		Success: true,
		Code:    0,
		Message: getStringValue(resp.Message),
	})
}

// 修改邮箱
func (h *UserHTTPHandler) ChangeEmail(ctx context.Context, c *app.RequestContext) {
	var req api.ChangeEmailReq
	if err := c.BindAndValidate(&req); err != nil {
		hlog.CtxErrorf(ctx, "参数绑定失败: %v", err)
		c.JSON(consts.StatusBadRequest, HTTPResponse{
			Success: false,
			Code:    400,
			Message: "参数错误",
			Error:   err.Error(),
		})
		return
	}

	resp, err := h.userService.ChangeEmail(ctx, &req)
	if err != nil {
		hlog.CtxErrorf(ctx, "修改邮箱失败: %v", err)
		c.JSON(consts.StatusInternalServerError, HTTPResponse{
			Success: false,
			Code:    500,
			Message: "修改邮箱失败",
			Error:   err.Error(),
		})
		return
	}

	if !resp.Success {
		c.JSON(consts.StatusBadRequest, HTTPResponse{
			Success: false,
			Code:    resp.Code,
			Message: getStringValue(resp.Message),
		})
		return
	}

	c.JSON(consts.StatusOK, HTTPResponse{
		Success: true,
		Code:    0,
		Message: getStringValue(resp.Message),
	})
}

// 修改手机号
func (h *UserHTTPHandler) ChangePhone(ctx context.Context, c *app.RequestContext) {
	var req api.ChangePhoneReq
	if err := c.BindAndValidate(&req); err != nil {
		hlog.CtxErrorf(ctx, "参数绑定失败: %v", err)
		c.JSON(consts.StatusBadRequest, HTTPResponse{
			Success: false,
			Code:    400,
			Message: "参数错误",
			Error:   err.Error(),
		})
		return
	}

	resp, err := h.userService.ChangePhone(ctx, &req)
	if err != nil {
		hlog.CtxErrorf(ctx, "修改手机号失败: %v", err)
		c.JSON(consts.StatusInternalServerError, HTTPResponse{
			Success: false,
			Code:    500,
			Message: "修改手机号失败",
			Error:   err.Error(),
		})
		return
	}

	if !resp.Success {
		c.JSON(consts.StatusBadRequest, HTTPResponse{
			Success: false,
			Code:    resp.Code,
			Message: getStringValue(resp.Message),
		})
		return
	}

	c.JSON(consts.StatusOK, HTTPResponse{
		Success: true,
		Code:    0,
		Message: getStringValue(resp.Message),
	})
}

// 获取用户资料
func (h *UserHTTPHandler) GetUserProfile(ctx context.Context, c *app.RequestContext) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(consts.StatusBadRequest, HTTPResponse{
			Success: false,
			Code:    400,
			Message: "ID格式错误",
			Error:   err.Error(),
		})
		return
	}

	req := &api.GetUserProfileReq{Id: id}
	resp, err := h.userService.GetUserProfile(ctx, req)
	if err != nil {
		hlog.CtxErrorf(ctx, "获取用户资料失败: %v", err)
		c.JSON(consts.StatusInternalServerError, HTTPResponse{
			Success: false,
			Code:    500,
			Message: "获取用户资料失败",
			Error:   err.Error(),
		})
		return
	}

	if !resp.Success {
		status := consts.StatusBadRequest
		if resp.Code == 404 {
			status = consts.StatusNotFound
		} else if resp.Code == 403 {
			status = consts.StatusForbidden
		}
		c.JSON(status, HTTPResponse{
			Success: false,
			Code:    resp.Code,
			Message: getStringValue(resp.Message),
		})
		return
	}

	c.JSON(consts.StatusOK, HTTPResponse{
		Success: true,
		Code:    0,
		Message: getStringValue(resp.Message),
		Data:    resp.User,
	})
}

// 用户登出
func (h *UserHTTPHandler) Logout(ctx context.Context, c *app.RequestContext) {
	req := &api.LogoutReq{}
	resp, err := h.userService.Logout(ctx, req)
	if err != nil {
		hlog.CtxErrorf(ctx, "用户登出失败: %v", err)
		c.JSON(consts.StatusInternalServerError, HTTPResponse{
			Success: false,
			Code:    500,
			Message: "登出失败",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(consts.StatusOK, HTTPResponse{
		Success: true,
		Code:    0,
		Message: getStringValue(resp.Message),
	})
}

// 获取用户状态
func (h *UserHTTPHandler) GetUserStatus(ctx context.Context, c *app.RequestContext) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(consts.StatusBadRequest, HTTPResponse{
			Success: false,
			Code:    400,
			Message: "ID格式错误",
			Error:   err.Error(),
		})
		return
	}

	req := &api.GetUserStatusReq{UserId: id}
	resp, err := h.userService.GetUserStatus(ctx, req)
	if err != nil {
		hlog.CtxErrorf(ctx, "获取用户状态失败: %v", err)
		c.JSON(consts.StatusInternalServerError, HTTPResponse{
			Success: false,
			Code:    500,
			Message: "获取用户状态失败",
			Error:   err.Error(),
		})
		return
	}

	if !resp.Success {
		status := consts.StatusBadRequest
		if resp.Code == 404 {
			status = consts.StatusNotFound
		} else if resp.Code == 403 {
			status = consts.StatusForbidden
		}
		c.JSON(status, HTTPResponse{
			Success: false,
			Code:    resp.Code,
			Message: getStringValue(resp.Message),
		})
		return
	}

	c.JSON(consts.StatusOK, HTTPResponse{
		Success: true,
		Code:    0,
		Message: getStringValue(resp.Message),
		Data: map[string]interface{}{
			"status":     resp.Status,
			"is_banned":  resp.IsBanned,
			"is_deleted": resp.IsDeleted,
		},
	})
}

// 管理员接口：封禁用户
func (h *UserHTTPHandler) BanUser(ctx context.Context, c *app.RequestContext) {
	var req api.BanUserReq
	if err := c.BindAndValidate(&req); err != nil {
		hlog.CtxErrorf(ctx, "参数绑定失败: %v", err)
		c.JSON(consts.StatusBadRequest, HTTPResponse{
			Success: false,
			Code:    400,
			Message: "参数错误",
			Error:   err.Error(),
		})
		return
	}

	resp, err := h.userService.BanUser(ctx, &req)
	if err != nil {
		hlog.CtxErrorf(ctx, "封禁用户失败: %v", err)
		c.JSON(consts.StatusInternalServerError, HTTPResponse{
			Success: false,
			Code:    500,
			Message: "封禁用户失败",
			Error:   err.Error(),
		})
		return
	}

	if !resp.Success {
		c.JSON(consts.StatusBadRequest, HTTPResponse{
			Success: false,
			Code:    resp.Code,
			Message: getStringValue(resp.Message),
		})
		return
	}

	c.JSON(consts.StatusOK, HTTPResponse{
		Success: true,
		Code:    0,
		Message: getStringValue(resp.Message),
		Data: map[string]interface{}{
			"banned_at":  resp.BannedAt,
			"ban_reason": resp.BanReason,
		},
	})
}

// 管理员接口：解封用户
func (h *UserHTTPHandler) UnbanUser(ctx context.Context, c *app.RequestContext) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(consts.StatusBadRequest, HTTPResponse{
			Success: false,
			Code:    400,
			Message: "ID格式错误",
			Error:   err.Error(),
		})
		return
	}

	req := &api.UnbanUserReq{UserId: id}
	resp, err := h.userService.UnbanUser(ctx, req)
	if err != nil {
		hlog.CtxErrorf(ctx, "解封用户失败: %v", err)
		c.JSON(consts.StatusInternalServerError, HTTPResponse{
			Success: false,
			Code:    500,
			Message: "解封用户失败",
			Error:   err.Error(),
		})
		return
	}

	if !resp.Success {
		c.JSON(consts.StatusBadRequest, HTTPResponse{
			Success: false,
			Code:    resp.Code,
			Message: getStringValue(resp.Message),
		})
		return
	}

	c.JSON(consts.StatusOK, HTTPResponse{
		Success: true,
		Code:    0,
		Message: getStringValue(resp.Message),
	})
}

// 管理员接口：删除用户
func (h *UserHTTPHandler) DeleteUser(ctx context.Context, c *app.RequestContext) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(consts.StatusBadRequest, HTTPResponse{
			Success: false,
			Code:    400,
			Message: "ID格式错误",
			Error:   err.Error(),
		})
		return
	}

	req := &api.DeleteUserReq{UserId: id}
	resp, err := h.userService.DeleteUser(ctx, req)
	if err != nil {
		hlog.CtxErrorf(ctx, "删除用户失败: %v", err)
		c.JSON(consts.StatusInternalServerError, HTTPResponse{
			Success: false,
			Code:    500,
			Message: "删除用户失败",
			Error:   err.Error(),
		})
		return
	}

	if !resp.Success {
		c.JSON(consts.StatusBadRequest, HTTPResponse{
			Success: false,
			Code:    resp.Code,
			Message: getStringValue(resp.Message),
		})
		return
	}

	c.JSON(consts.StatusOK, HTTPResponse{
		Success: true,
		Code:    0,
		Message: getStringValue(resp.Message),
		Data: map[string]interface{}{
			"deleted_at": resp.DeletedAt,
		},
	})
}

// 管理员接口：恢复用户
func (h *UserHTTPHandler) RestoreUser(ctx context.Context, c *app.RequestContext) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(consts.StatusBadRequest, HTTPResponse{
			Success: false,
			Code:    400,
			Message: "ID格式错误",
			Error:   err.Error(),
		})
		return
	}

	req := &api.RestoreUserReq{UserId: id}
	resp, err := h.userService.RestoreUser(ctx, req)
	if err != nil {
		hlog.CtxErrorf(ctx, "恢复用户失败: %v", err)
		c.JSON(consts.StatusInternalServerError, HTTPResponse{
			Success: false,
			Code:    500,
			Message: "恢复用户失败",
			Error:   err.Error(),
		})
		return
	}

	if !resp.Success {
		c.JSON(consts.StatusBadRequest, HTTPResponse{
			Success: false,
			Code:    resp.Code,
			Message: getStringValue(resp.Message),
		})
		return
	}

	c.JSON(consts.StatusOK, HTTPResponse{
		Success: true,
		Code:    0,
		Message: getStringValue(resp.Message),
	})
}

// 管理员接口：更新用户状态
func (h *UserHTTPHandler) UpdateUserStatus(ctx context.Context, c *app.RequestContext) {
	var req api.UpdateUserStatusReq
	if err := c.BindAndValidate(&req); err != nil {
		hlog.CtxErrorf(ctx, "参数绑定失败: %v", err)
		c.JSON(consts.StatusBadRequest, HTTPResponse{
			Success: false,
			Code:    400,
			Message: "参数错误",
			Error:   err.Error(),
		})
		return
	}

	resp, err := h.userService.UpdateUserStatus(ctx, &req)
	if err != nil {
		hlog.CtxErrorf(ctx, "更新用户状态失败: %v", err)
		c.JSON(consts.StatusInternalServerError, HTTPResponse{
			Success: false,
			Code:    500,
			Message: "更新用户状态失败",
			Error:   err.Error(),
		})
		return
	}

	if !resp.Success {
		c.JSON(consts.StatusBadRequest, HTTPResponse{
			Success: false,
			Code:    resp.Code,
			Message: getStringValue(resp.Message),
		})
		return
	}

	c.JSON(consts.StatusOK, HTTPResponse{
		Success: true,
		Code:    0,
		Message: getStringValue(resp.Message),
		Data: map[string]interface{}{
			"old_status": resp.OldStatus,
			"new_status": resp.NewStatus_,
			"updated_at": resp.UpdatedAt,
		},
	})
}

// 管理员接口：用户列表
func (h *UserHTTPHandler) ListUsers(ctx context.Context, c *app.RequestContext) {
	var req api.ListUsersReq

	//解析查询参数
	if page, err := strconv.Atoi(c.Query("page")); err == nil && page > 0 {
		req.Page = int32(page)
	} else {
		req.Page = 1
	}

	if pageSize, err := strconv.Atoi(c.Query("page_size")); err == nil && pageSize > 0 {
		if pageSize > 100 {
			pageSize = 100
		}
		req.PageSize = int32(pageSize)
	} else {
		req.PageSize = 20
	}

	if keyword := c.Query("keyword"); keyword != "" {
		req.Keyword = &keyword
	}

	if statusStr := c.Query("status"); statusStr != "" {
		if status, err := strconv.Atoi(statusStr); err == nil {
			statusVal := api.UserStatus(status)
			req.Status = &statusVal
		}
	}

	if minCreatedAtStr := c.Query("min_created_at"); minCreatedAtStr != "" {
		if minCreatedAt, err := strconv.ParseInt(minCreatedAtStr, 10, 64); err == nil {
			req.MinCreatedAt = &minCreatedAt
		}
	}

	if maxCreatedAtStr := c.Query("max_created_at"); maxCreatedAtStr != "" {
		if maxCreatedAt, err := strconv.ParseInt(maxCreatedAtStr, 10, 64); err == nil {
			req.MaxCreatedAt = &maxCreatedAt
		}
	}

	if orderBy := c.Query("order_by"); orderBy != "" {
		req.OrderBy = &orderBy
	}

	req.Desc = c.Query("desc") == "true"

	resp, err := h.userService.ListUsers(ctx, &req)
	if err != nil {
		hlog.CtxErrorf(ctx, "获取用户列表失败: %v", err)
		c.JSON(consts.StatusInternalServerError, HTTPResponse{
			Success: false,
			Code:    500,
			Message: "获取用户列表失败",
			Error:   err.Error(),
		})
		return
	}

	if !resp.Success {
		c.JSON(consts.StatusBadRequest, HTTPResponse{
			Success: false,
			Code:    resp.Code,
			Message: getStringValue(resp.Message),
		})
		return
	}

	c.JSON(consts.StatusOK, HTTPResponse{
		Success: true,
		Code:    0,
		Message: getStringValue(resp.Message),
		Data: utils.H{
			"total":     resp.Total,
			"page":      resp.Page,
			"page_size": resp.PageSize,
			"users":     resp.Users,
		},
	})
}

// 管理员接口：搜索用户
func (h *UserHTTPHandler) SearchUsers(ctx context.Context, c *app.RequestContext) {
	var req api.SearchUsersReq

	//解析查询参数
	if keyword := c.Query("keyword"); keyword != "" {
		req.Keyword = keyword
	}

	if page, err := strconv.Atoi(c.Query("page")); err == nil && page > 0 {
		req.Page = int32(page)
	} else {
		req.Page = 1
	}

	if pageSize, err := strconv.Atoi(c.Query("page_size")); err == nil && pageSize > 0 {
		if pageSize > 100 {
			pageSize = 100
		}
		req.PageSize = int32(pageSize)
	} else {
		req.PageSize = 20
	}

	resp, err := h.userService.SearchUsers(ctx, &req)
	if err != nil {
		hlog.CtxErrorf(ctx, "搜索用户失败: %v", err)
		c.JSON(consts.StatusInternalServerError, HTTPResponse{
			Success: false,
			Code:    500,
			Message: "搜索用户失败",
			Error:   err.Error(),
		})
		return
	}

	if !resp.Success {
		c.JSON(consts.StatusBadRequest, HTTPResponse{
			Success: false,
			Code:    resp.Code,
			Message: getStringValue(resp.Message),
		})
		return
	}

	c.JSON(consts.StatusOK, HTTPResponse{
		Success: true,
		Code:    0,
		Message: getStringValue(resp.Message),
		Data: utils.H{
			"total":     resp.Total,
			"page":      resp.Page,
			"page_size": resp.PageSize,
			"users":     resp.Users,
		},
	})
}

// 管理员接口：统计用户数量
func (h *UserHTTPHandler) CountUsers(ctx context.Context, c *app.RequestContext) {
	var req api.CountUsersReq

	if statusStr := c.Query("status"); statusStr != "" {
		if status, err := strconv.Atoi(statusStr); err == nil {
			statusVal := api.UserStatus(status)
			req.Status = &statusVal
		}
	}

	resp, err := h.userService.CountUsers(ctx, &req)
	if err != nil {
		hlog.CtxErrorf(ctx, "统计用户数量失败: %v", err)
		c.JSON(consts.StatusInternalServerError, HTTPResponse{
			Success: false,
			Code:    500,
			Message: "统计用户数量失败",
			Error:   err.Error(),
		})
		return
	}

	if !resp.Success {
		c.JSON(consts.StatusBadRequest, HTTPResponse{
			Success: false,
			Code:    resp.Code,
			Message: getStringValue(resp.Message),
		})
		return
	}

	c.JSON(consts.StatusOK, HTTPResponse{
		Success: true,
		Code:    0,
		Message: getStringValue(resp.Message),
		Data: map[string]interface{}{
			"count": resp.Count,
		},
	})
}

// 管理员接口：按状态统计用户
func (h *UserHTTPHandler) CountByStatus(ctx context.Context, c *app.RequestContext) {
	resp, err := h.userService.CountByStatus(ctx)
	if err != nil {
		hlog.CtxErrorf(ctx, "按状态统计用户失败: %v", err)
		c.JSON(consts.StatusInternalServerError, HTTPResponse{
			Success: false,
			Code:    500,
			Message: "按状态统计用户失败",
			Error:   err.Error(),
		})
		return
	}

	if !resp.Success {
		c.JSON(consts.StatusBadRequest, HTTPResponse{
			Success: false,
			Code:    resp.Code,
			Message: getStringValue(resp.Message),
		})
		return
	}

	c.JSON(consts.StatusOK, HTTPResponse{
		Success: true,
		Code:    0,
		Message: getStringValue(resp.Message),
		Data: map[string]interface{}{
			"counts": resp.Counts,
		},
	})
}
