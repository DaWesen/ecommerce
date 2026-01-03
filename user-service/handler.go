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
	"log"
	"strings"
	"time"
)

// UserServiceImpl implements the last service interface defined in the IDL.
type UserServiceImpl struct {
	userService service.UserService
	jwtManager  *jwt.JWTManager
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
		jwtManager:  jwtManager,
	}, nil
}

// 从请求中提取并验证token，返回包含用户信息的新上下文
func (s *UserServiceImpl) verifyTokenAndCreateContext(ctx context.Context, token string) (context.Context, error) {
	if token == "" {
		return ctx, fmt.Errorf("token为空")
	}

	//移除Bearer前缀
	token = strings.TrimPrefix(token, "Bearer ")
	token = strings.TrimSpace(token)

	//验证token
	claims, err := s.jwtManager.VerifyAccessToken(token)
	if err != nil {
		log.Printf("Token验证失败: %v", err)
		return ctx, fmt.Errorf("token无效或已过期")
	}

	//创建新的上下文，包含用户信息
	newCtx := context.WithValue(ctx, "user_id", claims.UserID)
	newCtx = context.WithValue(newCtx, "user_email", claims.Email)
	newCtx = context.WithValue(newCtx, "user_status", claims.Status)
	newCtx = context.WithValue(newCtx, "is_admin", claims.IsAdmin)

	return newCtx, nil
}

// Register implements the UserServiceImpl interface.
func (s *UserServiceImpl) Register(ctx context.Context, req *api.RegisterReq) (resp *api.RegisterResp, err error) {
	log.Printf("接收到注册用户请求：%s", req.Name)
	return s.userService.Register(ctx, req)
}

// Login implements the UserServiceImpl interface.
func (s *UserServiceImpl) Login(ctx context.Context, req *api.LoginReq) (resp *api.LoginResp, err error) {
	log.Printf("接收到登录用户请求")
	return s.userService.Login(ctx, req)
}

// GetUserProfile implements the UserServiceImpl interface.
func (s *UserServiceImpl) GetUserProfile(ctx context.Context, req *api.GetUserProfileReq) (resp *api.GetUserProfileResp, err error) {
	log.Printf("接收到获取用户资料请求,%d, token: %s", req.Id, req.Token)

	//验证token
	newCtx, err := s.verifyTokenAndCreateContext(ctx, req.Token)
	if err != nil {
		return &api.GetUserProfileResp{
			Code:    401,
			Success: false,
			Message: stringPtr("认证失败: " + err.Error()),
		}, nil
	}

	return s.userService.GetUserProfile(newCtx, req)
}

// UpdateUser implements the UserServiceImpl interface.
func (s *UserServiceImpl) UpdateUser(ctx context.Context, req *api.UpdateUserReq) (resp *api.UpdateUserResp, err error) {
	log.Printf("接收到更新用户请求, token: %s", req.Token)

	//验证token
	newCtx, err := s.verifyTokenAndCreateContext(ctx, req.Token)
	if err != nil {
		return &api.UpdateUserResp{
			Code:    401,
			Success: false,
			Message: stringPtr("认证失败: " + err.Error()),
		}, nil
	}

	return s.userService.UpdateUser(newCtx, req)
}

// ChangePassword implements the UserServiceImpl interface.
func (s *UserServiceImpl) ChangePassword(ctx context.Context, req *api.ChangePasswordReq) (resp *api.UpdateUserResp, err error) {
	log.Printf("接收到更改密码请求, token: %s", req.Token)

	//验证token
	newCtx, err := s.verifyTokenAndCreateContext(ctx, req.Token)
	if err != nil {
		return &api.UpdateUserResp{
			Code:    401,
			Success: false,
			Message: stringPtr("认证失败: " + err.Error()),
		}, nil
	}

	return s.userService.ChangePassword(newCtx, req)
}

// ChangeEmail implements the UserServiceImpl interface.
func (s *UserServiceImpl) ChangeEmail(ctx context.Context, req *api.ChangeEmailReq) (resp *api.UpdateUserResp, err error) {
	log.Printf("接收到更改邮箱请求, token: %s", req.Token)

	//验证token
	newCtx, err := s.verifyTokenAndCreateContext(ctx, req.Token)
	if err != nil {
		return &api.UpdateUserResp{
			Code:    401,
			Success: false,
			Message: stringPtr("认证失败: " + err.Error()),
		}, nil
	}

	return s.userService.ChangeEmail(newCtx, req)
}

// ChangePhone implements the UserServiceImpl interface.
func (s *UserServiceImpl) ChangePhone(ctx context.Context, req *api.ChangePhoneReq) (resp *api.UpdateUserResp, err error) {
	log.Printf("接收到更改手机号请求, token: %s", req.Token)

	//验证token
	newCtx, err := s.verifyTokenAndCreateContext(ctx, req.Token)
	if err != nil {
		return &api.UpdateUserResp{
			Code:    401,
			Success: false,
			Message: stringPtr("认证失败: " + err.Error()),
		}, nil
	}

	return s.userService.ChangePhone(newCtx, req)
}

// Logout implements the UserServiceImpl interface.
func (s *UserServiceImpl) Logout(ctx context.Context, req *api.LogoutReq) (resp *api.LogoutResp, err error) {
	log.Printf("接收到登出请求, token: %s", req.Token)

	//验证token
	newCtx, err := s.verifyTokenAndCreateContext(ctx, req.Token)
	if err != nil {
		return &api.LogoutResp{
			Code:    401,
			Success: false,
			Message: stringPtr("认证失败: " + err.Error()),
		}, nil
	}

	return s.userService.Logout(newCtx, req)
}

// GetUserStatus implements the UserServiceImpl interface.
func (s *UserServiceImpl) GetUserStatus(ctx context.Context, req *api.GetUserStatusReq) (resp *api.GetUserStatusResp, err error) {
	log.Printf("接收到获取用户状态请求：%d, token: %s", req.UserId, req.Token)

	//验证token
	newCtx, err := s.verifyTokenAndCreateContext(ctx, req.Token)
	if err != nil {
		return &api.GetUserStatusResp{
			Code:    401,
			Success: false,
			Message: stringPtr("认证失败: " + err.Error()),
		}, nil
	}

	return s.userService.GetUserStatus(newCtx, req)
}

// BanUser implements the UserServiceImpl interface.
func (s *UserServiceImpl) BanUser(ctx context.Context, req *api.BanUserReq) (resp *api.BanUserResp, err error) {
	log.Printf("接收到封禁用户请求：%d, token: %s", req.UserId, req.Token)

	//验证token
	newCtx, err := s.verifyTokenAndCreateContext(ctx, req.Token)
	if err != nil {
		return &api.BanUserResp{
			Code:    401,
			Success: false,
			Message: stringPtr("认证失败: " + err.Error()),
		}, nil
	}

	return s.userService.BanUser(newCtx, req)
}

// UnbanUser implements the UserServiceImpl interface.
func (s *UserServiceImpl) UnbanUser(ctx context.Context, req *api.UnbanUserReq) (resp *api.UnbanUserResp, err error) {
	log.Printf("接收到解封用户请求:%d, token: %s", req.UserId, req.Token)

	//验证token
	newCtx, err := s.verifyTokenAndCreateContext(ctx, req.Token)
	if err != nil {
		return &api.UnbanUserResp{
			Code:    401,
			Success: false,
			Message: stringPtr("认证失败: " + err.Error()),
		}, nil
	}

	return s.userService.UnbanUser(newCtx, req)
}

// DeleteUser implements the UserServiceImpl interface.
func (s *UserServiceImpl) DeleteUser(ctx context.Context, req *api.DeleteUserReq) (resp *api.DeleteUserResp, err error) {
	log.Printf("接收到删除用户请求:%d, token: %s", req.UserId, req.Token)

	//验证token
	newCtx, err := s.verifyTokenAndCreateContext(ctx, req.Token)
	if err != nil {
		return &api.DeleteUserResp{
			Code:    401,
			Success: false,
			Message: stringPtr("认证失败: " + err.Error()),
		}, nil
	}

	return s.userService.DeleteUser(newCtx, req)
}

// RestoreUser implements the UserServiceImpl interface.
func (s *UserServiceImpl) RestoreUser(ctx context.Context, req *api.RestoreUserReq) (resp *api.RestoreUserResp, err error) {
	log.Printf("接收到恢复用户请求:%d, token: %s", req.UserId, req.Token)

	//验证token
	newCtx, err := s.verifyTokenAndCreateContext(ctx, req.Token)
	if err != nil {
		return &api.RestoreUserResp{
			Code:    401,
			Success: false,
			Message: stringPtr("认证失败: " + err.Error()),
		}, nil
	}

	return s.userService.RestoreUser(newCtx, req)
}

// UpdateUserStatus implements the UserServiceImpl interface.
func (s *UserServiceImpl) UpdateUserStatus(ctx context.Context, req *api.UpdateUserStatusReq) (resp *api.UpdateUserStatusResp, err error) {
	log.Printf("接收到更新用户状态请求:%d, token: %s", req.UserId, req.Token)

	//验证token
	newCtx, err := s.verifyTokenAndCreateContext(ctx, req.Token)
	if err != nil {
		return &api.UpdateUserStatusResp{
			Code:    401,
			Success: false,
			Message: stringPtr("认证失败: " + err.Error()),
		}, nil
	}

	return s.userService.UpdateUserStatus(newCtx, req)
}

// ListUsers implements the UserServiceImpl interface.
func (s *UserServiceImpl) ListUsers(ctx context.Context, req *api.ListUsersReq) (resp *api.ListUsersResp, err error) {
	log.Printf("接收到列举用户请求, token: %s", req.Token)

	//验证token
	newCtx, err := s.verifyTokenAndCreateContext(ctx, req.Token)
	if err != nil {
		return &api.ListUsersResp{
			Code:    401,
			Success: false,
			Message: stringPtr("认证失败: " + err.Error()),
		}, nil
	}

	return s.userService.ListUsers(newCtx, req)
}

// SearchUsers implements the UserServiceImpl interface.
func (s *UserServiceImpl) SearchUsers(ctx context.Context, req *api.SearchUsersReq) (resp *api.SearchUsersResp, err error) {
	log.Printf("接收到搜索用户请求, token: %s", req.Token)

	//验证token
	newCtx, err := s.verifyTokenAndCreateContext(ctx, req.Token)
	if err != nil {
		return &api.SearchUsersResp{
			Code:    401,
			Success: false,
			Message: stringPtr("认证失败: " + err.Error()),
		}, nil
	}

	return s.userService.SearchUsers(newCtx, req)
}

// CountUsers implements the UserServiceImpl interface.
func (s *UserServiceImpl) CountUsers(ctx context.Context, req *api.CountUsersReq) (resp *api.CountUsersResp, err error) {
	log.Printf("接收到根据用户状态获取用户数量请求, token: %s", req.Token)

	//验证token
	newCtx, err := s.verifyTokenAndCreateContext(ctx, req.Token)
	if err != nil {
		return &api.CountUsersResp{
			Code:    401,
			Success: false,
			Message: stringPtr("认证失败: " + err.Error()),
		}, nil
	}

	return s.userService.CountUsers(newCtx, req)
}

// CountByStatus implements the UserServiceImpl interface.
func (s *UserServiceImpl) CountByStatus(ctx context.Context, req *api.CountByStatusReq) (resp *api.CountByStatusResp, err error) {
	log.Printf("接收到获取用户数量请求, token: %s", req.Token)

	//验证token
	newCtx, err := s.verifyTokenAndCreateContext(ctx, req.Token)
	if err != nil {
		return &api.CountByStatusResp{
			Code:    401,
			Success: false,
			Message: stringPtr("认证失败: " + err.Error()),
		}, nil
	}

	return s.userService.CountByStatus(newCtx)
}

// AdminUpdatePassword implements the UserServiceImpl interface.
func (s *UserServiceImpl) AdminUpdatePassword(ctx context.Context, req *api.UpdatePasswordReq) (resp *api.UpdatePasswordResp, err error) {
	log.Printf("接收到管理员更新用户密码请求:%d, token: %s", req.UserId, req.Token)

	//验证token
	newCtx, err := s.verifyTokenAndCreateContext(ctx, req.Token)
	if err != nil {
		return &api.UpdatePasswordResp{
			Code:    401,
			Success: false,
			Message: stringPtr("认证失败: " + err.Error()),
		}, nil
	}

	return s.userService.AdminUpdatePassword(newCtx, req)
}

// AdminUpdateEmail implements the UserServiceImpl interface.
func (s *UserServiceImpl) AdminUpdateEmail(ctx context.Context, req *api.UpdateEmailReq) (resp *api.UpdateEmailResp, err error) {
	log.Printf("接收到管理员更新用户邮箱请求:%d, token: %s", req.UserId, req.Token)

	//验证token
	newCtx, err := s.verifyTokenAndCreateContext(ctx, req.Token)
	if err != nil {
		return &api.UpdateEmailResp{
			Code:    401,
			Success: false,
			Message: stringPtr("认证失败: " + err.Error()),
		}, nil
	}

	return s.userService.AdminUpdateEmail(newCtx, req)
}

// AdminUpdatePhone implements the UserServiceImpl interface.
func (s *UserServiceImpl) AdminUpdatePhone(ctx context.Context, req *api.UpdatePhoneReq) (resp *api.UpdatePhoneResp, err error) {
	log.Printf("接收到管理员更新用户电话号请求:%d, token: %s", req.UserId, req.Token)

	//验证token
	newCtx, err := s.verifyTokenAndCreateContext(ctx, req.Token)
	if err != nil {
		return &api.UpdatePhoneResp{
			Code:    401,
			Success: false,
			Message: stringPtr("认证失败: " + err.Error()),
		}, nil
	}

	return s.userService.AdminUpdatePhone(newCtx, req)
}

// AdminUpdateUserProfile implements the UserServiceImpl interface.
func (s *UserServiceImpl) AdminUpdateUserProfile(ctx context.Context, req *api.UpdateUserProfileReq) (resp *api.UpdateUserProfileResp, err error) {
	name := "未知"
	if req.Name != nil {
		name = *req.Name
	}
	log.Printf("接收到管理员更新用户资料请求:%s, token: %s", name, req.Token)

	//验证token
	newCtx, err := s.verifyTokenAndCreateContext(ctx, req.Token)
	if err != nil {
		return &api.UpdateUserProfileResp{
			Code:    401,
			Success: false,
			Message: stringPtr("认证失败: " + err.Error()),
		}, nil
	}

	return s.userService.AdminUpdateUserProfile(newCtx, req)
}

// 辅助函数
func stringPtr(s string) *string {
	return &s
}
