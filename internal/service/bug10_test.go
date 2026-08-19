package service

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/jb843051627/beeyard/internal/model"
)

func TestBug10_ExportCSVUsesCSTTimezone(t *testing.T) {
	st, svc := newTestService(t)
	apiaryID := seedApiary(t, st)
	hiveID := seedHive(t, st, apiaryID, "H-010")

	cst := time.FixedZone("CST", 8*3600)
	recordedAt := time.Date(2024, 6, 15, 10, 0, 0, 0, cst)

	ctx := context.Background()
	_, err := svc.Readings.Record(ctx, hiveID, model.SensorTypeTemp, 25.5, recordedAt)
	if err != nil {
		t.Fatal(err)
	}

	from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	csv, err := svc.Readings.ExportCSV(ctx, apiaryID, from, to)
	if err != nil {
		t.Fatal(err)
	}

	if strings.Contains(csv, "02:00:00") {
		t.Fatalf("ExportCSV appears to use UTC (02:00:00) instead of CST (10:00:00);\n%s", csv)
	}
	if !strings.Contains(csv, "10:00:00") {
		t.Fatalf("ExportCSV should contain CST time '10:00:00';\n%s", csv)
	}
}
