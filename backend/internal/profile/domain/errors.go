package domain

import "errors"

var (
	ErrMother      = errors.New("the mother is not specified")
	ErrNoMother    = errors.New("there is no such mother")
	ErrBirth       = errors.New("the date of birth is incorrect")
	ErrIncome      = errors.New("the income cannot be negative")
	ErrWeeks       = errors.New("pregnancy weeks must be from 0 to 42")
	ErrStatus      = errors.New("unknown status")
	ErrTransition  = errors.New("this status change is not allowed")
	ErrNotFilled   = errors.New("the profile is not filled yet")
	ErrNotComplete = errors.New("fill the date of birth and the situation before sending")
)
