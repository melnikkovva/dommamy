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

func AddNew(motherId int, firstName string, birthDate time.Time, now time.Time) (*Child, error) {
	firstName = strings.TrimSpace(firstName)
	if firstName == "" {
		return nil, ErrorName
	}
	if birthDate.After(now) {
		return nil, ErrorBirth
	}
	return &Child{
		MotherId: motherId,
		FirstName: firstName,
		BirthDate: birthDate,
		CreatedAt: now,
	}, nil
}