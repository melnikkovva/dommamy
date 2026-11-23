package query

import (
	"context"
	"time"
)

type CourseStatsDTO struct {
	TotalCourses      int     `json:"total_courses"`
	CompletedCourses  int     `json:"completed_courses"`
	InProgressCourses int     `json:"in_progress_courses"`
	AverageProgress   float64 `json:"average_progress"`
	TotalTimeSpent    int     `json:"total_time_spent"`
}

type CourseResultDTO struct {
	CourseID         string   `json:"course_id"`
	NameCourse       string   `json:"name_course"`
	AccesDate        string   `json:"access_date"`
	DueDate          *string  `json:"due_date,omitempty"`
	CompletionDate   *string  `json:"completion_date,omitempty"`
	TimeSpent        int      `json:"time_spent"`
	Duration         int      `json:"duration"`
	CompletionStatus string   `json:"completion_status"`
	Progress         float64  `json:"progress"`
	AwardedScore     *float64 `json:"awarded_score,omitempty"`
	IsOverdue        bool     `json:"is_overdue"`
	ViewsCount       int      `json:"views_count"`
}

type IPRPlanDTO struct {
	PlanID         string        `json:"plan_id"`
	Title          string        `json:"title"`
	EmployeeID     string        `json:"employee_id"`
	MentorID       string        `json:"mentor_id"`
	MentorName     string        `json:"mentor_name"`
	StartDate      string        `json:"start_date"`
	EndDate        *string       `json:"end_date,omitempty"`
	Status         string        `json:"status"`
	Progress       float64       `json:"progress"`
	CompletedTasks int           `json:"completed_tasks"`
	TotalTasks     int           `json:"total_tasks"`
	Goals          []string      `json:"goals"`
	Stages         []IPRStageDTO `json:"stages"`
}

type IPRStageDTO struct {
	StageID     string       `json:"stage_id"`
	Title       string       `json:"title"`
	Description string       `json:"description"`
	StartDate   *string      `json:"start_date,omitempty"`
	EndDate     *string      `json:"end_date,omitempty"`
	Tasks       []IPRTaskDTO `json:"tasks"`
}

type IPRTaskDTO struct {
	TaskID          string  `json:"task_id"`
	TaskType        string  `json:"task_type"`
	Title           string  `json:"title"`
	Description     string  `json:"description"`
	Status          string  `json:"status"`
	IsOverdue       bool    `json:"is_overdue"`
	AbsoluteDueDate *string `json:"absolute_due_date,omitempty"`
	CompletionDate  *string `json:"completion_date,omitempty"`
}

type MotherView struct {
	ID          int             `json:"id"`
	FirstName   string          `json:"first_name"`
	LastName    string          `json:"last_name"`
	Phone       string          `json:"phone"`
	Email       string          `json:"email"`
	LearnUserID string          `json:"learn_user_id"`
	PhotoURL    string          `json:"photo_url,omitempty"`
	CreatedAt   time.Time       `json:"created_at"`
	CourseStats *CourseStatsDTO `json:"course_stats,omitempty"`
}

type LearnStats interface {
	GetCourseStats(ctx context.Context, learnUserID string) (CourseStatsDTO, error)
	GetCourseResults(ctx context.Context, learnUserID string) ([]CourseResultDTO, error)
	GetIPRPlans(ctx context.Context, learnUserID string) ([]IPRPlanDTO, error)
}

type ReadModel interface {
	GetByID(ctx context.Context, id int) (MotherView, error)
	Search(ctx context.Context, q string) ([]MotherView, error)
	Exists(ctx context.Context, id int) (bool, error)
}

type MotherRef struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
}

type MotherListReader interface {
	GetAllMothers(ctx context.Context) ([]MotherRef, error)
}
