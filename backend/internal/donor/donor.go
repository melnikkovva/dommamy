package donor

import (
	"errors"
	"strings"
	"time"
)

var ErrorName = errors.New("specify the name of the donor")
var ErrorKind = errors.New("unknown kind of donor")
var ErrorEmail = errors.New("the email is incorrect")
var ErrorAmount = errors.New("the amount must be greater than zero")
var ErrorDate = errors.New("the date is incorrect")
var ErrorNotFound = errors.New("donor not found")

const (
	KindPerson  = "person"
	KindCompany = "company"
)

type Donor struct {
	Id int
	Name string
	Kind string
	Phone string
	Email string
	Comment string
	CreatedAt time.Time
}

type Summary struct {
	Donor
	Total int 
	Count int 
}

type Donation struct {
	Id int
	DonorId int
	Amount int
	Date string 
	Comment string
	CreatedAt time.Time
}

func (d *Donor) Validate() error {
	d.Name = strings.TrimSpace(d.Name)
	d.Phone = strings.TrimSpace(d.Phone)
	d.Email = strings.TrimSpace(d.Email)
	d.Comment = strings.TrimSpace(d.Comment)
	if d.Name == "" {
		return ErrorName
	}
	if d.Kind != KindPerson && d.Kind != KindCompany {
		return ErrorKind
	}
	if d.Email != "" && !strings.Contains(d.Email, "@") {
		return ErrorEmail
	}
	return nil
}

func NewDonation(donorId, amount int, date, comment string) (*Donation, error) {
	if amount <= 0 {
		return nil, ErrorAmount
	}
	today := time.Now().Format("2006-01-02")
	date = strings.TrimSpace(date)
	if date == "" {
		date = today
	}
	if _, err := time.Parse("2006-01-02", date); err != nil || date > today {
		return nil, ErrorDate
	}
	return &Donation{
		DonorId: donorId,
		Amount: amount,
		Date: date,
		Comment: strings.TrimSpace(comment),
		CreatedAt: time.Now(),
	}, nil
}