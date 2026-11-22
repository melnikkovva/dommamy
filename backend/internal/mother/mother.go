package mother

import (
	"time"
	"strings"

	"backend/internal/learn"
)

type Mother struct {
	Id int
	FirstName string
	LastName string
	Phone string
	Email string
	LearnUserID string
	PhotoURL string
	CreatedAt time.Time
	CourseStats *learn.CourseStats
}

func CreateNew(firstName string, lastName string, phone string, email string) (*Mother, error) {
	firstName = strings.TrimSpace(firstName)
	lastName = strings.TrimSpace(lastName)
	if firstName == "" || lastName == "" {
		return nil, ErrorName
	}
	return &Mother{
		FirstName: firstName,
		LastName: lastName,
		Phone: phone,
		Email: email,
		CreatedAt: time.Now(),
	}, nil
}