package mysql

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"backend/internal/learn"
)

type Repo struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) *Repo {
	return &Repo{db: db}
}

func (r *Repo) SaveLearners(learners []learn.MotherLearnAccount) error {
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
		_, err := r.db.Exec(query,
			l.UserId, l.DepartmentId, l.FirstName, l.LastName, l.Email, l.Login, l.Phone,
			l.Status, addedDate, lastLoginDate)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Repo) SaveCourseResults(results []learn.CourseResult) error {
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
		var dueDate, completionDate interface{}
		if res.DueDate != nil {
			dueDate = *res.DueDate
		}
		if res.CompletionDate != nil {
			completionDate = *res.CompletionDate
		}

		var awardedScore interface{}
		if res.AwardedScore != nil {
			awardedScore = *res.AwardedScore
		}

		_, err := r.db.Exec(query,
			res.UserId, res.CourseId, res.NameCourse, "", res.AccesDate,
			dueDate, completionDate, res.TimeSpent, res.Duration, res.CompletionStatus,
			res.Progress, awardedScore, res.IsOverdue, res.ViewsCount)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Repo) GetCourseStats(learnUserID string) (*learn.CourseStats, error) {
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

	var stats learn.CourseStats
	err := r.db.QueryRow(query, learnUserID).Scan(
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

func (r *Repo) GetAllLearners() ([]learn.Learner, error) {
	query := `SELECT user_id, login, email, phone, first_name, last_name FROM lms_learners`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var learners []learn.Learner
	for rows.Next() {
		var l learn.Learner
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

func (r *Repo) GetCourseResults(learnUserID string) ([]learn.CourseResult, error) {
	query := `
		SELECT 
			user_id, course_id, course_title, access_date, due_date, 
			completion_date, time_spent, duration, completion_status, 
			progress, awarded_score, is_overdue, views_count
		FROM lms_course_results
		WHERE user_id = ?
		ORDER BY access_date DESC
	`

	rows, err := r.db.Query(query, learnUserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []learn.CourseResult
	for rows.Next() {
		var cr learn.CourseResult
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

func (r *Repo) SaveIPRPlans(plans []learn.IPRPlan) error {
	if len(plans) == 0 {
		return nil
	}

	r.db.Exec("DELETE FROM lms_ipr_tasks")
	r.db.Exec("DELETE FROM lms_ipr_stages")
	r.db.Exec("DELETE FROM lms_ipr_plans")

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
		var endDate interface{}
		if p.EndDate != nil {
			endDate = *p.EndDate
		}

		goalsJSON, _ := json.Marshal(p.Goals)

		_, err := r.db.Exec(planQuery,
			p.PlanID, p.Title, p.EmployeeID, p.MentorID, p.MentorName,
			p.StartDate, endDate, p.Status, p.Progress,
			p.CompletedTasks, p.TotalTasks, string(goalsJSON),
		)
		if err != nil {
			return fmt.Errorf("сохранение ИПР %s: %w", p.PlanID, err)
		}

		for _, s := range p.Stages {
			var sStart, sEnd interface{}
			if s.StartDate != nil {
				sStart = *s.StartDate
			}
			if s.EndDate != nil {
				sEnd = *s.EndDate
			}

			stageResult, err := r.db.Exec(stageQuery,
				p.PlanID, s.StageID, s.Title, s.Description, sStart, sEnd,
			)
			if err != nil {
				return fmt.Errorf("сохранение этапа %s: %w", s.StageID, err)
			}
			stageID, _ := stageResult.LastInsertId()

			for _, t := range s.Tasks {
				var dueDate, completionDate interface{}
				if t.AbsoluteDueDate != nil {
					dueDate = *t.AbsoluteDueDate
				}
				if t.CompletionDate != nil {
					completionDate = *t.CompletionDate
				}

				_, err := r.db.Exec(taskQuery,
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

func (r *Repo) GetIPRPlans(learnUserID string) ([]learn.IPRPlan, error) {
	query := `
		SELECT plan_id, title, employee_id, mentor_id, mentor_name,
			start_date, end_date, status, progress, completed_tasks, total_tasks, goals
		FROM lms_ipr_plans
		WHERE employee_id = ?
		ORDER BY start_date DESC
	`

	rows, err := r.db.Query(query, learnUserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var plans []learn.IPRPlan
	for rows.Next() {
		var p learn.IPRPlan
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

		stages, err := r.getIPRStages(p.PlanID)
		if err != nil {
			return nil, err
		}
		p.Stages = stages

		plans = append(plans, p)
	}
	return plans, rows.Err()
}

func (r *Repo) getIPRStages(planID string) ([]learn.IPRStage, error) {
	query := `
		SELECT id, stage_id, title, description, start_date, end_date
		FROM lms_ipr_stages
		WHERE plan_id = ?
		ORDER BY start_date ASC
	`

	rows, err := r.db.Query(query, planID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stages []learn.IPRStage
	for rows.Next() {
		var s learn.IPRStage
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

		tasks, err := r.getIPRTasks(stageDBID)
		if err != nil {
			return nil, err
		}
		s.Tasks = tasks

		stages = append(stages, s)
	}
	return stages, rows.Err()
}

func (r *Repo) getIPRTasks(stageID int64) ([]learn.IPRTask, error) {
	query := `
		SELECT task_id, task_type, title, description,
			status, is_overdue, absolute_due_date, completion_date
		FROM lms_ipr_tasks
		WHERE stage_id = ?
		ORDER BY absolute_due_date ASC
	`

	rows, err := r.db.Query(query, stageID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []learn.IPRTask
	for rows.Next() {
		var t learn.IPRTask
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