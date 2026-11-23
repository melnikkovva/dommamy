package help

import (
	"time"
)

type Kind string

const (
	Food Kind = "food"
	Clothes Kind = "clothes"
	Furniture Kind = "furniture"
	Stationery Kind = "stationery"
	Gifts Kind = "gifts"
	Consultation Kind = "consultation"
)

func (h Kind) RequireCheck() bool {
	return h == Food
}

type Item struct {
	Id int
	MotherId int
	Kind Kind
	Sum int
	CheckUrl string
	CreatedAt time.Time
} 

func AddItem(motherId int, kind Kind, sum int, checkUrl string, now time.Time) (*Item, error) {	
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
		CreatedAt: now,
	}, nil
}