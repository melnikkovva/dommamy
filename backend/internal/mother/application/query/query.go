package query

import (
	"context"
	"time"

	"backend/internal/mother/domain"
)

type MotherView struct {
	ID          int                    `json:"id"`
	FirstName   string                 `json:"first_name"`
	LastName    string                 `json:"last_name"`
	Phone       string                 `json:"phone"`
	Email       string                 `json:"email"`
	LearnUserID string                 `json:"learn_user_id"`
	PhotoURL    string                 `json:"photo_url,omitempty"`
	CreatedAt   time.Time              `json:"created_at"`
	CourseStats *domain.CourseStatsDTO `json:"course_stats,omitempty"`
}

type ReadModel interface {
	GetByID(ctx context.Context, id int) (MotherView, error)
	Search(ctx context.Context, q string) ([]MotherView, error)
	Exists(ctx context.Context, id int) (bool, error)
}
