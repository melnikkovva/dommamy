package storage

import (
	"crypto/rand"
	"encoding/hex"

	"backend/internal/document/domain"
)

func GenerateStoredName(contentType string) (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	ext := domain.GetExtensionByMimeType(contentType)
	return hex.EncodeToString(b) + ext, nil
}
