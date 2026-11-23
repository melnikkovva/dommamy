package auth

import (
	"net/http"
	"strings"
)

func PublicPath(method, path string) bool {
	for _, p := range publicPaths {
		if p.method == method && matchPath(p.pattern, path) {
			return true
		}
	}
	return false
}

type publicPath struct {
	method  string
	pattern string
}

var publicPaths = []publicPath{
	{"POST", "/api/auth/login"},
	{"POST", "/api/auth/me"},
	{"GET", "/health"},
}

func Middleware(secret []byte) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if PublicPath(r.Method, r.URL.Path) {
				next.ServeHTTP(w, r)
				return
			}

			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, `{"error":"missing authorization token"}`, http.StatusUnauthorized)
				return
			}

			var tokenString string
			if strings.HasPrefix(authHeader, "Bearer ") {
				tokenString = strings.TrimPrefix(authHeader, "Bearer ")
			} else {
				tokenString = authHeader
			}

			if tokenString == "" {
				http.Error(w, `{"error":"empty token"}`, http.StatusUnauthorized)
				return
			}

			claims, err := Parse(secret, tokenString)
			if err != nil {
				http.Error(w, `{"error":"invalid or expired token"}`, http.StatusUnauthorized)
				return
			}

			ctx := WithClaims(r.Context(), claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func matchPath(pattern, path string) bool {
	if pattern == path {
		return true
	}
	if strings.HasSuffix(pattern, "/*") {
		prefix := pattern[:len(pattern)-1]
		return strings.HasPrefix(path, prefix)
	}
	return false
}
