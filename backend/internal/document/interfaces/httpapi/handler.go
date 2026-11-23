package httpapi

import (
	"io"
	"log"
	"mime"
	"net/http"
	"path"
	"slices"
	"strings"

	"backend/internal/auth"
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
		mapper: httpx.NewErrorMapper(
			httpx.WithNotFound(domain.ErrNotFound),
			httpx.WithValidation(
				domain.ErrMother, domain.ErrNoMother, domain.ErrType,
				domain.ErrNeedHelp, domain.ErrHelpNotAllowed, domain.ErrNoHelp,
				domain.ErrHelpMismatch, domain.ErrDuplicate, domain.ErrSigned,
				domain.ErrNoFile, domain.ErrEmptyFile, domain.ErrTooBig,
				domain.ErrFileType, domain.ErrSignDate,
			),
		),
	}
}

type handler struct {
	createCmd    command.CreateHandler
	attachFile   command.AttachFileHandler
	downloadCmd  command.DownloadHandler
	signCmd      command.SignHandler
	reads        query.ReadModel
	mapper       *httpx.ErrorMapper
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
		CreatedAt: v.CreatedAt.Format(httpx.APITime),
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
	if !auth.RequireOwnershipOrAdmin(w, r, motherID) {
		return
	}
	body, err := httpx.ReadJSON[createBody](w, r)
	if err != nil {
		if httpx.IsMaxBytesExceeded(err) {
			httpx.WriteError(w, http.StatusRequestEntityTooLarge, "request body too large")
			return
		}
		httpx.WriteError(w, http.StatusBadRequest, "incorrect JSON")
		return
	}
	id, err := h.createCmd.Handle(r.Context(), command.Create{
		MotherID: motherID,
		Type:     body.Type,
		HelpID:   body.HelpID,
	})
	if err != nil {
		h.mapper.Handle(w, "create document", err)
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
	if !auth.RequireOwnershipOrAdmin(w, r, motherID) {
		return
	}
	items, err := h.reads.ByMother(r.Context(), motherID)
	if err != nil {
		h.mapper.Handle(w, "list documents", err)
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
	sniffed := http.DetectContentType(head[:n])
	if slices.Contains(domain.DeniedContentTypes(), sniffed) {
		httpx.WriteError(w, http.StatusUnprocessableEntity, domain.ErrFileType.Error())
		return
	}
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		log.Println("error seek upload:", err)
		httpx.WriteError(w, http.StatusInternalServerError, "something went wrong")
		return
	}

	ext := strings.ToLower(path.Ext(header.Filename))
	contentType := domain.AllowedExtensions()[ext]

	docID, err := h.attachFile.Handle(r.Context(), command.AttachFile{
		DocumentID:   id,
		Content:      file,
		Size:         header.Size,
		ContentType:  contentType,
		OriginalName: header.Filename,
	})
	if err != nil {
		h.mapper.Handle(w, "upload file", err)
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
		h.mapper.Handle(w, "download file", err)
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
	body, err := httpx.ReadJSON[signBody](w, r)
	if err != nil {
		if httpx.IsMaxBytesExceeded(err) {
			httpx.WriteError(w, http.StatusRequestEntityTooLarge, "request body too large")
			return
		}
		httpx.WriteError(w, http.StatusBadRequest, "incorrect JSON")
		return
	}
	docID, err := h.signCmd.Handle(r.Context(), command.Sign{
		DocumentID: id,
		SignedAt:   body.SignedAt,
	})
	if err != nil {
		h.mapper.Handle(w, "sign document", err)
		return
	}
	h.respondWith(w, r, http.StatusOK, docID)
}

func (h *handler) respondWith(w http.ResponseWriter, r *http.Request, status, id int) {
	v, err := h.reads.GetByID(r.Context(), id)
	if err != nil {
		h.mapper.Handle(w, "get document", err)
		return
	}
	httpx.WriteJSON(w, status, toResponse(v))
}
