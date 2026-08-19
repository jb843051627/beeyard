package store

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jb843051627/beeyard/internal/model"
)

// HiveStore 蜂箱持久化。
type HiveStore struct {
	db *sql.DB
}

func NewHiveStore(db *sql.DB) *HiveStore {
	return &HiveStore{db: db}
}

// GetByID 按 id 查蜂箱；不存在返回 (nil, ErrHiveNotFound)。
func (s *HiveStore) GetByID(ctx context.Context, id int64) (*model.Hive, error) {
	row := s.db.QueryRowContext(ctx,
		`SELECT id, apiary_id, code, status, installed_at, created_at FROM hives WHERE id = ?`, id)
	var h model.Hive
	err := row.Scan(&h.ID, &h.ApiaryID, &h.Code, &h.Status, &h.InstalledAt, &h.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("hive get by id: %w", err)
	}
	return &h, nil
}

// ListByApiary 列出蜂场下全部蜂箱。
func (s *HiveStore) ListByApiary(ctx context.Context, apiaryID int64) ([]*model.Hive, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, apiary_id, code, status, installed_at, created_at FROM hives WHERE apiary_id = ? ORDER BY id`, apiaryID)
	if err != nil {
		return nil, fmt.Errorf("hive list by apiary: %w", err)
	}
	defer rows.Close()
	var list []*model.Hive
	for rows.Next() {
		var h model.Hive
		if err := rows.Scan(&h.ID, &h.ApiaryID, &h.Code, &h.Status, &h.InstalledAt, &h.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, &h)
	}
	return list, rows.Err()
}

// Create 新建蜂箱。
func (s *HiveStore) Create(ctx context.Context, apiaryID int64, code, status string, installedAt interface{}) (int64, error) {
	res, err := s.db.ExecContext(ctx,
		`INSERT INTO hives(apiary_id, code, status, installed_at) VALUES(?,?,?,?)`,
		apiaryID, code, status, installedAt)
	if err != nil {
		return 0, fmt.Errorf("hive create: %w", err)
	}
	return res.LastInsertId()
}

// UpdateStatus 更新蜂箱状态。
func (s *HiveStore) UpdateStatus(ctx context.Context, id int64, status string) error {
	res, err := s.db.ExecContext(ctx, `UPDATE hives SET status = ? WHERE id = ?`, status, id)
	if err != nil {
		return fmt.Errorf("hive update status: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return model.ErrHiveNotFound
	}
	return nil
}

// Delete 删除蜂箱。
func (s *HiveStore) Delete(ctx context.Context, id int64) error {
	res, err := s.db.ExecContext(ctx, `DELETE FROM hives WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("hive delete: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return model.ErrHiveNotFound
	}
	return nil
}
