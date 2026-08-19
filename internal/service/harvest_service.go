package service

import (
	"context"
	"sort"

	"github.com/jb843051627/beeyard/internal/model"
	"github.com/jb843051627/beeyard/internal/store"
)

// HarvestService 采蜜业务。
type HarvestService struct {
	store *store.HarvestStore
}

func NewHarvestService(s *store.HarvestStore) *HarvestService {
	return &HarvestService{store: s}
}

// Record 登记采蜜。
func (s *HarvestService) Record(ctx context.Context, h *model.Harvest) (int64, error) {
	if h.AmountKg <= 0 {
		return 0, model.NewValidationError("amount_kg", "amount must be positive")
	}
	return s.store.Create(ctx, h)
}

// GenerateReport 生成采蜜报表，按采蜜量降序。
func (s *HarvestService) GenerateReport(ctx context.Context, apiaryID int64) ([]*model.Harvest, error) {
	list, err := s.store.ListByApiary(ctx, apiaryID)
	if err != nil {
		return nil, err
	}
	sorted := make([]*model.Harvest, len(list))
	copy(sorted, list)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].AmountKg > sorted[j].AmountKg
	})
	return sorted, nil
}

// ListByApiary 列出蜂场采蜜记录。
func (s *HarvestService) ListByApiary(ctx context.Context, apiaryID int64) ([]*model.Harvest, error) {
	return s.store.ListByApiary(ctx, apiaryID)
}

// GetByID 查采蜜记录。
func (s *HarvestService) GetByID(ctx context.Context, id int64) (*model.Harvest, error) {
	return s.store.GetByID(ctx, id)
}
