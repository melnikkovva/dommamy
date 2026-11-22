package user

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Repository interface {
	Save(u *User) error
	UserById(id int) (*User, error)
	UserByLogin(login string) (*User, error)
	ListUsers() ([]*User, error)
}

type Service struct {
	repo    Repository
	secret  []byte
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
	if _, err := s.repo.UserByLogin(login); err == nil {
		return nil, ErrorLoginTaken
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	u := &User{
		Login:     login,
		Password:  string(hashed),
		FirstName: firstName,
		LastName:  lastName,
		Role:      role,
		MotherId:  motherId,
		CreatedAt: time.Now(),
	}
	if err := u.Validate(); err != nil {
		return nil, err
	}
	if err := s.repo.Save(u); err != nil {
		return nil, err
	}
	return u, nil
}

func (s *Service) Login(login, password string) (string, *User, error) {
	u, err := s.repo.UserByLogin(login)
	if err != nil {
		return "", nil, ErrorInvalidCredentials
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

func (s *Service) generateToken(u *User) (string, error) {
	randomBytes := make([]byte, 16)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", err
	}
	expire := time.Now().Add(s.tokenTTL).Unix()
	randomHex := hex.EncodeToString(randomBytes)

	payload := string(u.Id) + "-" + string(rune(expire%100000)) + "-" + randomHex
	return payload, nil
}
