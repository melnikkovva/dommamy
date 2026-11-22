package mysql

import (
	"context"
	"database/sql"

	"backend/internal/lifecycle/application/query"
	"backend/internal/lifecycle/domain"
)

type ReadModel struct {
	db *sql.DB
}

func NewReadModel(db *sql.DB) *ReadModel {
	return &ReadModel{db: db}
}

func (r *ReadModel) GetStatus(ctx context.Context, motherID int) (query.StatusView, error) {
	view := query.StatusView{
		MotherID: motherID,
		Status:   domain.StatusReview,
		History:  []query.HistoryItemView{},
	}

	rows, err := r.db.QueryContext(
		ctx,
		"SELECT status, reason, created_at FROM mother_status_history WHERE mother_id = ? ORDER BY id DESC",
		motherID,
	)
	if err != nil {
		return view, err
	}
	defer rows.Close()

	for rows.Next() {
		var item query.HistoryItemView
		if err := rows.Scan(&item.Status, &item.Reason, &item.CreatedAt); err != nil {
			return view, err
		}
		view.History = append(view.History, item)
	}
	if err := rows.Err(); err != nil {
		return view, err
	}

	if len(view.History) > 0 {
		view.Status = view.History[0].Status
	}

	return view, nil
}

func (r *ReadModel) Currents(ctx context.Context) ([]query.CurrentItem, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT h.mother_id, h.status
		 FROM mother_status_history h
		 JOIN (SELECT MAX(id) AS last_id FROM mother_status_history GROUP BY mother_id) m
		   ON h.id = m.last_id`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []query.CurrentItem
	for rows.Next() {
		var item query.CurrentItem
		if err := rows.Scan(&item.MotherID, &item.Status); err != nil {
			return nil, err
		}
		result = append(result, item)
	}
	return result, rows.Err()
}
