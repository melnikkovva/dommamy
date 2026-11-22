package user

import "errors"

var ErrorLogin = errors.New("specify a login")
var ErrorPassword = errors.New("password must be at least 6 characters")
var ErrorName = errors.New("specify a first name")
var ErrorSurname = errors.New("specify a last name")
var ErrorNotFound = errors.New("user not found")
var ErrorInvalidCredentials = errors.New("invalid login or password")
var ErrorLoginTaken = errors.New("login is already taken")
