package mysql

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"backend/internal/learn/domain"
)

type WriteRepo struct {
	db *sql.DB
}

func NewWriteRepo(db *sql.DB) *WriteRepo {
	return &WriteRepo{db: db}
}

func (r *WriteRepo) SaveLearners(ctx context.Context, learners []domain.MotherLearnAccount) error {
	if len(learners) == 0 {
		return nil
	}

	query := `
		INSERT INTO lms_learners (user_id, department_id, first_name, last_name, email, login, phone, status, added_date, last_login_date)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE 
			department_id = VALUES(department_id),
			first_name = VALUES(first_name),
			last_name = VALUES(last_name),
			email = VALUES(email),
			login = VALUES(login),
			phone = VALUES(phone),
			status = VALUES(status),
			added_date = VALUES(added_date),
			last_login_date = VALUES(last_login_date)
	`

	for _, l := range learners {
		addedDate := l.AddedDate
		if addedDate.IsZero() || addedDate.Year() == 0 {
			addedDate = time.Now()
		}
		lastLoginDate := l.LastLoginDate
		if lastLoginDate.IsZero() || lastLoginDate.Year() == 0 {
			lastLoginDate = time.Now()
		}
		_, err := r.db.ExecContext(ctx, query,
			l.UserId, l.DepartmentId, l.FirstName, l.LastName, l.Email, l.Login, l.Phone,
			l.Status, addedDate, lastLoginDate)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *WriteRepo) SaveCourseResults(ctx context.Context, results []domain.CourseResult) error {
	if len(results) == 0 {
		return nil
	}

	query := `
		INSERT INTO lms_course_results (
			user_id, course_id, course_title, enrollment_id, access_date, 
			due_date, completion_date, time_spent, duration, completion_status, 
			progress, awarded_score, is_overdue, views_count
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE 
			course_title = VALUES(course_title),
			enrollment_id = VALUES(enrollment_id),
			access_date = VALUES(access_date),
			due_date = VALUES(due_date),
			completion_date = VALUES(completion_date),
			time_spent = VALUES(time_spent),
			duration = VALUES(duration),
			completion_status = VALUES(completion_status),
			progress = VALUES(progress),
			awarded_score = VALUES(awarded_score),
			is_overdue = VALUES(is_overdue),
			views_count = VALUES(views_count)
	`

	for _, res := range results {
		var dueDate sql.NullTime
		if res.DueDate != nil {
			dueDate = sql.NullTime{Time: *res.DueDate, Valid: true}
		}
		var completionDate sql.NullTime
		if res.CompletionDate != nil {
			completionDate = sql.NullTime{Time: *res.CompletionDate, Valid: true}
		}

		var awardedScore sql.NullFloat64
		if res.AwardedScore != nil {
			awardedScore = sql.NullFloat64{Float64: *res.AwardedScore, Valid: true}
		}

		_, err := r.db.ExecContext(ctx, query,
			res.UserId, res.CourseId, res.NameCourse, "", res.AccesDate,
			dueDate, completionDate, res.TimeSpent, res.Duration, res.CompletionStatus,
			res.Progress, awardedScore, res.IsOverdue, res.ViewsCount)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *WriteRepo) SaveIPRPlans(ctx context.Context, plans []domain.IPRPlan) error {
	if len(plans) == 0 {
		return nil
	}

	r.db.ExecContext(ctx, "DELETE FROM lms_ipr_tasks")
	r.db.ExecContext(ctx, "DELETE FROM lms_ipr_stages")
	r.db.ExecContext(ctx, "DELETE FROM lms_ipr_plans")

	planQuery := `
		INSERT INTO lms_ipr_plans (
			plan_id, title, employee_id, mentor_id, mentor_name,
			start_date, end_date, status, progress, completed_tasks, total_tasks, goals
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	stageQuery := `
		INSERT INTO lms_ipr_stages (
			plan_id, stage_id, title, description, start_date, end_date
		) VALUES (?, ?, ?, ?, ?, ?)
	`
	taskQuery := `
		INSERT INTO lms_ipr_tasks (
			stage_id, task_id, task_type, title, description,
			status, is_overdue, absolute_due_date, completion_date
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	for _, p := range plans {
		var endDate sql.NullTime
		if p.EndDate != nil {
			endDate = sql.NullTime{Time: *p.EndDate, Valid: true}
		}

		goalsJSON, _ := json.Marshal(p.Goals)

		_, err := r.db.ExecContext(ctx, planQuery,
			p.PlanID, p.Title, p.EmployeeID, p.MentorID, p.MentorName,
			p.StartDate, endDate, p.Status, p.Progress,
			p.CompletedTasks, p.TotalTasks, string(goalsJSON),
		)
		if err != nil {
			return fmt.Errorf("сохранение ИПР %s: %w", p.PlanID, err)
		}

		for _, s := range p.Stages {
			var sStart sql.NullTime
			if s.StartDate != nil {
				sStart = sql.NullTime{Time: *s.StartDate, Valid: true}
			}
			var sEnd sql.NullTime
			if s.EndDate != nil {
				sEnd = sql.NullTime{Time: *s.EndDate, Valid: true}
			}

			stageResult, err := r.db.ExecContext(ctx, stageQuery,
				p.PlanID, s.StageID, s.Title, s.Description, sStart, sEnd,
			)
			if err != nil {
				return fmt.Errorf("сохранение этапа %s: %w", s.StageID, err)
			}
			stageID, _ := stageResult.LastInsertId()

			for _, t := range s.Tasks {
				var dueDate sql.NullTime
				if t.AbsoluteDueDate != nil {
					dueDate = sql.NullTime{Time: *t.AbsoluteDueDate, Valid: true}
				}
				var completionDate sql.NullTime
				if t.CompletionDate != nil {
					completionDate = sql.NullTime{Time: *t.CompletionDate, Valid: true}
				}

				_, err := r.db.ExecContext(ctx, taskQuery,
					stageID, t.TaskID, t.TaskType, t.Title, t.Description,
					t.Status, t.IsOverdue, dueDate, completionDate,
				)
				if err != nil {
					return fmt.Errorf("сохранение задачи %s: %w", t.TaskID, err)
				}
			}
		}
	}
	return nil
}
