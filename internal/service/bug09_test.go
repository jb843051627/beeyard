package service

import (
	"context"
	"errors"
	"testing"

	"github.com/jb843051627/beeyard/internal/model"
)

// TestBug09_CompleteInspectionValidatesStatus 验证：
// 对已完成巡检再次调用 Complete 应返回 ValidationError，
// 而非静默覆盖完成状态。
func TestBug09_CompleteInspectionValidatesStatus(t *testing.T) {
	st, svc := newTestService(t)
	apiaryID := seedApiary(t, st)
	hiveID := seedHive(t, st, apiaryID, "H-009")

	ctx := context.Background()

	// 创建巡检并先完成
	insID := seedInspection(t, st, hiveID, model.InspectionStatusCompleted)

	// 再次完成 → 应报错
	err := svc.Inspections.Complete(ctx, insID, "second time")
	if err == nil {
		t.Fatal("expected error when completing already-completed inspection, got nil")
	}
	var ve *model.ValidationError
	if !errors.As(err, &ve) {
		t.Fatalf("expected ValidationError for duplicate completion, got: %v", err)
	}
}
