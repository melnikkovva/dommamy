package app

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type ConnectionContainer struct {
	db       *sql.DB
	filesDir string
}

func NewConnectionContainer(ctx context.Context, cfg Config) (*ConnectionContainer, error) {
	db, err := sql.Open("mysql", cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}
	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("база недоступна: %w", err)
	}

	return &ConnectionContainer{db: db, filesDir: cfg.FilesDir}, nil
}

func (c *ConnectionContainer) DB() *sql.DB     { return c.db }
func (c *ConnectionContainer) FilesDir() string { return c.filesDir }

func (c *ConnectionContainer) Close() error { return c.db.Close() }