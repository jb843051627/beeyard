package store

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jb843051627/beeyard/internal/model"
)

// InspectionStore 巡检持久化。
type InspectionStore struct {
	db *sql.DB
}

func NewInspectionStore(db *sql.DB) *InspectionStore {
	return &InspectionStore{db: db}
}

// Create 新建巡检。
func (s *InspectionStore) Create(ctx context.Context, hiveID int64, scheduledAt interface{}, notes string) (int64, error) {
	res, err := s.db.ExecContext(ctx,
		`INSERT INTO inspections(hive_id, scheduled_at, notes) VALUES(?,?,?)`, hiveID, scheduledAt, notes)
	if err != nil {
		return 0, fmt.Errorf("inspection create: %w", err)
	}
	return res.LastInsertId()
}

// GetByID 查巡检；不存在返回 (nil, ErrInspectionNotFound)。
func (s *InspectionStore) GetByID(ctx context.Context, id int64) (*model.Inspection, error) {
	row := s.db.QueryRowContext(ctx,
		`SELECT id, hive_id, scheduled_at, status, completed_at, notes, created_at FROM inspections WHERE id = ?`, id)
	var ins model.Inspection
	var completed sql.NullTime
	err := row.Scan(&ins.ID, &ins.HiveID, &ins.ScheduledAt, &ins.Status, &completed, &ins.Notes, &ins.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, model.ErrInspectionNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("inspection get: %w", err)
	}
	if completed.Valid {
		t := completed.Time
		ins.CompletedAt = &t
	}
	return &ins, nil
}

// ListByHive 列出蜂箱巡检。
func (s *InspectionStore) ListByHive(ctx context.Context, hiveID int64) ([]*model.Inspection, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, hive_id, scheduled_at, status, completed_at, notes, created_at FROM inspections WHERE hive_id = ? ORDER BY scheduled_at DESC`, hiveID)
	if err != nil {
		return nil, fmt.Errorf("inspection list by hive: %w", err)
	}
	defer rows.Close()
	var list []*model.Inspection
	for rows.Next() {
		var ins model.Inspection
		var completed sql.NullTime
		if err := rows.Scan(&ins.ID, &ins.HiveID, &ins.ScheduledAt, &ins.Status, &completed, &ins.Notes, &ins.CreatedAt); err != nil {
			return nil, err
		}
		if completed.Valid {
			t := completed.Time
			ins.CompletedAt = &t
		}
		list = append(list, &ins)
	}
	return list, rows.Err()
}

// Complete 完成巡检。
func (s *InspectionStore) Complete(ctx context.Context, id int64, completedAt interface{}) error {
	res, err := s.db.ExecContext(ctx,
		`UPDATE inspections SET status = ?, completed_at = ? WHERE id = ?`, model.InspectionStatusCompleted, completedAt, id)
	if err != nil {
		return fmt.Errorf("inspection complete: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return model.ErrInspectionNotFound
	}
	return nil
}

// ListOverdue 列出逾期未完成巡检。
func (s *InspectionStore) ListOverdue(ctx context.Context, before interface{}) ([]*model.Inspection, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, hive_id, scheduled_at, status, completed_at, notes, created_at FROM inspections WHERE status = ? AND scheduled_at < ? ORDER BY scheduled_at ASC`, model.InspectionStatusPending, before)
	if err != nil {
		return nil, fmt.Errorf("inspection list overdue: %w", err)
	}
	defer rows.Close()
	var list []*model.Inspection
	for rows.Next() {
		var ins model.Inspection
		var completed sql.NullTime
		if err := rows.Scan(&ins.ID, &ins.HiveID, &ins.ScheduledAt, &ins.Status, &completed, &ins.Notes, &ins.CreatedAt); err != nil {
			return nil, err
		}
		if completed.Valid {
			t := completed.Time
			ins.CompletedAt = &t
		}
		list = append(list, &ins)
	}
	return list, rows.Err()
}


