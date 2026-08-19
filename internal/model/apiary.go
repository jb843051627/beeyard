package model

import "time"

// Apiary 蜂场。
type Apiary struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	Location     string    `json:"location"`
	EstablishedAt time.Time `json:"established_at"`
	CreatedAt    time.Time `json:"created_at"`
}

// HiveStatus 蜂箱状态机。
const (
	HiveStatusActive      = "active"
	HiveStatusMaintenance = "maintenance"
	HiveStatusInactive    = "inactive"
)

// Hive 蜂箱，隶属于蜂场。
type Hive struct {
	ID          int64     `json:"id"`
	ApiaryID    int64     `json:"apiary_id"`
	Code        string    `json:"code"`
	Status      string    `json:"status"`
	InstalledAt time.Time `json:"installed_at"`
	CreatedAt   time.Time `json:"created_at"`
}

// QueenStatus 蜂王状态机。
const (
	QueenStatusActive    = "active"
	QueenStatusSuperseded = "superseded"
	QueenStatusDead      = "dead"
)

// Queen 蜂王，绑定到蜂箱。
type Queen struct {
	ID        int64     `json:"id"`
	HiveID    int64     `json:"hive_id"`
	Breed     string    `json:"breed"`
	Status    string    `json:"status"`
	MarkedAt  time.Time `json:"marked_at"`
	CreatedAt time.Time `json:"created_at"`
}
