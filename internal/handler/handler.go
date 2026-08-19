package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/jb843051627/beeyard/internal/model"
	"github.com/jb843051627/beeyard/internal/service"
)

// Handler HTTP 路由容器。
type Handler struct {
	Hives       *HiveHandler
	Readings    *ReadingHandler
	Alerts      *AlertHandler
	Inspections *InspectionHandler
	Harvests    *HarvestHandler
	Transfers   *TransferHandler
	Maintenance *MaintenanceHandler
}

func NewHandler(svc *service.Service) *Handler {
	return &Handler{
		Hives:       NewHiveHandler(svc.Hives),
		Readings:    NewReadingHandler(svc.Readings),
		Alerts:      NewAlertHandler(svc.Alerts),
		Inspections: NewInspectionHandler(svc.Inspections),
		Harvests:    NewHarvestHandler(svc.Harvests),
		Transfers:   NewTransferHandler(svc.Transfers),
		Maintenance: NewMaintenanceHandler(svc.Maintenance),
	}
}

// Routes 注册全部 HTTP 路由（Go 1.22 增强路由）。
func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/hives/{id}", h.Hives.Get)
	mux.HandleFunc("GET /api/apiaries/{apiary_id}/hives", h.Hives.ListByApiary)
	mux.HandleFunc("POST /api/hives", h.Hives.Create)
	mux.HandleFunc("PUT /api/hives/{id}/status", h.Hives.UpdateStatus)
	mux.HandleFunc("POST /api/readings", h.Readings.Record)
	mux.HandleFunc("POST /api/readings/batch", h.Readings.BatchIngest)
	mux.HandleFunc("GET /api/readings/export", h.Readings.ExportCSV)
	mux.HandleFunc("POST /api/alerts", h.Alerts.Create)
	mux.HandleFunc("GET /api/apiaries/{apiary_id}/alerts", h.Alerts.ListAlerts)
	mux.HandleFunc("PUT /api/alerts/{id}/ack", h.Alerts.Acknowledge)
	mux.HandleFunc("POST /api/inspections", h.Inspections.Schedule)
	mux.HandleFunc("PUT /api/inspections/{id}/complete", h.Inspections.Complete)
	mux.HandleFunc("GET /api/hives/{id}/inspections", h.Inspections.ListByHive)
	mux.HandleFunc("POST /api/harvests", h.Harvests.Record)
	mux.HandleFunc("GET /api/apiaries/{apiary_id}/harvests/report", h.Harvests.Report)
	mux.HandleFunc("POST /api/transfers", h.Transfers.Create)
	mux.HandleFunc("PUT /api/transfers/{id}/assign-queen/{queen_id}", h.Transfers.AssignQueen)
	mux.HandleFunc("PUT /api/transfers/{id}/complete", h.Transfers.Complete)
	mux.HandleFunc("POST /api/maintenance/batch", h.Maintenance.BatchCreate)
	mux.HandleFunc("PUT /api/maintenance/{id}/complete", h.Maintenance.Complete)
	mux.Handle("GET /", http.FileServer(http.Dir("web")))
	return mux
}

func writeJSON(w http.ResponseWriter, code int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

// writeError 统一错误映射：ValidationError→400，NotFound→404，其余→500。
func writeError(w http.ResponseWriter, err error) {
	var ve *model.ValidationError
	if errors.As(err, &ve) {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": ve.Error(), "field": ve.Field})
		return
	}
	for _, sentinel := range []error{
		model.ErrHiveNotFound, model.ErrQueenNotFound, model.ErrAlertNotFound,
		model.ErrInspectionNotFound, model.ErrMaintenanceNotFound,
		model.ErrHarvestNotFound, model.ErrTransferNotFound, model.ErrReadingNotFound,
	} {
		if errors.Is(err, sentinel) {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
			return
		}
	}
	writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
}

func writeCreated(w http.ResponseWriter, id int64) {
	writeJSON(w, http.StatusCreated, map[string]int64{"id": id})
}
