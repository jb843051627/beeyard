package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/jb843051627/beeyard/internal/model"
)

// TestBug06_AssignQueenRejectsMissingQueen 验证：
// 为不存在的蜂王 id 执行转场分配时，service 应返回可识别错误，
// 而非拿到 nil 蜂王解引用 panic。
func TestBug06_AssignQueenRejectsMissingQueen(t *testing.T) {
	st, svc := newTestService(t)
	apiaryID := seedApiary(t, st)
	hiveID := seedHive(t, st, apiaryID, "H-006")

	// 创建转场请求
	tid, err := svc.Transfers.Create(context.Background(), hiveID, apiaryID, apiaryID+1, time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatal(err)
	}

	// 为不存在的蜂王分配 → 应报错而非 nil panic
	err = svc.Transfers.AssignQueen(context.Background(), tid, 999999)
	if err == nil {
		t.Fatal("expected error for non-existent queen, got nil")
	}
	if !errors.Is(err, model.ErrQueenNotFound) {
		t.Fatalf("expected errors.Is(err, ErrQueenNotFound), got: %v", err)
	}
}
