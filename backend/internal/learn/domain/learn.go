package domain

import "time"

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
	UserId           string
	CourseId         string
	NameCourse       string
	AccesDate        time.Time
	DueDate          *time.Time
	CompletionDate   *time.Time
	TimeSpent        int
	Duration         int
	CompletionStatus string
	Progress         float64
	AwardedScore     *float64
	IsOverdue        bool
	ViewsCount       int
}

type CourseStats struct {
	TotalCourses      int
	CompletedCourses  int
	InProgressCourses int
	AverageProgress   float64
	TotalTimeSpent    int
}

type Learner struct {
	UserId    string
	Login     string
	Email     string
	Phone     string
	FirstName string
	LastName  string
}

type IPRPlan struct {
	PlanID         string
	Title          string
	EmployeeID     string
	MentorID       string
	MentorName     string
	StartDate      time.Time
	EndDate        *time.Time
	Status         string
	Progress       float64
	CompletedTasks int
	TotalTasks     int
	Goals          []string
	Stages         []IPRStage
}

type IPRStage struct {
	StageID     string
	Title       string
	Description string
	StartDate   *time.Time
	EndDate     *time.Time
	Tasks       []IPRTask
}

type IPRTask struct {
	TaskID          string
	TaskType        string
	Title           string
	Description     string
	Status          string
	IsOverdue       bool
	AbsoluteDueDate *time.Time
	CompletionDate  *time.Time
}
