package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/jb843051627/beeyard/internal/model"
	"github.com/jb843051627/beeyard/internal/service"
)

type HarvestHandler struct {
	svc *service.HarvestService
}

func NewHarvestHandler(s *service.HarvestService) *HarvestHandler {
	return &HarvestHandler{svc: s}
}

func (h *HarvestHandler) Record(w http.ResponseWriter, r *http.Request) {
	var hv model.Harvest
	if err := json.NewDecoder(r.Body).Decode(&hv); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid body"})
		return
	}
	id, err := h.svc.Record(r.Context(), &hv)
	if err != nil {
		writeError(w, err)
		return
	}
	writeCreated(w, id)
}

func (h *HarvestHandler) Report(w http.ResponseWriter, r *http.Request) {
	aid, _ := strconv.ParseInt(r.PathValue("apiary_id"), 10, 64)
	report, err := h.svc.GenerateReport(r.Context(), aid)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, report)
}
