package model

import "time"

// InspectionStatus 巡检状态机。
const (
	InspectionStatusPending   = "pending"
	InspectionStatusCompleted = "completed"
	InspectionStatusOverdue   = "overdue"
)

// Inspection 巡检记录。
type Inspection struct {
	ID           int64      `json:"id"`
	HiveID       int64      `json:"hive_id"`
	ScheduledAt  time.Time  `json:"scheduled_at"`
	Status       string     `json:"status"`
	CompletedAt  *time.Time `json:"completed_at,omitempty"`
	Notes        string     `json:"notes"`
	CreatedAt    time.Time  `json:"created_at"`
}

// TaskStatus 维护任务状态机。
const (
	TaskStatusPending   = "pending"
	TaskStatusCompleted = "completed"
)

// MaintenanceTask 维护任务。
type MaintenanceTask struct {
	ID          int64      `json:"id"`
	HiveID      int64      `json:"hive_id"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	ScheduledFor time.Time  `json:"scheduled_for"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
}
