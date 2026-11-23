package auth

import (
	"context"
	"net/http"
	"slices"
)

func RequireRole(roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, err := FromContext(r.Context())
			if err != nil {
				http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
				return
			}
			if !slices.Contains(roles, claims.Role) {
				http.Error(w, `{"error":"forbidden"}`, http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func RequireOwnership(r *http.Request, ownerID int) bool {
	claims, err := FromContext(r.Context())
	if err != nil {
		return false
	}
	if claims.Role == "admin" {
		return true
	}
	return claims.UserID == ownerID || claims.MotherID == ownerID
}

func RequireOwnershipOrAdmin(w http.ResponseWriter, r *http.Request, ownerID int) bool {
	claims, err := FromContext(r.Context())
	if err != nil {
		http.Error(w, `{"error":"forbidden"}`, http.StatusForbidden)
		return false
	}
	if claims.Role == "admin" {
		return true
	}
	if claims.UserID == ownerID || claims.MotherID == ownerID {
		return true
	}
	http.Error(w, `{"error":"forbidden"}`, http.StatusForbidden)
	return false
}

func CurrentUser(ctx context.Context) (Claims, error) {
	return FromContext(ctx)
}
