package child

import "errors"

var ErrorName = errors.New("specify the date of birth")
var ErrorBirth = errors.New("the date of birth cannot be in the future")
var ErrorMother = errors.New("the mother is not specified")
var ErrorNoMother = errors.New("there is no such mother")
