package donor

import "errors"

var ErrorName = errors.New("specify the name of the donor")
var ErrorKind = errors.New("unknown kind of donor")
var ErrorEmail = errors.New("the email is incorrect")
var ErrorAmount = errors.New("the amount must be greater than zero")
var ErrorDate = errors.New("the date is incorrect")
var ErrorNotFound = errors.New("donor not found")
