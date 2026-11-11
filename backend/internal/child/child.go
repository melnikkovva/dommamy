package child

import (
	"errors"
	"strings"
	"time"
)

var ErrorName = errors.New("specify the date of birth")
var ErrorBirth = errors.New("the date of birth cannot be in the future")
var ErrorMother = errors.New("the mother is not specified")
var ErrorNoMother = errors.New("there is no such mother")

type Child struct {
	Id int
	MotherId int
	FirstName string
	BirthDate time.Time
	CreatedAt time.Time
}

func AddNew(motherId int, firstName string, birthDate time.Time) (*Child, error) {
	firstName = strings.TrimSpace(firstName)
	if firstName == "" {
		return nil, ErrorName
	}
	if birthDate.After(time.Now()) {
		return nil, ErrorBirth
	}
	return &Child{
		MotherId: motherId,
		FirstName: firstName,
		BirthDate: birthDate,
		CreatedAt: time.Now(),
	}, nil
}