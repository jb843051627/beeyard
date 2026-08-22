package service

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/jb843051627/beeyard/internal/cache"
	"github.com/jb843051627/beeyard/internal/model"
	"github.com/jb843051627/beeyard/internal/store"
)

func newTestService(t *testing.T) (*store.Store, *Service) {
	t.Helper()
	dbPath := filepath.Join(t.TempDir(), "test.db")
	st, err := store.NewStore(dbPath)
	if err != nil {
		t.Fatalf("new store: %v", err)
	}
	t.Cleanup(func() { st.Close() })
	rc := cache.NewReadingCache()
	return st, NewService(st, rc)
}

func seedApiary(t *testing.T, st *store.Store) int64 {
	t.Helper()
	id, err := st.SeedApiary("test-apiary", "test-location")
	if err != nil {
		t.Fatalf("seed apiary: %v", err)
	}
	return id
}

func seedHive(t *testing.T, st *store.Store, apiaryID int64, code string) int64 {
	t.Helper()
	id, err := st.Hives.Create(context.Background(), apiaryID, code, model.HiveStatusActive, time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("create hive %s: %v", code, err)
	}
	return id
}

func seedQueen(t *testing.T, st *store.Store, hiveID int64, breed string) int64 {
	t.Helper()
	id, err := st.Queens.Create(context.Background(), hiveID, breed, model.QueenStatusActive, time.Date(2024, 3, 5, 0, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("create queen: %v", err)
	}
	return id
}

func seedAlert(t *testing.T, st *store.Store, apiaryID, hiveID int64, level, msg string) int64 {
	t.Helper()
	id, err := st.Alerts.Create(context.Background(), &model.Alert{
		ApiaryID: apiaryID, HiveID: hiveID, Level: level, Message: msg, Status: model.AlertStatusActive,
	})
	if err != nil {
		t.Fatalf("create alert: %v", err)
	}
	return id
}

func seedHarvest(t *testing.T, st *store.Store, apiaryID, hiveID int64, amount float64) {
	t.Helper()
	_, err := st.Harvests.Create(context.Background(), &model.Harvest{
		ApiaryID: apiaryID, HiveID: hiveID, AmountKg: amount, HarvestedAt: time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatalf("create harvest: %v", err)
	}
}

func seedInspection(t *testing.T, st *store.Store, hiveID int64, status string) int64 {
	t.Helper()
	ctx := context.Background()
	id, err := st.Inspections.Create(ctx, hiveID, time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC), "routine")
	if err != nil {
		t.Fatalf("create inspection: %v", err)
	}
	if status == model.InspectionStatusCompleted {
		if err := st.Inspections.Complete(ctx, id, time.Now()); err != nil {
			t.Fatalf("complete inspection: %v", err)
		}
	}
	return id
}