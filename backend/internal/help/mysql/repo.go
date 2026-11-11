package mysql

import (
	"context"
	"database/sql"

	"backend/internal/help"
)

type Repo struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) *Repo {
	return &Repo{db: db}
}

func (r *Repo) Save(ctx context.Context, item *help.Item) error {
	res, err := r.db.ExecContext(ctx,
		"INSERT INTO help_items (mother_id, kind, `sum`, check_url, created_at) VALUES (?, ?, ?, ?, ?)",
		item.MotherId, string(item.Kind), item.Sum, item.CheckUrl, item.CreatedAt,
	)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	item.Id = int(id)
	return nil
}

func (r *Repo) List(ctx context.Context) ([]*help.Item, error) {
	rows, err := r.db.QueryContext(ctx,
		"SELECT id, mother_id, kind, `sum`, check_url, created_at FROM help_items ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := []*help.Item{}
	for rows.Next() {
		item := &help.Item{}
		var kind string
		err := rows.Scan(&item.Id, &item.MotherId, &kind, &item.Sum, &item.CheckUrl, &item.CreatedAt)
		if err != nil {
			return nil, err
		}
		item.Kind = help.Kind(kind)
		result = append(result, item)
	}
	return result, rows.Err()
}

func (r *Repo) ById(ctx context.Context, id int) (*help.Item, error) {
	item := &help.Item{}
	var kind string
	err := r.db.QueryRowContext(ctx,
		"SELECT id, mother_id, kind, `sum`, check_url, created_at FROM help_items WHERE id = ?", id,
	).Scan(&item.Id, &item.MotherId, &kind, &item.Sum, &item.CheckUrl, &item.CreatedAt)
	if err != nil {
		return nil, err
	}
	item.Kind = help.Kind(kind)
	return item, nil
}