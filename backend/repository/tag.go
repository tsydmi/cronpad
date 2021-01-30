package repository

import "errors"

type Tag struct {
	ID          string `json:"id,omitempty" bson:"_id"`
	Name        string `json:"name" validate:"required,max=256"`
	Description string `json:"description" validate:"max=256"`
	Color       string `json:"color" validate:"required,max=16"`
	ParentID    string `json:"parent,omitempty" validate:"max=36" bson:"parent,omitempty"`
	ProjectID   string `json:"project,omitempty" validate:"max=36" bson:"project,omitempty"`
	Basic       bool   `json:"basic"`
}

func (t *Tag) PrepareReceivedProjectTag() {
	t.ID = ""
	t.Basic = false
}

func (t *Tag) PrepareReceivedBaseTag() {
	t.ID = ""
	t.Basic = true
	t.ProjectID = ""
	t.ParentID = ""
}

func (t *Tag) PrepareToSend() {
}

func (t *Tag) IsValid() []error {
	var errs []error

	if !t.Basic && len(t.ProjectID) == 0 {
		errs = append(errs, errors.New("project is required"))
	}

	return errs
}
