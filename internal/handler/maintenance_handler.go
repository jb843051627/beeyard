package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/jb843051627/beeyard/internal/model"
	"github.com/jb843051627/beeyard/internal/service"
)

type MaintenanceHandler struct {
	svc *service.MaintenanceService
}

func NewMaintenanceHandler(s *service.MaintenanceService) *MaintenanceHandler {
	return &MaintenanceHandler{svc: s}
}

func (h *MaintenanceHandler) BatchCreate(w http.ResponseWriter, r *http.Request) {
	var tasks []model.MaintenanceTask
	if err := json.NewDecoder(r.Body).Decode(&tasks); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid body"})
		return
	}
	total, err := h.svc.BatchCreate(r.Context(), tasks)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]int64{"created": total})
}

func (h *MaintenanceHandler) Complete(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err := h.svc.Complete(r.Context(), id); err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "completed"})
}


