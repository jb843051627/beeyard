package store

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jb843051627/beeyard/internal/model"
)

// TransferStore 转场持久化。
type TransferStore struct {
	db *sql.DB
}

func NewTransferStore(db *sql.DB) *TransferStore {
	return &TransferStore{db: db}
}

// Create 新建转场请求。
func (s *TransferStore) Create(ctx context.Context, hiveID, fromApiary, toApiary int64, requestedAt interface{}) (int64, error) {
	res, err := s.db.ExecContext(ctx,
		`INSERT INTO transfers(hive_id, from_apiary, to_apiary, status, requested_at) VALUES(?,?,?,?,?)`,
		hiveID, fromApiary, toApiary, model.TransferStatusPending, requestedAt)
	if err != nil {
		return 0, fmt.Errorf("transfer create: %w", err)
	}
	return res.LastInsertId()
}

// GetByID 查转场；不存在返回 (nil, ErrTransferNotFound)。
func (s *TransferStore) GetByID(ctx context.Context, id int64) (*model.Transfer, error) {
	row := s.db.QueryRowContext(ctx,
		`SELECT id, hive_id, from_apiary, to_apiary, status, requested_at, completed_at, created_at FROM transfers WHERE id = ?`, id)
	var t model.Transfer
	var completed sql.NullTime
	err := row.Scan(&t.ID, &t.HiveID, &t.FromApiary, &t.ToApiary, &t.Status, &t.RequestedAt, &completed, &t.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, model.ErrTransferNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("transfer get: %w", err)
	}
	if completed.Valid {
		tm := completed.Time
		t.CompletedAt = &tm
	}
	return &t, nil
}

// UpdateStatus 更新转场状态。
func (s *TransferStore) UpdateStatus(ctx context.Context, id int64, status string, completedAt interface{}) error {
	res, err := s.db.ExecContext(ctx,
		`UPDATE transfers SET status = ?, completed_at = ? WHERE id = ?`, status, completedAt, id)
	if err != nil {
		return fmt.Errorf("transfer update status: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return model.ErrTransferNotFound
	}
	return nil
}

// ListByDateRange 按日期范围列转场。
func (s *TransferStore) ListByDateRange(ctx context.Context, from, to interface{}) ([]*model.Transfer, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, hive_id, from_apiary, to_apiary, status, requested_at, completed_at, created_at FROM transfers WHERE requested_at >= ? AND requested_at <= ? ORDER BY requested_at ASC`, from, to)
	if err != nil {
		return nil, fmt.Errorf("transfer list by date range: %w", err)
	}
	defer rows.Close()
	var list []*model.Transfer
	for rows.Next() {
		var t model.Transfer
		var completed sql.NullTime
		if err := rows.Scan(&t.ID, &t.HiveID, &t.FromApiary, &t.ToApiary, &t.Status, &t.RequestedAt, &completed, &t.CreatedAt); err != nil {
			return nil, err
		}
		if completed.Valid {
			tm := completed.Time
			t.CompletedAt = &tm
		}
		list = append(list, &t)
	}
	return list, rows.Err()
}

// ListByApiary 列出蜂场相关转场。
func (s *TransferStore) ListByApiary(ctx context.Context, apiaryID int64) ([]*model.Transfer, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, hive_id, from_apiary, to_apiary, status, requested_at, completed_at, created_at FROM transfers WHERE from_apiary = ? OR to_apiary = ? ORDER BY requested_at DESC`, apiaryID, apiaryID)
	if err != nil {
		return nil, fmt.Errorf("transfer list by apiary: %w", err)
	}
	defer rows.Close()
	var list []*model.Transfer
	for rows.Next() {
		var t model.Transfer
		var completed sql.NullTime
		if err := rows.Scan(&t.ID, &t.HiveID, &t.FromApiary, &t.ToApiary, &t.Status, &t.RequestedAt, &completed, &t.CreatedAt); err != nil {
			return nil, err
		}
		if completed.Valid {
			tm := completed.Time
			t.CompletedAt = &tm
		}
		list = append(list, &t)
	}
	return list, rows.Err()
}
