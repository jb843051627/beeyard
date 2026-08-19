package store

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jb843051627/beeyard/internal/model"
)

// MaintenanceStore 维护任务持久化。
type MaintenanceStore struct {
	db *sql.DB
}

func NewMaintenanceStore(db *sql.DB) *MaintenanceStore {
	return &MaintenanceStore{db: db}
}

// Create 新建维护任务。
func (s *MaintenanceStore) Create(ctx context.Context, hiveID int64, description string, scheduledFor interface{}) (int64, error) {
	res, err := s.db.ExecContext(ctx,
		`INSERT INTO maintenance_tasks(hive_id, description, scheduled_for) VALUES(?,?,?)`, hiveID, description, scheduledFor)
	if err != nil {
		return 0, fmt.Errorf("maintenance create: %w", err)
	}
	return res.LastInsertId()
}

// GetByID 查维护任务；不存在返回 (nil, ErrMaintenanceNotFound)。
func (s *MaintenanceStore) GetByID(ctx context.Context, id int64) (*model.MaintenanceTask, error) {
	row := s.db.QueryRowContext(ctx,
		`SELECT id, hive_id, description, status, scheduled_for, completed_at, created_at FROM maintenance_tasks WHERE id = ?`, id)
	var m model.MaintenanceTask
	var completed sql.NullTime
	err := row.Scan(&m.ID, &m.HiveID, &m.Description, &m.Status, &m.ScheduledFor, &completed, &m.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, model.ErrMaintenanceNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("maintenance get: %w", err)
	}
	if completed.Valid {
		t := completed.Time
		m.CompletedAt = &t
	}
	return &m, nil
}

// ListPending 列出待完成维护任务。
func (s *MaintenanceStore) ListPending(ctx context.Context, hiveID int64) ([]*model.MaintenanceTask, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, hive_id, description, status, scheduled_for, completed_at, created_at FROM maintenance_tasks WHERE hive_id = ? AND status = ? ORDER BY scheduled_for ASC`, hiveID, model.TaskStatusPending)
	if err != nil {
		return nil, fmt.Errorf("maintenance list pending: %w", err)
	}
	defer rows.Close()
	var list []*model.MaintenanceTask
	for rows.Next() {
		var m model.MaintenanceTask
		var completed sql.NullTime
		if err := rows.Scan(&m.ID, &m.HiveID, &m.Description, &m.Status, &m.ScheduledFor, &completed, &m.CreatedAt); err != nil {
			return nil, err
		}
		if completed.Valid {
			t := completed.Time
			m.CompletedAt = &t
		}
		list = append(list, &m)
	}
	return list, rows.Err()
}

// Complete 完成维护任务。
func (s *MaintenanceStore) Complete(ctx context.Context, id int64, completedAt interface{}) error {
	res, err := s.db.ExecContext(ctx,
		`UPDATE maintenance_tasks SET status = ?, completed_at = ? WHERE id = ?`, model.TaskStatusCompleted, completedAt, id)
	if err != nil {
		return fmt.Errorf("maintenance complete: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return model.ErrMaintenanceNotFound
	}
	return nil
}

// Reschedule 改期维护任务。
func (s *MaintenanceStore) Reschedule(ctx context.Context, id int64, scheduledFor interface{}) error {
	res, err := s.db.ExecContext(ctx,
		`UPDATE maintenance_tasks SET scheduled_for = ? WHERE id = ? AND status = ?`, scheduledFor, id, model.TaskStatusPending)
	if err != nil {
		return fmt.Errorf("maintenance reschedule: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return model.ErrMaintenanceNotFound
	}
	return nil
}

// BatchCreate 批量创建维护任务（单事务，逐条出错即回滚）。
func (s *MaintenanceStore) BatchCreate(ctx context.Context, tasks []model.MaintenanceTask) (int64, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("begin tx: %w", err)
	}
	var total int64
	for _, t := range tasks {
		res, err := tx.ExecContext(ctx,
			`INSERT INTO maintenance_tasks(hive_id, description, scheduled_for) VALUES(?,?,?)`, t.HiveID, t.Description, t.ScheduledFor)
		if err != nil {
			tx.Rollback()
			return 0, fmt.Errorf("batch create row: %w", err)
		}
		n, _ := res.RowsAffected()
		total += n
	}
	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("commit batch: %w", err)
	}
	return total, nil
}
