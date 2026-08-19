package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/jb843051627/beeyard/internal/service"
)

type TransferHandler struct {
	svc *service.TransferService
}

func NewTransferHandler(s *service.TransferService) *TransferHandler {
	return &TransferHandler{svc: s}
}

func (h *TransferHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		HiveID     int64 `json:"hive_id"`
		FromApiary int64 `json:"from_apiary"`
		ToApiary   int64 `json:"to_apiary"`
	}
	_ = json.NewDecoder(r.Body).Decode(&req)
	id, err := h.svc.Create(r.Context(), req.HiveID, req.FromApiary, req.ToApiary, time.Now())
	if err != nil {
		writeError(w, err)
		return
	}
	writeCreated(w, id)
}

func (h *TransferHandler) AssignQueen(w http.ResponseWriter, r *http.Request) {
	tid, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	qid, _ := strconv.ParseInt(r.PathValue("queen_id"), 10, 64)
	if err := h.svc.AssignQueen(r.Context(), tid, qid); err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "queen assigned"})
}

func (h *TransferHandler) Complete(w http.ResponseWriter, r *http.Request) {
	tid, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err := h.svc.Complete(r.Context(), tid); err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "completed"})
}
