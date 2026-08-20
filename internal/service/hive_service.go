package service

import (
	"context"
	"time"

	"github.com/jb843051627/beeyard/internal/model"
	"github.com/jb843051627/beeyard/internal/store"
)

// HiveService 蜂箱业务。
type HiveService struct {
	store *store.HiveStore
}

func NewHiveService(s *store.HiveStore) *HiveService {
	return &HiveService{store: s}
}

// Get 按 id 取蜂箱；不存在时返回 ErrHiveNotFound，调用方可 errors.Is 区分。
func (s *HiveService) Get(ctx context.Context, id int64) (*model.Hive, error) {
	return s.store.GetByID(ctx, id)
}

// Create 新建蜂箱（含状态校验）。
func (s *HiveService) Create(ctx context.Context, apiaryID int64, code, status string, installedAt time.Time) (int64, error) {
	if status == "" {
		status = model.HiveStatusActive
	}
	switch status {
	case model.HiveStatusActive, model.HiveStatusMaintenance, model.HiveStatusInactive:
	default:
		return 0, model.NewValidationError("status", "invalid hive status: "+status)
	}
	return s.store.Create(ctx, apiaryID, code, status, installedAt)
}

// ListByApiary 列出蜂场蜂箱。
func (s *HiveService) ListByApiary(ctx context.Context, apiaryID int64) ([]*model.Hive, error) {
	return s.store.ListByApiary(ctx, apiaryID)
}

// UpdateStatus 更新蜂箱状态（含合法性校验）。
func (s *HiveService) UpdateStatus(ctx context.Context, id int64, status string) error {
	switch status {
	case model.HiveStatusActive, model.HiveStatusMaintenance, model.HiveStatusInactive:
	default:
		return model.NewValidationError("status", "invalid hive status: "+status)
	}
	return s.store.UpdateStatus(ctx, id, status)
}

// Delete 删除蜂箱。
func (s *HiveService) Delete(ctx context.Context, id int64) error {
	return s.store.Delete(ctx, id)
}
