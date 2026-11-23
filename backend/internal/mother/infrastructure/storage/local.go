package storage

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type LocalStorage struct {
	dir string
}

func NewLocalStorage(dir string) *LocalStorage {
	if err := os.MkdirAll(dir, 0755); err != nil {
		panic(fmt.Sprintf("photo storage: cannot create directory %s: %v", dir, err))
	}
	return &LocalStorage{dir: dir}
}

func (s *LocalStorage) Upload(_ context.Context, data []byte, _ string) (string, error) {
	ext := ".jpg"
	randomBytes := make([]byte, 16)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", err
	}
	fileName := hex.EncodeToString(randomBytes) + ext

	dest := filepath.Join(s.dir, fileName)
	f, err := os.OpenFile(dest, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return "", err
	}
	defer f.Close()

	if _, err := io.Copy(f, io.LimitReader(&fixedReader{data}, int64(len(data)))); err != nil {
		os.Remove(dest)
		return "", err
	}

	return "/api/photos/" + fileName, nil
}

type fixedReader struct{ data []byte }

func (r *fixedReader) Read(p []byte) (int, error) {
	if len(r.data) == 0 {
		return 0, io.EOF
	}
	n := copy(p, r.data)
	r.data = r.data[n:]
	return n, nil
}
