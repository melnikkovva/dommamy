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
	learnUserID string
	photoURL  string
	createdAt time.Time
}

func NewMother(firstName, lastName, phone, email string) (Mother, error) {
	firstName = strings.TrimSpace(firstName)
	lastName = strings.TrimSpace(lastName)
	if firstName == "" || lastName == "" {
		return Mother{}, ErrName
	}
	return Mother{
		firstName: firstName,
		lastName:  lastName,
		phone:     phone,
		email:     email,
		createdAt: time.Now(),
	}, nil
}

func (m *Mother) SetID(id int) {
	m.id = id
}

func (m *Mother) SetLearnUserID(learnUserID string) {
	m.learnUserID = learnUserID
}

func (m *Mother) SetPhotoURL(photoURL string) {
	m.photoURL = photoURL
}

func (m Mother) ID() int              { return m.id }
func (m Mother) FirstName() string    { return m.firstName }
func (m Mother) LastName() string     { return m.lastName }
func (m Mother) Phone() string        { return m.phone }
func (m Mother) Email() string        { return m.email }
func (m Mother) LearnUserID() string  { return m.learnUserID }
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
		ID:          m.id,
		FirstName:   m.firstName,
		LastName:    m.lastName,
		Phone:       m.phone,
		Email:       m.email,
		LearnUserID: m.learnUserID,
		PhotoURL:    m.photoURL,
		CreatedAt:   m.createdAt,
	}
}

func restoreMother(s MotherSnapshot) Mother {
	return Mother{
		id:          s.ID,
		firstName:   s.FirstName,
		lastName:    s.LastName,
		phone:       s.Phone,
		email:       s.Email,
		learnUserID: s.LearnUserID,
		photoURL:    s.PhotoURL,
		createdAt:   s.CreatedAt,
	}
}
