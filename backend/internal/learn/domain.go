package learn

import (
	"errors"
	"time"
)

var ErrInvalidData = errors.New("invalid learn data")

type MotherLearnAccount struct {
	UserId string 
	FirstName string
	LastName string
	Email string
	Login string
	Status string
	LastLoginDate time.Time
	AddedDate time.Time
}

type CourseResult struct {
	UserId string
	CourseId string
	NameCourse string
	AccesDate time.Time
	DueDate *time.Time
	CompletionDate *time.Time 
	TimeSpent int
	Duration int
	CompletionStatus string
	Progress float64
	AwardedScore *float64
	IsOverdue bool
	ViewsCount int
}

func LearnAccountRepository interface  {
	SaveLearnAccounts(learners []MotherLearnAccount)
	SaveCourseResult(courses []CourseResult)
}

func LearnApiClient interface {
	FetchAllLearnAccounts() ([]MotherLearnAccount, error)
	FetchAllCourseResult() ([]CourseResult, error)
}

