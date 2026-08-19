package store

import (
	"context"
	"database/sql"
	"fmt"
	"sync"

	"github.com/jb843051627/beeyard/internal/model"
)

// AlertStore 告警持久化，带按蜂场缓存。
type AlertStore struct {
	db    *sql.DB
	cache map[int64][]*model.Alert
	mu    sync.RWMutex
}

func NewAlertStore(db *sql.DB) *AlertStore {
	return &AlertStore{db: db, cache: make(map[int64][]*model.Alert)}
}

func (s *AlertStore) invalidate(apiaryID int64) {
	s.mu.Lock()
	delete(s.cache, apiaryID)
	s.mu.Unlock()
}

// Create 新建告警并失效缓存。
func (s *AlertStore) Create(ctx context.Context, a *model.Alert) (int64, error) {
	res, err := s.db.ExecContext(ctx,
		`INSERT INTO alerts(apiary_id, hive_id, level, message, status) VALUES(?,?,?,?,?)`,
		a.ApiaryID, a.HiveID, a.Level, a.Message, a.Status)
	if err != nil {
		return 0, fmt.Errorf("alert create: %w", err)
	}
	s.invalidate(a.ApiaryID)
	return res.LastInsertId()
}

// GetByID 查告警；不存在返回 (nil, ErrAlertNotFound)。
func (s *AlertStore) GetByID(ctx context.Context, id int64) (*model.Alert, error) {
	row := s.db.QueryRowContext(ctx,
		`SELECT id, apiary_id, hive_id, level, message, status, created_at FROM alerts WHERE id = ?`, id)
	var a model.Alert
	err := row.Scan(&a.ID, &a.ApiaryID, &a.HiveID, &a.Level, &a.Message, &a.Status, &a.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, model.ErrAlertNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("alert get: %w", err)
	}
	return &a, nil
}

// ListByApiary 列出蜂场告警；命中缓存时返回副本。
func (s *AlertStore) ListByApiary(ctx context.Context, apiaryID int64) ([]*model.Alert, error) {
	s.mu.RLock()
	cached, ok := s.cache[apiaryID]
	s.mu.RUnlock()
	if ok {
		return cached, nil
	}
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, apiary_id, hive_id, level, message, status, created_at FROM alerts WHERE apiary_id = ? ORDER BY id`, apiaryID)
	if err != nil {
		return nil, fmt.Errorf("alert list by apiary: %w", err)
	}
	defer rows.Close()
	var list []*model.Alert
	for rows.Next() {
		var a model.Alert
		if err := rows.Scan(&a.ID, &a.ApiaryID, &a.HiveID, &a.Level, &a.Message, &a.Status, &a.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, &a)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	s.mu.Lock()
	s.cache[apiaryID] = list
	s.mu.Unlock()
	return list, nil
}

// ListActive 列出未确认告警。
func (s *AlertStore) ListActive(ctx context.Context, apiaryID int64) ([]*model.Alert, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, apiary_id, hive_id, level, message, status, created_at FROM alerts WHERE apiary_id = ? AND status = ? ORDER BY created_at DESC`, apiaryID, model.AlertStatusActive)
	if err != nil {
		return nil, fmt.Errorf("alert list active: %w", err)
	}
	defer rows.Close()
	var list []*model.Alert
	for rows.Next() {
		var a model.Alert
		if err := rows.Scan(&a.ID, &a.ApiaryID, &a.HiveID, &a.Level, &a.Message, &a.Status, &a.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, &a)
	}
	return list, rows.Err()
}

// Acknowledge 确认告警。
func (s *AlertStore) Acknowledge(ctx context.Context, id int64) error {
	res, err := s.db.ExecContext(ctx, `UPDATE alerts SET status = ? WHERE id = ?`, model.AlertStatusAcknowledged, id)
	if err != nil {
		return fmt.Errorf("alert acknowledge: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return model.ErrAlertNotFound
	}
	var a model.Alert
	_ = s.db.QueryRowContext(ctx, `SELECT apiary_id FROM alerts WHERE id = ?`, id).Scan(&a.ApiaryID)
	s.invalidate(a.ApiaryID)
	return nil
}
