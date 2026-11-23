package domain

import (
	"path"
	"strings"
	"time"
)

type Type string

const (
	TypeQuestionnaire  Type = "questionnaire"
	TypeConsent        Type = "consent"
	TypeResponsibility Type = "responsibility"
	TypeOther          Type = "other"
)

var validTypes = []Type{TypeQuestionnaire, TypeConsent, TypeResponsibility, TypeOther}

func (t Type) Valid() bool {
	for _, v := range validTypes {
		if v == t {
			return true
		}
	}
	return false
}

type Status string

const (
	StatusNotSigned Status = "not_signed"
	StatusSigned    Status = "signed"
)

const MaxFileSize = 10 << 20 // 10 MB

var allowedExtensions = map[string]string{
	".pdf":  "application/pdf",
	".jpg":  "image/jpeg",
	".jpeg": "image/jpeg",
	".png":  "image/png",
	".doc":  "application/msword",
	".docx": "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
	".xls":  "application/vnd.ms-excel",
	".xlsx": "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
	".ppt":  "application/vnd.ms-powerpoint",
	".pptx": "application/vnd.openxmlformats-officedocument.presentationml.presentation",
	".zip":  "application/zip",
	".txt":  "text/plain",
	".csv":  "text/csv",
}

var deniedContentTypes = []string{
	"text/html",
	"application/x-executable",
	"application/x-msdownload",
	"application/octet-stream",
}

func AllowedExtensions() map[string]string {
	return allowedExtensions
}

func DeniedContentTypes() []string {
	return deniedContentTypes
}

type Document struct {
	id           int
	motherID     int
	docType      Type
	helpID       int
	status       Status
	signedAt     string
	fileName     string
	originalName string
	contentType  string
	size         int64
	createdAt    time.Time
}

func NewDocument(motherID int, docType Type, helpID int, now time.Time) (Document, error) {
	if motherID <= 0 {
		return Document{}, ErrMother
	}
	if !docType.Valid() {
		return Document{}, ErrType
	}
	if docType == TypeResponsibility && helpID <= 0 {
		return Document{}, ErrNeedHelp
	}
	if docType != TypeResponsibility && helpID != 0 {
		return Document{}, ErrHelpNotAllowed
	}
	return Document{
		motherID:  motherID,
		docType:   docType,
		helpID:    helpID,
		status:    StatusNotSigned,
		createdAt: now,
	}, nil
}


func (d Document) ID() int           { return d.id }
func (d Document) MotherID() int     { return d.motherID }
func (d Document) Type() Type        { return d.docType }
func (d Document) HelpID() int       { return d.helpID }
func (d Document) Status() Status    { return d.status }
func (d Document) SignedAt() string  { return d.signedAt }
func (d Document) FileName() string  { return d.fileName }
func (d Document) OriginalName() string { return d.originalName }
func (d Document) ContentType() string { return d.contentType }
func (d Document) Size() int64       { return d.size }
func (d Document) CreatedAt() time.Time { return d.createdAt }

func (d Document) HasFile() bool {
	return d.fileName != ""
}

func (d *Document) AttachFile(storedName, originalName, contentType string, size int64) error {
	if d.status == StatusSigned {
		return ErrSigned
	}
	if size <= 0 {
		return ErrEmptyFile
	}
	if size > MaxFileSize {
		return ErrTooBig
	}
	ext := strings.ToLower(path.Ext(cleanName(originalName)))
	if mimeType, ok := allowedExtensions[ext]; ok {
		d.contentType = mimeType
	} else {
		return ErrFileType
	}
	d.fileName = storedName
	d.originalName = cleanName(originalName)
	d.size = size
	return nil
}

func (d *Document) Sign(date, today string) error {
	if d.status == StatusSigned {
		return ErrSigned
	}
	if !d.HasFile() {
		return ErrNoFile
	}
	date = strings.TrimSpace(date)
	if date == "" {
		date = today
	}
	if _, err := time.Parse("2006-01-02", date); err != nil || date > today {
		return ErrSignDate
	}
	d.status = StatusSigned
	d.signedAt = date
	return nil
}

func GetExtensionByMimeType(mimeType string) string {
	for ext, mt := range allowedExtensions {
		if mt == mimeType {
			return ext
		}
	}
	return ".bin"
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

type Snapshot struct {
	ID           int
	MotherID     int
	DocType      Type
	HelpID       int
	Status       Status
	SignedAt     string
	FileName     string
	OriginalName string
	ContentType  string
	Size         int64
	CreatedAt    time.Time
}

func (d Document) Snapshot() Snapshot {
	return Snapshot{
		ID:           d.id,
		MotherID:     d.motherID,
		DocType:      d.docType,
		HelpID:       d.helpID,
		Status:       d.status,
		SignedAt:     d.signedAt,
		FileName:     d.fileName,
		OriginalName: d.originalName,
		ContentType:  d.contentType,
		Size:         d.size,
		CreatedAt:    d.createdAt,
	}
}

func Restore(s Snapshot) Document {
	return Document{
		id:           s.ID,
		motherID:     s.MotherID,
		docType:      s.DocType,
		helpID:       s.HelpID,
		status:       s.Status,
		signedAt:     s.SignedAt,
		fileName:     s.FileName,
		originalName: s.OriginalName,
		contentType:  s.ContentType,
		size:         s.Size,
		createdAt:    s.CreatedAt,
	}
}
