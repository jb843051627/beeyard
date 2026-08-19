package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jb843051627/beeyard/internal/cache"
	"github.com/jb843051627/beeyard/internal/model"
	"github.com/jb843051627/beeyard/internal/store"
)

// beeyardTZ 蜂场所在地时区（UTC+8，固定偏移，不依赖系统时区）。
var beeyardTZ = time.FixedZone("CST", 8*3600)

// ReadingService 读数业务。
type ReadingService struct {
	store *store.ReadingStore
	cache *cache.ReadingCache
}

func NewReadingService(s *store.ReadingStore, c *cache.ReadingCache) *ReadingService {
	return &ReadingService{store: s, cache: c}
}

// Record 单条上报，同步更新缓存。
func (s *ReadingService) Record(ctx context.Context, hiveID int64, sensorType string, value float64, recordedAt time.Time) (int64, error) {
	switch sensorType {
	case model.SensorTypeTemp, model.SensorTypeHumidity, model.SensorTypeWeight:
	default:
		return 0, model.NewValidationError("sensor_type", "invalid sensor type: "+sensorType)
	}
	id, err := s.store.Record(ctx, hiveID, sensorType, value, recordedAt)
	if err != nil {
		return 0, err
	}
	s.cache.Update(&model.SensorReading{ID: id, HiveID: hiveID, SensorType: sensorType, Value: value, RecordedAt: recordedAt})
	return id, nil
}

// BatchIngest 批量导入；逐条检查 ctx 取消，取消后立即返回已处理量与错误。
func (s *ReadingService) BatchIngest(ctx context.Context, batch []model.ReadingBatch) (int64, error) {
	var total int64
	for _, r := range batch {
		if err := ctx.Err(); err != nil {
			return total, fmt.Errorf("batch ingest cancelled: %w", err)
		}
		id, err := s.store.Record(ctx, r.HiveID, r.SensorType, r.Value, r.RecordedAt)
		if err != nil {
			return total, err
		}
		s.cache.Update(&model.SensorReading{ID: id, HiveID: r.HiveID, SensorType: r.SensorType, Value: r.Value, RecordedAt: r.RecordedAt})
		total++
	}
	return total, nil
}

// GetLatest 取蜂箱某传感器最新读数（先查缓存，miss 再查库）。
func (s *ReadingService) GetLatest(ctx context.Context, hiveID int64, sensorType string) (*model.SensorReading, error) {
	if r, ok := s.cache.Get(hiveID, sensorType); ok {
		return r, nil
	}
	return s.store.LatestByHive(ctx, hiveID, sensorType)
}

// ListByHive 按蜂箱列读数。
func (s *ReadingService) ListByHive(ctx context.Context, hiveID int64, limit int) ([]*model.SensorReading, error) {
	return s.store.ListByHive(ctx, hiveID, limit)
}

// ExportCSV 导出蜂场读数 CSV；时间用蜂场时区（CST UTC+8）格式化。
func (s *ReadingService) ExportCSV(ctx context.Context, apiaryID int64, from, to time.Time) (string, error) {
	readings, err := s.store.ListByApiaryAndTimeRange(ctx, apiaryID, from, to)
	if err != nil {
		return "", err
	}
	var sb strings.Builder
	sb.WriteString("id,hive_id,sensor_type,value,timestamp\n")
	for _, r := range readings {
		ts := r.RecordedAt.UTC().Format("2006-01-02 15:04:05")
		sb.WriteString(fmt.Sprintf("%d,%d,%s,%.2f,%s\n",
			r.ID, r.HiveID, r.SensorType, r.Value, ts))
	}
	return sb.String(), nil
}

// EvaluateThresholds 检查读数是否超阈值并返回待创建告警。
func (s *ReadingService) EvaluateThresholds(ctx context.Context, hiveID int64, apiaryID int64) ([]*model.Alert, error) {
	var alerts []*model.Alert
	if r, ok := s.cache.Get(hiveID, model.SensorTypeTemp); ok {
		if r.Value > 38.0 {
			alerts = append(alerts, &model.Alert{
				ApiaryID: apiaryID, HiveID: hiveID, Level: model.AlertLevelCritical,
				Message: fmt.Sprintf("temperature %.1f exceeds 38°C", r.Value), Status: model.AlertStatusActive,
			})
		}
	}
	if r, ok := s.cache.Get(hiveID, model.SensorTypeHumidity); ok {
		if r.Value < 40.0 {
			alerts = append(alerts, &model.Alert{
				ApiaryID: apiaryID, HiveID: hiveID, Level: model.AlertLevelWarning,
				Message: fmt.Sprintf("humidity %.1f below 40%%", r.Value), Status: model.AlertStatusActive,
			})
		}
	}
	return alerts, nil
}
