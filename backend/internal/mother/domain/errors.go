package domain

import "errors"

var (
	ErrName     = errors.New("need name and surname")
	ErrNotFound = errors.New("mother not found")
)
