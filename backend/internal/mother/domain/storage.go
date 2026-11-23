package domain

import "context"

type PhotoStorage interface {
	Upload(ctx context.Context, data []byte, contentType string) (string, error)
}
