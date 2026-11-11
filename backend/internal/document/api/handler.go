package api

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"mime"
	"net/http"
	"strconv"
	"time"

	"backend/internal/document"
)

type Handler struct {
	s *document.Service
}

func New(s *document.Service) *Handler {
	return &Handler{s: s}
}

func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/mother/{id}/documents", h.create)
	mux.HandleFunc("GET /api/mother/{id}/documents", h.list)
	mux.HandleFunc("POST /api/documents/{id}/file", h.upload)
	mux.HandleFunc("GET /api/documents/{id}/file", h.download)
	mux.HandleFunc("PUT /api/documents/{id}/sign", h.sign)
}

type response struct {
	Id int `json:"id"`
	MotherId int `json:"mother_id"`
	Type string `json:"type"`
	HelpId int `json:"help_id"`
	Status string `json:"status"`
	SignedAt string `json:"signed_at"`
	HasFile bool`json:"has_file"`
	OriginalName string `json:"original_name"`
	ContentType string `json:"content_type"`
	Size int64`json:"size"`
	CreatedAt string `json:"created_at"`
}

func toResponse(d *document.Document) response {
	return response{
		Id: d.Id, MotherId: d.MotherId, Type: d.Type, HelpId: d.HelpId,
		Status: d.Status, SignedAt: d.SignedAt, HasFile: d.HasFile(),
		OriginalName: d.OriginalName, ContentType: d.ContentType, Size: d.Size,
		CreatedAt: d.CreatedAt.Format(time.RFC3339),
	}
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	motherId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "incorrect id"})
		return
	}
	var req struct {
		Type string `json:"type"`
		HelpId int `json:"help_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "incorrect JSON"})
		return
	}
	d, err := h.s.Create(motherId, req.Type, req.HelpId)
	if err != nil {
		fail(w, "create document", err)
		return
	}
	writeJSON(w, http.StatusCreated, toResponse(d))
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	motherId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "incorrect id"})
		return
	}
	items, err := h.s.ByMother(motherId)
	if err != nil {
		fail(w, "list documents", err)
		return
	}
	result := []response{}
	for _, d := range items {
		result = append(result, toResponse(d))
	}
	writeJSON(w, http.StatusOK, result)
}

func (h *Handler) upload(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "incorrect id"})
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, document.MaxFileSize+(1<<20))
	file, header, err := r.FormFile("file")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "send the file in the field 'file' (max 10 MB)"})
		return
	}
	defer file.Close()

	head := make([]byte, 512)
	n, _ := file.Read(head)
	contentType := http.DetectContentType(head[:n])
	if _, err := file.Seek(0, io.SeekStart); err != nil { 
		log.Println("error seek upload:", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "something went wrong"})
		return
	}

	d, err := h.s.AttachFile(id, file, header.Size, contentType, header.Filename)
	if err != nil {
		fail(w, "upload file", err)
		return
	}
	writeJSON(w, http.StatusOK, toResponse(d))
}

func (h *Handler) download(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "incorrect id"})
		return
	}
	d, f, err := h.s.Open(id)
	if err != nil {
		fail(w, "download file", err)
		return
	}
	defer f.Close()

	w.Header().Set("Content-Type", d.ContentType)
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Content-Disposition", mime.FormatMediaType("inline", map[string]string{"filename": d.OriginalName}))
	io.Copy(w, f)
}

func (h *Handler) sign(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "incorrect id"})
		return
	}
	var req struct {
		SignedAt string `json:"signed_at"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	d, err := h.s.Sign(id, req.SignedAt)
	if err != nil {
		fail(w, "sign document", err)
		return
	}
	writeJSON(w, http.StatusOK, toResponse(d))
}

func fail(w http.ResponseWriter, where string, err error) {
	switch {
	case errors.Is(err, document.ErrorNotFound):
		writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
	case errors.Is(err, document.ErrorMother), errors.Is(err, document.ErrorNoMother),
		errors.Is(err, document.ErrorType), errors.Is(err, document.ErrorNeedHelp),
		errors.Is(err, document.ErrorHelpNotAllowed), errors.Is(err, document.ErrorNoHelp),
		errors.Is(err, document.ErrorHelpMismatch), errors.Is(err, document.ErrorDuplicate),
		errors.Is(err, document.ErrorSigned), errors.Is(err, document.ErrorNoFile),
		errors.Is(err, document.ErrorEmptyFile), errors.Is(err, document.ErrorTooBig),
		errors.Is(err, document.ErrorFileType), errors.Is(err, document.ErrorSignDate):
		writeJSON(w, http.StatusUnprocessableEntity, map[string]string{"error": err.Error()})
	default:
		log.Println("error", where+":", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "something went wrong"})
	}
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}