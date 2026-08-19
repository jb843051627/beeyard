package service

import (
	"context"
	"testing"
)

func TestBug08_GenerateReportSortDoesNotPolluteCache(t *testing.T) {
	st, svc := newTestService(t)
	apiaryID := seedApiary(t, st)
	hiveID := seedHive(t, st, apiaryID, "H-008")

	ctx := context.Background()

	seedHarvest(t, st, apiaryID, hiveID, 5.0)
	seedHarvest(t, st, apiaryID, hiveID, 15.0)
	seedHarvest(t, st, apiaryID, hiveID, 10.0)

	first, err := st.Harvests.ListByApiary(ctx, apiaryID)
	if err != nil {
		t.Fatal(err)
	}
	if len(first) != 3 {
		t.Fatalf("expected 3 harvests, got %d", len(first))
	}

	_, err = svc.Harvests.GenerateReport(ctx, apiaryID)
	if err != nil {
		t.Fatal(err)
	}

	second, err := st.Harvests.ListByApiary(ctx, apiaryID)
	if err != nil {
		t.Fatal(err)
	}
	for i := 1; i < len(second); i++ {
		if second[i].ID < second[i-1].ID {
			t.Fatalf("cache polluted: expected id ascending after report sort, "+
				"got id[%d]=%d < id[%d]=%d", i, second[i].ID, i-1, second[i-1].ID)
		}
	}
}
