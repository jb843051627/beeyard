package service

import (
	"context"
	"testing"
	"time"

	"github.com/jb843051627/beeyard/internal/model"
)

// TestBug07_BatchCreateRollbackOnError 验证：
// 批量创建维护任务时，若某条出错应回滚整批——
// 数据库中不应残留部分任务。
func TestBug07_BatchCreateRollbackOnError(t *testing.T) {
	st, svc := newTestService(t)
	apiaryID := seedApiary(t, st)
	hiveID := seedHive(t, st, apiaryID, "H-007")

	tasks := []model.MaintenanceTask{
		{HiveID: hiveID, Description: "task-1", ScheduledFor: time.Date(2024, 7, 1, 0, 0, 0, 0, time.UTC)},
		{HiveID: -1, Description: "bad-hive", ScheduledFor: time.Date(2024, 7, 2, 0, 0, 0, 0, time.UTC)},
		{HiveID: hiveID, Description: "task-3", ScheduledFor: time.Date(2024, 7, 3, 0, 0, 0, 0, time.UTC)},
	}

	_, err := svc.Maintenance.BatchCreate(context.Background(), tasks)
	if err == nil {
		t.Fatal("expected error from batch with invalid hive_id, got nil")
	}

	// 验证数据库中没有残留——坏批次应回滚
	pending, err := svc.Maintenance.ListPending(context.Background(), hiveID)
	if err != nil {
		t.Fatal(err)
	}
	for _, m := range pending {
		if m.Description == "task-1" || m.Description == "task-3" {
			t.Fatalf("found residual task %q after failed batch (expected rollback)", m.Description)
		}
	}
}
