package service

import (
	"context"
	"errors"
	"testing"

	"github.com/jb843051627/beeyard/internal/model"
)

func TestBug01_HiveGetByIDReturnsError(t *testing.T) {
	st, svc := newTestService(t)
	apiaryID := seedApiary(t, st)
	hiveID := seedHive(t, st, apiaryID, "H-001")

	ctx := context.Background()

	h, err := svc.Hives.Get(ctx, hiveID)
	if err != nil || h == nil {
		t.Fatalf("get existing hive %d: err=%v h=%v", hiveID, err, h)
	}

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
