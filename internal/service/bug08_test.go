package service

import (
	"context"
	"testing"
)

// TestBug08_GenerateReportSortDoesNotPolluteCache 验证：
// service 对采蜜报表排序不会污染 store 缓存——
// 排序后再查 store 应返回原始 id 升序。
func TestBug08_GenerateReportSortDoesNotPolluteCache(t *testing.T) {
	st, svc := newTestService(t)
	apiaryID := seedApiary(t, st)
	hiveID := seedHive(t, st, apiaryID, "H-008")

	ctx := context.Background()

	// 创建 3 条采蜜记录（amount 各不同，id 递增）
	seedHarvest(t, st, apiaryID, hiveID, 5.0)
	seedHarvest(t, st, apiaryID, hiveID, 15.0)
	seedHarvest(t, st, apiaryID, hiveID, 10.0)

	// 第一次从 store 获取（id 升序）
	first, err := st.Harvests.ListByApiary(ctx, apiaryID)
	if err != nil {
		t.Fatal(err)
	}
	if len(first) != 3 {
		t.Fatalf("expected 3 harvests, got %d", len(first))
	}

	// 调 service 排序（按 amount 降序）
	_, err = svc.Harvests.GenerateReport(ctx, apiaryID)
	if err != nil {
		t.Fatal(err)
	}

	// 再次从 store 获取——应仍为 id 升序
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
