package main

import (
	"context"
	api "ecommerce/user-service/kitex_gen/api"
)

// UserServiceImpl implements the last service interface defined in the IDL.
type UserServiceImpl struct{}

// Register implements the UserServiceImpl interface.
func (s *UserServiceImpl) Register(ctx context.Context, req *api.RegisterReq) (resp *api.RegisterResp, err error) {
	// TODO: Your code here...
	return
}

// Login implements the UserServiceImpl interface.
func (s *UserServiceImpl) Login(ctx context.Context, req *api.LoginReq) (resp *api.LoginResp, err error) {
	// TODO: Your code here...
	return
}

// UpdateUser implements the UserServiceImpl interface.
func (s *UserServiceImpl) UpdateUser(ctx context.Context, req *api.UpdateUserReq) (resp *api.UpdateUserResp, err error) {
	// TODO: Your code here...
	return
}

// ChangePassword implements the UserServiceImpl interface.
func (s *UserServiceImpl) ChangePassword(ctx context.Context, req *api.ChangePasswordReq) (resp *api.UpdateUserResp, err error) {
	// TODO: Your code here...
	return
}

// ChangeEmail implements the UserServiceImpl interface.
func (s *UserServiceImpl) ChangeEmail(ctx context.Context, req *api.ChangeEmailReq) (resp *api.UpdateUserResp, err error) {
	// TODO: Your code here...
	return
}

// ChangePhone implements the UserServiceImpl interface.
func (s *UserServiceImpl) ChangePhone(ctx context.Context, req *api.ChangePhoneReq) (resp *api.UpdateUserResp, err error) {
	// TODO: Your code here...
	return
}

// GetUserProfile implements the UserServiceImpl interface.
func (s *UserServiceImpl) GetUserProfile(ctx context.Context, req *api.GetUserProfileReq) (resp *api.GetUserProfileResp, err error) {
	// TODO: Your code here...
	return
}

// Logout implements the UserServiceImpl interface.
func (s *UserServiceImpl) Logout(ctx context.Context, req *api.LogoutReq) (resp *api.LogoutResp, err error) {
	// TODO: Your code here...
	return
}

// GetUserStatus implements the UserServiceImpl interface.
func (s *UserServiceImpl) GetUserStatus(ctx context.Context, req *api.GetUserStatusReq) (resp *api.GetUserStatusResp, err error) {
	// TODO: Your code here...
	return
}

// BanUser implements the UserServiceImpl interface.
func (s *UserServiceImpl) BanUser(ctx context.Context, req *api.BanUserReq) (resp *api.BanUserResp, err error) {
	// TODO: Your code here...
	return
}

// UnbanUser implements the UserServiceImpl interface.
func (s *UserServiceImpl) UnbanUser(ctx context.Context, req *api.UnbanUserReq) (resp *api.UnbanUserResp, err error) {
	// TODO: Your code here...
	return
}

// DeleteUser implements the UserServiceImpl interface.
func (s *UserServiceImpl) DeleteUser(ctx context.Context, req *api.DeleteUserReq) (resp *api.DeleteUserResp, err error) {
	// TODO: Your code here...
	return
}

// RestoreUser implements the UserServiceImpl interface.
func (s *UserServiceImpl) RestoreUser(ctx context.Context, req *api.RestoreUserReq) (resp *api.RestoreUserResp, err error) {
	// TODO: Your code here...
	return
}

// UpdateUserStatus implements the UserServiceImpl interface.
func (s *UserServiceImpl) UpdateUserStatus(ctx context.Context, req *api.UpdateUserStatusReq) (resp *api.UpdateUserStatusResp, err error) {
	// TODO: Your code here...
	return
}

// ListUsers implements the UserServiceImpl interface.
func (s *UserServiceImpl) ListUsers(ctx context.Context, req *api.ListUsersReq) (resp *api.ListUsersResp, err error) {
	// TODO: Your code here...
	return
}

// SearchUsers implements the UserServiceImpl interface.
func (s *UserServiceImpl) SearchUsers(ctx context.Context, req *api.SearchUsersReq) (resp *api.SearchUsersResp, err error) {
	// TODO: Your code here...
	return
}

// CountUsers implements the UserServiceImpl interface.
func (s *UserServiceImpl) CountUsers(ctx context.Context, req *api.CountUsersReq) (resp *api.CountUsersResp, err error) {
	// TODO: Your code here...
	return
}

// CountByStatus implements the UserServiceImpl interface.
func (s *UserServiceImpl) CountByStatus(ctx context.Context) (resp *api.CountByStatusResp, err error) {
	// TODO: Your code here...
	return
}

// AdminUpdatePassword implements the UserServiceImpl interface.
func (s *UserServiceImpl) AdminUpdatePassword(ctx context.Context, req *api.UpdatePasswordReq) (resp *api.UpdatePasswordResp, err error) {
	// TODO: Your code here...
	return
}

// AdminUpdateEmail implements the UserServiceImpl interface.
func (s *UserServiceImpl) AdminUpdateEmail(ctx context.Context, req *api.UpdateEmailReq) (resp *api.UpdateEmailResp, err error) {
	// TODO: Your code here...
	return
}

// AdminUpdatePhone implements the UserServiceImpl interface.
func (s *UserServiceImpl) AdminUpdatePhone(ctx context.Context, req *api.UpdatePhoneReq) (resp *api.UpdatePhoneResp, err error) {
	// TODO: Your code here...
	return
}

// AdminUpdateUserProfile implements the UserServiceImpl interface.
func (s *UserServiceImpl) AdminUpdateUserProfile(ctx context.Context, req *api.UpdateUserProfileReq) (resp *api.UpdateUserProfileResp, err error) {
	// TODO: Your code here...
	return
}
