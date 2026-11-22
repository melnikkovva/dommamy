package mother

import "errors"

var ErrorName = errors.New("need name and surname")
var ErrorNotFound = errors.New("mother not found")