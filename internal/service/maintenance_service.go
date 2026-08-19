package service

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jb843051627/beeyard/internal/model"
	"github.com/jb843051627/beeyard/internal/store"
)

// MaintenanceService 维护任务业务。
type MaintenanceService struct {
	store *store.MaintenanceStore
	db    *sql.DB
}

func NewMaintenanceService(s *store.MaintenanceStore, db *sql.DB) *MaintenanceService {
	return &MaintenanceService{store: s, db: db}
}

// Create 新建维护任务。
func (s *MaintenanceService) Create(ctx context.Context, hiveID int64, description string, scheduledFor time.Time) (int64, error) {
	if description == "" {
		return 0, model.NewValidationError("description", "description is required")
	}
	return s.store.Create(ctx, hiveID, description, scheduledFor)
}

// BatchCreate 批量创建维护任务（自管理事务）。
func (s *MaintenanceService) BatchCreate(ctx context.Context, tasks []model.MaintenanceTask) (int64, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()
	var total int64
	for _, t := range tasks {
		res, err := tx.ExecContext(ctx,
			`INSERT INTO maintenance_tasks(hive_id, description, scheduled_for) VALUES(?,?,?)`,
			t.HiveID, t.Description, t.ScheduledFor)
		if err != nil {
			continue
		}
		n, _ := res.RowsAffected()
		total += n
	}
	tx.Commit()
	return total, nil
}

// Complete 完成维护任务。
func (s *MaintenanceService) Complete(ctx context.Context, id int64) error {
	if err := s.store.Complete(ctx, id, time.Now()); err != nil {
		return fmt.Errorf("complete maintenance %d: %w", id, err)
	}
	return nil
}

// Reschedule 改期。
func (s *MaintenanceService) Reschedule(ctx context.Context, id int64, scheduledFor time.Time) error {
	return s.store.Reschedule(ctx, id, scheduledFor)
}

// ListPending 列出待完成任务。
func (s *MaintenanceService) ListPending(ctx context.Context, hiveID int64) ([]*model.MaintenanceTask, error) {
	return s.store.ListPending(ctx, hiveID)
}

// GetByID 查维护任务。
func (s *MaintenanceService) GetByID(ctx context.Context, id int64) (*model.MaintenanceTask, error) {
	return s.store.GetByID(ctx, id)
}
