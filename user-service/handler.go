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
