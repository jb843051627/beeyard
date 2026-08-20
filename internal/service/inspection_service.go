package service

import (
	"context"
	"time"

	"github.com/jb843051627/beeyard/internal/model"
	"github.com/jb843051627/beeyard/internal/store"
)

// InspectionService 巡检业务。
type InspectionService struct {
	store *store.InspectionStore
}

func NewInspectionService(s *store.InspectionStore) *InspectionService {
	return &InspectionService{store: s}
}

// Schedule 排期巡检。
func (s *InspectionService) Schedule(ctx context.Context, hiveID int64, scheduledAt time.Time, notes string) (int64, error) {
	if scheduledAt.IsZero() {
		return 0, model.NewValidationError("scheduled_at", "scheduled time is required")
	}
	return s.store.Create(ctx, hiveID, scheduledAt, notes)
}

// Complete 完成巡检；先校验当前状态合法再标记完成，避免重复完成静默覆盖完成时间。
func (s *InspectionService) Complete(ctx context.Context, id int64, notes string) error {
	ins, err := s.store.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if ins.Status == model.InspectionStatusCompleted {
		return model.NewValidationError("status", "inspection already completed")
	}
	return s.store.Complete(ctx, id, time.Now())
}

// ListByHive 列出蜂箱巡检。
func (s *InspectionService) ListByHive(ctx context.Context, hiveID int64) ([]*model.Inspection, error) {
	return s.store.ListByHive(ctx, hiveID)
}

// ListOverdue 列出逾期巡检。
func (s *InspectionService) ListOverdue(ctx context.Context, before time.Time) ([]*model.Inspection, error) {
	return s.store.ListOverdue(ctx, before)
}

// GetByID 查巡检。
func (s *InspectionService) GetByID(ctx context.Context, id int64) (*model.Inspection, error) {
	return s.store.GetByID(ctx, id)
}
