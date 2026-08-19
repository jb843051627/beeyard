package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/jb843051627/beeyard/internal/model"
)

func TestBug06_AssignQueenRejectsMissingQueen(t *testing.T) {
	st, svc := newTestService(t)
	apiaryID := seedApiary(t, st)
	hiveID := seedHive(t, st, apiaryID, "H-006")

	tid, err := svc.Transfers.Create(context.Background(), hiveID, apiaryID, apiaryID+1, time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatal(err)
	}

	err = svc.Transfers.AssignQueen(context.Background(), tid, 999999)
	if err == nil {
		t.Fatal("expected error for non-existent queen, got nil")
	}
	if !errors.Is(err, model.ErrQueenNotFound) {
		t.Fatalf("expected errors.Is(err, ErrQueenNotFound), got: %v", err)
	}
}
