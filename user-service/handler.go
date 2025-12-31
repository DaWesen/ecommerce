package main

import (
	"context"
	"ecommerce/user-service/internal/repository"
	"ecommerce/user-service/internal/service"
	api "ecommerce/user-service/kitex_gen/api"
	"ecommerce/user-service/pkg/config"
	"ecommerce/user-service/pkg/database"
	"ecommerce/user-service/pkg/jwt"
	"fmt"
	"time"
)

// UserServiceImpl implements the last service interface defined in the IDL.
type UserServiceImpl struct {
	userService service.UserService
}

// 创建处理器
func NewUserServiceImpl() (*UserServiceImpl, error) {
	//加载配置
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("加载配置失败: %v", err)
	}
	//创建数据库连接
	db, _, err := database.NewDatabase(&cfg.Database)
	if err != nil {
		return nil, fmt.Errorf("数据库连接失败：%v", err)
	}
	//根据配置文件创建 JWT 配置
	jwtCfg := jwt.Config{
		SecretKey:     cfg.JWT.Secret,
		Issuer:        "ecommerce-user-service",
		AccessExpire:  time.Duration(cfg.JWT.ExpireHours) * time.Hour,
		RefreshExpire: 7 * 24 * time.Hour,
		Algorithm:     "HS256",
	}
	//创建 JWT 管理器
	jwtManager := jwt.NewJWTManager(jwtCfg)
	//创建 repository 和 service
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo, jwtManager)
	return &UserServiceImpl{
		userService: userService,
	}, nil
}

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
