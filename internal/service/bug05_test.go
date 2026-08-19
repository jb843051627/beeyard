package service

import (
	"context"
	"testing"
	"time"

	"github.com/jb843051627/beeyard/internal/model"
)

// TestBug05_BatchIngestRespectsContextCancellation 验证：
// 传入已取消的 context 时，BatchIngest 应立即返回 0 条已处理 + 错误，
// 而非无视取消继续处理全部数据。
func TestBug05_BatchIngestRespectsContextCancellation(t *testing.T) {
	st, svc := newTestService(t)
	apiaryID := seedApiary(t, st)
	hiveID := seedHive(t, st, apiaryID, "H-005")

	batch := make([]model.ReadingBatch, 50)
	for i := range batch {
		batch[i] = model.ReadingBatch{
			HiveID:     hiveID,
			SensorType: model.SensorTypeTemp,
			Value:      float64(i),
			RecordedAt: time.Date(2024, 6, 1, 0, i, 0, 0, time.UTC),
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	total, err := svc.Readings.BatchIngest(ctx, batch)
	if err == nil {
		t.Fatal("expected error from cancelled context, got nil (context cancellation ignored)")
	}
	if total > 0 {
		t.Fatalf("expected 0 ingested with cancelled context, got %d", total)
	}
}
