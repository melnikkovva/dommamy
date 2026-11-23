package domain

import "errors"

var (
	ErrNotFound            = errors.New("event not found")
	ErrTitle               = errors.New("specify the title")
	ErrLimit               = errors.New("the limit cannot be negative")
	ErrNoMother            = errors.New("there is no such mother")
	ErrMotherNotFound      = errors.New("mother not found")
	ErrAlreadyJoined       = errors.New("the mother is already registered")
	ErrFull                = errors.New("no free places")
	ErrNotRegistered       = errors.New("the mother is not registered")
	ErrInvalidStatus       = errors.New("unknown status")
	ErrInvalidTransition   = errors.New("invalid status transition")
)
