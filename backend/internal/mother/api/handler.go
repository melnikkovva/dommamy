package api

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"backend/internal/mother"
)

type MotherHandler struct {
	s *mother.Service
}

func NewMotherHandler(s *mother.Service) *MotherHandler{
	return &MotherHandler{s: s}
}

var coursesPathRe = regexp.MustCompile(`^/api/mother/([^/]+)/courses$`)

func (h *MotherHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/mother", h.create)
	mux.HandleFunc("GET /api/mother", h.search)
	mux.HandleFunc("GET /api/mother/", h.coursesRouter)
}

type motherRequest struct {
	FirstName string `json:"first_name"`
	LastName string `json:"last_name"`
	Phone string `json:"phone"`
	Email string `json:"email"`
}

type courseStatsResponse struct {
	TotalCourses int `json:"total_courses"`
	CompletedCourses int `json:"completed_courses"`
	InProgressCourses int `json:"in_progress_courses"`
	AverageProgress float64 `json:"average_progress"`
	TotalTimeSpent int `json:"total_time_spent"`
}

type motherResponse struct {
	Id int `json:"id"`
	FirstName string `json:"first_name"`
	LastName string `json:"last_name"`
	Phone string `json:"phone"`
	Email string `json:"email"`
	LearnUserID string `json:"learn_user_id"`
	CreatedAt time.Time `json:"created_at"`
	CourseStats *courseStatsResponse `json:"course_stats,omitempty"`
}

func handleToMotherResponse(m *mother.Mother) motherResponse {
	resp := motherResponse{
		Id: m.Id, 
		FirstName: m.FirstName, 
		LastName: m.LastName,
		Phone: m.Phone, 
		Email: m.Email,
		LearnUserID: m.LearnUserID,
		CreatedAt: m.CreatedAt,
	}
	if m.CourseStats != nil {
		resp.CourseStats = &courseStatsResponse{
			TotalCourses: m.CourseStats.TotalCourses,
			CompletedCourses: m.CourseStats.CompletedCourses,
			InProgressCourses: m.CourseStats.InProgressCourses,
			AverageProgress: m.CourseStats.AverageProgress,
			TotalTimeSpent: m.CourseStats.TotalTimeSpent,
		}
	}
	return resp
}

func (h *MotherHandler) create(w http.ResponseWriter, r *http.Request) {
	var req motherRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendJSON(w, http.StatusBadRequest, map[string]string{"error": "incorrect JSON"})
		return
	}

	m, err := h.s.Create(req.FirstName, req.LastName, req.Phone, req.Email)
	if err != nil {
		sendJSON(w, http.StatusUnprocessableEntity, map[string]string{"error": err.Error()})
		return
	}
	sendJSON(w, http.StatusCreated, handleToMotherResponse(m))
}

func (h *MotherHandler) search(w http.ResponseWriter, r *http.Request) {
	items, err := h.s.Search(r.URL.Query().Get("q"))
	if err != nil {
		sendJSON(w, http.StatusInternalServerError, map[string]string{"error": "server error"})
		return
	}
	result := []motherResponse{}
	for _, m := range items {
		result = append(result, handleToMotherResponse(m))
	}
	sendJSON(w, http.StatusOK, result)
}

type courseResultResponse struct {
	CourseId         string    `json:"course_id"`
	NameCourse       string    `json:"name_course"`
	AccesDate        time.Time `json:"access_date"`
	DueDate          *time.Time `json:"due_date,omitempty"`
	CompletionDate   *time.Time `json:"completion_date,omitempty"`
	TimeSpent        int       `json:"time_spent"`
	Duration         int       `json:"duration"`
	CompletionStatus string    `json:"completion_status"`
	Progress         float64   `json:"progress"`
	AwardedScore     *float64  `json:"awarded_score,omitempty"`
	IsOverdue        bool      `json:"is_overdue"`
	ViewsCount       int       `json:"views_count"`
}

func (h *MotherHandler) coursesRouter(w http.ResponseWriter, r *http.Request) {
	matches := coursesPathRe.FindStringSubmatch(r.URL.Path)
	if matches == nil {
		http.NotFound(w, nil)
		return
	}

	motherId, err := strconv.Atoi(matches[1])
	if err != nil {
		sendJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
		return
	}

	courses, err := h.s.GetCourses(motherId)
	if err != nil {
		sendJSON(w, http.StatusInternalServerError, map[string]string{"error": "server error"})
		return
	}

	result := make([]courseResultResponse, len(courses))
	for i, c := range courses {
		result[i] = courseResultResponse{
			CourseId:         c.CourseId,
			NameCourse:       c.NameCourse,
			AccesDate:        c.AccesDate,
			DueDate:          c.DueDate,
			CompletionDate:   c.CompletionDate,
			TimeSpent:        c.TimeSpent,
			Duration:         c.Duration,
			CompletionStatus: c.CompletionStatus,
			Progress:         c.Progress,
			AwardedScore:     c.AwardedScore,
			IsOverdue:        c.IsOverdue,
			ViewsCount:       c.ViewsCount,
		}
	}
	sendJSON(w, http.StatusOK, result)
}

func sendJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

