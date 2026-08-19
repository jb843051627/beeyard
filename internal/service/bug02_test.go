package service

import (
	"context"
	"testing"

	"github.com/jb843051627/beeyard/internal/model"
)

// TestBug02_ListAlertsSortDoesNotPolluteCache 验证：
// service 对告警列表排序不会污染 store 缓存——
// 排序后再查 store 应返回原始 id 升序，而非被排序后的顺序。
func TestBug02_ListAlertsSortDoesNotPolluteCache(t *testing.T) {
	st, svc := newTestService(t)
	apiaryID := seedApiary(t, st)
	hiveID := seedHive(t, st, apiaryID, "H-002")

	ctx := context.Background()

	// 创建 3 条告警（id 递增 = created_at 递增）
	for i := 0; i < 3; i++ {
		seedAlert(t, st, apiaryID, hiveID, model.AlertLevelInfo, "alert-"+string(rune('A'+i)))
	}

	// 第一次从 store 获取原始顺序（id 升序）
	first, err := st.Alerts.ListByApiary(ctx, apiaryID)
	if err != nil {
		t.Fatal(err)
	}
	if len(first) != 3 {
		t.Fatalf("expected 3 alerts, got %d", len(first))
	}

	// 调 service 排序（按 created_at 倒序）
	_, err = svc.Alerts.ListAlerts(ctx, apiaryID)
	if err != nil {
		t.Fatal(err)
	}

	// 再次从 store 获取——应仍为 id 升序（未被排序污染）
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
