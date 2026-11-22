package domain

import "errors"

var (
	ErrMother     = errors.New("the mother is not specified")
	ErrNoMother   = errors.New("there is no such mother")
	ErrStatus     = errors.New("unknown status")
	ErrTransition = errors.New("this status change is not allowed")
	ErrReason     = errors.New("specify the reason")
)
