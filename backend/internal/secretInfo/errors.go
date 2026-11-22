package secretInfo

import "errors"

var ErrorMother = errors.New("the mother is not specified")
var ErrorNoMother = errors.New("there is no such mother")
var ErrorNotFound = errors.New("secret info not filled")
