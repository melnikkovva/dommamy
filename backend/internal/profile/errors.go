package profile

import "errors"

var ErrorMother = errors.New("the mother is not specified")
var ErrorNoMother = errors.New("there is no such mother")
var ErrorBirth = errors.New("the date of birth is incorrect")
var ErrorIncome = errors.New("the income cannot be negative")
var ErrorWeeks = errors.New("pregnancy weeks must be from 0 to 42")
var ErrorStatus = errors.New("unknown status")
var ErrorTransition = errors.New("this status change is not allowed")
var ErrorNotFilled = errors.New("the profile is not filled yet")
var ErrorNotComplete = errors.New("fill the date of birth and the situation before sending")
