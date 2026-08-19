package service

import (
	"context"
	"testing"
	"time"
)

func TestBug02_ListAlertsSortDoesNotPolluteCache(t *testing.T) {
	st, svc := newTestService(t)
	apiaryID := seedApiary(t, st)
	hiveID := seedHive(t, st, apiaryID, "H-002")

	ctx := context.Background()
	db := st.DB()

	for i := 0; i < 3; i++ {
		_, err := db.ExecContext(ctx,
			`INSERT INTO alerts(apiary_id, hive_id, level, message, status, created_at) VALUES(?,?,?,?,?,?)`,
			apiaryID, hiveID, "info", "alert", "active",
			time.Date(2024, 6, 1, 0, i, 0, 0, time.UTC))
		if err != nil {
			t.Fatal(err)
		}
	}

	first, err := st.Alerts.ListByApiary(ctx, apiaryID)
	if err != nil {
		t.Fatal(err)
	}
	if len(first) != 3 {
		t.Fatalf("expected 3 alerts, got %d", len(first))
	}

	_, err = svc.Alerts.ListAlerts(ctx, apiaryID)
	if err != nil {
		t.Fatal(err)
	}

	second, err := st.Alerts.ListByApiary(ctx, apiaryID)
	if err != nil {
		t.Fatal(err)
	}
	for i := 1; i < len(second); i++ {
		if second[i].ID < second[i-1].ID {
			t.Fatalf("cache polluted: expected id ascending after service sort, "+
				"got id[%d]=%d < id[%d]=%d", i, second[i].ID, i-1, second[i-1].ID)
		}
	}
}
