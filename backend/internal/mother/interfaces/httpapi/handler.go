package httpapi

import (
	"context"
	"io"
	"log"
	"net/http"

	"backend/internal/auth"
	mothercommand "backend/internal/mother/application/command"
	motherquery "backend/internal/mother/application/query"
	"backend/internal/mother/domain"
	"backend/internal/shared/httpx"
)

type LearnLinkRepository interface {
	Link(ctx context.Context, motherID int, learnUserID string) error
	GetByMotherID(ctx context.Context, motherID int) (string, error)
	Unlink(ctx context.Context, motherID int) error
}

type LearnerReader interface {
	GetAllLearners(ctx context.Context) ([]LearnerDTO, error)
}

type LearnerDTO struct {
	UserID    string `json:"user_id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

type Handler struct {
	createHandler   mothercommand.CreateHandler
	uploadPhotoCmd  mothercommand.UploadPhotoHandler
	softDeleteCmd   mothercommand.SoftDeleteHandler
	getCourses      motherquery.GetCoursesHandler
	getIPR          motherquery.GetIPRHandler
	reads           motherquery.ReadModel
	learnLinks      LearnLinkRepository
	learnerReader   LearnerReader
	mapper          *httpx.ErrorMapper
}

func NewHandler(
	create mothercommand.CreateHandler,
	uploadPhoto mothercommand.UploadPhotoHandler,
	softDelete mothercommand.SoftDeleteHandler,
	getCourses motherquery.GetCoursesHandler,
	getIPR motherquery.GetIPRHandler,
	reads motherquery.ReadModel,
	learnLinks LearnLinkRepository,
	learnerReader LearnerReader,
) *Handler {
	return &Handler{
		createHandler:  create,
		uploadPhotoCmd: uploadPhoto,
		softDeleteCmd:  softDelete,
		getCourses:     getCourses,
		getIPR:         getIPR,
		reads:          reads,
		learnLinks:     learnLinks,
		learnerReader:  learnerReader,
		mapper: httpx.NewErrorMapper(
			httpx.WithValidation(domain.ErrName),
			httpx.WithNotFound(domain.ErrNotFound),
		),
	}
}

func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/mother", h.create)
	mux.HandleFunc("GET /api/mother", h.search)
	mux.HandleFunc("DELETE /api/mother/{id}", h.softDelete)
	mux.HandleFunc("POST /api/mother/{id}/photo", h.uploadPhoto)
	mux.HandleFunc("GET /api/mother/{id}/courses", h.courses)
	mux.HandleFunc("GET /api/mother/{id}/ipr", h.ipr)
	mux.HandleFunc("PUT /api/mother/{id}/learn-link", h.linkLearn)
	mux.HandleFunc("DELETE /api/mother/{id}/learn-link", h.unlinkLearn)
	mux.HandleFunc("GET /api/learners", h.listLearners)
}

func (h *Handler) softDelete(w http.ResponseWriter, r *http.Request) {
	motherID, err := httpx.PathID(r, "id")
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "incorrect id")
		return
	}
	err = h.softDeleteCmd.Handle(r.Context(), mothercommand.SoftDelete{MotherID: motherID})
	if err != nil {
		h.mapper.Handle(w, "soft delete mother", err)
		return
	}
	httpx.WriteJSON(w, http.StatusNoContent, nil)
}

type motherRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	req, err := httpx.ReadJSON[motherRequest](w, r)
	if err != nil {
		if httpx.IsMaxBytesExceeded(err) {
			httpx.WriteError(w, http.StatusRequestEntityTooLarge, "request body too large")
			return
		}
		httpx.WriteError(w, http.StatusBadRequest, "incorrect JSON")
		return
	}

	snapshot, err := h.createHandler.Handle(r.Context(), mothercommand.Create(req))
	if err != nil {
		h.mapper.Handle(w, "create mother", err)
		return
	}

	httpx.WriteJSON(w, http.StatusCreated, toCreateResponse(snapshot))
}

type motherSearchResponse struct {
	ID            int    `json:"id"`
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	Phone         string `json:"phone"`
	Email         string `json:"email"`
	LearnUserID   string `json:"learn_user_id"`
	PhotoURL      string `json:"photo_url,omitempty"`
	CreatedAt     string `json:"created_at"`
	ProfileStatus string `json:"profile_status"`
	MotherStatus  string `json:"mother_status"`
}

func toSearchResponse(v motherquery.MotherView) motherSearchResponse {
	return motherSearchResponse{
		ID: v.ID, FirstName: v.FirstName, LastName: v.LastName,
		Phone: v.Phone, Email: v.Email, LearnUserID: v.LearnUserID,
		PhotoURL: v.PhotoURL, CreatedAt: v.CreatedAt.Format(httpx.APITime),
		ProfileStatus: v.ProfileStatus, MotherStatus: v.MotherStatus,
	}
}

func (h *Handler) search(w http.ResponseWriter, r *http.Request) {
	filter := motherquery.MotherFilter{
		Q:             r.URL.Query().Get("q"),
		Phone:         r.URL.Query().Get("phone"),
		ProfileStatus: r.URL.Query().Get("profileStatus"),
		MotherStatus:  r.URL.Query().Get("motherStatus"),
	}
	views, err := h.reads.Search(r.Context(), filter)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "server error")
		return
	}
	result := make([]motherSearchResponse, 0, len(views))
	for _, v := range views {
		result = append(result, toSearchResponse(v))
	}
	httpx.WriteJSON(w, http.StatusOK, result)
}

type courseResultResponse struct {
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

func toCourseResultResponse(cr motherquery.CourseResultDTO) courseResultResponse {
	return courseResultResponse{
		CourseID: cr.CourseID, NameCourse: cr.NameCourse, AccesDate: cr.AccesDate,
		DueDate: cr.DueDate, CompletionDate: cr.CompletionDate,
		TimeSpent: cr.TimeSpent, Duration: cr.Duration,
		CompletionStatus: cr.CompletionStatus, Progress: cr.Progress,
		AwardedScore: cr.AwardedScore, IsOverdue: cr.IsOverdue, ViewsCount: cr.ViewsCount,
	}
}

func (h *Handler) courses(w http.ResponseWriter, r *http.Request) {
	motherID, err := httpx.PathID(r, "id")
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "incorrect id")
		return
	}
	if !auth.RequireOwnershipOrAdmin(w, r, motherID) {
		return
	}
	results, err := h.getCourses.Handle(r.Context(), motherquery.GetCourses{MotherID: motherID})
	if err != nil {
		h.mapper.Handle(w, "get courses", err)
		return
	}
	resp := make([]courseResultResponse, 0, len(results))
	for _, cr := range results {
		resp = append(resp, toCourseResultResponse(cr))
	}
	httpx.WriteJSON(w, http.StatusOK, resp)
}

type iprTaskResponse struct {
	TaskID          string  `json:"task_id"`
	TaskType        string  `json:"task_type"`
	Title           string  `json:"title"`
	Description     string  `json:"description"`
	Status          string  `json:"status"`
	IsOverdue       bool    `json:"is_overdue"`
	AbsoluteDueDate *string `json:"absolute_due_date,omitempty"`
	CompletionDate  *string `json:"completion_date,omitempty"`
}

type iprStageResponse struct {
	StageID     string          `json:"stage_id"`
	Title       string          `json:"title"`
	Description string          `json:"description"`
	StartDate   *string         `json:"start_date,omitempty"`
	EndDate     *string         `json:"end_date,omitempty"`
	Tasks       []iprTaskResponse `json:"tasks"`
}

type iprPlanResponse struct {
	PlanID         string          `json:"plan_id"`
	Title          string          `json:"title"`
	EmployeeID     string          `json:"employee_id"`
	MentorID       string          `json:"mentor_id"`
	MentorName     string          `json:"mentor_name"`
	StartDate      string          `json:"start_date"`
	EndDate        *string         `json:"end_date,omitempty"`
	Status         string          `json:"status"`
	Progress       float64         `json:"progress"`
	CompletedTasks int             `json:"completed_tasks"`
	TotalTasks     int             `json:"total_tasks"`
	Goals          []string        `json:"goals"`
	Stages         []iprStageResponse `json:"stages"`
}

func toIPRPlanResponse(p motherquery.IPRPlanDTO) iprPlanResponse {
	stages := make([]iprStageResponse, 0, len(p.Stages))
	for _, s := range p.Stages {
		tasks := make([]iprTaskResponse, 0, len(s.Tasks))
		for _, t := range s.Tasks {
			tasks = append(tasks, iprTaskResponse{
				TaskID: t.TaskID, TaskType: t.TaskType, Title: t.Title,
				Description: t.Description, Status: t.Status, IsOverdue: t.IsOverdue,
				AbsoluteDueDate: t.AbsoluteDueDate, CompletionDate: t.CompletionDate,
			})
		}
		stages = append(stages, iprStageResponse{
			StageID: s.StageID, Title: s.Title, Description: s.Description,
			StartDate: s.StartDate, EndDate: s.EndDate, Tasks: tasks,
		})
	}
	return iprPlanResponse{
		PlanID: p.PlanID, Title: p.Title, EmployeeID: p.EmployeeID,
		MentorID: p.MentorID, MentorName: p.MentorName,
		StartDate: p.StartDate, EndDate: p.EndDate,
		Status: p.Status, Progress: p.Progress,
		CompletedTasks: p.CompletedTasks, TotalTasks: p.TotalTasks,
		Goals: p.Goals, Stages: stages,
	}
}

func (h *Handler) ipr(w http.ResponseWriter, r *http.Request) {
	motherID, err := httpx.PathID(r, "id")
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "incorrect id")
		return
	}
	if !auth.RequireOwnershipOrAdmin(w, r, motherID) {
		return
	}
	plans, err := h.getIPR.Handle(r.Context(), motherquery.GetIPR{MotherID: motherID})
	if err != nil {
		h.mapper.Handle(w, "get ipr", err)
		return
	}
	resp := make([]iprPlanResponse, 0, len(plans))
	for _, p := range plans {
		resp = append(resp, toIPRPlanResponse(p))
	}
	httpx.WriteJSON(w, http.StatusOK, resp)
}

type createResponse struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
}

func toCreateResponse(s domain.MotherSnapshot) createResponse {
	return createResponse{
		ID:        s.ID,
		FirstName: s.FirstName,
		LastName:  s.LastName,
		Phone:     s.Phone,
		Email:     s.Email,
		CreatedAt: s.CreatedAt.Format(httpx.APITime),
	}
}

const maxPhotoSize = 5 << 20 // 5 MB

var allowedPhotoTypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/webp": true,
}

func (h *Handler) uploadPhoto(w http.ResponseWriter, r *http.Request) {
	motherID, err := httpx.PathID(r, "id")
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "incorrect id")
		return
	}
	if !auth.RequireOwnershipOrAdmin(w, r, motherID) {
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxPhotoSize+(1<<20))
	file, _, err := r.FormFile("photo")
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "send the photo in the field 'photo' (max 5 MB)")
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
		httpx.WriteError(w, http.StatusInternalServerError, "something went wrong")
		return
	}
	if !allowedPhotoTypes[contentType] {
		httpx.WriteError(w, http.StatusUnprocessableEntity, "only JPEG, PNG or WebP images are allowed")
		return
	}

	data, err := io.ReadAll(io.LimitReader(file, maxPhotoSize))
	if err != nil {
		log.Println("error reading photo:", err)
		httpx.WriteError(w, http.StatusInternalServerError, "something went wrong")
		return
	}

	photoURL, err := h.uploadPhotoCmd.Handle(r.Context(), mothercommand.UploadPhoto{
		MotherID:    motherID,
		Data:        data,
		ContentType: contentType,
	})
	if err != nil {
		h.mapper.Handle(w, "upload photo", err)
		return
	}
 httpx.WriteJSON(w, http.StatusOK, map[string]string{"photo_url": photoURL})
}

type linkLearnRequest struct {
	LearnUserID string `json:"learn_user_id"`
}

func (h *Handler) linkLearn(w http.ResponseWriter, r *http.Request) {
	motherID, err := httpx.PathID(r, "id")
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "incorrect id")
		return
	}
	req, err := httpx.ReadJSON[linkLearnRequest](w, r)
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "incorrect JSON")
		return
	}
	if req.LearnUserID == "" {
		httpx.WriteError(w, http.StatusUnprocessableEntity, "specify learn_user_id")
		return
	}
	if err := h.learnLinks.Link(r.Context(), motherID, req.LearnUserID); err != nil {
		log.Printf("link learn error: %v", err)
		httpx.WriteError(w, http.StatusInternalServerError, "server error")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) unlinkLearn(w http.ResponseWriter, r *http.Request) {
	motherID, err := httpx.PathID(r, "id")
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "incorrect id")
		return
	}
	if err := h.learnLinks.Unlink(r.Context(), motherID); err != nil {
		log.Printf("unlink learn error: %v", err)
		httpx.WriteError(w, http.StatusInternalServerError, "server error")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) listLearners(w http.ResponseWriter, r *http.Request) {
	if h.learnerReader == nil {
		httpx.WriteJSON(w, http.StatusOK, []LearnerDTO{})
		return
	}
	learners, err := h.learnerReader.GetAllLearners(r.Context())
	if err != nil {
		log.Printf("list learners error: %v", err)
		httpx.WriteError(w, http.StatusInternalServerError, "server error")
		return
	}
	dtos := make([]LearnerDTO, len(learners))
	for i, l := range learners {
		dtos[i] = LearnerDTO{
			UserID:    l.UserID,
			FirstName: l.FirstName,
			LastName:  l.LastName,
			Email:     l.Email,
		}
	}
	httpx.WriteJSON(w, http.StatusOK, dtos)
}
