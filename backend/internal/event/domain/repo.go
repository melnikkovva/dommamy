package domain

import "context"

type Repository interface {
	Add(ctx context.Context, e Event) (id int, err error)
	ByID(ctx context.Context, id int) (Event, error)
	CountParticipants(ctx context.Context, eventID int) (int, error)
	AddParticipant(ctx context.Context, eventID int, motherID int) error
	HasParticipation(ctx context.Context, eventID int, motherID int) (bool, error)
	GetParticipation(ctx context.Context, eventID int, motherID int) (ParticipationStatus, error)
	SetParticipationStatus(ctx context.Context, eventID int, motherID int, status ParticipationStatus) error
}

type EventType struct {
	ID   int
	Name string
}
