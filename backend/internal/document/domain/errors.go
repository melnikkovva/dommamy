package domain

import "errors"

var (
	ErrMother           = errors.New("the mother is not specified")
	ErrNoMother         = errors.New("there is no such mother")
	ErrType             = errors.New("unknown type of document")
	ErrNeedHelp         = errors.New("the statement of responsibility must be linked to a given help")
	ErrHelpNotAllowed   = errors.New("only the statement of responsibility can be linked to a help")
	ErrNoHelp           = errors.New("there is no such help")
	ErrHelpMismatch     = errors.New("this help was given to another mother")
	ErrDuplicate        = errors.New("this help already has a statement")
	ErrNotFound         = errors.New("document not found")
	ErrSigned           = errors.New("the document is already signed")
	ErrNoFile           = errors.New("upload the scan first")
	ErrEmptyFile        = errors.New("the file is empty")
	ErrTooBig           = errors.New("the file is too big (max 10 MB)")
	ErrFileType         = errors.New("only PDF, JPG, PNG, DOC, DOCX, XLS, XLSX, PPT, PPTX, ZIP, TXT, CSV files are allowed")
	ErrSignDate         = errors.New("the signing date is incorrect")
	ErrFileStorage      = errors.New("file storage error")
)
