package document

import "errors"

var ErrorMother = errors.New("the mother is not specified")
var ErrorNoMother = errors.New("there is no such mother")
var ErrorType = errors.New("unknown type of document")
var ErrorNeedHelp = errors.New("the statement of responsibility must be linked to a given help")
var ErrorHelpNotAllowed = errors.New("only the statement of responsibility can be linked to a help")
var ErrorNoHelp = errors.New("there is no such help")
var ErrorHelpMismatch = errors.New("this help was given to another mother")
var ErrorDuplicate = errors.New("this help already has a statement")
var ErrorNotFound = errors.New("document not found")
var ErrorSigned = errors.New("the document is already signed")
var ErrorNoFile = errors.New("upload the scan first")
var ErrorEmptyFile = errors.New("the file is empty")
var ErrorTooBig = errors.New("the file is too big (max 10 MB)")
var ErrorFileType = errors.New("only PDF, JPG, PNG, DOC, DOCX, XLS, XLSX, PPT, PPTX, ZIP, TXT, CSV files are allowed")
var ErrorSignDate = errors.New("the signing date is incorrect")
