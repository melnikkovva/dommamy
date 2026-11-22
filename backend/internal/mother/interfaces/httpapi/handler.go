package httpapi

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	mothercommand "backend/internal/mother/application/command"
	motherquery "backend/internal/mother/application/query"
	"backend/internal/mother/domain"
	"backend/internal/shared/httpx"
)

type LearnReader interface {
	GetCourseResults(ctx context.Context, learnUserID string) ([]domain.CourseResultDTO, error)
	GetIPRPlans(ctx context.Context, learnUserID string) ([]domain.IPRPlanDTO, error)
}

type MotherLookup interface {
	ByID(ctx context.Context, id int) (domain.MotherSnapshot, error)
}

type Handler struct {
	createHandler mothercommand.CreateHandler
	updatePhoto   mothercommand.UpdatePhotoHandler
	reads         motherquery.ReadModel
	lookup        MotherLookup
	learnReader   LearnReader
	photoDir      string
}

func NewHandler(
	create mothercommand.CreateHandler,
	updatePhoto mothercommand.UpdatePhotoHandler,
	reads motherquery.ReadModel,
	lookup MotherLookup,
	learnReader LearnReader,
	photoDir string,
) *Handler {
	return &Handler{
		createHandler: create,
		updatePhoto:   updatePhoto,
		reads:         reads,
		lookup:        lookup,
		learnReader:   learnReader,
		photoDir:      photoDir,
	}
}

func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/mother", h.create)
	mux.HandleFunc("GET /api/mother", h.search)
	mux.HandleFunc("POST /api/mother/{id}/photo", h.uploadPhoto)
	mux.HandleFunc("GET /api/mother/{id}/courses", h.courses)
	mux.HandleFunc("GET /api/mother/{id}/ipr", h.ipr)
}

type motherRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	req, err := httpx.ReadJSON[motherRequest](r)
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "incorrect JSON")
		return
	}

	snapshot, err := h.createHandler.Handle(r.Context(), mothercommand.Create(req))
	if err != nil {
		fail(w, "create mother", err)
		return
	}

	httpx.WriteJSON(w, http.StatusCreated, toCreateResponse(snapshot))
}

func (h *Handler) search(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	views, err := h.reads.Search(r.Context(), q)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "server error")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, views)
}

func (h *Handler) courses(w http.ResponseWriter, r *http.Request) {
	motherID, err := httpx.PathID(r, "id")
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "incorrect id")
		return
	}
	snapshot, err := h.lookup.ByID(r.Context(), motherID)
	if err != nil {
		fail(w, "get mother", err)
		return
	}
	if snapshot.LearnUserID == "" || h.learnReader == nil {
		httpx.WriteJSON(w, http.StatusOK, []domain.CourseResultDTO{})
		return
	}
	results, err := h.learnReader.GetCourseResults(r.Context(), snapshot.LearnUserID)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "server error")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, results)
}

func (h *Handler) ipr(w http.ResponseWriter, r *http.Request) {
	motherID, err := httpx.PathID(r, "id")
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "incorrect id")
		return
	}
	snapshot, err := h.lookup.ByID(r.Context(), motherID)
	if err != nil {
		fail(w, "get mother", err)
		return
	}
	if snapshot.LearnUserID == "" || h.learnReader == nil {
		httpx.WriteJSON(w, http.StatusOK, []domain.IPRPlanDTO{})
		return
	}
	plans, err := h.learnReader.GetIPRPlans(r.Context(), snapshot.LearnUserID)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "server error")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, plans)
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
		CreatedAt: s.CreatedAt.Format("2006-01-02T15:04:05"),
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

	r.Body = http.MaxBytesReader(w, r.Body, maxPhotoSize+(1<<20))
	file, header, err := r.FormFile("photo")
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
		httpx.WriteError(w, http.StatusInternalServerError, "something went wrong")
		return
	}
	if _, err := io.Copy(f, file); err != nil {
		f.Close()
		os.Remove(filepath.Join(h.photoDir, fileName))
		log.Println("error copying photo:", err)
		httpx.WriteError(w, http.StatusInternalServerError, "something went wrong")
		return
	}
	f.Close()

	photoURL := "/api/photos/" + fileName
	if err := h.updatePhoto.Handle(r.Context(), mothercommand.UpdatePhoto{
		MotherID: motherID,
		PhotoURL: photoURL,
	}); err != nil {
		fail(w, "update photo", err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, map[string]string{"photo_url": photoURL})
}

func fail(w http.ResponseWriter, where string, err error) {
	switch {
	case errors.Is(err, domain.ErrName):
		httpx.WriteError(w, http.StatusUnprocessableEntity, err.Error())
	case errors.Is(err, domain.ErrNotFound):
		httpx.WriteError(w, http.StatusNotFound, err.Error())
	default:
		log.Println("error "+where+":", err)
		httpx.WriteError(w, http.StatusInternalServerError, "something went wrong")
	}
}
