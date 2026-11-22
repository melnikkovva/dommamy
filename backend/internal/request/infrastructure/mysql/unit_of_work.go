package mysql

import (
	"context"
	"database/sql"

	"backend/internal/request/application/command"
	"backend/internal/request/domain"
)

func NewUnitOfWork(db *sql.DB) command.UnitOfWork { return &unitOfWork{db: db} }

type unitOfWork struct{ db *sql.DB }

func (u *unitOfWork) Do(ctx context.Context, fn func(ctx context.Context, repo domain.Repository) error) error {
	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback() 

	if err := fn(ctx, NewRepository(tx)); err != nil {
		return err
	}
	return tx.Commit()
}