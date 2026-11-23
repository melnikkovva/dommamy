package domain

import "context"

type APIClient interface {
	FetchAllLearnAccounts(ctx context.Context) ([]MotherLearnAccount, error)
	FetchAllCourseResults(ctx context.Context) ([]CourseResult, error)
	FetchAllIPR(ctx context.Context) ([]IPRPlan, error)
}

type MotherDirectory interface {
	GetAllMothers(ctx context.Context) ([]MotherRef, error)
}

type MotherRef struct {
	Id        int
	FirstName string
	LastName  string
	Phone     string
	Email     string
}
