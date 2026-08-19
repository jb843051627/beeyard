package model

import "time"

// AlertLevel 告警级别。
const (
	AlertLevelInfo     = "info"
	AlertLevelWarning  = "warning"
	AlertLevelCritical = "critical"
)

// AlertStatus 告警状态机。
const (
	AlertStatusActive        = "active"
	AlertStatusAcknowledged = "acknowledged"
)

// Alert 异常告警。
type Alert struct {
	ID        int64     `json:"id"`
	ApiaryID  int64     `json:"apiary_id"`
	HiveID    int64     `json:"hive_id"`
	Level     string    `json:"level"`
	Message   string    `json:"message"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}
