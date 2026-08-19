package store

import (
	"context"
	"fmt"
)

// CountHivesByApiary 统计蜂场蜂箱数。
func (s *HiveStore) CountHivesByApiary(ctx context.Context, apiaryID int64) (int64, error) {
	var n int64
	err := s.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM hives WHERE apiary_id = ?`, apiaryID).Scan(&n)
	if err != nil {
		return 0, fmt.Errorf("count hives: %w", err)
	}
	return n, nil
}

// CountActiveAlerts 统计蜂场活跃告警数。
func (s *AlertStore) CountActiveAlerts(ctx context.Context, apiaryID int64) (int64, error) {
	var n int64
	err := s.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM alerts WHERE apiary_id = ? AND status = ?`, apiaryID, "active").Scan(&n)
	if err != nil {
		return 0, fmt.Errorf("count active alerts: %w", err)
	}
	return n, nil
}

// CountPendingInspections 统计蜂箱待巡检数。
func (s *InspectionStore) CountPendingInspections(ctx context.Context, hiveID int64) (int64, error) {
	var n int64
	err := s.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM inspections WHERE hive_id = ? AND status = ?`, hiveID, "pending").Scan(&n)
	if err != nil {
		return 0, fmt.Errorf("count pending inspections: %w", err)
	}
	return n, nil
}

// SumHarvestAmount 统计蜂场累计采蜜量。
func (s *HarvestStore) SumHarvestAmount(ctx context.Context, apiaryID int64) (float64, error) {
	var total float64
	err := s.db.QueryRowContext(ctx,
		`SELECT COALESCE(SUM(amount_kg), 0) FROM harvests WHERE apiary_id = ?`, apiaryID).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("sum harvest: %w", err)
	}
	return total, nil
}

// AvgReadingBySensor 取蜂箱某传感器平均读数。
func (s *ReadingStore) AvgReadingBySensor(ctx context.Context, hiveID int64, sensorType string) (float64, error) {
	var avg float64
	err := s.db.QueryRowContext(ctx,
		`SELECT COALESCE(AVG(value), 0) FROM readings WHERE hive_id = ? AND sensor_type = ?`, hiveID, sensorType).Scan(&avg)
	if err != nil {
		return 0, fmt.Errorf("avg reading: %w", err)
	}
	return avg, nil
}
