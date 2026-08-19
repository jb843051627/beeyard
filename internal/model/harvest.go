package model

import "time"

// Harvest 采蜜记录。
type Harvest struct {
	ID          int64     `json:"id"`
	ApiaryID    int64     `json:"apiary_id"`
	HiveID      int64     `json:"hive_id"`
	AmountKg    float64   `json:"amount_kg"`
	HarvestedAt time.Time `json:"harvested_at"`
	CreatedAt   time.Time `json:"created_at"`
}

// TransferStatus 转场状态机。
const (
	TransferStatusPending    = "pending"
	TransferStatusInTransit  = "in_transit"
	TransferStatusCompleted  = "completed"
)

// Transfer 转场记录。
type Transfer struct {
	ID          int64     `json:"id"`
	HiveID      int64     `json:"hive_id"`
	FromApiary  int64     `json:"from_apiary"`
	ToApiary    int64     `json:"to_apiary"`
	Status      string    `json:"status"`
	RequestedAt time.Time `json:"requested_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}
