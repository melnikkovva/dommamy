package storage

import (
	"errors"
	"io"
	"os"
	"path/filepath"
)

type Local struct {
	dir string
}

func NewLocal(dir string) (*Local, error) {
	if err := os.MkdirAll(dir, 0o750); err != nil {
		return nil, err
	}
	return &Local{dir: dir}, nil
}

func (l *Local) path(name string) (string, error) {
	if name == "" || filepath.Base(name) != name {
		return "", errors.New("bad file name")
	}
	return filepath.Join(l.dir, name), nil
}

func (l *Local) Put(name string, content io.Reader) error {
	p, err := l.path(name)
	if err != nil {
		return err
	}
	f, err := os.OpenFile(p, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
	if err != nil {
		return err
	}
	if _, err := io.Copy(f, content); err != nil {
		f.Close()
		os.Remove(p) 
		return err
	}
	return f.Close()
}

func (l *Local) Open(name string) (io.ReadCloser, error) {
	p, err := l.path(name)
	if err != nil {
		return nil, err
	}
	return os.Open(p)
}

func (l *Local) Delete(name string) error {
	p, err := l.path(name)
	if err != nil {
		return err
	}
	err = os.Remove(p)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return err
}