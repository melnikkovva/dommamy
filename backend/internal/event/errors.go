package event

import "errors"

var ErrorTitle = errors.New("specify the title")
var ErrorLimit = errors.New("the limit cannot be negative")
var ErrorNotFound = errors.New("event not found")
var ErrorNoMother = errors.New("there is no such mother")
var ErrorAlreadyJoined = errors.New("the mother is already registered")
var ErrorFull = errors.New("no free places")
var ErrorNotRegistered = errors.New("the mother is not registered")
var ErrorStatus = errors.New("unknown status")
