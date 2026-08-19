package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/jb843051627/beeyard/internal/service"
)

type InspectionHandler struct {
	svc *service.InspectionService
}

func NewInspectionHandler(s *service.InspectionService) *InspectionHandler {
	return &InspectionHandler{svc: s}
}

func (h *InspectionHandler) Schedule(w http.ResponseWriter, r *http.Request) {
	var req struct {
		HiveID      int64     `json:"hive_id"`
		ScheduledAt time.Time `json:"scheduled_at"`
		Notes       string    `json:"notes"`
	}
	_ = json.NewDecoder(r.Body).Decode(&req)
	id, err := h.svc.Schedule(r.Context(), req.HiveID, req.ScheduledAt, req.Notes)
	if err != nil {
		writeError(w, err)
		return
	}
	writeCreated(w, id)
}

// Complete 完成巡检。错误经 writeError 映射：ValidationError→400（bug-009 修复面）。
func (h *InspectionHandler) Complete(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	var req struct {
		Notes string `json:"notes"`
	}
	_ = json.NewDecoder(r.Body).Decode(&req)
	if err := h.svc.Complete(r.Context(), id, req.Notes); err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "completed"})
}

func (h *InspectionHandler) ListByHive(w http.ResponseWriter, r *http.Request) {
	hid, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	list, err := h.svc.ListByHive(r.Context(), hid)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, list)
}
