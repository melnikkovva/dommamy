package query

import (
	"context"
)

type GetCourses struct {
	MotherID int
}

type GetCoursesHandler interface {
	Handle(ctx context.Context, q GetCourses) ([]CourseResultDTO, error)
}

func NewGetCoursesHandler(
	reads ReadModel,
	learnReader LearnStats,
) GetCoursesHandler {
	return &getCoursesHandler{
		reads:       reads,
		learnReader: learnReader,
	}
}

type getCoursesHandler struct {
	reads       ReadModel
	learnReader LearnStats
}

func (h *getCoursesHandler) Handle(ctx context.Context, q GetCourses) ([]CourseResultDTO, error) {
	view, err := h.reads.GetByID(ctx, q.MotherID)
	if err != nil {
		return nil, err
	}
	if view.LearnUserID == "" || h.learnReader == nil {
		return nil, nil
	}
	return h.learnReader.GetCourseResults(ctx, view.LearnUserID)
}

type GetIPR struct {
	MotherID int
}

type GetIPRHandler interface {
	Handle(ctx context.Context, q GetIPR) ([]IPRPlanDTO, error)
}

func NewGetIPRHandler(
	reads ReadModel,
	learnReader LearnStats,
) GetIPRHandler {
	return &getIPRHandler{
		reads:       reads,
		learnReader: learnReader,
	}
}

type getIPRHandler struct {
	reads       ReadModel
	learnReader LearnStats
}

func (h *getIPRHandler) Handle(ctx context.Context, q GetIPR) ([]IPRPlanDTO, error) {
	view, err := h.reads.GetByID(ctx, q.MotherID)
	if err != nil {
		return nil, err
	}
	if view.LearnUserID == "" || h.learnReader == nil {
		return nil, nil
	}
	return h.learnReader.GetIPRPlans(ctx, view.LearnUserID)
}
