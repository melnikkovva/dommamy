package mysql

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	"backend/internal/learn/domain"
)

type ReadModel struct {
	db *sql.DB
}

func NewReadModel(db *sql.DB) *ReadModel {
	return &ReadModel{db: db}
}

func (r *ReadModel) GetCourseStats(ctx context.Context, learnUserID string) (*domain.CourseStats, error) {
	query := `
		SELECT 
			COUNT(*) as total_courses,
			COALESCE(SUM(CASE WHEN completion_status = 'completed' THEN 1 ELSE 0 END), 0) as completed_courses,
			COALESCE(SUM(CASE WHEN completion_status != 'completed' AND progress > 0 THEN 1 ELSE 0 END), 0) as in_progress_courses,
			COALESCE(AVG(progress), 0) as average_progress,
			COALESCE(SUM(time_spent), 0) as total_time_spent
		FROM lms_course_results
		WHERE user_id = ?
	`

	var stats domain.CourseStats
	err := r.db.QueryRowContext(ctx, query, learnUserID).Scan(
		&stats.TotalCourses,
		&stats.CompletedCourses,
		&stats.InProgressCourses,
		&stats.AverageProgress,
		&stats.TotalTimeSpent,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return &stats, nil
	}
	if err != nil {
		return nil, err
	}
	return &stats, nil
}

func (r *ReadModel) GetAllLearners(ctx context.Context) ([]domain.Learner, error) {
	query := `SELECT user_id, login, email, phone, first_name, last_name FROM lms_learners`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var learners []domain.Learner
	for rows.Next() {
		var l domain.Learner
		var login, email, phone sql.NullString
		err := rows.Scan(&l.UserId, &login, &email, &phone, &l.FirstName, &l.LastName)
		if err != nil {
			return nil, err
		}
		if login.Valid {
			l.Login = login.String
		}
		if email.Valid {
			l.Email = email.String
		}
		if phone.Valid {
			l.Phone = phone.String
		}
		learners = append(learners, l)
	}
	return learners, rows.Err()
}

func (r *ReadModel) GetCourseResults(ctx context.Context, learnUserID string) ([]domain.CourseResult, error) {
	query := `
		SELECT 
			user_id, course_id, course_title, access_date, due_date, 
			completion_date, time_spent, duration, completion_status, 
			progress, awarded_score, is_overdue, views_count
		FROM lms_course_results
		WHERE user_id = ?
		ORDER BY access_date DESC
	`

	rows, err := r.db.QueryContext(ctx, query, learnUserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []domain.CourseResult
	for rows.Next() {
		var cr domain.CourseResult
		var dueDate, completionDate sql.NullTime
		var awardedScore sql.NullFloat64

		err := rows.Scan(
			&cr.UserId, &cr.CourseId, &cr.NameCourse, &cr.AccesDate,
			&dueDate, &completionDate, &cr.TimeSpent, &cr.Duration,
			&cr.CompletionStatus, &cr.Progress, &awardedScore,
			&cr.IsOverdue, &cr.ViewsCount,
		)
		if err != nil {
			return nil, err
		}

		if dueDate.Valid {
			t := dueDate.Time
			cr.DueDate = &t
		}
		if completionDate.Valid {
			t := completionDate.Time
			cr.CompletionDate = &t
		}
		if awardedScore.Valid {
			s := awardedScore.Float64
			cr.AwardedScore = &s
		}

		results = append(results, cr)
	}

	return results, rows.Err()
}

func (r *ReadModel) GetIPRPlans(ctx context.Context, learnUserID string) ([]domain.IPRPlan, error) {
	query := `
		SELECT plan_id, title, employee_id, mentor_id, mentor_name,
			start_date, end_date, status, progress, completed_tasks, total_tasks, goals
		FROM lms_ipr_plans
		WHERE employee_id = ?
		ORDER BY start_date DESC
	`

	rows, err := r.db.QueryContext(ctx, query, learnUserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var plans []domain.IPRPlan
	for rows.Next() {
		var p domain.IPRPlan
		var endDate sql.NullTime
		var goalsJSON sql.NullString

		err := rows.Scan(
			&p.PlanID, &p.Title, &p.EmployeeID, &p.MentorID, &p.MentorName,
			&p.StartDate, &endDate, &p.Status, &p.Progress,
			&p.CompletedTasks, &p.TotalTasks, &goalsJSON,
		)
		if err != nil {
			return nil, err
		}

		if endDate.Valid {
			t := endDate.Time
			p.EndDate = &t
		}
		if goalsJSON.Valid {
			json.Unmarshal([]byte(goalsJSON.String), &p.Goals)
		}

		stages, err := r.getIPRStages(ctx, p.PlanID)
		if err != nil {
			return nil, err
		}
		p.Stages = stages

		plans = append(plans, p)
	}
	return plans, rows.Err()
}

func (r *ReadModel) getIPRStages(ctx context.Context, planID string) ([]domain.IPRStage, error) {
	query := `
		SELECT id, stage_id, title, description, start_date, end_date
		FROM lms_ipr_stages
		WHERE plan_id = ?
		ORDER BY start_date ASC
	`

	rows, err := r.db.QueryContext(ctx, query, planID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stages []domain.IPRStage
	for rows.Next() {
		var s domain.IPRStage
		var stageDBID int64
		var sStart, sEnd sql.NullTime

		err := rows.Scan(&stageDBID, &s.StageID, &s.Title, &s.Description, &sStart, &sEnd)
		if err != nil {
			return nil, err
		}

		if sStart.Valid {
			t := sStart.Time
			s.StartDate = &t
		}
		if sEnd.Valid {
			t := sEnd.Time
			s.EndDate = &t
		}

		tasks, err := r.getIPRTasks(ctx, stageDBID)
		if err != nil {
			return nil, err
		}
		s.Tasks = tasks

		stages = append(stages, s)
	}
	return stages, rows.Err()
}

func (r *ReadModel) getIPRTasks(ctx context.Context, stageID int64) ([]domain.IPRTask, error) {
	query := `
		SELECT task_id, task_type, title, description,
			status, is_overdue, absolute_due_date, completion_date
		FROM lms_ipr_tasks
		WHERE stage_id = ?
		ORDER BY absolute_due_date ASC
	`

	rows, err := r.db.QueryContext(ctx, query, stageID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []domain.IPRTask
	for rows.Next() {
		var t domain.IPRTask
		var dueDate, completionDate sql.NullTime

		err := rows.Scan(&t.TaskID, &t.TaskType, &t.Title, &t.Description,
			&t.Status, &t.IsOverdue, &dueDate, &completionDate)
		if err != nil {
			return nil, err
		}

		if dueDate.Valid {
			d := dueDate.Time
			t.AbsoluteDueDate = &d
		}
		if completionDate.Valid {
			d := completionDate.Time
			t.CompletionDate = &d
		}

		tasks = append(tasks, t)
	}
	return tasks, rows.Err()
}
