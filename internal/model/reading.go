package model

import "time"

// SensorType 传感器类型。
const (
	SensorTypeTemp     = "temperature"
	SensorTypeHumidity = "humidity"
	SensorTypeWeight   = "weight"
)

// SensorReading 传感器读数。
type SensorReading struct {
	ID         int64     `json:"id"`
	HiveID     int64     `json:"hive_id"`
	SensorType string    `json:"sensor_type"`
	Value      float64   `json:"value"`
	RecordedAt time.Time `json:"recorded_at"`
	CreatedAt  time.Time `json:"created_at"`
}

// ReadingBatch 批量导入载体。
type ReadingBatch struct {
	HiveID     int64     `json:"hive_id"`
	SensorType string    `json:"sensor_type"`
	Value      float64   `json:"value"`
	RecordedAt time.Time `json:"recorded_at"`
}
