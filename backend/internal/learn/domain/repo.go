package domain

import "context"

type WriteRepository interface {
	SaveLearners(ctx context.Context, learners []MotherLearnAccount) error
	SaveCourseResults(ctx context.Context, results []CourseResult) error
	SaveIPRPlans(ctx context.Context, plans []IPRPlan) error
}

type ReadRepository interface {
	GetCourseStats(ctx context.Context, learnUserID string) (*CourseStats, error)
	GetCourseResults(ctx context.Context, learnUserID string) ([]CourseResult, error)
	GetIPRPlans(ctx context.Context, learnUserID string) ([]IPRPlan, error)
	GetAllLearners(ctx context.Context) ([]Learner, error)
}
