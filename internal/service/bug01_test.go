package service

import (
	"context"
	"errors"
	"testing"

	"github.com/jb843051627/beeyard/internal/model"
)

// TestBug01_HiveGetByIDReturnsError 验证：不存在的蜂箱 id 查询时，
// store 返回 (nil, ErrHiveNotFound) 且 service 正确传播错误，
// 调用方不会拿到 nil 蜂箱而解引用 panic。
func TestBug01_HiveGetByIDReturnsError(t *testing.T) {
	st, svc := newTestService(t)
	apiaryID := seedApiary(t, st)
	hiveID := seedHive(t, st, apiaryID, "H-001")

	ctx := context.Background()

	// 存在的蜂箱应正常返回
	h, err := svc.Hives.Get(ctx, hiveID)
	if err != nil || h == nil {
		t.Fatalf("get existing hive %d: err=%v h=%v", hiveID, err, h)
	}

	// 不存在的蜂箱应返回可识别的错误，而非 (nil, nil)
	missing, err := svc.Hives.Get(ctx, 999999)
	if err == nil {
		t.Fatalf("expected error for non-existent hive, got hive=%v err=nil", missing)
	}
	if missing != nil {
		t.Fatalf("expected nil hive for non-existent id, got %v", missing)
	}
	if !errors.Is(err, model.ErrHiveNotFound) {
		t.Fatalf("expected errors.Is(err, ErrHiveNotFound), got: %v", err)
	}
}
