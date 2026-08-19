package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/jb843051627/beeyard/internal/service"
)

type HiveHandler struct {
	svc *service.HiveService
}

func NewHiveHandler(s *service.HiveService) *HiveHandler { return &HiveHandler{svc: s} }

func (h *HiveHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	hive, err := h.svc.Get(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, hive)
}

func (h *HiveHandler) ListByApiary(w http.ResponseWriter, r *http.Request) {
	aid, _ := strconv.ParseInt(r.PathValue("apiary_id"), 10, 64)
	list, err := h.svc.ListByApiary(r.Context(), aid)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, list)
}

func (h *HiveHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ApiaryID    int64     `json:"apiary_id"`
		Code        string    `json:"code"`
		Status      string    `json:"status"`
		InstalledAt time.Time `json:"installed_at"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid body"})
		return
	}
	id, err := h.svc.Create(r.Context(), req.ApiaryID, req.Code, req.Status, req.InstalledAt)
	if err != nil {
		writeError(w, err)
		return
	}
	writeCreated(w, id)
}

func (h *HiveHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	var req struct {
		Status string `json:"status"`
	}
	_ = json.NewDecoder(r.Body).Decode(&req)
	if err := h.svc.UpdateStatus(r.Context(), id, req.Status); err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "updated"})
}


