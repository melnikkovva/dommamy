package mother

import (
	"errors"
	"time"
	"strings"

	"backend/internal/learn"
)

var ErrorName = errors.New("need name and surname")
var ErrorNotFound = errors.New("mother not found")

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