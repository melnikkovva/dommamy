package document

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"path"
	"strings"
	"time"
)

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

const (
	TypeQuestionnaire  = "questionnaire"
	TypeConsent = "consent"
	TypeResponsibility = "responsibility"
	TypeOther = "other"
)

const (
	StatusNotSigned = "not_signed"
	StatusSigned = "signed"
)

const MaxFileSize = 10 << 20 

var extensions = map[string]string{
	"application/pdf": ".pdf",
	"image/jpeg": ".jpg",
	"image/png": ".png",
	"application/msword": ".doc",
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document": ".docx",
	"application/vnd.ms-excel": ".xls",
	"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet": ".xlsx",
	"application/vnd.ms-powerpoint": ".ppt",
	"application/vnd.openxmlformats-officedocument.presentationml.presentation": ".pptx",
	"application/zip": ".zip",
	"application/x-zip-compressed": ".zip",
	"text/plain": ".txt",
	"text/csv": ".csv",
}

type Document struct {
	Id int
	MotherId int
	Type string
	HelpId int 
	Status string
	SignedAt string
	FileName string 
	OriginalName string 
	ContentType string
	Size int64
	CreatedAt time.Time
}

func ValidType(t string) bool {
	return t == TypeQuestionnaire || t == TypeConsent || t == TypeResponsibility || t == TypeOther
}

func AddNew(motherId int, docType string, helpId int) (*Document, error) {
	if motherId <= 0 {
		return nil, ErrorMother
	}
	if !ValidType(docType) {
		return nil, ErrorType
	}
	if docType == TypeResponsibility && helpId <= 0 {
		return nil, ErrorNeedHelp
	}
	if docType != TypeResponsibility && helpId != 0 {
		return nil, ErrorHelpNotAllowed
	}
	return &Document{
		MotherId: motherId,
		Type: docType,
		HelpId: helpId,
		Status: StatusNotSigned,
		CreatedAt: time.Now(),
	}, nil
}

func (d *Document) HasFile() bool {
	return d.FileName != ""
}

func (d *Document) AttachFile(storedName, originalName, contentType string, size int64) error {
	if d.Status == StatusSigned {
		return ErrorSigned
	}
	if size <= 0 {
		return ErrorEmptyFile
	}
	if size > MaxFileSize {
		return ErrorTooBig
	}
	if _, ok := extensions[contentType]; !ok {
		return ErrorFileType
	}
	d.FileName = storedName
	d.OriginalName = cleanName(originalName)
	d.ContentType = contentType
	d.Size = size
	return nil
}

func (d *Document) Sign(date string) error {
	if d.Status == StatusSigned {
		return ErrorSigned
	}
	if !d.HasFile() {
		return ErrorNoFile
	}
	date = strings.TrimSpace(date)
	today := time.Now().Format("2006-01-02")
	if date == "" {
		date = today
	}
	if _, err := time.Parse("2006-01-02", date); err != nil || date > today {
		return ErrorSignDate
	}
	d.Status = StatusSigned
	d.SignedAt = date
	return nil
}

func NewStoredName(contentType string) (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b) + extensions[contentType], nil
}

func cleanName(name string) string {
	name = path.Base(strings.ReplaceAll(name, "\\", "/"))
	if name == "." || name == "/" || name == "" {
		return "file"
	}
	runes := []rune(name)
	if len(runes) > 200 {
		runes = runes[:200]
	}
	return string(runes)
}