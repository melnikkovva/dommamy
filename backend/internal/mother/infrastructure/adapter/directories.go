package adapter

import (
	"context"

	learndomain "backend/internal/learn/domain"
	motherquery "backend/internal/mother/application/query"
)

type MotherDirectory struct {
	reads ReadModel
}

type ReadModel interface {
	Exists(ctx context.Context, id int) (bool, error)
}

func NewMotherDirectory(reads ReadModel) *MotherDirectory {
	return &MotherDirectory{reads: reads}
}

func (d *MotherDirectory) Exists(ctx context.Context, id int) (bool, error) {
	return d.reads.Exists(ctx, id)
}

type LearnStatsAdapter struct {
	reads LearnStatsReads
}

type LearnStatsReads interface {
	GetCourseStats(ctx context.Context, learnUserID string) (*learndomain.CourseStats, error)
	GetCourseResults(ctx context.Context, learnUserID string) ([]learndomain.CourseResult, error)
	GetIPRPlans(ctx context.Context, learnUserID string) ([]learndomain.IPRPlan, error)
}

func NewLearnStatsAdapter(reads LearnStatsReads) *LearnStatsAdapter {
	return &LearnStatsAdapter{reads: reads}
}

func (a *LearnStatsAdapter) GetCourseStats(ctx context.Context, learnUserID string) (motherquery.CourseStatsDTO, error) {
	stats, err := a.reads.GetCourseStats(ctx, learnUserID)
	if err != nil {
		return motherquery.CourseStatsDTO{}, err
	}
	return motherquery.CourseStatsDTO{
		TotalCourses:      stats.TotalCourses,
		CompletedCourses:  stats.CompletedCourses,
		InProgressCourses: stats.InProgressCourses,
		AverageProgress:   stats.AverageProgress,
		TotalTimeSpent:    stats.TotalTimeSpent,
	}, nil
}

func (a *LearnStatsAdapter) GetCourseResults(ctx context.Context, learnUserID string) ([]motherquery.CourseResultDTO, error) {
	results, err := a.reads.GetCourseResults(ctx, learnUserID)
	if err != nil {
		return nil, err
	}
	dtos := make([]motherquery.CourseResultDTO, len(results))
	for i, r := range results {
		dtos[i] = motherquery.CourseResultDTO{
			CourseID:         r.CourseId,
			NameCourse:       r.NameCourse,
			AccesDate:        r.AccesDate.Format("2006-01-02T15:04"),
			TimeSpent:        r.TimeSpent,
			Duration:         r.Duration,
			CompletionStatus: r.CompletionStatus,
			Progress:         r.Progress,
			IsOverdue:        r.IsOverdue,
			ViewsCount:       r.ViewsCount,
		}
		if r.DueDate != nil {
			s := r.DueDate.Format("2006-01-02T15:04")
			dtos[i].DueDate = &s
		}
		if r.CompletionDate != nil {
			s := r.CompletionDate.Format("2006-01-02T15:04")
			dtos[i].CompletionDate = &s
		}
		if r.AwardedScore != nil {
			dtos[i].AwardedScore = r.AwardedScore
		}
	}
	return dtos, nil
}

func (a *LearnStatsAdapter) GetIPRPlans(ctx context.Context, learnUserID string) ([]motherquery.IPRPlanDTO, error) {
	plans, err := a.reads.GetIPRPlans(ctx, learnUserID)
	if err != nil {
		return nil, err
	}
	dtos := make([]motherquery.IPRPlanDTO, len(plans))
	for i, p := range plans {
		dtos[i] = motherquery.IPRPlanDTO{
			PlanID:         p.PlanID,
			Title:          p.Title,
			EmployeeID:     p.EmployeeID,
			MentorID:       p.MentorID,
			MentorName:     p.MentorName,
			StartDate:      p.StartDate.Format("2006-01-02T15:04"),
			Status:         p.Status,
			Progress:       p.Progress,
			CompletedTasks: p.CompletedTasks,
			TotalTasks:     p.TotalTasks,
			Goals:          p.Goals,
		}
		if p.EndDate != nil {
			s := p.EndDate.Format("2006-01-02T15:04")
			dtos[i].EndDate = &s
		}
		stages := make([]motherquery.IPRStageDTO, len(p.Stages))
		for j, s := range p.Stages {
			stages[j] = motherquery.IPRStageDTO{
				StageID:     s.StageID,
				Title:       s.Title,
				Description: s.Description,
			}
			if s.StartDate != nil {
				str := s.StartDate.Format("2006-01-02T15:04")
				stages[j].StartDate = &str
			}
			if s.EndDate != nil {
				str := s.EndDate.Format("2006-01-02T15:04")
				stages[j].EndDate = &str
			}
			tasks := make([]motherquery.IPRTaskDTO, len(s.Tasks))
			for k, t := range s.Tasks {
				tasks[k] = motherquery.IPRTaskDTO{
					TaskID:    t.TaskID,
					TaskType:  t.TaskType,
					Title:     t.Title,
					Description: t.Description,
					Status:    t.Status,
					IsOverdue: t.IsOverdue,
				}
				if t.AbsoluteDueDate != nil {
					str := t.AbsoluteDueDate.Format("2006-01-02T15:04")
					tasks[k].AbsoluteDueDate = &str
				}
				if t.CompletionDate != nil {
					str := t.CompletionDate.Format("2006-01-02T15:04")
					tasks[k].CompletionDate = &str
				}
			}
			stages[j].Tasks = tasks
		}
		dtos[i].Stages = stages
	}
	return dtos, nil
}
