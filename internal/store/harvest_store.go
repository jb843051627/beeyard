package store

import (
	"context"
	"database/sql"
	"fmt"
	"sync"

	"github.com/jb843051627/beeyard/internal/model"
)

// HarvestStore 采蜜记录持久化，带按蜂场缓存。
type HarvestStore struct {
	db    *sql.DB
	cache map[int64][]*model.Harvest
	mu    sync.RWMutex
}

func NewHarvestStore(db *sql.DB) *HarvestStore {
	return &HarvestStore{db: db, cache: make(map[int64][]*model.Harvest)}
}

func (s *HarvestStore) invalidate(apiaryID int64) {
	s.mu.Lock()
	delete(s.cache, apiaryID)
	s.mu.Unlock()
}

// Create 新建采蜜记录并失效缓存。
func (s *HarvestStore) Create(ctx context.Context, h *model.Harvest) (int64, error) {
	res, err := s.db.ExecContext(ctx,
		`INSERT INTO harvests(apiary_id, hive_id, amount_kg, harvested_at) VALUES(?,?,?,?)`,
		h.ApiaryID, h.HiveID, h.AmountKg, h.HarvestedAt)
	if err != nil {
		return 0, fmt.Errorf("harvest create: %w", err)
	}
	s.invalidate(h.ApiaryID)
	return res.LastInsertId()
}

// GetByID 查采蜜记录；不存在返回 (nil, ErrHarvestNotFound)。
func (s *HarvestStore) GetByID(ctx context.Context, id int64) (*model.Harvest, error) {
	row := s.db.QueryRowContext(ctx,
		`SELECT id, apiary_id, hive_id, amount_kg, harvested_at, created_at FROM harvests WHERE id = ?`, id)
	var h model.Harvest
	err := row.Scan(&h.ID, &h.ApiaryID, &h.HiveID, &h.AmountKg, &h.HarvestedAt, &h.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, model.ErrHarvestNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("harvest get: %w", err)
	}
	return &h, nil
}

// ListByApiary 列出蜂场采蜜记录；命中缓存时返回副本。
func (s *HarvestStore) ListByApiary(ctx context.Context, apiaryID int64) ([]*model.Harvest, error) {
	s.mu.RLock()
	cached, ok := s.cache[apiaryID]
	s.mu.RUnlock()
	if ok {
		out := make([]*model.Harvest, len(cached))
		copy(out, cached)
		return out, nil
	}
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, apiary_id, hive_id, amount_kg, harvested_at, created_at FROM harvests WHERE apiary_id = ? ORDER BY id`, apiaryID)
	if err != nil {
		return nil, fmt.Errorf("harvest list by apiary: %w", err)
	}
	defer rows.Close()
	var list []*model.Harvest
	for rows.Next() {
		var h model.Harvest
		if err := rows.Scan(&h.ID, &h.ApiaryID, &h.HiveID, &h.AmountKg, &h.HarvestedAt, &h.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, &h)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	s.mu.Lock()
	s.cache[apiaryID] = list
	s.mu.Unlock()
	out := make([]*model.Harvest, len(list))
	copy(out, list)
	return out, nil
}
