package domain

import (
	"strings"
	"time"
)

type Mother struct {
	id        int
	firstName string
	lastName  string
	phone     string
	email     string
	photoURL  string
	createdAt time.Time
}

func NewMother(firstName, lastName, phone, email string, now time.Time) (Mother, error) {
	firstName = strings.TrimSpace(firstName)
	lastName = strings.TrimSpace(lastName)
	email = strings.TrimSpace(strings.ToLower(email))
	if firstName == "" || lastName == "" {
		return Mother{}, ErrName
	}
	return Mother{
		firstName: firstName,
		lastName:  lastName,
		phone:     phone,
		email:     email,
		createdAt: now,
	}, nil
}

func (m *Mother) SetPhoto(photoURL string) {
	m.photoURL = photoURL
}

func (m Mother) ID() int              { return m.id }
func (m Mother) FirstName() string    { return m.firstName }
func (m Mother) LastName() string     { return m.lastName }
func (m Mother) Phone() string        { return m.phone }
func (m Mother) Email() string        { return m.email }
func (m Mother) PhotoURL() string     { return m.photoURL }
func (m Mother) CreatedAt() time.Time { return m.createdAt }

type MotherSnapshot struct {
	ID          int
	FirstName   string
	LastName    string
	Phone       string
	Email       string
	LearnUserID string
	PhotoURL    string
	CreatedAt   time.Time
}

func (m Mother) Snapshot() MotherSnapshot {
	return MotherSnapshot{
		ID:        m.id,
		FirstName: m.firstName,
		LastName:  m.lastName,
		Phone:     m.phone,
		Email:     m.email,
		PhotoURL:  m.photoURL,
		CreatedAt: m.createdAt,
	}
}

func (s MotherSnapshot) Restore() Mother {
	return Mother{
		id:        s.ID,
		firstName: s.FirstName,
		lastName:  s.LastName,
		phone:     s.Phone,
		email:     s.Email,
		photoURL:  s.PhotoURL,
	}
}
