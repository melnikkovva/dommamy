package mysql

import (
	"database/sql"
	"time"

	"backend/internal/lms"
)

type Repo struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) *Repo {
	return &Repo{db: db}
}

func (r *Repo) SaveLearners(learners []lms.Learner) error {
	if len(learners) == 0 {
		return nil
	}

	query := `
		INSERT INTO lms_learners (user_id, department_id, first_name, last_name, email, login, status, added_date, last_login_date)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE 
			department_id = VALUES(department_id),
			first_name = VALUES(first_name),
			last_name = VALUES(last_name),
			email = VALUES(email),
			login = VALUES(login),
			status = VALUES(status),
			added_date = VALUES(added_date),
			last_login_date = VALUES(last_login_date)
	`

	for _, l := range learners {
		_, err := r.db.Exec(query, 
			l.UserID, l.DepartmentID, l.FirstName, l.LastName, l.Email, l.Login, 
			l.Status, l.AddedDate, l.LastLoginDate)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Repo) SaveCourseResults(results []lms.CourseResult) error {
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
			res.UserID, res.CourseID, res.CourseTitle, res.EnrollmentID, res.AccessDate,
			dueDate, completionDate, res.TimeSpent, res.Duration, res.CompletionStatus,
			res.Progress, awardedScore, res.IsOverdue, res.ViewsCount)
		if err != nil {
			return err
		}
	}
	return nil
}