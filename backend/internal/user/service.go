package user

import (
	"errors"
	"time"

	"backend/internal/auth"
	"golang.org/x/crypto/bcrypt"
)

const fakeHash = "$2a$10$njfjPpPf/ElcAXhzlE1Use8QDUG7m6urKLk9aCM1pX1Mxm3UJsZwG"

type Repository interface {
	Save(u *User) error
	UserById(id int) (*User, error)
	UserByLogin(login string) (*User, error)
	ListUsers() ([]*User, error)
	DeleteUser(id int) error
}

type Service struct {
	repo     Repository
	secret   []byte
	tokenTTL time.Duration
}

func NewService(repo Repository, secret string, tokenTTL time.Duration) *Service {
	return &Service{
		repo:     repo,
		secret:   []byte(secret),
		tokenTTL: tokenTTL,
	}
}

func (s *Service) Create(login, password, firstName, lastName, role string, motherId int) (*User, error) {
	if len(password) < 6 {
		return nil, ErrorPassword
	}

	u := &User{
		Login:     login,
		Password:  password,
		FirstName: firstName,
		LastName:  lastName,
		Role:      role,
		MotherId:  motherId,
		CreatedAt: time.Now(),
	}
	if err := u.Validate(); err != nil {
		return nil, err
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	u.Password = string(hashed)

	if err := s.repo.Save(u); err != nil {
		return nil, err
	}
	return u, nil
}

func (s *Service) Login(login, password string) (string, *User, error) {
	u, err := s.repo.UserByLogin(login)
	if err != nil {
		if errors.Is(err, ErrorNotFound) {
			_ = bcrypt.CompareHashAndPassword([]byte(fakeHash), []byte(password))
			return "", nil, ErrorInvalidCredentials
		}
		return "", nil, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return "", nil, ErrorInvalidCredentials
	}
	token, err := s.generateToken(u)
	if err != nil {
		return "", nil, err
	}
	return token, u, nil
}

func (s *Service) GetUser(id int) (*User, error) {
	return s.repo.UserById(id)
}

func (s *Service) ListUsers() ([]*User, error) {
	return s.repo.ListUsers()
}

func (s *Service) DeleteUser(id int) error {
	return s.repo.DeleteUser(id)
}

func (s *Service) generateToken(u *User) (string, error) {
	return auth.Generate(s.secret, u.Id, u.Role, u.MotherId, s.tokenTTL)
}
