package service

import (
	"context"
	"fmt"
	"sort"

	"github.com/jb843051627/beeyard/internal/model"
	"github.com/jb843051627/beeyard/internal/store"
)

// AlertService 告警业务。
type AlertService struct {
	store *store.AlertStore
}

func NewAlertService(s *store.AlertStore) *AlertService {
	return &AlertService{store: s}
}

// Create 新建告警；校验级别后入库，错误用 %w 包装以保留 errors.Is 链。
func (s *AlertService) Create(ctx context.Context, a *model.Alert) (int64, error) {
	if a.Status == "" {
		a.Status = model.AlertStatusActive
	}
	switch a.Level {
	case model.AlertLevelInfo, model.AlertLevelWarning, model.AlertLevelCritical:
	default:
		return 0, model.NewValidationError("level", "invalid alert level: "+a.Level)
	}
	id, err := s.store.Create(ctx, a)
	if err != nil {
		return 0, fmt.Errorf("create alert: %w", err)
	}
	return id, nil
}

// ListAlerts 按蜂场列告警，按创建时间倒序。
func (s *AlertService) ListAlerts(ctx context.Context, apiaryID int64) ([]*model.Alert, error) {
	alerts, err := s.store.ListByApiary(ctx, apiaryID)
	if err != nil {
		return nil, err
	}
	sorted := make([]*model.Alert, len(alerts))
	copy(sorted, alerts)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].CreatedAt.After(sorted[j].CreatedAt)
	})
	return sorted, nil
}

// Acknowledge 确认告警。
func (s *AlertService) Acknowledge(ctx context.Context, id int64) error {
	return s.store.Acknowledge(ctx, id)
}

// ListActive 列出未确认告警。
func (s *AlertService) ListActive(ctx context.Context, apiaryID int64) ([]*model.Alert, error) {
	return s.store.ListActive(ctx, apiaryID)
}

// GetByID 查告警。
func (s *AlertService) GetByID(ctx context.Context, id int64) (*model.Alert, error) {
	return s.store.GetByID(ctx, id)
}
