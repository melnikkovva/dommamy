package lifecycle

import (
	"errors"
	"strings"
	"time"
)

var ErrorMother = errors.New("the mother is not specified")
var ErrorNoMother = errors.New("there is no such mother")
var ErrorStatus = errors.New("unknown status")
var ErrorTransition = errors.New("this status change is not allowed")
var ErrorReason = errors.New("specify the reason")

const (
	StatusReview = "review"
	StatusActive = "active"
	StatusInactive = "inactive"
)

type Change struct {
	Id int
	MotherId int
	Status string
	Reason string
	CreatedAt time.Time
}

var transitions = map[string][]string{
	StatusReview: {StatusActive, StatusInactive},
	StatusActive: {StatusInactive},
	StatusInactive: {StatusReview},
}

func CanChange(from, to string) bool {
	for _, next := range transitions[from] {
		if next == to {
			return true
		}
	}
	return false
}

func ValidStatus(s string) bool {
	return s == StatusReview || s == StatusActive || s == StatusInactive
}

func NewChange(motherId int, status, reason string) (*Change, error) {
	if motherId <= 0 {
		return nil, ErrorMother
	}
	if !ValidStatus(status) {
		return nil, ErrorStatus
	}
	reason = strings.TrimSpace(reason)
	if status == StatusInactive && reason == "" {
		return nil, ErrorReason
	}
	return &Change{
		MotherId: motherId,
		Status: status,
		Reason: reason,
		CreatedAt: time.Now(),
	}, nil
}