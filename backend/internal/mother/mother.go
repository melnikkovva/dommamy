package mother

import (
	"errors"
	"time"
	"strings"
)

var ErrorName = errors.New("need name and surname")
var ErrorNotFound = errors.New("mother not found")

type Mother struct {
	Id int
	FirstName string
	LastName string
	Phone string
	CreatedAt time.Time
}

func CreateNew(firstName string, lastName string, phone string) (*Mother, error) {
	firstName = strings.TrimSpace(firstName)
	lastName = strings.TrimSpace(lastName)
	if firstName == "" || lastName == "" {
		return nil, ErrorName
	}
	return &Mother{
		FirstName: firstName,
		LastName: lastName,
		Phone: phone,
		CreatedAt: time.Now(),
	}, nil
}