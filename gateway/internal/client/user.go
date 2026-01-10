package client

import (
	"context"
	"ecommerce/user-service/kitex_gen/api"
)

// Register 用户注册
func (uc *UserClient) Register(ctx context.Context, req *api.RegisterReq) (*api.RegisterResp, error) {
	return uc.client.Register(ctx, req)
}

// Login 用户登录
func (uc *UserClient) Login(ctx context.Context, req *api.LoginReq) (*api.LoginResp, error) {
	return uc.client.Login(ctx, req)
}

// UpdateUser 更新用户信息
func (uc *UserClient) UpdateUser(ctx context.Context, req *api.UpdateUserReq) (*api.UpdateUserResp, error) {
	return uc.client.UpdateUser(ctx, req)
}

// ChangePassword 修改密码
func (uc *UserClient) ChangePassword(ctx context.Context, req *api.ChangePasswordReq) (*api.UpdateUserResp, error) {
	return uc.client.ChangePassword(ctx, req)
}

// ChangeEmail 修改邮箱
func (uc *UserClient) ChangeEmail(ctx context.Context, req *api.ChangeEmailReq) (*api.UpdateUserResp, error) {
	return uc.client.ChangeEmail(ctx, req)
}

// ChangePhone 修改手机
func (uc *UserClient) ChangePhone(ctx context.Context, req *api.ChangePhoneReq) (*api.UpdateUserResp, error) {
	return uc.client.ChangePhone(ctx, req)
}

// GetUserProfile 获取用户信息
func (uc *UserClient) GetUserProfile(ctx context.Context, req *api.GetUserProfileReq) (*api.GetUserProfileResp, error) {
	return uc.client.GetUserProfile(ctx, req)
}

// Logout 用户登出
func (uc *UserClient) Logout(ctx context.Context, req *api.LogoutReq) (*api.LogoutResp, error) {
	return uc.client.Logout(ctx, req)
}

// GetUserStatus 获取用户状态
func (uc *UserClient) GetUserStatus(ctx context.Context, req *api.GetUserStatusReq) (*api.GetUserStatusResp, error) {
	return uc.client.GetUserStatus(ctx, req)
}

// BanUser 封禁用户（管理员）
func (uc *UserClient) BanUser(ctx context.Context, req *api.BanUserReq) (*api.BanUserResp, error) {
	return uc.client.BanUser(ctx, req)
}

// UnbanUser 解封用户（管理员）
func (uc *UserClient) UnbanUser(ctx context.Context, req *api.UnbanUserReq) (*api.UnbanUserResp, error) {
	return uc.client.UnbanUser(ctx, req)
}

// DeleteUser 删除用户（管理员）
func (uc *UserClient) DeleteUser(ctx context.Context, req *api.DeleteUserReq) (*api.DeleteUserResp, error) {
	return uc.client.DeleteUser(ctx, req)
}

// RestoreUser 恢复用户（管理员）
func (uc *UserClient) RestoreUser(ctx context.Context, req *api.RestoreUserReq) (*api.RestoreUserResp, error) {
	return uc.client.RestoreUser(ctx, req)
}

// UpdateUserStatus 更新用户状态（管理员）
func (uc *UserClient) UpdateUserStatus(ctx context.Context, req *api.UpdateUserStatusReq) (*api.UpdateUserStatusResp, error) {
	return uc.client.UpdateUserStatus(ctx, req)
}

// ListUsers 用户列表（管理员）
func (uc *UserClient) ListUsers(ctx context.Context, req *api.ListUsersReq) (*api.ListUsersResp, error) {
	return uc.client.ListUsers(ctx, req)
}

// SearchUsers 搜索用户（管理员）
func (uc *UserClient) SearchUsers(ctx context.Context, req *api.SearchUsersReq) (*api.SearchUsersResp, error) {
	return uc.client.SearchUsers(ctx, req)
}

// CountUsers 统计用户数量
func (uc *UserClient) CountUsers(ctx context.Context, req *api.CountUsersReq) (*api.CountUsersResp, error) {
	return uc.client.CountUsers(ctx, req)
}

// CountByStatus 按状态统计用户
func (uc *UserClient) CountByStatus(ctx context.Context, req *api.CountByStatusReq) (*api.CountByStatusResp, error) {
	return uc.client.CountByStatus(ctx, req)
}

// AdminUpdatePassword 管理员修改密码
func (uc *UserClient) AdminUpdatePassword(ctx context.Context, req *api.UpdatePasswordReq) (*api.UpdatePasswordResp, error) {
	return uc.client.AdminUpdatePassword(ctx, req)
}

// AdminUpdateEmail 管理员修改邮箱
func (uc *UserClient) AdminUpdateEmail(ctx context.Context, req *api.UpdateEmailReq) (*api.UpdateEmailResp, error) {
	return uc.client.AdminUpdateEmail(ctx, req)
}

// AdminUpdatePhone 管理员修改手机
func (uc *UserClient) AdminUpdatePhone(ctx context.Context, req *api.UpdatePhoneReq) (*api.UpdatePhoneResp, error) {
	return uc.client.AdminUpdatePhone(ctx, req)
}

// AdminUpdateUserProfile 管理员更新用户信息
func (uc *UserClient) AdminUpdateUserProfile(ctx context.Context, req *api.UpdateUserProfileReq) (*api.UpdateUserProfileResp, error) {
	return uc.client.AdminUpdateUserProfile(ctx, req)
}
