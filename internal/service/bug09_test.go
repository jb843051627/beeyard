package service

import (
	"context"
	"errors"
	"testing"

	"github.com/jb843051627/beeyard/internal/model"
)

func TestBug09_CompleteInspectionValidatesStatus(t *testing.T) {
	st, svc := newTestService(t)
	apiaryID := seedApiary(t, st)
	hiveID := seedHive(t, st, apiaryID, "H-009")

	ctx := context.Background()

	insID := seedInspection(t, st, hiveID, model.InspectionStatusCompleted)

	err := svc.Inspections.Complete(ctx, insID, "second time")
	if err == nil {
		t.Fatal("expected error when completing already-completed inspection, got nil")
	}
	var ve *model.ValidationError
	if !errors.As(err, &ve) {
		t.Fatalf("expected ValidationError for duplicate completion, got: %v", err)
	}
}
