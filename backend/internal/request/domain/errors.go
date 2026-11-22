package domain

import "errors"

var (
	ErrNotFound = errors.New("request not found")
	ErrMotherNotSpecified = errors.New("the mother is not specified")
	ErrMotherNotFound = errors.New("there is no such mother")
	ErrInvalidKind = errors.New("unknown kind of help")
	ErrInvalidSource = errors.New("unknown source")
	ErrInvalidStatus = errors.New("unknown status")
	ErrTransition = errors.New("this status change is not allowed")
	ErrCommentRequired = errors.New("specify the comment for rejection")
	ErrUseFulfill = errors.New("a request becomes issued only after linking the given help")
	ErrHelpNotFound = errors.New("there is no such help")
	ErrHelpMismatch = errors.New("this help was given to another mother")
)