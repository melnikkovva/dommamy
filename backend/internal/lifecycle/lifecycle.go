package lifecycle

import (
	"strings"
	"time"
)

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