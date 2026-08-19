package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/jb843051627/beeyard/internal/model"
	"github.com/jb843051627/beeyard/internal/service"
)

type ReadingHandler struct {
	svc *service.ReadingService
}

func NewReadingHandler(s *service.ReadingService) *ReadingHandler { return &ReadingHandler{svc: s} }

func (h *ReadingHandler) Record(w http.ResponseWriter, r *http.Request) {
	var req struct {
		HiveID     int64     `json:"hive_id"`
		SensorType string    `json:"sensor_type"`
		Value      float64   `json:"value"`
		RecordedAt time.Time `json:"recorded_at"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid body"})
		return
	}
	id, err := h.svc.Record(r.Context(), req.HiveID, req.SensorType, req.Value, req.RecordedAt)
	if err != nil {
		writeError(w, err)
		return
	}
	writeCreated(w, id)
}

// BatchIngest 批量导入。
func (h *ReadingHandler) BatchIngest(w http.ResponseWriter, r *http.Request) {
	var batch []model.ReadingBatch
	if err := json.NewDecoder(r.Body).Decode(&batch); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid body"})
		return
	}
	total, err := h.svc.BatchIngest(r.Context(), batch)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]int64{"ingested": total})
}

func (h *ReadingHandler) ExportCSV(w http.ResponseWriter, r *http.Request) {
	apiaryID, _ := strconv.ParseInt(r.URL.Query().Get("apiary_id"), 10, 64)
	from, _ := time.Parse(time.RFC3339, r.URL.Query().Get("from"))
	to, _ := time.Parse(time.RFC3339, r.URL.Query().Get("to"))
	csv, err := h.svc.ExportCSV(r.Context(), apiaryID, from, to)
	if err != nil {
		writeError(w, err)
		return
	}
	w.Header().Set("Content-Type", "text/csv")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(csv))
}
