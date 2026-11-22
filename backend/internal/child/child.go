package child

import (
	"strings"
	"time"
)

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