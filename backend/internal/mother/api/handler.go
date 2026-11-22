package api

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"time"

	"backend/internal/mother"
)

type MotherHandler struct {
	s *mother.Service
	photoDir string
}

func NewMotherHandler(s *mother.Service, photoDir string) *MotherHandler{
	return &MotherHandler{s: s, photoDir: photoDir}
}

var coursesPathRe = regexp.MustCompile(`^/api/mother/([^/]+)/courses$`)
var iprPathRe = regexp.MustCompile(`^/api/mother/([^/]+)/ipr$`)

func (h *MotherHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/mother", h.create)
	mux.HandleFunc("GET /api/mother", h.search)
	mux.HandleFunc("POST /api/mother/{id}/photo", h.uploadPhoto)
	mux.HandleFunc("GET /api/mother/", h.learnRouter)
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
	PhotoURL string `json:"photo_url,omitempty"`
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
		PhotoURL: m.PhotoURL,
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

func (h *MotherHandler) learnRouter(w http.ResponseWriter, r *http.Request) {
	if matches := coursesPathRe.FindStringSubmatch(r.URL.Path); matches != nil {
		h.handleCourses(w, r, matches)
		return
	}
	if matches := iprPathRe.FindStringSubmatch(r.URL.Path); matches != nil {
		h.handleIPR(w, r, matches)
		return
	}
	http.NotFound(w, nil)
}

func (h *MotherHandler) handleCourses(w http.ResponseWriter, r *http.Request, matches []string) {
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

type iprTaskResponse struct {
	TaskID          string     `json:"task_id"`
	TaskType        string     `json:"task_type"`
	Title           string     `json:"title"`
	Description     string     `json:"description"`
	Status          string     `json:"status"`
	IsOverdue       bool       `json:"is_overdue"`
	AbsoluteDueDate *time.Time `json:"absolute_due_date,omitempty"`
	CompletionDate  *time.Time `json:"completion_date,omitempty"`
}

type iprStageResponse struct {
	StageID     string          `json:"stage_id"`
	Title       string          `json:"title"`
	Description string          `json:"description"`
	StartDate   *time.Time      `json:"start_date,omitempty"`
	EndDate     *time.Time      `json:"end_date,omitempty"`
	Tasks       []iprTaskResponse `json:"tasks"`
}

type iprPlanResponse struct {
	PlanID         string            `json:"plan_id"`
	Title          string            `json:"title"`
	EmployeeID     string            `json:"employee_id"`
	MentorID       string            `json:"mentor_id"`
	MentorName     string            `json:"mentor_name"`
	StartDate      time.Time         `json:"start_date"`
	EndDate        *time.Time        `json:"end_date,omitempty"`
	Status         string            `json:"status"`
	Progress       float64           `json:"progress"`
	CompletedTasks int               `json:"completed_tasks"`
	TotalTasks     int               `json:"total_tasks"`
	Goals          []string          `json:"goals"`
	Stages         []iprStageResponse `json:"stages"`
}

func (h *MotherHandler) handleIPR(w http.ResponseWriter, r *http.Request, matches []string) {
	motherId, err := strconv.Atoi(matches[1])
	if err != nil {
		sendJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
		return
	}

	plans, err := h.s.GetIPR(motherId)
	if err != nil {
		sendJSON(w, http.StatusInternalServerError, map[string]string{"error": "server error"})
		return
	}

	result := make([]iprPlanResponse, len(plans))
	for i, p := range plans {
		var stages []iprStageResponse
		for _, s := range p.Stages {
			var tasks []iprTaskResponse
			for _, t := range s.Tasks {
				tasks = append(tasks, iprTaskResponse{
					TaskID:          t.TaskID,
					TaskType:        t.TaskType,
					Title:           t.Title,
					Description:     t.Description,
					Status:          t.Status,
					IsOverdue:       t.IsOverdue,
					AbsoluteDueDate: t.AbsoluteDueDate,
					CompletionDate:  t.CompletionDate,
				})
			}
			stages = append(stages, iprStageResponse{
				StageID:     s.StageID,
				Title:       s.Title,
				Description: s.Description,
				StartDate:   s.StartDate,
				EndDate:     s.EndDate,
				Tasks:       tasks,
			})
		}
		result[i] = iprPlanResponse{
			PlanID:         p.PlanID,
			Title:          p.Title,
			EmployeeID:     p.EmployeeID,
			MentorID:       p.MentorID,
			MentorName:     p.MentorName,
			StartDate:      p.StartDate,
			EndDate:        p.EndDate,
			Status:         p.Status,
			Progress:       p.Progress,
			CompletedTasks: p.CompletedTasks,
			TotalTasks:     p.TotalTasks,
			Goals:          p.Goals,
			Stages:         stages,
		}
	}
	sendJSON(w, http.StatusOK, result)
}

func sendJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

const maxPhotoSize = 5 << 20 // 5 MB

var allowedPhotoTypes = map[string]bool{
	"image/jpeg": true,
	"image/png": true,
	"image/webp": true,
}

func (h *MotherHandler) uploadPhoto(w http.ResponseWriter, r *http.Request) {
	motherId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		sendJSON(w, http.StatusBadRequest, map[string]string{"error": "incorrect id"})
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxPhotoSize+(1<<20))
	file, header, err := r.FormFile("photo")
	if err != nil {
		sendJSON(w, http.StatusBadRequest, map[string]string{"error": "send the photo in the field 'photo' (max 5 MB)"})
		return
	}
	defer file.Close()

	contentType := http.DetectContentType(make([]byte, 0))
	buf := make([]byte, 512)
	n, _ := file.Read(buf)
	if n > 0 {
		contentType = http.DetectContentType(buf[:n])
	}
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		log.Println("error seek photo:", err)
		sendJSON(w, http.StatusInternalServerError, map[string]string{"error": "something went wrong"})
		return
	}
	if !allowedPhotoTypes[contentType] {
		sendJSON(w, http.StatusUnprocessableEntity, map[string]string{"error": "only JPEG, PNG or WebP images are allowed"})
		return
	}

	ext := filepath.Ext(header.Filename)
	if ext == "" {
		ext = ".jpg"
	}
	randomBytes := make([]byte, 16)
	rand.Read(randomBytes)
	fileName := hex.EncodeToString(randomBytes) + ext

	f, err := os.OpenFile(filepath.Join(h.photoDir, fileName), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Println("error saving photo:", err)
		sendJSON(w, http.StatusInternalServerError, map[string]string{"error": "something went wrong"})
		return
	}
	if _, err := io.Copy(f, file); err != nil {
		f.Close()
		os.Remove(filepath.Join(h.photoDir, fileName))
		log.Println("error copying photo:", err)
		sendJSON(w, http.StatusInternalServerError, map[string]string{"error": "something went wrong"})
		return
	}
	f.Close()

	photoURL := "/api/photos/" + fileName
	if err := h.s.UpdatePhotoURL(motherId, photoURL); err != nil {
		sendJSON(w, http.StatusUnprocessableEntity, map[string]string{"error": err.Error()})
		return
	}
	sendJSON(w, http.StatusOK, map[string]string{"photo_url": photoURL})
}

