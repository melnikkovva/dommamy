package help

import (
	"errors"
	"time"
)

type Kind string

const (
	Food Kind = "food"
	Clothes Kind = "clothes"
)

func (h Kind) RequireCheck() bool {
	return h == Food
}

var ErrorCheck = errors.New("Need check for food")
var ErrorSum = errors.New("Incorrect sum")
var ErrorMum = errors.New("Incorrect mum")
var ErrorKind = errors.New("Unknown kind of help")
var ErrorNotFound = errors.New("help not found")

type Item struct {
	Id int
	MotherId int
	Kind Kind
	Sum int
	CheckUrl string
	CreatedAt time.Time
} 

func AddItem(motherId int, kind Kind, sum int, checkUrl string) (*Item, error) {	
	if sum < 0 {
		return nil, ErrorSum
	}

	if kind.RequireCheck() {
		if sum == 0 {
			return nil, ErrorSum
		}
		if checkUrl == "" {
			return nil, ErrorSum
		}
	}

	return &Item {
		MotherId: motherId,
		Kind: kind,
		Sum: sum,
		CheckUrl: checkUrl,
		CreatedAt: time.Now(),
	}, nil
}