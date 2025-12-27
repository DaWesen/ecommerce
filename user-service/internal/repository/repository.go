package repository

import (
	"context"
	"ecommerce/user-service/internal/model"
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"
)

type UserRepository interface {
	//基础CRUD
	Create(ctx context.Context, user *model.User) error
	FindByID(ctx context.Context, id int64) (*model.User, error)
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	FindByPhone(ctx context.Context, phone string) (*model.User, error)
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id int64) error
	//验证邮箱或手机号的使用
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	ExistsByPhone(ctx context.Context, phone string) (bool, error)
	//状态更新
	UpdateStatus(ctx context.Context, id int64, status model.UserStatus) error
	UpdateLastLogin(ctx context.Context, id int64) error
	//封禁相关操作
	BanUser(ctx context.Context, id int64, reason string) error
	UnbanUser(ctx context.Context, id int64) error
	//软删除相关操作
	SoftDelete(ctx context.Context, id int64) error
	RestoreUser(ctx context.Context, id int64) error
	//状态检查
	IsUserActive(ctx context.Context, id int64) (bool, error)
	IsUserBanned(ctx context.Context, id int64) (bool, error)
	IsUserDeleted(ctx context.Context, id int64) (bool, error)
	//带状态过滤的查询
	FindActiveByID(ctx context.Context, id int64) (*model.User, error)
	FindActiveByEmail(ctx context.Context, email string) (*model.User, error)
	FindActiveByPhone(ctx context.Context, phone string) (*model.User, error)
	//管理员查询
	FindAllByID(ctx context.Context, id int64) (*model.User, error)
	//更新特定字段
	UpdatePassword(ctx context.Context, id int64, newPassword string) error
	UpdateEmail(ctx context.Context, id int64, newEmail string) error
	UpdatePhone(ctx context.Context, id int64, newPhone string) error
	UpdateProfile(ctx context.Context, id int64, updates map[string]interface{}) error
	//分页查询
	ListUsers(ctx context.Context, page, pageSize int, filters ...func(*gorm.DB) *gorm.DB) ([]*model.User, int64, error)
	SearchUsers(ctx context.Context, keyword string, page, pageSize int) ([]*model.User, int64, error)
	//统计
	CountUsers(ctx context.Context, status *model.UserStatus) (int64, error)
	CountByStatus(ctx context.Context) (map[model.UserStatus]int64, error)
}

type userRepositoryImpl struct {
	db *gorm.DB
}

// 创建用户实例
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepositoryImpl{db: db}
}

// 创建用户
func (r *userRepositoryImpl) Create(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

// 根据ID查找用户
func (r *userRepositoryImpl) FindByID(ctx context.Context, id int64) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// 根据邮箱查找用户
func (r *userRepositoryImpl) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).
		Where("email = ?", email).
		First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// 根据电话号码查找用户
func (r *userRepositoryImpl) FindByPhone(ctx context.Context, phone string) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).
		Where("phone = ?", phone).
		First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// 更新用户
func (r *userRepositoryImpl) Update(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

// 删除用户
func (r *userRepositoryImpl) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&model.User{}, id).Error
}

// 验证邮箱的使用
func (r *userRepositoryImpl) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.User{}).
		Where("email = ?", email).
		Count(&count).Error
	return count > 0, err
}

// 验证手机号的使用
func (r *userRepositoryImpl) ExistsByPhone(ctx context.Context, phone string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.User{}).
		Where("phone = ?", phone).
		Count(&count).Error
	return count > 0, err
}

// 更新用户状态
func (r *userRepositoryImpl) UpdateStatus(ctx context.Context, id int64, status model.UserStatus) error {
	return r.db.WithContext(ctx).Model(&model.User{}).
		Where("id = ?", id).
		Update("status", status).Error
}

// 更新最后登录的时间
func (r *userRepositoryImpl) UpdateLastLogin(ctx context.Context, id int64) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&model.User{}).
		Where("id = ?", id).
		Update("last_login", now).Error
}

// 封禁用户
func (r *userRepositoryImpl) BanUser(ctx context.Context, id int64, reason string) error {
	return r.db.WithContext(ctx).Model(&model.User{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":     model.UserStatusBANNED,
			"ban_reason": reason,
		}).Error
}

// 解封用户
func (r *userRepositoryImpl) UnbanUser(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Model(&model.User{}).
		Where("id = ? AND status = ?", id, model.UserStatusBANNED).
		Updates(map[string]interface{}{
			"status":     model.UserStatusACTIVE,
			"ban_reason": "",
		}).Error
}

// 软删除用户
func (r *userRepositoryImpl) SoftDelete(ctx context.Context, id int64) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&model.User{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":     model.UserStatusDELETED,
			"daleted_at": now,
		}).Error
}

// 恢复已删除的用户
func (r *userRepositoryImpl) RestoreUser(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Model(&model.User{}).
		Where("id = ? AND status = ?", id, model.UserStatusDELETED).
		Updates(map[string]interface{}{
			"status":     model.UserStatusACTIVE,
			"deleted_at": nil,
		}).Error
}

// 检察用户是否活跃
func (r *userRepositoryImpl) IsUserActive(ctx context.Context, id int64) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.User{}).
		Where("id = ? AND status = ?", id, model.UserStatusACTIVE).
		Count(&count).Error
	return count > 0, err
}

// 检查用户是否被封禁
func (r *userRepositoryImpl) IsUserBanned(ctx context.Context, id int64) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.User{}).
		Where("id = ? AND status = ?", id, model.UserStatusBANNED).
		Count(&count).Error
	return count > 0, err
}

// 检查用户是否被删除
func (r *userRepositoryImpl) IsUserDeleted(ctx context.Context, id int64) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.User{}).
		Where("id = ? AND status = ?", id, model.UserStatusDELETED).
		Count(&count).Error
	return count > 0, err
}

// 通过id查找活跃用户
func (r *userRepositoryImpl) FindActiveByID(ctx context.Context, id int64) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).
		Where("id = ? AND status = ?", id, model.UserStatusACTIVE).
		First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil

}

// 通过邮箱查找活跃用户
func (r *userRepositoryImpl) FindActiveByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).
		Where("email = ? AND status = ?", email, model.UserStatusACTIVE).
		First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// 通过手机号查找活跃用户
func (r *userRepositoryImpl) FindActiveByPhone(ctx context.Context, phone string) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).
		Where("phone = ? AND status = ?", phone, model.UserStatusACTIVE).
		First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// 管理员查询
func (r *userRepositoryImpl) FindAllByID(ctx context.Context, id int64) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		Unscoped(). //软删除
		First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// 更新密码
func (r *userRepositoryImpl) UpdatePassword(ctx context.Context, id int64, newPassword string) error {
	return r.db.WithContext(ctx).Model(&model.User{}).
		Where("id = ?", id).
		Update("password", newPassword).Error
}

// 更新邮箱
func (r *userRepositoryImpl) UpdateEmail(ctx context.Context, id int64, newEmail string) error {
	return r.db.WithContext(ctx).Model(&model.User{}).
		Where("id = ?", id).
		Update("email", newEmail).Error
}

// 更新手机号
func (r *userRepositoryImpl) UpdatePhone(ctx context.Context, id int64, newPhone string) error {
	return r.db.WithContext(ctx).Model(&model.User{}).
		Where("id = ?", id).
		Update("phone", newPhone).Error
}

// 更新用户资料
func (r *userRepositoryImpl) UpdateProfile(ctx context.Context, id int64, updates map[string]interface{}) error {
	return r.db.WithContext(ctx).Model(&model.User{}).
		Where("id = ?", id).
		Updates(updates).Error
}

// 分页查询用户
func (r *userRepositoryImpl) ListUsers(ctx context.Context, page, pageSize int,
	filters ...func(*gorm.DB) *gorm.DB) ([]*model.User, int64, error) {
	db := r.db.WithContext(ctx).Model(&model.User{})
	for _, filter := range filters {
		db = filter(db)
	}
	//计算总数
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	//分页查询
	var users []*model.User
	offset := (page - 1) * pageSize
	err := db.Offset(offset).
		Limit(pageSize).
		Order("created_at DESC").
		Find(&users).Error
	return users, total, err
}

// 搜索用户
func (r *userRepositoryImpl) SearchUsers(ctx context.Context, keyword string,
	page, pageSize int) ([]*model.User, int64, error) {
	db := r.db.WithContext(ctx).Model(&model.User{})
	//模糊搜索
	if keyword != "" {
		keyword = "%" + strings.TrimSpace(keyword) + "%"
		db = db.Where("name LIKE ? OR email LIKE ? OR phone LIKE ?", keyword, keyword, keyword)
	}
	//计算总数
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	//分页查询
	var users []*model.User
	offset := (page - 1) * pageSize
	err := db.Offset(offset).
		Limit(pageSize).
		Order("created_at DESC").
		Find(&users).Error
	return users, total, err
}

// 统计用户数量
func (r *userRepositoryImpl) CountUsers(ctx context.Context,
	status *model.UserStatus) (int64, error) {
	db := r.db.WithContext(ctx).Model(&model.User{})

	if status != nil {
		db = db.Where("status = ?", *status)
	}

	var count int64
	err := db.Count(&count).Error
	return count, err
}

// 按状态统计用户数量
func (r *userRepositoryImpl) CountByStatus(ctx context.Context) (map[model.UserStatus]int64,
	error) {
	type StatusCount struct {
		Status model.UserStatus
		Count  int64
	}
	var results []StatusCount
	err := r.db.WithContext(ctx).Model(&model.User{}).
		Select("status, COUNT(*) as count").
		Group("status").
		Scan(&results).Error
	if err != nil {
		return nil, err
	}
	//转换为map
	countMap := make(map[model.UserStatus]int64)
	for _, result := range results {
		countMap[result.Status] = result.Count
	}

	return countMap, nil
}
