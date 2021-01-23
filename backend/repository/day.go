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

func (d *Day) PrepareReceived() {
	d.ID = ""
	d.UserID = ""
}

func (d *Day) PrepareToSend() {
	d.UserID = ""
}

type Event struct {
	ID        string    `json:"id,omitempty" bson:"_id"`
	Name      string    `json:"name" validate:"required,max=256"`
	Start     time.Time `json:"start" validate:"required"`
	End       time.Time `json:"end" validate:"required"`
	TagID     string    `json:"tag" validate:"required,max=36"`
	ProjectID string    `json:"project" validate:"max=36"`
	Timed     bool      `json:"timed" validate:"required"`
}

func (e *Event) PrepareReceived() {
	e.ID = ""
}

func (e *Event) PrepareToSend() {
	// Event doesn't have anything to clean before sending
}

func (e *Event) Copy() Event {
	return Event{
		ID:        e.ID,
		Name:      e.Name,
		Start:     e.Start,
		End:       e.End,
		TagID:     e.TagID,
		ProjectID: e.ProjectID,
		Timed:     e.Timed,
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
