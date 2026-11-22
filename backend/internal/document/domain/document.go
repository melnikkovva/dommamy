package domain

import (
	"crypto/rand"
	"encoding/hex"
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

var extensions = map[string]string{
	"application/pdf": ".pdf",
	"image/jpeg":      ".jpg",
	"image/png":       ".png",
	"application/msword": ".doc",
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document": ".docx",
	"application/vnd.ms-excel": ".xls",
	"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet": ".xlsx",
	"application/vnd.ms-powerpoint": ".ppt",
	"application/vnd.openxmlformats-officedocument.presentationml.presentation": ".pptx",
	"application/zip": ".zip",
	"application/x-zip-compressed": ".zip",
	"text/plain": ".txt",
	"text/csv":   ".csv",
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

func NewDocument(motherID int, docType Type, helpID int) (Document, error) {
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
		createdAt: time.Now(),
	}, nil
}

func (d *Document) SetID(id int) {
	d.id = id
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
	if _, ok := extensions[contentType]; !ok {
		return ErrFileType
	}
	d.fileName = storedName
	d.originalName = cleanName(originalName)
	d.contentType = contentType
	d.size = size
	return nil
}

func (d *Document) Sign(date string) error {
	if d.status == StatusSigned {
		return ErrSigned
	}
	if !d.HasFile() {
		return ErrNoFile
	}
	date = strings.TrimSpace(date)
	today := time.Now().Format("2006-01-02")
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

func NewStoredName(contentType string) (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	ext := extensions[contentType]
	if ext == "" {
		ext = ".bin"
	}
	return hex.EncodeToString(b) + ext, nil
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
