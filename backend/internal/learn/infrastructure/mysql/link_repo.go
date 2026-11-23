package mysql

import (
	"context"
	"database/sql"
	"errors"
)

type MotherLearnLinkRepo struct {
	db *sql.DB
}

func NewMotherLearnLinkRepo(db *sql.DB) *MotherLearnLinkRepo {
	return &MotherLearnLinkRepo{db: db}
}

func (r *MotherLearnLinkRepo) Link(ctx context.Context, motherID int, learnUserID string) error {
	_, err := r.db.ExecContext(
		ctx,
		`INSERT INTO mother_learn_link (mother_id, learn_user_id) VALUES (?, ?)
		 ON DUPLICATE KEY UPDATE learn_user_id = VALUES(learn_user_id)`,
		motherID, learnUserID,
	)
	return err
}

func (r *MotherLearnLinkRepo) GetByMotherID(ctx context.Context, motherID int) (string, error) {
	var learnUserID string
	err := r.db.QueryRowContext(ctx,
		`SELECT learn_user_id FROM mother_learn_link WHERE mother_id = ?`, motherID,
	).Scan(&learnUserID)
	if errors.Is(err, sql.ErrNoRows) {
		return "", nil
	}
	return learnUserID, err
}

func (r *MotherLearnLinkRepo) IsLinked(ctx context.Context, motherID int) (bool, error) {
	var count int
	err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(1) FROM mother_learn_link WHERE mother_id = ?`, motherID,
	).Scan(&count)
	return count > 0, err
}
