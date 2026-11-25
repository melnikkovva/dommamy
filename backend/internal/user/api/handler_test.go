package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"backend/internal/auth"
	"backend/internal/user"
)

type deleteTestRepo struct {
	deletedID int
	deleteErr error
}

func (r *deleteTestRepo) Save(*user.User) error { return nil }

func (r *deleteTestRepo) UserById(int) (*user.User, error) {
	return nil, user.ErrorNotFound
}

func (r *deleteTestRepo) UserByLogin(string) (*user.User, error) {
	return nil, user.ErrorNotFound
}

func (r *deleteTestRepo) ListUsers() ([]*user.User, error) { return nil, nil }

func (r *deleteTestRepo) DeleteUser(id int) error {
	r.deletedID = id
	return r.deleteErr
}

func TestDeleteUserRequiresAdmin(t *testing.T) {
	repo := &deleteTestRepo{}
	handler := New(user.NewService(repo, "secret", time.Hour))
	mux := http.NewServeMux()
	handler.Register(mux)
	request := httptest.NewRequest(http.MethodDelete, "/api/users/42", nil)
	request = request.WithContext(auth.WithClaims(context.Background(), auth.Claims{UserID: 1, Role: user.RoleUser}))
	response := httptest.NewRecorder()

	mux.ServeHTTP(response, request)

	if response.Code != http.StatusForbidden {
		t.Fatalf("expected status %d, got %d", http.StatusForbidden, response.Code)
	}
	if repo.deletedID != 0 {
		t.Fatalf("expected no user to be deleted, got id %d", repo.deletedID)
	}
}

func TestAdminCanSoftDeleteAnotherUser(t *testing.T) {
	repo := &deleteTestRepo{}
	handler := New(user.NewService(repo, "secret", time.Hour))
	mux := http.NewServeMux()
	handler.Register(mux)
	request := httptest.NewRequest(http.MethodDelete, "/api/users/42", nil)
	request = request.WithContext(auth.WithClaims(context.Background(), auth.Claims{UserID: 1, Role: user.RoleAdmin}))
	response := httptest.NewRecorder()

	mux.ServeHTTP(response, request)

	if response.Code != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d", http.StatusNoContent, response.Code)
	}
	if repo.deletedID != 42 {
		t.Fatalf("expected user 42 to be deleted, got id %d", repo.deletedID)
	}
}

func TestAdminCannotDeleteOwnAccount(t *testing.T) {
	repo := &deleteTestRepo{}
	handler := New(user.NewService(repo, "secret", time.Hour))
	mux := http.NewServeMux()
	handler.Register(mux)
	request := httptest.NewRequest(http.MethodDelete, "/api/users/42", nil)
	request = request.WithContext(auth.WithClaims(context.Background(), auth.Claims{UserID: 42, Role: user.RoleAdmin}))
	response := httptest.NewRecorder()

	mux.ServeHTTP(response, request)

	if response.Code != http.StatusConflict {
		t.Fatalf("expected status %d, got %d", http.StatusConflict, response.Code)
	}
	if repo.deletedID != 0 {
		t.Fatalf("expected no user to be deleted, got id %d", repo.deletedID)
	}
}
