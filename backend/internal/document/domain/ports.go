package domain

import "io"

type FileStorage interface {
	Put(name string, content io.Reader) error
	Open(name string) (io.ReadCloser, error)
	Delete(name string) error
}
