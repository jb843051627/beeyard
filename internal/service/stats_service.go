package service

import (
	"context"
	"fmt"

	"github.com/jb843051627/beeyard/internal/cache"
	"github.com/jb843051627/beeyard/internal/model"
	"github.com/jb843051627/beeyard/internal/store"
)

// StatsService 蜂场统计与聚合。
type StatsService struct {
	hives       *store.HiveStore
	readings    *store.ReadingStore
	alerts      *store.AlertStore
	inspections *store.InspectionStore
	harvests    *store.HarvestStore
	cache       *cache.ReadingCache
}

func NewStatsService(st *store.Store, rc *cache.ReadingCache) *StatsService {
	return &StatsService{
		hives:       st.Hives,
		readings:    st.Readings,
		alerts:      st.Alerts,
		inspections: st.Inspections,
		harvests:    st.Harvests,
		cache:       rc,
	}
}

// ApiaryOverview 蜂场概览聚合体。
type ApiaryOverview struct {
	HiveCount        int                    `json:"hive_count"`
	ActiveAlerts     int                    `json:"active_alerts"`
	PendingInspects  int                    `json:"pending_inspections"`
	LatestReadings   map[int64]map[string]float64 `json:"latest_readings"`
}

// GetApiaryOverview 聚合蜂场概览：蜂箱数、活跃告警数、待巡检数、各蜂箱最新读数。
func (s *StatsService) GetApiaryOverview(ctx context.Context, apiaryID int64) (*ApiaryOverview, error) {
	hives, err := s.hives.ListByApiary(ctx, apiaryID)
	if err != nil {
		return nil, fmt.Errorf("list hives for overview: %w", err)
	}
	activeAlerts, err := s.alerts.ListActive(ctx, apiaryID)
	if err != nil {
		return nil, fmt.Errorf("list active alerts for overview: %w", err)
	}
	overview := &ApiaryOverview{
		HiveCount:    len(hives),
		ActiveAlerts: len(activeAlerts),
		LatestReadings: make(map[int64]map[string]float64),
	}
	for _, h := range hives {
		snap := s.cache.Snapshot()
		if m, ok := snap[h.ID]; ok {
			readings := make(map[string]float64)
			for sensorType, r := range m {
				readings[sensorType] = r.Value
			}
			overview.LatestReadings[h.ID] = readings
		}
	}
	totalPending := 0
	for _, h := range hives {
		list, err := s.inspections.ListByHive(ctx, h.ID)
		if err != nil {
			continue
		}
		for _, ins := range list {
			if ins.Status == model.InspectionStatusPending {
				totalPending++
			}
		}
	}
	overview.PendingInspects = totalPending
	return overview, nil
}

// GetHiveTemperatureTrend 取蜂箱最近 N 条温度读数（倒序）。
func (s *StatsService) GetHiveTemperatureTrend(ctx context.Context, hiveID int64, limit int) ([]*model.SensorReading, error) {
	readings, err := s.readings.ListByHive(ctx, hiveID, limit)
	if err != nil {
		return nil, err
	}
	var trend []*model.SensorReading
	for _, r := range readings {
		if r.SensorType == model.SensorTypeTemp {
			trend = append(trend, r)
		}
	}
	return trend, nil
}

// GetTotalYield 统计蜂场累计采蜜量（kg）。
func (s *StatsService) GetTotalYield(ctx context.Context, apiaryID int64) (float64, error) {
	list, err := s.harvests.ListByApiary(ctx, apiaryID)
	if err != nil {
		return 0, err
	}
	var total float64
	for _, h := range list {
		total += h.AmountKg
	}
	return total, nil
}
