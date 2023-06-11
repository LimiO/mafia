package db

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type Config struct {
	DBName string
}

type Manager struct {
	DB *sql.DB
}

func NewManager(cfg *Config) (*Manager, error) {
	db, err := sql.Open("sqlite3", cfg.DBName)
	if err != nil {
		return nil, fmt.Errorf("failed to open sqlite3 db %q: %v", cfg.DBName, err)
	}
	mgr := &Manager{DB: db}
	if err = mgr.CreateUserTable(); err != nil {
		return nil, fmt.Errorf("failed to create table user: %v", err)
	}

	if err = mgr.CreateStatsTable(); err != nil {
		return nil, fmt.Errorf("failed to create table stats: %v", err)
	}
	return &Manager{DB: db}, nil
}

func (m *Manager) Close() error {
	return m.DB.Close()
}
