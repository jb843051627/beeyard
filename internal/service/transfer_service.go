package service

import (
	"context"
	"fmt"
	"time"

	"github.com/jb843051627/beeyard/internal/model"
	"github.com/jb843051627/beeyard/internal/store"
)

// TransferService 转场业务。
type TransferService struct {
	store  *store.TransferStore
	queens *store.QueenStore
}

func NewTransferService(s *store.TransferStore, q *store.QueenStore) *TransferService {
	return &TransferService{store: s, queens: q}
}

// Create 新建转场请求。
func (s *TransferService) Create(ctx context.Context, hiveID, fromApiary, toApiary int64, requestedAt time.Time) (int64, error) {
	if fromApiary == toApiary {
		return 0, model.NewValidationError("to_apiary", "source and destination apiary must differ")
	}
	return s.store.Create(ctx, hiveID, fromApiary, toApiary, requestedAt)
}

// AssignQueen 在转场途中为蜂箱分配蜂王；校验蜂王存在。
func (s *TransferService) AssignQueen(ctx context.Context, transferID, queenID int64) error {
	q, err := s.queens.GetByID(ctx, queenID)
	if err != nil {
		return fmt.Errorf("get queen %d for assign: %w", queenID, err)
	}
	_ = q
	return s.store.UpdateStatus(ctx, transferID, model.TransferStatusInTransit, nil)
}

// Complete 完成转场。
func (s *TransferService) Complete(ctx context.Context, transferID int64) error {
	return s.store.UpdateStatus(ctx, transferID, model.TransferStatusCompleted, time.Now())
}

// ListByDateRange 按日期范围列转场。
func (s *TransferService) ListByDateRange(ctx context.Context, from, to time.Time) ([]*model.Transfer, error) {
	return s.store.ListByDateRange(ctx, from, to)
}

// GetByID 查转场。
func (s *TransferService) GetByID(ctx context.Context, id int64) (*model.Transfer, error) {
	return s.store.GetByID(ctx, id)
}
