package store

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jb843051627/beeyard/internal/model"
)

// QueenStore 蜂王持久化。
type QueenStore struct {
	db *sql.DB
}

func NewQueenStore(db *sql.DB) *QueenStore {
	return &QueenStore{db: db}
}

// GetByID 按 id 查蜂王；不存在返回 (nil, ErrQueenNotFound)。
func (s *QueenStore) GetByID(ctx context.Context, id int64) (*model.Queen, error) {
	row := s.db.QueryRowContext(ctx,
		`SELECT id, hive_id, breed, status, marked_at, created_at FROM queens WHERE id = ?`, id)
	var q model.Queen
	err := row.Scan(&q.ID, &q.HiveID, &q.Breed, &q.Status, &q.MarkedAt, &q.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, model.ErrQueenNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("queen get by id: %w", err)
	}
	return &q, nil
}

// ListByHive 列出蜂箱下蜂王。
func (s *QueenStore) ListByHive(ctx context.Context, hiveID int64) ([]*model.Queen, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, hive_id, breed, status, marked_at, created_at FROM queens WHERE hive_id = ? ORDER BY id`, hiveID)
	if err != nil {
		return nil, fmt.Errorf("queen list by hive: %w", err)
	}
	defer rows.Close()
	var list []*model.Queen
	for rows.Next() {
		var q model.Queen
		if err := rows.Scan(&q.ID, &q.HiveID, &q.Breed, &q.Status, &q.MarkedAt, &q.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, &q)
	}
	return list, rows.Err()
}

// Create 新建蜂王。
func (s *QueenStore) Create(ctx context.Context, hiveID int64, breed, status string, markedAt interface{}) (int64, error) {
	res, err := s.db.ExecContext(ctx,
		`INSERT INTO queens(hive_id, breed, status, marked_at) VALUES(?,?,?,?)`,
		hiveID, breed, status, markedAt)
	if err != nil {
		return 0, fmt.Errorf("queen create: %w", err)
	}
	return res.LastInsertId()
}

// UpdateStatus 更新蜂王状态。
func (s *QueenStore) UpdateStatus(ctx context.Context, id int64, status string) error {
	res, err := s.db.ExecContext(ctx, `UPDATE queens SET status = ? WHERE id = ?`, status, id)
	if err != nil {
		return fmt.Errorf("queen update status: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return model.ErrQueenNotFound
	}
	return nil
}
