package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/jb843051627/beeyard/internal/model"
	"github.com/jb843051627/beeyard/internal/service"
)

type AlertHandler struct {
	svc *service.AlertService
}

func NewAlertHandler(s *service.AlertService) *AlertHandler { return &AlertHandler{svc: s} }

// Create 新建告警。
func (h *AlertHandler) Create(w http.ResponseWriter, r *http.Request) {
	var a model.Alert
	if err := json.NewDecoder(r.Body).Decode(&a); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid body"})
		return
	}
	id, err := h.svc.Create(r.Context(), &a)
	if err != nil {
		writeError(w, err)
		return
	}
	writeCreated(w, id)
}

func (h *AlertHandler) ListAlerts(w http.ResponseWriter, r *http.Request) {
	aid, _ := strconv.ParseInt(r.PathValue("apiary_id"), 10, 64)
	list, err := h.svc.ListAlerts(r.Context(), aid)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, list)
}

func (h *AlertHandler) Acknowledge(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err := h.svc.Acknowledge(r.Context(), id); err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "acknowledged"})
}
