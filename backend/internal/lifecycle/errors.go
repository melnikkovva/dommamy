package lifecycle

import "errors"

var ErrorMother = errors.New("the mother is not specified")
var ErrorNoMother = errors.New("there is no such mother")
var ErrorStatus = errors.New("unknown status")
var ErrorTransition = errors.New("this status change is not allowed")
var ErrorReason = errors.New("specify the reason")
