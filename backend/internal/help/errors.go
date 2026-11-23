package help

import "errors"

var ErrorCheck = errors.New("need check for food")
var ErrorSum = errors.New("incorrect sum")
var ErrorMum = errors.New("incorrect mum")
var ErrorKind = errors.New("unknown kind of help")
var ErrorNotFound = errors.New("help not found")
