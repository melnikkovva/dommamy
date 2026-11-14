package learn

import (
	"errors"
	"time"
)

var ErrInvalidData = errors.New("invalid learn data")

type MotherLearnAccount struct {
	UserId        string 
	DepartmentId  string
	FirstName     string
	LastName      string
	Email         string
	Login         string
	Phone         string
	Status        string
	LastLoginDate time.Time
	AddedDate     time.Time
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

type CourseStats struct {
	TotalCourses int
	CompletedCourses int
	InProgressCourses int
	AverageProgress float64
	TotalTimeSpent int
}

type Learner struct {
	UserId    string
	Login     string
	Email     string
	Phone     string
	FirstName string
	LastName  string
}

type LearnAccountRepository interface  {
	SaveLearnAccounts(learners []MotherLearnAccount) error
	SaveCourseResult(courses []CourseResult) error
}

type LearnApiClient interface {
	FetchAllLearnAccounts() ([]MotherLearnAccount, error)
	FetchAllCourseResults() ([]CourseResult, error)
}

type StatsRepository interface {
	GetCourseStats(learnUserID string) (*CourseStats, error)
	GetCourseResults(learnUserID string) ([]CourseResult, error)
}

