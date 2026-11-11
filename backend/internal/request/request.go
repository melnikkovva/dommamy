package request

import (
	"errors"
	"strings"
	"time"
)

var ErrorMother = errors.New("the mother is not specified")
var ErrorNoMother = errors.New("there is no such mother")
var ErrorKind = errors.New("unknown kind of help")
var ErrorSource = errors.New("unknown source")
var ErrorNotFound = errors.New("request not found")
var ErrorStatus = errors.New("unknown status")
var ErrorTransition = errors.New("this status change is not allowed")
var ErrorComment = errors.New("specify the comment for rejection")
var ErrorUseIssue = errors.New("a request becomes issued only after linking the given help")
var ErrorNoHelp = errors.New("there is no such help")
var ErrorHelpMismatch = errors.New("this help was given to another mother")

const (
	StatusNew = "new"
	StatusInWork = "in_work"
	StatusIssued = "issued"
	StatusRejected = "rejected"
)

const (
	SourceStaff  = "staff"
	SourceMother = "mother" 
)

type Request struct {
	Id int
	MotherId int
	Kind string
	Description string
	Status string
	Comment string
	Source string
	HelpId int 
	CreatedAt time.Time
	UpdatedAt time.Time
}

var transitions = map[string][]string{
	StatusNew: {StatusInWork, StatusRejected},
	StatusInWork: {StatusRejected},
	StatusRejected: {StatusNew},
}

func CanChange(from, to string) bool {
	for _, next := range transitions[from] {
		if next == to {
			return true
		}
	}
	return false
}

func ValidStatus(s string) bool {
	return s == StatusNew || s == StatusInWork || s == StatusIssued || s == StatusRejected
}

func ValidKind(k string) bool {
	return k == "food" || k == "clothes" || k == "furniture" || k == "other"
}

func AddNew(motherId int, kind, description, source string) (*Request, error) {
	if motherId <= 0 {
		return nil, ErrorMother
	}
	if !ValidKind(kind) {
		return nil, ErrorKind
	}
	if source != SourceStaff && source != SourceMother {
		return nil, ErrorSource
	}
	now := time.Now()
	return &Request{
		MotherId: motherId,
		Kind: kind,
		Description: strings.TrimSpace(description),
		Status: StatusNew,
		Source: source,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func (r *Request) ChangeStatus(to, comment string) error {
	if to == StatusIssued {
		return ErrorUseIssue
	}
	if !ValidStatus(to) {
		return ErrorStatus
	}
	if !CanChange(r.Status, to) {
		return ErrorTransition
	}
	comment = strings.TrimSpace(comment)
	if to == StatusRejected && comment == "" {
		return ErrorComment
	}
	r.Status = to
	if comment != "" {
		r.Comment = comment
	}
	r.UpdatedAt = time.Now()
	return nil
}

func (r *Request) Fulfill(helpId int) error {
	if helpId <= 0 {
		return ErrorNoHelp
	}
	if r.Status != StatusNew && r.Status != StatusInWork {
		return ErrorTransition
	}
	r.Status = StatusIssued
	r.HelpId = helpId
	r.UpdatedAt = time.Now()
	return nil
}