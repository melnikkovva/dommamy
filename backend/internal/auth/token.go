package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrorInvalidToken = errors.New("invalid or expired token")
	ErrorNoToken      = errors.New("missing authorization token")
)

type Claims struct {
	UserID   int    `json:"uid"`
	Role     string `json:"role"`
	MotherID int    `json:"mid,omitempty"`
	jwt.RegisteredClaims
}

type ContextKey string

const ClaimsKey ContextKey = "auth_claims"

func Generate(secret []byte, userID int, role string, motherID int, ttl time.Duration) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID:   userID,
		Role:     role,
		MotherID: motherID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(now),
			ID:        randomID(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

func Parse(secret []byte, tokenString string) (Claims, error) {
	var claims Claims
	_, err := jwt.ParseWithClaims(tokenString, &claims,
		func(token *jwt.Token) (any, error) {
			return secret, nil
		},
		jwt.WithValidMethods([]string{"HS256"}),
	)
	if err != nil {
		return claims, ErrorInvalidToken
	}
	return claims, nil
}

func FromContext(ctx context.Context) (Claims, error) {
	v := ctx.Value(ClaimsKey)
	c, ok := v.(Claims)
	if !ok {
		return Claims{}, ErrorInvalidToken
	}
	return c, nil
}

func WithClaims(ctx context.Context, c Claims) context.Context {
	return context.WithValue(ctx, ClaimsKey, c)
}

func randomID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "unknown"
	}
	return hex.EncodeToString(b)
}
