package domain

import "context"

type CourseStatsDTO struct {
	TotalCourses     int     `json:"total_courses"`
	CompletedCourses int     `json:"completed_courses"`
	InProgressCourses int    `json:"in_progress_courses"`
	AverageProgress  float64 `json:"average_progress"`
	TotalTimeSpent   int     `json:"total_time_spent"`
}

type CourseResultDTO struct {
	CourseID         string  `json:"course_id"`
	NameCourse       string  `json:"name_course"`
	AccesDate        string  `json:"access_date"`
	DueDate          *string `json:"due_date,omitempty"`
	CompletionDate   *string `json:"completion_date,omitempty"`
	TimeSpent        int     `json:"time_spent"`
	Duration         int     `json:"duration"`
	CompletionStatus string  `json:"completion_status"`
	Progress         float64 `json:"progress"`
	AwardedScore     *float64 `json:"awarded_score,omitempty"`
	IsOverdue        bool    `json:"is_overdue"`
	ViewsCount       int     `json:"views_count"`
}

type IPRPlanDTO struct {
	PlanID         string       `json:"plan_id"`
	Title          string       `json:"title"`
	EmployeeID     string       `json:"employee_id"`
	MentorID       string       `json:"mentor_id"`
	MentorName     string       `json:"mentor_name"`
	StartDate      string       `json:"start_date"`
	EndDate        *string      `json:"end_date,omitempty"`
	Status         string       `json:"status"`
	Progress       float64      `json:"progress"`
	CompletedTasks int          `json:"completed_tasks"`
	TotalTasks     int          `json:"total_tasks"`
	Goals          []string     `json:"goals"`
	Stages         []IPRStageDTO `json:"stages"`
}

type IPRStageDTO struct {
	StageID     string        `json:"stage_id"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	StartDate   *string       `json:"start_date,omitempty"`
	EndDate     *string       `json:"end_date,omitempty"`
	Tasks       []IPRTaskDTO  `json:"tasks"`
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

type LearnStats interface {
	GetCourseStats(ctx context.Context, learnUserID string) (CourseStatsDTO, error)
	GetCourseResults(ctx context.Context, learnUserID string) ([]CourseResultDTO, error)
	GetIPRPlans(ctx context.Context, learnUserID string) ([]IPRPlanDTO, error)
}
