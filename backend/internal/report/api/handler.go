package api

import (
	"encoding/csv"
	"net/http"
	"strings"
	"time"

	"backend/internal/report"
	"backend/internal/shared/httpx"
)

type Handler struct {
	s      *report.Service
	mapper *httpx.ErrorMapper
}

func New(s *report.Service) *Handler {
	return &Handler{
		s: s,
		mapper: httpx.NewErrorMapper(
			httpx.WithNotFound(report.ErrorUnknown),
			httpx.WithValidation(report.ErrorPeriod),
		),
	}
}

func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/reports/{name}", h.get)
	mux.HandleFunc("GET /api/reports/{name}/csv", h.download)
}

type totalResponse struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

type reportResponse struct {
	Name string `json:"name"`
	Title string `json:"title"`
	Headers []string `json:"headers"`
	Rows [][]string  `json:"rows"`
	Totals []totalResponse `json:"totals"`
}

func filterFrom(r *http.Request) report.Filter {
	q := r.URL.Query()
	return report.Filter{From: q.Get("from"), To: q.Get("to")}
}

func (h *Handler) get(w http.ResponseWriter, r *http.Request) {
	rep, err := h.s.Build(r.PathValue("name"), filterFrom(r))
	if err != nil {
		h.mapper.Handle(w, "get report", err)
		return
	}
	totals := []totalResponse{}
	for _, t := range rep.Totals {
		totals = append(totals, totalResponse{Label: t.Label, Value: t.Value})
	}
	httpx.WriteJSON(w, http.StatusOK, reportResponse{Name: rep.Name, Title: rep.Title, Headers: rep.Headers, Rows: rep.Rows, Totals: totals})
}

func (h *Handler) download(w http.ResponseWriter, r *http.Request) {
	rep, err := h.s.Build(r.PathValue("name"), filterFrom(r))
	if err != nil {
		h.mapper.Handle(w, "download report", err)
		return
	}
	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w.Header().Set("Content-Disposition", `attachment; filename="report-`+rep.Name+"-"+time.Now().Format("2006-01-02")+`.csv"`)
	w.Write([]byte("\xEF\xBB\xBF"))
	cw := csv.NewWriter(w)
	cw.Comma = ';'
	cw.Write(rep.Headers)
	for _, row := range rep.Rows {
		safe := make([]string, len(row))
		for i, cell := range row {
			safe[i] = safeCell(cell)
		}
		cw.Write(safe)
	}
	cw.Flush()
}

func safeCell(s string) string {
	if s == "" {
		return s
	}
	switch s[0] {
	case '=', '@', '\t', '\r':
		return "'" + s
	case '+', '-':
		if !looksLikePhone(s) {
			return "'" + s
		}
	}
	return s
}

func looksLikePhone(s string) bool {
	for _, c := range s {
		if !strings.ContainsRune("+-0123456789 ()", c) {
			return false
		}
	}
	return true
}