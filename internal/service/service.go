package service

import (
	"github.com/jb843051627/beeyard/internal/cache"
	"github.com/jb843051627/beeyard/internal/store"
)

// Service 业务容器，聚合各子服务。
type Service struct {
	Hives       *HiveService
	Readings    *ReadingService
	Alerts      *AlertService
	Inspections *InspectionService
	Harvests    *HarvestService
	Transfers   *TransferService
	Maintenance *MaintenanceService
	Stats       *StatsService
}

// NewService 构造业务容器。
func NewService(st *store.Store, rc *cache.ReadingCache) *Service {
	return &Service{
		Hives:       NewHiveService(st.Hives),
		Readings:    NewReadingService(st.Readings, rc),
		Alerts:      NewAlertService(st.Alerts),
		Inspections: NewInspectionService(st.Inspections),
		Harvests:    NewHarvestService(st.Harvests),
		Transfers:   NewTransferService(st.Transfers, st.Queens),
		Maintenance: NewMaintenanceService(st.Maintenance, st.DB()),
		Stats:       NewStatsService(st, rc),
	}
}
