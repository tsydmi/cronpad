package repository

import (
	"errors"
	"time"
)

type Day struct {
	ID     string    `json:"id,omitempty" bson:"_id"`
	Events []Event   `json:"events"`
	Date   time.Time `json:"date"`
	UserID string    `json:"user,omitempty"`
}

func (t *Day) PrepareReceived() {
	t.ID = ""
	t.UserID = ""
}

func (t *Day) PrepareToSend() {
	t.UserID = ""
}

type Event struct {
	ID        string    `json:"id,omitempty" bson:"_id"`
	Name      string    `json:"name" validate:"required"`
	Start     time.Time `json:"start" validate:"required"`
	End       time.Time `json:"end" validate:"required"`
	TagID     string    `json:"tag" validate:"required"`
	ProjectID string    `json:"project"`
	Timed     bool      `json:"timed" validate:"required"`
}

func (t *Event) PrepareReceived() {
	t.ID = ""
}

func (t *Event) PrepareToSend() {
	// Event doesn't have anything to clean before sending
}

func (t *Event) Copy() Event {
	return Event{
		ID:        t.ID,
		Name:      t.Name,
		Start:     t.Start,
		End:       t.End,
		TagID:     t.TagID,
		ProjectID: t.ProjectID,
		Timed:     t.Timed,
	}
}

func (e *Event) IsValid() []error {
	var errs []error

	if !e.Start.IsZero() && !e.End.IsZero() {
		if e.Start.After(e.End) {
			errs = append(errs, errors.New("wrong date rage"))
		}
		if e.Start.Day() != e.End.Day() || e.Start.Month() != e.End.Month() || e.Start.Year() != e.End.Year() {
			errs = append(errs, errors.New("event cannot last more than one day"))
		}
	}

	return errs
}
