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
	Name      string    `json:"name"`
	Start     time.Time `json:"start"`
	End       time.Time `json:"end"`
	TagID     string    `json:"tag"`
	ProjectID string    `json:"project"`
	Timed     bool      `json:"timed"`
}

func (t *Event) PrepareReceived() {
	t.ID = ""
}

func (t *Event) PrepareToSend() {
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
