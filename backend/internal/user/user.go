package user

import (
	"strings"
	"time"
)

const (
	RoleAdmin = "admin"
	RoleUser  = "user"
	RoleMom   = "mom"
)

type User struct {
	Id        int
	Login     string
	Password  string 
	FirstName string
	LastName  string
	Role      string
	MotherId  int 
	CreatedAt time.Time
}

func (u *User) Validate() error {
	u.Login = strings.TrimSpace(u.Login)
	u.FirstName = strings.TrimSpace(u.FirstName)
	u.LastName = strings.TrimSpace(u.LastName)
	if u.Login == "" {
		return ErrorLogin
	}
	if len(u.Password) < 6 {
		return ErrorPassword
	}
	if u.FirstName == "" {
		return ErrorName
	}
	if u.LastName == "" {
		return ErrorSurname
	}
	if u.Role != RoleAdmin && u.Role != RoleUser && u.Role != RoleMom {
		u.Role = RoleUser
	}
	return nil
}
