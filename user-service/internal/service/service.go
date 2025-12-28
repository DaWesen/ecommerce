package service

import (
	"context"
	"ecommerce/user-service/internal/repository"
	"ecommerce/user-service/kitex_gen/api"
)

type UserService interface {
	Register(ctx context.Context, req *api.RegisterReq) (*api.RegisterResp, error)
	Login(ctx context.Context, req *api.LoginReq) (*api.LoginResp, error)
	UpdateUser(ctx context.Context, req *api.UpdateUserReq) (*api.UpdateUserResp, error)
	ChangePassword(ctx context.Context, req *api.ChangePasswordReq) (*api.UpdateUserResp, error)
	ChangeEmail(ctx context.Context, req *api.ChangeEmailReq) (*api.UpdateUserResp, error)
	ChangePhone(ctx context.Context, req *api.ChangePhoneReq) (*api.UpdateUserResp, error)
	GetUserProfile(ctx context.Context, req *api.GetUserProfileReq) *api.GetUserProfileResp
	Logout(ctx context.Context, req *api.LogoutReq) (*api.LogoutResp, error)
	GetUserStatus(ctx context.Context, req *api.GetUserStatusReq) (*api.GetUserStatusResp, error)
	//基本用户管理操作
	BanUser(ctx context.Context, req *api.BanUserReq) (*api.BanUserResp, error)
	UnbanUser(ctx context.Context, req *api.UnbanUserReq) (*api.UnbanUserResp, error)
	DeleteUser(ctx context.Context, req *api.DeleteUserReq) (*api.DeleteUserResp, error)
	RestoreUser(ctx context.Context, req *api.RestoreUserReq) (*api.RestoreUserResp, error)
	UpdateUserStatus(ctx context.Context, req *api.UpdateUserStatusReq) (*api.UpdateUserStatusResp, error)
	//搜索和列表
	ListUsers(ctx context.Context, req *api.ListUsersReq) (*api.ListUsersResp, error)
	SearchUsers(ctx context.Context, req *api.SearchUsersReq) (*api.SearchUsersResp, error)
	CountUsers(ctx context.Context, req *api.CountUsersReq) (*api.CountUsersResp, error)
	CountByStatus(ctx context.Context) (*api.CountByStatusResp, error)
	//管理员特殊操作
	AdminUpdatePassword(ctx context.Context, req *api.UpdatePasswordReq) (*api.UpdatePasswordResp, error)
	AdminUpdateEmail(ctx context.Context, req *api.UpdateEmailReq) (*api.UpdateEmailResp, error)
	AdminUpdatePhone(ctx context.Context, req *api.UpdatePhoneReq) (*api.UpdatePhoneResp, error)
	AdminUpdateUserProfile(ctx context.Context, req *api.UpdateUserProfileReq) (*api.UpdateUserProfileResp, error)
}

type userServiceImpl struct {
	userRepo repository.UserRepository
}

// 创建用户实例
func NewUserService(userRepo repository.UserRepository) UserService {
	return &userServiceImpl{
		userRepo: userRepo,
	}
}

const (
	ErrCodeSuccess       = 0
	ErrCodeBadRequest    = 400
	ErrCodeUnauthorized  = 401
	ErrCodeForbidden     = 403
	ErrCodeNotFound      = 404
	ErrCodeInternalError = 500
	//其他常量
	MinPasswordLength = 6
	MaxLoginAttempts  = 5
)
