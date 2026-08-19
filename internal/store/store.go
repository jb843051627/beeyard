package store

import (
	"database/sql"
	"fmt"
	"strings"

	_ "modernc.org/sqlite"
)

// Store 持久化容器，持有各实体 store 与共享读数缓存。
type Store struct {
	db        *sql.DB
	Hives     *HiveStore
	Queens    *QueenStore
	Readings  *ReadingStore
	Alerts    *AlertStore
	Inspections *InspectionStore
	Harvests  *HarvestStore
	Transfers *TransferStore
	Maintenance *MaintenanceStore
}

// NewStore 打开 SQLite 文件库并执行建表迁移。
func NewStore(dsn string) (*Store, error) {
	if !strings.Contains(dsn, "_pragma") {
		if strings.Contains(dsn, "?") {
			dsn += "&_pragma=foreign_keys(1)"
		} else {
			dsn += "?_pragma=foreign_keys(1)"
		}
	}
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("open sqlite: %w", err)
	}
	if _, err := db.Exec("PRAGMA journal_mode=WAL;"); err != nil {
		db.Close()
		return nil, fmt.Errorf("set journal_mode: %w", err)
	}
	if _, err := db.Exec("PRAGMA foreign_keys=ON;"); err != nil {
		db.Close()
		return nil, fmt.Errorf("set foreign_keys: %w", err)
	}
	s := &Store{db: db}
	if err := s.migrate(); err != nil {
		db.Close()
		return nil, err
	}
	s.Hives = NewHiveStore(db)
	s.Queens = NewQueenStore(db)
	s.Readings = NewReadingStore(db)
	s.Alerts = NewAlertStore(db)
	s.Inspections = NewInspectionStore(db)
	s.Harvests = NewHarvestStore(db)
	s.Transfers = NewTransferStore(db)
	s.Maintenance = NewMaintenanceStore(db)
	return s, nil
}

// Close 关闭底层连接。
func (s *Store) Close() error { return s.db.Close() }

// DB 暴露底层连接，供需要直接事务的操作使用。
func (s *Store) DB() *sql.DB { return s.db }

func (s *Store) migrate() error {
	stmts := []string{
		`CREATE TABLE IF NOT EXISTS apiaries (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			location TEXT NOT NULL DEFAULT '',
			established_at DATETIME NOT NULL,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE IF NOT EXISTS hives (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			apiary_id INTEGER NOT NULL,
			code TEXT NOT NULL UNIQUE,
			status TEXT NOT NULL DEFAULT 'active',
			installed_at DATETIME NOT NULL,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY(apiary_id) REFERENCES apiaries(id)
		);`,
		`CREATE TABLE IF NOT EXISTS queens (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			hive_id INTEGER NOT NULL,
			breed TEXT NOT NULL,
			status TEXT NOT NULL DEFAULT 'active',
			marked_at DATETIME NOT NULL,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY(hive_id) REFERENCES hives(id)
		);`,
		`CREATE TABLE IF NOT EXISTS readings (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			hive_id INTEGER NOT NULL,
			sensor_type TEXT NOT NULL,
			value REAL NOT NULL,
			recorded_at DATETIME NOT NULL,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY(hive_id) REFERENCES hives(id)
		);`,
		`CREATE INDEX IF NOT EXISTS idx_readings_hive ON readings(hive_id, recorded_at DESC);`,
		`CREATE TABLE IF NOT EXISTS alerts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			apiary_id INTEGER NOT NULL,
			hive_id INTEGER NOT NULL,
			level TEXT NOT NULL,
			message TEXT NOT NULL,
			status TEXT NOT NULL DEFAULT 'active',
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE IF NOT EXISTS inspections (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			hive_id INTEGER NOT NULL,
			scheduled_at DATETIME NOT NULL,
			status TEXT NOT NULL DEFAULT 'pending',
			completed_at DATETIME,
			notes TEXT NOT NULL DEFAULT '',
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE IF NOT EXISTS harvests (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			apiary_id INTEGER NOT NULL,
			hive_id INTEGER NOT NULL,
			amount_kg REAL NOT NULL,
			harvested_at DATETIME NOT NULL,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE IF NOT EXISTS transfers (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			hive_id INTEGER NOT NULL,
			from_apiary INTEGER NOT NULL,
			to_apiary INTEGER NOT NULL,
			status TEXT NOT NULL DEFAULT 'pending',
			requested_at DATETIME NOT NULL,
			completed_at DATETIME,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE IF NOT EXISTS maintenance_tasks (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			hive_id INTEGER NOT NULL CHECK(hive_id > 0),
			description TEXT NOT NULL,
			status TEXT NOT NULL DEFAULT 'pending',
			scheduled_for DATETIME NOT NULL,
			completed_at DATETIME,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE IF NOT EXISTS apiaries_seed (n INTEGER);`,
	}
	for _, stmt := range stmts {
		if _, err := s.db.Exec(stmt); err != nil {
			return fmt.Errorf("migrate: %w (stmt=%s)", err, stmt)
		}
	}
	return nil
}

// SeedApiary 创建一个蜂场并返回 id，用于测试与冒烟启动。
func (s *Store) SeedApiary(name, location string) (int64, error) {
	res, err := s.db.Exec(
		`INSERT INTO apiaries(name, location, established_at) VALUES(?,?,?)`,
		name, location, "2024-03-01T00:00:00Z",
	)
	if err != nil {
		return 0, fmt.Errorf("seed apiary: %w", err)
	}
	return res.LastInsertId()
}
