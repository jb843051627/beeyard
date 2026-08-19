package cache

import (
	"sync"

	"github.com/jb843051627/beeyard/internal/model"
)

// ReadingCache 最新读数缓存，按蜂箱+传感器类型索引。
type ReadingCache struct {
	mu     sync.RWMutex
	latest map[int64]map[string]*model.SensorReading
}

func NewReadingCache() *ReadingCache {
	return &ReadingCache{latest: make(map[int64]map[string]*model.SensorReading)}
}

// Update 写入最新读数（写操作，必须持写锁）。
func (c *ReadingCache) Update(r *model.SensorReading) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.latest[r.HiveID] == nil {
		c.latest[r.HiveID] = make(map[string]*model.SensorReading)
	}
	c.latest[r.HiveID][r.SensorType] = r
}

// Get 取蜂箱某传感器最新读数。
func (c *ReadingCache) Get(hiveID int64, sensorType string) (*model.SensorReading, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if m, ok := c.latest[hiveID]; ok {
		if r, ok := m[sensorType]; ok {
			return r, true
		}
	}
	return nil, false
}

// Snapshot 返回全部最新读数副本。
func (c *ReadingCache) Snapshot() map[int64]map[string]*model.SensorReading {
	c.mu.RLock()
	defer c.mu.RUnlock()
	out := make(map[int64]map[string]*model.SensorReading, len(c.latest))
	for k, v := range c.latest {
		nv := make(map[string]*model.SensorReading, len(v))
		for sk, sv := range v {
			nv[sk] = sv
		}
		out[k] = nv
	}
	return out
}
