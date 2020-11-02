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
	ID    string    `json:"id,omitempty" bson:"_id"`
	Name  string    `json:"name"`
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
	TagID string    `json:"tag"`
	Timed bool      `json:"timed"`
}

func (t *Event) PrepareReceived() {
	t.ID = ""
}

func (t *Event) PrepareToSend() {
}
