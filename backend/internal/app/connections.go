package app

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"

	documentstorage "backend/internal/document/storage"
)

type ConnectionContainer struct {
	db       *sql.DB
	files    *documentstorage.Local
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

	files, err := documentstorage.NewLocal(cfg.FilesDir) 
	if err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("папка для файлов недоступна: %w", err)
	}

	return &ConnectionContainer{db: db, files: files, filesDir: cfg.FilesDir}, nil
}

func (c *ConnectionContainer) DB() *sql.DB                  { return c.db }
func (c *ConnectionContainer) Files() *documentstorage.Local { return c.files }
func (c *ConnectionContainer) FilesDir() string              { return c.filesDir }

func (c *ConnectionContainer) Close() error { return c.db.Close() }