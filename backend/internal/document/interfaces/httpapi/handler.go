package httpapi

import (
	"errors"
	"io"
	"log"
	"mime"
	"net/http"
	"slices"
	"time"

	"backend/internal/document/application/command"
	"backend/internal/document/application/query"
	"backend/internal/document/domain"
	"backend/internal/shared/httpx"
)

func NewHandler(
	create command.CreateHandler,
	attachFile command.AttachFileHandler,
	download command.DownloadHandler,
	sign command.SignHandler,
	reads query.ReadModel,
) httpx.Router {
 return &handler{
		createCmd:    create,
		attachFile:   attachFile,
		downloadCmd:  download,
		signCmd:      sign,
		reads:        reads,
	}
}

type handler struct {
	createCmd    command.CreateHandler
	attachFile   command.AttachFileHandler
	downloadCmd  command.DownloadHandler
	signCmd      command.SignHandler
	reads        query.ReadModel
}

func (h *handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/mother/{id}/documents", h.create)
	mux.HandleFunc("GET /api/mother/{id}/documents", h.list)
	mux.HandleFunc("POST /api/documents/{id}/file", h.upload)
	mux.HandleFunc("GET /api/documents/{id}/file", h.download)
	mux.HandleFunc("PUT /api/documents/{id}/sign", h.signDocument)
}

type documentResponse struct {
	ID           int    `json:"id"`
	MotherID     int    `json:"mother_id"`
	Type         string `json:"type"`
	HelpID       int    `json:"help_id"`
	Status       string `json:"status"`
	SignedAt     string `json:"signed_at"`
	HasFile      bool   `json:"has_file"`
	OriginalName string `json:"original_name"`
	ContentType  string `json:"content_type"`
	Size         int64  `json:"size"`
	CreatedAt    string `json:"created_at"`
}

func toResponse(v query.DocumentView) documentResponse {
	return documentResponse{
		ID: v.ID, MotherID: v.MotherID, Type: string(v.Type), HelpID: v.HelpID,
		Status: string(v.Status), SignedAt: v.SignedAt, HasFile: v.HasFile,
		OriginalName: v.OriginalName, ContentType: v.ContentType, Size: v.Size,
		CreatedAt: v.CreatedAt.Format(time.RFC3339),
	}
}

func toResponses(items []query.DocumentView) []documentResponse {
	result := make([]documentResponse, 0, len(items))
	for _, v := range items {
		result = append(result, toResponse(v))
	}
	return result
}

type createBody struct {
	Type   string `json:"type"`
	HelpID int    `json:"help_id"`
}

type signBody struct {
	SignedAt string `json:"signed_at"`
}

func (h *handler) create(w http.ResponseWriter, r *http.Request) {
	motherID, err := httpx.PathID(r, "id")
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "incorrect id")
		return
	}
	body, err := httpx.ReadJSON[createBody](r)
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "incorrect JSON")
		return
	}
	id, err := h.createCmd.Handle(r.Context(), command.Create{
		MotherID: motherID,
		Type:     body.Type,
		HelpID:   body.HelpID,
	})
	if err != nil {
		fail(w, "create document", err)
		return
	}
	h.respondWith(w, r, http.StatusCreated, id)
}

func (h *handler) list(w http.ResponseWriter, r *http.Request) {
	motherID, err := httpx.PathID(r, "id")
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "incorrect id")
		return
	}
	items, err := h.reads.ByMother(r.Context(), motherID)
	if err != nil {
		fail(w, "list documents", err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, toResponses(items))
}

func (h *handler) upload(w http.ResponseWriter, r *http.Request) {
	id, err := httpx.PathID(r, "id")
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "incorrect id")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, domain.MaxFileSize+(1<<20))
	file, header, err := r.FormFile("file")
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "send the file in the field 'file' (max 10 MB)")
		return
	}
	defer file.Close()

	head := make([]byte, 512)
	n, _ := file.Read(head)
	contentType := http.DetectContentType(head[:n])
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		log.Println("error seek upload:", err)
		httpx.WriteError(w, http.StatusInternalServerError, "something went wrong")
		return
	}

	docID, err := h.attachFile.Handle(r.Context(), command.AttachFile{
		DocumentID:   id,
		Content:      file,
		Size:         header.Size,
		ContentType:  contentType,
		OriginalName: header.Filename,
	})
	if err != nil {
		fail(w, "upload file", err)
		return
	}
	h.respondWith(w, r, http.StatusOK, docID)
}

func (h *handler) download(w http.ResponseWriter, r *http.Request) {
	id, err := httpx.PathID(r, "id")
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "incorrect id")
		return
	}

	df, err := h.downloadCmd.Handle(r.Context(), command.Download{DocumentID: id})
	if err != nil {
		fail(w, "download file", err)
		return
	}
	defer df.Reader.Close()

	w.Header().Set("Content-Type", df.ContentType)
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Content-Disposition", mime.FormatMediaType("inline", map[string]string{"filename": df.OriginalName}))
	io.Copy(w, df.Reader)
}

func (h *handler) signDocument(w http.ResponseWriter, r *http.Request) {
	id, err := httpx.PathID(r, "id")
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "incorrect id")
		return
	}
	body, err := httpx.ReadJSON[signBody](r)
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "incorrect JSON")
		return
	}
	docID, err := h.signCmd.Handle(r.Context(), command.Sign{
		DocumentID: id,
		SignedAt:   body.SignedAt,
	})
	if err != nil {
		fail(w, "sign document", err)
		return
	}
	h.respondWith(w, r, http.StatusOK, docID)
}

func (h *handler) respondWith(w http.ResponseWriter, r *http.Request, status, id int) {
	v, err := h.reads.GetByID(r.Context(), id)
	if err != nil {
		fail(w, "get document", err)
		return
	}
	httpx.WriteJSON(w, status, toResponse(v))
}

var validationErrors = []error{
	domain.ErrMother, domain.ErrNoMother, domain.ErrType,
	domain.ErrNeedHelp, domain.ErrHelpNotAllowed, domain.ErrNoHelp,
	domain.ErrHelpMismatch, domain.ErrDuplicate, domain.ErrSigned,
	domain.ErrNoFile, domain.ErrEmptyFile, domain.ErrTooBig,
	domain.ErrFileType, domain.ErrSignDate,
}

func fail(w http.ResponseWriter, where string, err error) {
	switch {
	case errors.Is(err, domain.ErrNotFound):
		httpx.WriteError(w, http.StatusNotFound, err.Error())
	case slices.ContainsFunc(validationErrors, func(e error) bool {
		return errors.Is(err, e)
	}):
		httpx.WriteError(w, http.StatusUnprocessableEntity, err.Error())
	default:
		log.Println("error", where+":", err)
		httpx.WriteError(w, http.StatusInternalServerError, "something went wrong")
	}
}
