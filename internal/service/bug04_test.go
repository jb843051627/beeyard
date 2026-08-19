package service

import (
	"context"
	"errors"
	"testing"

	"github.com/jb843051627/beeyard/internal/model"
)

func TestBug04_AlertCreateValidatesLevel(t *testing.T) {
	st, svc := newTestService(t)
	apiaryID := seedApiary(t, st)
	hiveID := seedHive(t, st, apiaryID, "H-004")

	_, err := svc.Alerts.Create(context.Background(), &model.Alert{
		ApiaryID: apiaryID, HiveID: hiveID, Level: "bogus", Message: "test",
	})
	if err == nil {
		t.Fatal("expected error for invalid alert level 'bogus', got nil")
	}
	var ve *model.ValidationError
	if !errors.As(err, &ve) {
		t.Fatalf("expected errors.As to detect ValidationError, got: %v", err)
	}
	if ve.Field != "level" {
		t.Fatalf("expected field=level, got %s", ve.Field)
	}
}
