package store

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jb843051627/beeyard/internal/model"
)

// ReadingStore 传感器读数持久化。
type ReadingStore struct {
	db *sql.DB
}

func NewReadingStore(db *sql.DB) *ReadingStore {
	return &ReadingStore{db: db}
}

// Record 单条入库。
func (s *ReadingStore) Record(ctx context.Context, hiveID int64, sensorType string, value float64, recordedAt interface{}) (int64, error) {
	res, err := s.db.ExecContext(ctx,
		`INSERT INTO readings(hive_id, sensor_type, value, recorded_at) VALUES(?,?,?,?)`,
		hiveID, sensorType, value, recordedAt)
	if err != nil {
		return 0, fmt.Errorf("reading record: %w", err)
	}
	return res.LastInsertId()
}

// BatchInsert 批量入库（单事务）。
func (s *ReadingStore) BatchInsert(ctx context.Context, batch []model.ReadingBatch) (int64, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("begin tx: %w", err)
	}
	var total int64
	for _, r := range batch {
		res, err := tx.ExecContext(ctx,
			`INSERT INTO readings(hive_id, sensor_type, value, recorded_at) VALUES(?,?,?,?)`,
			r.HiveID, r.SensorType, r.Value, r.RecordedAt)
		if err != nil {
			tx.Rollback()
			return 0, fmt.Errorf("batch insert row: %w", err)
		}
		n, _ := res.RowsAffected()
		total += n
	}
	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("commit batch: %w", err)
	}
	return total, nil
}

// ListByHive 按蜂箱列读数（倒序）。
func (s *ReadingStore) ListByHive(ctx context.Context, hiveID int64, limit int) ([]*model.SensorReading, error) {
	if limit <= 0 {
		limit = 100
	}
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, hive_id, sensor_type, value, recorded_at, created_at FROM readings WHERE hive_id = ? ORDER BY recorded_at DESC LIMIT ?`, hiveID, limit)
	if err != nil {
		return nil, fmt.Errorf("reading list by hive: %w", err)
	}
	defer rows.Close()
	var list []*model.SensorReading
	for rows.Next() {
		var r model.SensorReading
		if err := rows.Scan(&r.ID, &r.HiveID, &r.SensorType, &r.Value, &r.RecordedAt, &r.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, &r)
	}
	return list, rows.Err()
}

// LatestByHive 取蜂箱最新读数。
func (s *ReadingStore) LatestByHive(ctx context.Context, hiveID int64, sensorType string) (*model.SensorReading, error) {
	row := s.db.QueryRowContext(ctx,
		`SELECT id, hive_id, sensor_type, value, recorded_at, created_at FROM readings WHERE hive_id = ? AND sensor_type = ? ORDER BY recorded_at DESC LIMIT 1`, hiveID, sensorType)
	var r model.SensorReading
	err := row.Scan(&r.ID, &r.HiveID, &r.SensorType, &r.Value, &r.RecordedAt, &r.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, model.ErrReadingNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("reading latest: %w", err)
	}
	return &r, nil
}

// ListByApiaryAndTimeRange 按蜂场与时间范围列读数（关联蜂箱表）。
func (s *ReadingStore) ListByApiaryAndTimeRange(ctx context.Context, apiaryID int64, from, to interface{}) ([]*model.SensorReading, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT r.id, r.hive_id, r.sensor_type, r.value, r.recorded_at, r.created_at
		 FROM readings r JOIN hives h ON r.hive_id = h.id
		 WHERE h.apiary_id = ? AND r.recorded_at >= ? AND r.recorded_at <= ?
		 ORDER BY r.recorded_at ASC`, apiaryID, from, to)
	if err != nil {
		return nil, fmt.Errorf("reading list by apiary range: %w", err)
	}
	defer rows.Close()
	var list []*model.SensorReading
	for rows.Next() {
		var r model.SensorReading
		if err := rows.Scan(&r.ID, &r.HiveID, &r.SensorType, &r.Value, &r.RecordedAt, &r.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, &r)
	}
	return list, rows.Err()
}
