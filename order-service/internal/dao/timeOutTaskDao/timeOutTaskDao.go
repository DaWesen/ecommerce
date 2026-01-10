package timeOutTaskDao

import (
	"context"
	"ecommerce/order-service/internal/dao/interfaces"
	"ecommerce/order-service/internal/model"
	"time"

	"gorm.io/gorm"
)

type TimeoutTaskRepository struct {
	db *gorm.DB
}

func NewTimeOutTaskRepository(db *gorm.DB) interfaces.ITimeoutTaskRepository {
	return &TimeoutTaskRepository{db: db}
}

// 创建超时任务
func (r *TimeoutTaskRepository) Create(ctx context.Context, task *model.TimeoutTask) error {
	return r.db.WithContext(ctx).Create(task).Error
}

// 更新超时任务
func (r *TimeoutTaskRepository) Update(ctx context.Context, task *model.TimeoutTask) error {
	return r.db.WithContext(ctx).Save(task).Error
}

// 根据任务ID查询
func (r *TimeoutTaskRepository) FindByTaskID(ctx context.Context, taskID string) (*model.TimeoutTask, error) {
	var task model.TimeoutTask
	err := r.db.WithContext(ctx).Where("task_id = ?", taskID).First(&task).Error
	return &task, err
}

// 查询已过期的超时任务
func (r *TimeoutTaskRepository) FindExpiredTasks(ctx context.Context, taskType string, limit int) ([]*model.TimeoutTask, error) {
	var tasks []*model.TimeoutTask
	now := time.Now()

	db := r.db.WithContext(ctx).
		Where("status = ? AND expire_time < ?", model.TaskStatusPending, now)

	if taskType != "" {
		db = db.Where("type = ?", taskType)
	}

	err := db.Limit(limit).Order("expire_time ASC").Find(&tasks).Error
	return tasks, err
}

// 删除超时任务
func (r *TimeoutTaskRepository) Delete(ctx context.Context, taskID string) error {
	return r.db.WithContext(ctx).Where("task_id = ?", taskID).Delete(&model.TimeoutTask{}).Error
}

// 更新超时任务状态
func (r *TimeoutTaskRepository) UpdateStatus(ctx context.Context, taskID, status string) error {
	return r.db.WithContext(ctx).Model(&model.TimeoutTask{}).
		Where("task_id = ?", taskID).
		Updates(map[string]interface{}{
			"status":     status,
			"updated_at": time.Now(),
		}).Error
}

// 增加重试次数
func (r *TimeoutTaskRepository) IncrementRetryCount(ctx context.Context, taskID string) error {
	return r.db.WithContext(ctx).Model(&model.TimeoutTask{}).
		Where("task_id = ?", taskID).
		UpdateColumn("retry_count", gorm.Expr("retry_count + ?", 1)).Error
}
