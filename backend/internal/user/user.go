package user

import (
	"errors"
	"strings"
	"time"
)

var ErrorLogin = errors.New("specify a login")
var ErrorPassword = errors.New("password must be at least 6 characters")
var ErrorName = errors.New("specify a first name")
var ErrorSurname = errors.New("specify a last name")
var ErrorNotFound = errors.New("user not found")
var ErrorInvalidCredentials = errors.New("invalid login or password")
var ErrorLoginTaken = errors.New("login is already taken")

const (
	RoleAdmin = "admin"
	RoleUser  = "user"
)

type User struct {
	Id        int
	Login     string
	Password  string 
	FirstName string
	LastName  string
	Role      string
	CreatedAt time.Time
}

func (u *User) Validate() error {
	u.Login = strings.TrimSpace(u.Login)
	u.FirstName = strings.TrimSpace(u.FirstName)
	u.LastName = strings.TrimSpace(u.LastName)
	if u.Login == "" {
		return ErrorLogin
	}
	if u.Password == "" {
		return ErrorPassword
	}
	if u.FirstName == "" {
		return ErrorName
	}
	if u.LastName == "" {
		return ErrorSurname
	}
	if u.Role != RoleAdmin && u.Role != RoleUser {
		u.Role = RoleUser
	}
	return nil
}
