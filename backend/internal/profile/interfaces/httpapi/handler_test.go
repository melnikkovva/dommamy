package httpapi

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"backend/internal/auth"
	"backend/internal/profile/application/command"
	"backend/internal/profile/domain"
)

type profileSaveTestRepository struct {
	saved domain.Profile
}

func (r *profileSaveTestRepository) Save(_ context.Context, p domain.Profile) error {
	r.saved = p
	return nil
}

func (r *profileSaveTestRepository) ByMother(context.Context, int) (domain.Profile, error) {
	return domain.Profile{}, domain.ErrNotFilled
}

func (r *profileSaveTestRepository) UpdateStatus(context.Context, int, domain.Status) error {
	return nil
}

type profileSaveTestMothers struct{}

func (profileSaveTestMothers) Exists(context.Context, int) (bool, error) {
	return true, nil
}

func TestSaveProfilePersistsProblemType(t *testing.T) {
	repo := &profileSaveTestRepository{}
	save := command.NewSaveHandler(repo, profileSaveTestMothers{}, func() time.Time {
		return time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC)
	})
	handler := NewHandler(save, nil, nil)
	mux := http.NewServeMux()
	handler.Register(mux)

	request := httptest.NewRequest(http.MethodPut, "/api/mother/12/profile", strings.NewReader(`{"problem_type":"housing"}`))
	request = request.WithContext(auth.WithClaims(request.Context(), auth.Claims{UserID: 12, MotherID: 12, Role: "mom"}))
	response := httptest.NewRecorder()
	mux.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d: %s", http.StatusOK, response.Code, response.Body.String())
	}
	if got := repo.saved.ProblemType(); got != "housing" {
		t.Fatalf("expected problem type housing to be persisted, got %q", got)
	}
}