package repository

import "time"

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
