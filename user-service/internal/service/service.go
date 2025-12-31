package service

import (
	"context"
	"fmt"
	"time"

	"ecommerce/user-service/internal/model"
	"ecommerce/user-service/internal/repository"
	"ecommerce/user-service/kitex_gen/api"
	"ecommerce/user-service/pkg/jwt"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService interface {
	Register(ctx context.Context, req *api.RegisterReq) (*api.RegisterResp, error)
	Login(ctx context.Context, req *api.LoginReq) (*api.LoginResp, error)
	UpdateUser(ctx context.Context, req *api.UpdateUserReq) (*api.UpdateUserResp, error)
	ChangePassword(ctx context.Context, req *api.ChangePasswordReq) (*api.UpdateUserResp, error)
	ChangeEmail(ctx context.Context, req *api.ChangeEmailReq) (*api.UpdateUserResp, error)
	ChangePhone(ctx context.Context, req *api.ChangePhoneReq) (*api.UpdateUserResp, error)
	GetUserProfile(ctx context.Context, req *api.GetUserProfileReq) (*api.GetUserProfileResp, error)
	Logout(ctx context.Context, req *api.LogoutReq) (*api.LogoutResp, error)
	GetUserStatus(ctx context.Context, req *api.GetUserStatusReq) (*api.GetUserStatusResp, error)
	BanUser(ctx context.Context, req *api.BanUserReq) (*api.BanUserResp, error)
	UnbanUser(ctx context.Context, req *api.UnbanUserReq) (*api.UnbanUserResp, error)
	DeleteUser(ctx context.Context, req *api.DeleteUserReq) (*api.DeleteUserResp, error)
	RestoreUser(ctx context.Context, req *api.RestoreUserReq) (*api.RestoreUserResp, error)
	UpdateUserStatus(ctx context.Context, req *api.UpdateUserStatusReq) (*api.UpdateUserStatusResp, error)
	ListUsers(ctx context.Context, req *api.ListUsersReq) (*api.ListUsersResp, error)
	SearchUsers(ctx context.Context, req *api.SearchUsersReq) (*api.SearchUsersResp, error)
	CountUsers(ctx context.Context, req *api.CountUsersReq) (*api.CountUsersResp, error)
	CountByStatus(ctx context.Context) (*api.CountByStatusResp, error)
	AdminUpdatePassword(ctx context.Context, req *api.UpdatePasswordReq) (*api.UpdatePasswordResp, error)
	AdminUpdateEmail(ctx context.Context, req *api.UpdateEmailReq) (*api.UpdateEmailResp, error)
	AdminUpdatePhone(ctx context.Context, req *api.UpdatePhoneReq) (*api.UpdatePhoneResp, error)
	AdminUpdateUserProfile(ctx context.Context, req *api.UpdateUserProfileReq) (*api.UpdateUserProfileResp, error)
}

type userServiceImpl struct {
	userRepo   repository.UserRepository
	jwtManager *jwt.JWTManager
}

func NewUserService(userRepo repository.UserRepository, jwtManager *jwt.JWTManager) UserService {
	return &userServiceImpl{
		userRepo:   userRepo,
		jwtManager: jwtManager,
	}
}

// 检查是否是管理员
func (s *userServiceImpl) isAdmin(ctx context.Context) bool {
	isAdmin, ok := ctx.Value("is_admin").(bool)
	return ok && isAdmin
}

// 获取当前用户ID
func (s *userServiceImpl) getCurrentUserID(ctx context.Context) (int64, bool) {
	userID, ok := ctx.Value("user_id").(int64)
	return userID, ok
}

// 检查用户是否有权限访问其他用户信息
func (s *userServiceImpl) canAccessUserInfo(ctx context.Context, targetUserID int64) bool {
	currentUserID, hasUser := s.getCurrentUserID(ctx)
	isAdmin := s.isAdmin(ctx)

	//管理员可以访问任何用户信息
	if isAdmin {
		return true
	}

	//登录用户可以访问自己的信息
	if hasUser && currentUserID == targetUserID {
		return true
	}
	return true
}

// 指针转化
func stringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func int32Ptr(i int32) *int32 {
	return &i
}

func int64Ptr(i int64) *int64 {
	return &i
}

// 哈希加密
func hashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// 验证
func verifyPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func (s *userServiceImpl) Register(ctx context.Context, req *api.RegisterReq) (*api.RegisterResp, error) {
	if req.Name == "" || req.Email == "" || req.Password == "" || req.Phone == "" {
		return &api.RegisterResp{
			Code:    400,
			Success: false,
			Message: stringPtr("所有字段都必须填写"),
		}, nil
	}

	if len(req.Password) < 6 {
		return &api.RegisterResp{
			Code:    400,
			Success: false,
			Message: stringPtr("密码长度不能少于6位"),
		}, nil
	}

	exists, err := s.userRepo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return &api.RegisterResp{
			Code:    500,
			Success: false,
			Message: stringPtr("服务器内部错误"),
		}, err
	}
	if exists {
		return &api.RegisterResp{
			Code:    400,
			Success: false,
			Message: stringPtr("邮箱已被注册"),
		}, nil
	}

	exists, err = s.userRepo.ExistsByPhone(ctx, req.Phone)
	if err != nil {
		return &api.RegisterResp{
			Code:    500,
			Success: false,
			Message: stringPtr("服务器内部错误"),
		}, err
	}
	if exists {
		return &api.RegisterResp{
			Code:    400,
			Success: false,
			Message: stringPtr("手机号已被注册"),
		}, nil
	}

	hashedPassword, err := hashPassword(req.Password)
	if err != nil {
		return &api.RegisterResp{
			Code:    500,
			Success: false,
			Message: stringPtr("密码加密失败"),
		}, err
	}

	now := time.Now().Unix()
	user := &model.User{
		Name:      req.Name,
		Email:     req.Email,
		Password:  hashedPassword,
		Phone:     req.Phone,
		Status:    model.UserStatusACTIVE,
		CreatedAt: time.Unix(now, 0),
		UpdatedAt: time.Unix(now, 0),
	}

	err = s.userRepo.Create(ctx, user)
	if err != nil {
		return &api.RegisterResp{
			Code:    500,
			Success: false,
			Message: stringPtr("创建用户失败"),
		}, err
	}

	token, err := s.jwtManager.GenerateAccessToken(user.ID, user.Email, user.Name, int32(user.Status), false)
	if err != nil {
		return &api.RegisterResp{
			Code:    500,
			Success: false,
			Message: stringPtr("生成令牌失败"),
		}, err
	}

	return &api.RegisterResp{
		Id:      user.ID,
		Token:   token,
		Success: true,
		Code:    0,
	}, nil
}

func (s *userServiceImpl) Login(ctx context.Context, req *api.LoginReq) (*api.LoginResp, error) {
	if req.Password == "" {
		return &api.LoginResp{
			Code:    400,
			Success: false,
			Message: stringPtr("密码不能为空"),
		}, nil
	}

	var user *model.User
	var err error

	if *req.Email != "" && req.Email != nil {
		user, err = s.userRepo.FindByEmail(ctx, *req.Email)
	} else if req.Phone != "" {
		user, err = s.userRepo.FindByPhone(ctx, req.Phone)
	} else {
		return &api.LoginResp{
			Code:    400,
			Success: false,
			Message: stringPtr("邮箱或手机号必须提供一个"),
		}, nil
	}

	if err != nil {
		return &api.LoginResp{
			Code:    500,
			Success: false,
			Message: stringPtr("服务器内部错误"),
		}, err
	}

	if user == nil {
		return &api.LoginResp{
			Code:    404,
			Success: false,
			Message: stringPtr("用户不存在"),
		}, nil
	}

	if user.Status != model.UserStatusACTIVE && user.Status != model.UserStatusPOWER {
		return &api.LoginResp{
			Code:    403,
			Success: false,
			Message: stringPtr("账号已被封禁或删除"),
		}, nil
	}

	if !verifyPassword(user.Password, req.Password) {
		return &api.LoginResp{
			Code:    401,
			Success: false,
			Message: stringPtr("密码错误"),
		}, nil
	}

	token, err := s.jwtManager.GenerateAccessToken(user.ID, user.Email, user.Name, int32(user.Status), user.Status == model.UserStatusPOWER)
	if err != nil {
		return &api.LoginResp{
			Code:    500,
			Success: false,
			Message: stringPtr("生成令牌失败"),
		}, err
	}

	now := time.Now().Unix()
	user.LastLogin = now
	err = s.userRepo.Update(ctx, user)
	if err != nil {
		return &api.LoginResp{
			Code:    500,
			Success: false,
			Message: stringPtr("更新登录时间失败"),
		}, err
	}

	return &api.LoginResp{
		Id:      user.ID,
		Token:   token,
		Success: true,
		Code:    0,
	}, nil
}

func (s *userServiceImpl) UpdateUser(ctx context.Context, req *api.UpdateUserReq) (*api.UpdateUserResp, error) {
	userID, ok := s.getCurrentUserID(ctx)
	if !ok {
		return &api.UpdateUserResp{
			Code:    401,
			Success: false,
			Message: stringPtr("未授权访问"),
		}, nil
	}

	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return &api.UpdateUserResp{
			Code:    500,
			Success: false,
			Message: stringPtr("服务器内部错误"),
		}, err
	}
	if user == nil {
		return &api.UpdateUserResp{
			Code:    404,
			Success: false,
			Message: stringPtr("用户不存在"),
		}, nil
	}

	updates := make(map[string]interface{})
	if req.Name != nil && *req.Name != "" {
		updates["name"] = *req.Name
	}
	if req.Avatar != nil {
		updates["avatar"] = *req.Avatar
	}
	if req.Bio != nil {
		updates["bio"] = *req.Bio
	}
	if req.Gender != nil {
		updates["gender"] = fmt.Sprintf("%d", *req.Gender)
	}
	updates["updated_at"] = time.Now()

	if len(updates) > 0 {
		err = s.userRepo.UpdateProfile(ctx, userID, updates)
		if err != nil {
			return &api.UpdateUserResp{
				Code:    500,
				Success: false,
				Message: stringPtr("更新用户失败"),
			}, err
		}
	}

	return &api.UpdateUserResp{
		Success: true,
		Code:    0,
	}, nil
}

func (s *userServiceImpl) ChangePassword(ctx context.Context, req *api.ChangePasswordReq) (*api.UpdateUserResp, error) {
	userID, ok := s.getCurrentUserID(ctx)
	if !ok {
		return &api.UpdateUserResp{
			Code:    401,
			Success: false,
			Message: stringPtr("未授权访问"),
		}, nil
	}

	if req.OldPassword == "" || req.NewPassword_ == "" {
		return &api.UpdateUserResp{
			Code:    400,
			Success: false,
			Message: stringPtr("新旧密码不能为空"),
		}, nil
	}

	if len(req.NewPassword_) < 6 {
		return &api.UpdateUserResp{
			Code:    400,
			Success: false,
			Message: stringPtr("新密码长度不能少于6位"),
		}, nil
	}

	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return &api.UpdateUserResp{
			Code:    500,
			Success: false,
			Message: stringPtr("服务器内部错误"),
		}, err
	}
	if user == nil {
		return &api.UpdateUserResp{
			Code:    404,
			Success: false,
			Message: stringPtr("用户不存在"),
		}, nil
	}

	if !verifyPassword(user.Password, req.OldPassword) {
		return &api.UpdateUserResp{
			Code:    401,
			Success: false,
			Message: stringPtr("旧密码错误"),
		}, nil
	}

	hashedPassword, err := hashPassword(req.NewPassword_)
	if err != nil {
		return &api.UpdateUserResp{
			Code:    500,
			Success: false,
			Message: stringPtr("密码加密失败"),
		}, err
	}

	err = s.userRepo.UpdatePassword(ctx, userID, hashedPassword)
	if err != nil {
		return &api.UpdateUserResp{
			Code:    500,
			Success: false,
			Message: stringPtr("更新密码失败"),
		}, err
	}

	return &api.UpdateUserResp{
		Success: true,
		Code:    0,
	}, nil
}

func (s *userServiceImpl) ChangeEmail(ctx context.Context, req *api.ChangeEmailReq) (*api.UpdateUserResp, error) {
	userID, ok := s.getCurrentUserID(ctx)
	if !ok {
		return &api.UpdateUserResp{
			Code:    401,
			Success: false,
			Message: stringPtr("未授权访问"),
		}, nil
	}

	if req.NewEmail_ == "" || req.Password == "" {
		return &api.UpdateUserResp{
			Code:    400,
			Success: false,
			Message: stringPtr("新邮箱和密码不能为空"),
		}, nil
	}

	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return &api.UpdateUserResp{
			Code:    500,
			Success: false,
			Message: stringPtr("服务器内部错误"),
		}, err
	}
	if user == nil {
		return &api.UpdateUserResp{
			Code:    404,
			Success: false,
			Message: stringPtr("用户不存在"),
		}, nil
	}

	if !verifyPassword(user.Password, req.Password) {
		return &api.UpdateUserResp{
			Code:    401,
			Success: false,
			Message: stringPtr("密码错误"),
		}, nil
	}

	exists, err := s.userRepo.ExistsByEmail(ctx, req.NewEmail_)
	if err != nil {
		return &api.UpdateUserResp{
			Code:    500,
			Success: false,
			Message: stringPtr("服务器内部错误"),
		}, err
	}
	if exists {
		return &api.UpdateUserResp{
			Code:    400,
			Success: false,
			Message: stringPtr("邮箱已被使用"),
		}, nil
	}

	err = s.userRepo.UpdateEmail(ctx, userID, req.NewEmail_)
	if err != nil {
		return &api.UpdateUserResp{
			Code:    500,
			Success: false,
			Message: stringPtr("更新邮箱失败"),
		}, err
	}

	return &api.UpdateUserResp{
		Success: true,
		Code:    0,
	}, nil
}

func (s *userServiceImpl) ChangePhone(ctx context.Context, req *api.ChangePhoneReq) (*api.UpdateUserResp, error) {
	userID, ok := s.getCurrentUserID(ctx)
	if !ok {
		return &api.UpdateUserResp{
			Code:    401,
			Success: false,
			Message: stringPtr("未授权访问"),
		}, nil
	}

	if req.NewPhone_ == "" || req.Password == "" {
		return &api.UpdateUserResp{
			Code:    400,
			Success: false,
			Message: stringPtr("新手机号和密码不能为空"),
		}, nil
	}

	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return &api.UpdateUserResp{
			Code:    500,
			Success: false,
			Message: stringPtr("服务器内部错误"),
		}, err
	}
	if user == nil {
		return &api.UpdateUserResp{
			Code:    404,
			Success: false,
			Message: stringPtr("用户不存在"),
		}, nil
	}

	if !verifyPassword(user.Password, req.Password) {
		return &api.UpdateUserResp{
			Code:    401,
			Success: false,
			Message: stringPtr("密码错误"),
		}, nil
	}

	exists, err := s.userRepo.ExistsByPhone(ctx, req.NewPhone_)
	if err != nil {
		return &api.UpdateUserResp{
			Code:    500,
			Success: false,
			Message: stringPtr("服务器内部错误"),
		}, err
	}
	if exists {
		return &api.UpdateUserResp{
			Code:    400,
			Success: false,
			Message: stringPtr("手机号已被使用"),
		}, nil
	}

	err = s.userRepo.UpdatePhone(ctx, userID, req.NewPhone_)
	if err != nil {
		return &api.UpdateUserResp{
			Code:    500,
			Success: false,
			Message: stringPtr("更新手机号失败"),
		}, err
	}

	return &api.UpdateUserResp{
		Success: true,
		Code:    0,
	}, nil
}

func (s *userServiceImpl) GetUserProfile(ctx context.Context, req *api.GetUserProfileReq) (*api.GetUserProfileResp, error) {
	if req.Id <= 0 {
		return &api.GetUserProfileResp{
			Code:    400,
			Success: false,
			Message: stringPtr("用户ID不能为空"),
		}, nil
	}

	//检查权限:允许查看公开的用户信息
	if !s.canAccessUserInfo(ctx, req.Id) {
		return &api.GetUserProfileResp{
			Code:    403,
			Success: false,
			Message: stringPtr("无权查看该用户信息"),
		}, nil
	}

	user, err := s.userRepo.FindByID(ctx, req.Id)
	if err != nil {
		return &api.GetUserProfileResp{
			Code:    500,
			Success: false,
			Message: stringPtr("服务器内部错误"),
		}, err
	}
	if user == nil {
		return &api.GetUserProfileResp{
			Code:    404,
			Success: false,
			Message: stringPtr("用户不存在"),
		}, nil
	}

	//非管理员查看时，隐藏敏感信息
	isAdmin := s.isAdmin(ctx)
	currentUserID, hasUser := s.getCurrentUserID(ctx)
	isSelf := hasUser && currentUserID == req.Id

	genderInt := int32(0)
	if user.Gender != "" {
		fmt.Sscanf(user.Gender, "%d", &genderInt)
	}

	safeUser := &api.SafeUser{
		Id:        user.ID,
		Name:      user.Name,
		Avatar:    stringPtr(user.Avatar),
		Bio:       stringPtr(user.Bio),
		Gender:    &genderInt,
		CreatedAt: user.CreatedAt.Unix(),
		Status:    api.UserStatus(user.Status),
	}
	//只有管理员或用户自己可以看到邮箱和手机号
	if isAdmin || isSelf {
		safeUser.Email = user.Email
		safeUser.Phone = user.Phone
		safeUser.UpdatedAt = user.UpdatedAt.Unix()
		safeUser.LastLogin = int64Ptr(user.LastLogin)
	}

	return &api.GetUserProfileResp{
		Success: true,
		Code:    0,
		User:    safeUser,
	}, nil
}

func (s *userServiceImpl) Logout(ctx context.Context, req *api.LogoutReq) (*api.LogoutResp, error) {
	return &api.LogoutResp{
		Success: true,
		Code:    0,
	}, nil
}

func (s *userServiceImpl) GetUserStatus(ctx context.Context, req *api.GetUserStatusReq) (*api.GetUserStatusResp, error) {
	if req.UserId <= 0 {
		return &api.GetUserStatusResp{
			Code:    400,
			Success: false,
			Message: stringPtr("用户ID不能为空"),
		}, nil
	}

	//检查权限：只有管理员或用户自己可以查看状态
	currentUserID, hasUser := s.getCurrentUserID(ctx)
	isAdmin := s.isAdmin(ctx)
	isSelf := hasUser && currentUserID == req.UserId

	if !isAdmin && !isSelf {
		return &api.GetUserStatusResp{
			Code:    403,
			Success: false,
			Message: stringPtr("无权查看该用户状态"),
		}, nil
	}

	user, err := s.userRepo.FindAllByID(ctx, req.UserId)
	if err != nil {
		return &api.GetUserStatusResp{
			Code:    500,
			Success: false,
			Message: stringPtr("服务器内部错误"),
		}, err
	}
	if user == nil {
		return &api.GetUserStatusResp{
			Code:    404,
			Success: false,
			Message: stringPtr("用户不存在"),
		}, nil
	}

	isBanned := user.Status == model.UserStatusBANNED
	isDeleted := user.Status == model.UserStatusDELETED

	return &api.GetUserStatusResp{
		Success:   true,
		Code:      0,
		Status:    api.UserStatus(user.Status),
		IsBanned:  isBanned,
		IsDeleted: isDeleted,
	}, nil
}

func (s *userServiceImpl) BanUser(ctx context.Context, req *api.BanUserReq) (*api.BanUserResp, error) {
	if !s.isAdmin(ctx) {
		return &api.BanUserResp{
			Code:    403,
			Success: false,
			Message: stringPtr("需要管理员权限"),
		}, nil
	}

	if req.UserId <= 0 {
		return &api.BanUserResp{
			Code:    400,
			Success: false,
			Message: stringPtr("用户ID不能为空"),
		}, nil
	}

	if req.Reason == "" {
		return &api.BanUserResp{
			Code:    400,
			Success: false,
			Message: stringPtr("封禁原因不能为空"),
		}, nil
	}

	user, err := s.userRepo.FindByID(ctx, req.UserId)
	if err != nil {
		return &api.BanUserResp{
			Code:    500,
			Success: false,
			Message: stringPtr("服务器内部错误"),
		}, err
	}
	if user == nil {
		return &api.BanUserResp{
			Code:    404,
			Success: false,
			Message: stringPtr("用户不存在"),
		}, nil
	}

	if user.Status == model.UserStatusBANNED {
		return &api.BanUserResp{
			Code:    400,
			Success: false,
			Message: stringPtr("用户已被封禁"),
		}, nil
	}

	err = s.userRepo.BanUser(ctx, req.UserId, req.Reason)
	if err != nil {
		return &api.BanUserResp{
			Code:    500,
			Success: false,
			Message: stringPtr("封禁用户失败"),
		}, err
	}

	now := time.Now().Unix()
	return &api.BanUserResp{
		Success:   true,
		Code:      0,
		BannedAt:  &now,
		BanReason: &req.Reason,
	}, nil
}

func (s *userServiceImpl) UnbanUser(ctx context.Context, req *api.UnbanUserReq) (*api.UnbanUserResp, error) {
	if !s.isAdmin(ctx) {
		return &api.UnbanUserResp{
			Code:    403,
			Success: false,
			Message: stringPtr("需要管理员权限"),
		}, nil
	}

	if req.UserId <= 0 {
		return &api.UnbanUserResp{
			Code:    400,
			Success: false,
			Message: stringPtr("用户ID不能为空"),
		}, nil
	}

	user, err := s.userRepo.FindByID(ctx, req.UserId)
	if err != nil {
		return &api.UnbanUserResp{
			Code:    500,
			Success: false,
			Message: stringPtr("服务器内部错误"),
		}, err
	}
	if user == nil {
		return &api.UnbanUserResp{
			Code:    404,
			Success: false,
			Message: stringPtr("用户不存在"),
		}, nil
	}

	if user.Status != model.UserStatusBANNED {
		return &api.UnbanUserResp{
			Code:    400,
			Success: false,
			Message: stringPtr("用户未被封禁"),
		}, nil
	}

	err = s.userRepo.UnbanUser(ctx, req.UserId)
	if err != nil {
		return &api.UnbanUserResp{
			Code:    500,
			Success: false,
			Message: stringPtr("解封用户失败"),
		}, err
	}

	return &api.UnbanUserResp{
		Success: true,
		Code:    0,
	}, nil
}

func (s *userServiceImpl) DeleteUser(ctx context.Context, req *api.DeleteUserReq) (*api.DeleteUserResp, error) {
	if !s.isAdmin(ctx) {
		return &api.DeleteUserResp{
			Code:    403,
			Success: false,
			Message: stringPtr("需要管理员权限"),
		}, nil
	}

	if req.UserId <= 0 {
		return &api.DeleteUserResp{
			Code:    400,
			Success: false,
			Message: stringPtr("用户ID不能为空"),
		}, nil
	}

	user, err := s.userRepo.FindByID(ctx, req.UserId)
	if err != nil {
		return &api.DeleteUserResp{
			Code:    500,
			Success: false,
			Message: stringPtr("服务器内部错误"),
		}, err
	}
	if user == nil {
		return &api.DeleteUserResp{
			Code:    404,
			Success: false,
			Message: stringPtr("用户不存在"),
		}, nil
	}

	if user.Status == model.UserStatusDELETED {
		return &api.DeleteUserResp{
			Code:    400,
			Success: false,
			Message: stringPtr("用户已被删除"),
		}, nil
	}

	err = s.userRepo.SoftDelete(ctx, req.UserId)
	if err != nil {
		return &api.DeleteUserResp{
			Code:    500,
			Success: false,
			Message: stringPtr("删除用户失败"),
		}, err
	}

	now := time.Now().Unix()
	return &api.DeleteUserResp{
		Success:   true,
		Code:      0,
		DeletedAt: &now,
	}, nil
}

func (s *userServiceImpl) RestoreUser(ctx context.Context, req *api.RestoreUserReq) (*api.RestoreUserResp, error) {
	if !s.isAdmin(ctx) {
		return &api.RestoreUserResp{
			Code:    403,
			Success: false,
			Message: stringPtr("需要管理员权限"),
		}, nil
	}

	if req.UserId <= 0 {
		return &api.RestoreUserResp{
			Code:    400,
			Success: false,
			Message: stringPtr("用户ID不能为空"),
		}, nil
	}

	user, err := s.userRepo.FindAllByID(ctx, req.UserId)
	if err != nil {
		return &api.RestoreUserResp{
			Code:    500,
			Success: false,
			Message: stringPtr("服务器内部错误"),
		}, err
	}
	if user == nil {
		return &api.RestoreUserResp{
			Code:    404,
			Success: false,
			Message: stringPtr("用户不存在"),
		}, nil
	}

	if user.Status != model.UserStatusDELETED {
		return &api.RestoreUserResp{
			Code:    400,
			Success: false,
			Message: stringPtr("用户未被删除"),
		}, nil
	}

	err = s.userRepo.RestoreUser(ctx, req.UserId)
	if err != nil {
		return &api.RestoreUserResp{
			Code:    500,
			Success: false,
			Message: stringPtr("恢复用户失败"),
		}, err
	}

	return &api.RestoreUserResp{
		Success: true,
		Code:    0,
	}, nil
}

func (s *userServiceImpl) UpdateUserStatus(ctx context.Context, req *api.UpdateUserStatusReq) (*api.UpdateUserStatusResp, error) {
	if !s.isAdmin(ctx) {
		return &api.UpdateUserStatusResp{
			Code:    403,
			Success: false,
			Message: stringPtr("需要管理员权限"),
		}, nil
	}

	if req.UserId <= 0 {
		return &api.UpdateUserStatusResp{
			Code:    400,
			Success: false,
			Message: stringPtr("用户ID不能为空"),
		}, nil
	}

	user, err := s.userRepo.FindAllByID(ctx, req.UserId)
	if err != nil {
		return &api.UpdateUserStatusResp{
			Code:    500,
			Success: false,
			Message: stringPtr("服务器内部错误"),
		}, err
	}
	if user == nil {
		return &api.UpdateUserStatusResp{
			Code:    404,
			Success: false,
			Message: stringPtr("用户不存在"),
		}, nil
	}

	oldStatus := user.Status
	err = s.userRepo.UpdateStatus(ctx, req.UserId, model.UserStatus(req.Status))
	if err != nil {
		return &api.UpdateUserStatusResp{
			Code:    500,
			Success: false,
			Message: stringPtr("更新用户状态失败"),
		}, err
	}

	now := time.Now().Unix()
	return &api.UpdateUserStatusResp{
		Success:    true,
		Code:       0,
		OldStatus:  api.UserStatus(oldStatus),
		NewStatus_: req.Status,
		UpdatedAt:  now,
	}, nil
}

func (s *userServiceImpl) ListUsers(ctx context.Context, req *api.ListUsersReq) (*api.ListUsersResp, error) {
	//允许所有登录用户查看用户列表（但只返回公开信息）
	currentUserID, hasUser := s.getCurrentUserID(ctx)
	if !hasUser {
		return &api.ListUsersResp{
			Code:    401,
			Success: false,
			Message: stringPtr("需要登录"),
		}, nil
	}

	isAdmin := s.isAdmin(ctx)

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}
	if req.PageSize > 100 {
		req.PageSize = 100
	}

	var filters []func(*gorm.DB) *gorm.DB

	//非管理员只能查看活跃用户
	if !isAdmin {
		filters = append(filters, func(db *gorm.DB) *gorm.DB {
			return db.Where("status IN (?)", []model.UserStatus{model.UserStatusACTIVE, model.UserStatusPOWER})
		})
	} else if req.Status != nil {
		//管理员可以按状态筛选
		filters = append(filters, func(db *gorm.DB) *gorm.DB {
			return db.Where("status = ?", *req.Status)
		})
	}

	if req.Keyword != nil && *req.Keyword != "" {
		keyword := "%" + *req.Keyword + "%"
		filters = append(filters, func(db *gorm.DB) *gorm.DB {
			return db.Where("name LIKE ?", keyword)
		})
	}

	if req.MinCreatedAt != nil {
		filters = append(filters, func(db *gorm.DB) *gorm.DB {
			return db.Where("created_at >= ?", time.Unix(*req.MinCreatedAt, 0))
		})
	}

	if req.MaxCreatedAt != nil {
		filters = append(filters, func(db *gorm.DB) *gorm.DB {
			return db.Where("created_at <= ?", time.Unix(*req.MaxCreatedAt, 0))
		})
	}

	orderBy := "created_at"
	if req.OrderBy != nil && *req.OrderBy != "" {
		orderBy = *req.OrderBy
	}

	orderClause := orderBy
	if req.Desc {
		orderClause += " DESC"
	} else {
		orderClause += " ASC"
	}

	filters = append(filters, func(db *gorm.DB) *gorm.DB {
		return db.Order(orderClause)
	})

	users, total, err := s.userRepo.ListUsers(ctx, int(req.Page), int(req.PageSize), filters...)
	if err != nil {
		return &api.ListUsersResp{
			Code:    500,
			Success: false,
			Message: stringPtr("获取用户列表失败"),
		}, err
	}

	safeUsers := make([]*api.SafeUser, len(users))
	for i, user := range users {
		genderInt := int32(0)
		if user.Gender != "" {
			fmt.Sscanf(user.Gender, "%d", &genderInt)
		}

		safeUser := &api.SafeUser{
			Id:        user.ID,
			Name:      user.Name,
			Avatar:    stringPtr(user.Avatar),
			Bio:       stringPtr(user.Bio),
			Gender:    &genderInt,
			CreatedAt: user.CreatedAt.Unix(),
			Status:    api.UserStatus(user.Status),
		}

		//只有管理员或用户自己可以看到邮箱和手机号
		if isAdmin || (hasUser && currentUserID == user.ID) {
			safeUser.Email = user.Email
			safeUser.Phone = user.Phone
			safeUser.UpdatedAt = user.UpdatedAt.Unix()
			safeUser.LastLogin = int64Ptr(user.LastLogin)
		}

		safeUsers[i] = safeUser
	}

	return &api.ListUsersResp{
		Success:  true,
		Code:     0,
		Total:    int32(total),
		Page:     req.Page,
		PageSize: req.PageSize,
		Users:    safeUsers,
	}, nil
}

func (s *userServiceImpl) SearchUsers(ctx context.Context, req *api.SearchUsersReq) (*api.SearchUsersResp, error) {
	//允许所有登录用户搜索用户（但只返回公开信息）
	_, hasUser := s.getCurrentUserID(ctx)
	if !hasUser {
		return &api.SearchUsersResp{
			Code:    401,
			Success: false,
			Message: stringPtr("需要登录"),
		}, nil
	}

	currentUserID, _ := s.getCurrentUserID(ctx)
	isAdmin := s.isAdmin(ctx)

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}
	if req.PageSize > 100 {
		req.PageSize = 100
	}

	//管理员可以搜索所有字段（名字、邮箱、手机号），非管理员只能搜索名字
	if isAdmin {
		//管理员：使用原有的SearchUsers方法搜索所有字段
		users, total, err := s.userRepo.SearchUsers(ctx, req.Keyword, int(req.Page), int(req.PageSize))
		if err != nil {
			return &api.SearchUsersResp{
				Code:    500,
				Success: false,
				Message: stringPtr("搜索用户失败"),
			}, err
		}

		safeUsers := make([]*api.SafeUser, len(users))
		for i, user := range users {
			genderInt := int32(0)
			if user.Gender != "" {
				fmt.Sscanf(user.Gender, "%d", &genderInt)
			}

			safeUser := &api.SafeUser{
				Id:        user.ID,
				Name:      user.Name,
				Avatar:    stringPtr(user.Avatar),
				Bio:       stringPtr(user.Bio),
				Gender:    &genderInt,
				CreatedAt: user.CreatedAt.Unix(),
				Status:    api.UserStatus(user.Status),
			}

			//管理员可以看到邮箱和手机号
			safeUser.Email = user.Email
			safeUser.Phone = user.Phone
			safeUser.UpdatedAt = user.UpdatedAt.Unix()
			safeUser.LastLogin = int64Ptr(user.LastLogin)

			safeUsers[i] = safeUser
		}

		return &api.SearchUsersResp{
			Success:  true,
			Code:     0,
			Total:    int32(total),
			Page:     req.Page,
			PageSize: req.PageSize,
			Users:    safeUsers,
		}, nil
	} else {
		//非管理员：只能按名字搜索，使用ListUsers方法并添加名字过滤
		var filters []func(*gorm.DB) *gorm.DB

		//非管理员只能查看活跃用户
		filters = append(filters, func(db *gorm.DB) *gorm.DB {
			return db.Where("status IN (?)", []model.UserStatus{model.UserStatusACTIVE, model.UserStatusPOWER})
		})
		if req.Keyword != "" {
			keyword := "%" + req.Keyword + "%"
			filters = append(filters, func(db *gorm.DB) *gorm.DB {
				return db.Where("name LIKE ?", keyword)
			})
		}

		users, total, err := s.userRepo.ListUsers(ctx, int(req.Page), int(req.PageSize), filters...)
		if err != nil {
			return &api.SearchUsersResp{
				Code:    500,
				Success: false,
				Message: stringPtr("搜索用户失败"),
			}, err
		}

		safeUsers := make([]*api.SafeUser, len(users))
		for i, user := range users {
			genderInt := int32(0)
			if user.Gender != "" {
				fmt.Sscanf(user.Gender, "%d", &genderInt)
			}

			safeUser := &api.SafeUser{
				Id:        user.ID,
				Name:      user.Name,
				Avatar:    stringPtr(user.Avatar),
				Bio:       stringPtr(user.Bio),
				Gender:    &genderInt,
				CreatedAt: user.CreatedAt.Unix(),
				Status:    api.UserStatus(user.Status),
			}

			//只有用户自己可以看到邮箱和手机号
			if currentUserID == user.ID {
				safeUser.Email = user.Email
				safeUser.Phone = user.Phone
				safeUser.UpdatedAt = user.UpdatedAt.Unix()
				safeUser.LastLogin = int64Ptr(user.LastLogin)
			}

			safeUsers[i] = safeUser
		}

		return &api.SearchUsersResp{
			Success:  true,
			Code:     0,
			Total:    int32(total),
			Page:     req.Page,
			PageSize: req.PageSize,
			Users:    safeUsers,
		}, nil
	}
}

func (s *userServiceImpl) CountUsers(ctx context.Context, req *api.CountUsersReq) (*api.CountUsersResp, error) {
	//允许所有登录用户统计用户数量
	_, hasUser := s.getCurrentUserID(ctx)
	if !hasUser {
		return &api.CountUsersResp{
			Code:    401,
			Success: false,
			Message: stringPtr("需要登录"),
		}, nil
	}

	isAdmin := s.isAdmin(ctx)

	var status *model.UserStatus
	if req.Status != nil {
		//非管理员只能统计活跃用户数量
		if !isAdmin && *req.Status != api.UserStatus_ACTIVE && *req.Status != api.UserStatus_POWER {
			return &api.CountUsersResp{
				Code:    403,
				Success: false,
				Message: stringPtr("无权统计该状态用户"),
			}, nil
		}
		s := model.UserStatus(*req.Status)
		status = &s
	} else if !isAdmin {
		//非管理员只能统计活跃用户
		activeStatus := model.UserStatusACTIVE
		status = &activeStatus
	}

	count, err := s.userRepo.CountUsers(ctx, status)
	if err != nil {
		return &api.CountUsersResp{
			Code:    500,
			Success: false,
			Message: stringPtr("统计用户失败"),
		}, err
	}

	return &api.CountUsersResp{
		Success: true,
		Code:    0,
		Count:   count,
	}, nil
}

func (s *userServiceImpl) CountByStatus(ctx context.Context) (*api.CountByStatusResp, error) {
	//只有管理员可以按状态统计
	if !s.isAdmin(ctx) {
		return &api.CountByStatusResp{
			Code:    403,
			Success: false,
			Message: stringPtr("需要管理员权限"),
		}, nil
	}

	countMap, err := s.userRepo.CountByStatus(ctx)
	if err != nil {
		return &api.CountByStatusResp{
			Code:    500,
			Success: false,
			Message: stringPtr("按状态统计失败"),
		}, err
	}

	apiCountMap := make(map[api.UserStatus]int64)
	for status, count := range countMap {
		apiCountMap[api.UserStatus(status)] = count
	}

	return &api.CountByStatusResp{
		Success: true,
		Code:    0,
		Counts:  apiCountMap,
	}, nil
}

func (s *userServiceImpl) AdminUpdatePassword(ctx context.Context, req *api.UpdatePasswordReq) (*api.UpdatePasswordResp, error) {
	if !s.isAdmin(ctx) {
		return &api.UpdatePasswordResp{
			Code:    403,
			Success: false,
			Message: stringPtr("需要管理员权限"),
		}, nil
	}

	if req.UserId <= 0 {
		return &api.UpdatePasswordResp{
			Code:    400,
			Success: false,
			Message: stringPtr("用户ID不能为空"),
		}, nil
	}

	if req.NewPassword_ == "" {
		return &api.UpdatePasswordResp{
			Code:    400,
			Success: false,
			Message: stringPtr("新密码不能为空"),
		}, nil
	}

	if len(req.NewPassword_) < 6 {
		return &api.UpdatePasswordResp{
			Code:    400,
			Success: false,
			Message: stringPtr("密码长度不能少于6位"),
		}, nil
	}

	user, err := s.userRepo.FindByID(ctx, req.UserId)
	if err != nil {
		return &api.UpdatePasswordResp{
			Code:    500,
			Success: false,
			Message: stringPtr("服务器内部错误"),
		}, err
	}
	if user == nil {
		return &api.UpdatePasswordResp{
			Code:    404,
			Success: false,
			Message: stringPtr("用户不存在"),
		}, nil
	}

	hashedPassword, err := hashPassword(req.NewPassword_)
	if err != nil {
		return &api.UpdatePasswordResp{
			Code:    500,
			Success: false,
			Message: stringPtr("密码加密失败"),
		}, err
	}

	err = s.userRepo.UpdatePassword(ctx, req.UserId, hashedPassword)
	if err != nil {
		return &api.UpdatePasswordResp{
			Code:    500,
			Success: false,
			Message: stringPtr("更新密码失败"),
		}, err
	}

	return &api.UpdatePasswordResp{
		Success: true,
		Code:    0,
	}, nil
}

func (s *userServiceImpl) AdminUpdateEmail(ctx context.Context, req *api.UpdateEmailReq) (*api.UpdateEmailResp, error) {
	if !s.isAdmin(ctx) {
		return &api.UpdateEmailResp{
			Code:    403,
			Success: false,
			Message: stringPtr("需要管理员权限"),
		}, nil
	}

	if req.UserId <= 0 {
		return &api.UpdateEmailResp{
			Code:    400,
			Success: false,
			Message: stringPtr("用户ID不能为空"),
		}, nil
	}

	if req.NewEmail_ == "" {
		return &api.UpdateEmailResp{
			Code:    400,
			Success: false,
			Message: stringPtr("新邮箱不能为空"),
		}, nil
	}

	user, err := s.userRepo.FindByID(ctx, req.UserId)
	if err != nil {
		return &api.UpdateEmailResp{
			Code:    500,
			Success: false,
			Message: stringPtr("服务器内部错误"),
		}, err
	}
	if user == nil {
		return &api.UpdateEmailResp{
			Code:    404,
			Success: false,
			Message: stringPtr("用户不存在"),
		}, nil
	}

	exists, err := s.userRepo.ExistsByEmail(ctx, req.NewEmail_)
	if err != nil {
		return &api.UpdateEmailResp{
			Code:    500,
			Success: false,
			Message: stringPtr("服务器内部错误"),
		}, err
	}
	if exists {
		return &api.UpdateEmailResp{
			Code:    400,
			Success: false,
			Message: stringPtr("邮箱已被使用"),
		}, nil
	}

	err = s.userRepo.UpdateEmail(ctx, req.UserId, req.NewEmail_)
	if err != nil {
		return &api.UpdateEmailResp{
			Code:    500,
			Success: false,
			Message: stringPtr("更新邮箱失败"),
		}, err
	}

	return &api.UpdateEmailResp{
		Success: true,
		Code:    0,
	}, nil
}

func (s *userServiceImpl) AdminUpdatePhone(ctx context.Context, req *api.UpdatePhoneReq) (*api.UpdatePhoneResp, error) {
	if !s.isAdmin(ctx) {
		return &api.UpdatePhoneResp{
			Code:    403,
			Success: false,
			Message: stringPtr("需要管理员权限"),
		}, nil
	}

	if req.UserId <= 0 {
		return &api.UpdatePhoneResp{
			Code:    400,
			Success: false,
			Message: stringPtr("用户ID不能为空"),
		}, nil
	}

	if req.NewPhone_ == "" {
		return &api.UpdatePhoneResp{
			Code:    400,
			Success: false,
			Message: stringPtr("新手机号不能为空"),
		}, nil
	}

	user, err := s.userRepo.FindByID(ctx, req.UserId)
	if err != nil {
		return &api.UpdatePhoneResp{
			Code:    500,
			Success: false,
			Message: stringPtr("服务器内部错误"),
		}, err
	}
	if user == nil {
		return &api.UpdatePhoneResp{
			Code:    404,
			Success: false,
			Message: stringPtr("用户不存在"),
		}, nil
	}

	exists, err := s.userRepo.ExistsByPhone(ctx, req.NewPhone_)
	if err != nil {
		return &api.UpdatePhoneResp{
			Code:    500,
			Success: false,
			Message: stringPtr("服务器内部错误"),
		}, err
	}
	if exists {
		return &api.UpdatePhoneResp{
			Code:    400,
			Success: false,
			Message: stringPtr("手机号已被使用"),
		}, nil
	}

	err = s.userRepo.UpdatePhone(ctx, req.UserId, req.NewPhone_)
	if err != nil {
		return &api.UpdatePhoneResp{
			Code:    500,
			Success: false,
			Message: stringPtr("更新手机号失败"),
		}, err
	}

	return &api.UpdatePhoneResp{
		Success: true,
		Code:    0,
	}, nil
}

func (s *userServiceImpl) AdminUpdateUserProfile(ctx context.Context, req *api.UpdateUserProfileReq) (*api.UpdateUserProfileResp, error) {
	if !s.isAdmin(ctx) {
		return &api.UpdateUserProfileResp{
			Code:    403,
			Success: false,
			Message: stringPtr("需要管理员权限"),
		}, nil
	}

	if req.UserId <= 0 {
		return &api.UpdateUserProfileResp{
			Code:    400,
			Success: false,
			Message: stringPtr("用户ID不能为空"),
		}, nil
	}

	user, err := s.userRepo.FindByID(ctx, req.UserId)
	if err != nil {
		return &api.UpdateUserProfileResp{
			Code:    500,
			Success: false,
			Message: stringPtr("服务器内部错误"),
		}, err
	}
	if user == nil {
		return &api.UpdateUserProfileResp{
			Code:    404,
			Success: false,
			Message: stringPtr("用户不存在"),
		}, nil
	}

	updates := make(map[string]interface{})
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Avatar != nil {
		updates["avatar"] = *req.Avatar
	}
	if req.Bio != nil {
		updates["bio"] = *req.Bio
	}
	if req.Gender != nil {
		updates["gender"] = fmt.Sprintf("%d", *req.Gender)
	}
	updates["updated_at"] = time.Now()

	if len(updates) > 0 {
		err = s.userRepo.UpdateProfile(ctx, req.UserId, updates)
		if err != nil {
			return &api.UpdateUserProfileResp{
				Code:    500,
				Success: false,
				Message: stringPtr("更新用户资料失败"),
			}, err
		}
	}

	updatedUser, err := s.userRepo.FindByID(ctx, req.UserId)
	if err != nil {
		return &api.UpdateUserProfileResp{
			Code:    500,
			Success: false,
			Message: stringPtr("获取更新后用户信息失败"),
		}, err
	}

	genderInt := int32(0)
	if updatedUser.Gender != "" {
		fmt.Sscanf(updatedUser.Gender, "%d", &genderInt)
	}

	safeUser := &api.SafeUser{
		Id:        updatedUser.ID,
		Name:      updatedUser.Name,
		Email:     updatedUser.Email,
		Phone:     updatedUser.Phone,
		Avatar:    stringPtr(updatedUser.Avatar),
		Bio:       stringPtr(updatedUser.Bio),
		Gender:    &genderInt,
		CreatedAt: updatedUser.CreatedAt.Unix(),
		UpdatedAt: updatedUser.UpdatedAt.Unix(),
		Status:    api.UserStatus(updatedUser.Status),
		LastLogin: int64Ptr(updatedUser.LastLogin),
	}

	return &api.UpdateUserProfileResp{
		Success: true,
		Code:    0,
		User:    safeUser,
	}, nil
}
