package event

import (
	"strings"
	"time"
)

type ParticipationStatus string

const (
	StatusRegistered ParticipationStatus = "registered"
	StatusCame ParticipationStatus = "came"
	StatusAbsent ParticipationStatus = "absent"
)

type Type struct {
	Id int
	Name string
}

type Event struct {
	Id int
	TypeId int
	Title string
	StartsAt time.Time
	Place string
	MaxPeople int
}

type Participation struct {
	EventId  int
	MotherId int
	Status ParticipationStatus
}

func AddNew(typeId int, title string, startsAt time.Time, place string, maxPeople int) (*Event, error){
	title = strings.TrimSpace(title)
	if title == "" {
		return nil, ErrorTitle
	}
	if maxPeople < 0 {
		return nil, ErrorLimit
	}
	return &Event{
		TypeId: typeId, 
		Title: title, 
		StartsAt: startsAt,
		Place: strings.TrimSpace(place), 
		MaxPeople: maxPeople,
	}, nil
}

