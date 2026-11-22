package query

import (
	"context"

	"backend/internal/learn/domain"
)

type ReadModel interface {
	GetCourseStats(ctx context.Context, learnUserID string) (*domain.CourseStats, error)
	GetCourseResults(ctx context.Context, learnUserID string) ([]domain.CourseResult, error)
	GetIPRPlans(ctx context.Context, learnUserID string) ([]domain.IPRPlan, error)
	GetAllLearners(ctx context.Context) ([]domain.Learner, error)
}
