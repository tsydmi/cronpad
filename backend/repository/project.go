package repository

import (
	"errors"
	"time"
)

type Project struct {
	ID          string     `json:"id,omitempty" bson:"_id"`
	Name        string     `json:"name" validate:"required,max=256"`
	Description string     `json:"description" validate:"required,max=256"`
	Users       []string   `json:"users"`
	Start       *time.Time `json:"start,omitempty" bson:"start,omitempty"`
	End         *time.Time `json:"end,omitempty" bson:"end,omitempty"`
}

func (t *Project) PrepareReceived() {
	t.ID = ""
}

func (t *Project) PrepareToSend() {
	t.Users = nil
}

func (t *Project) IsValid() []error {
	var errs []error

	if t.Start == nil && t.End != nil {
		errs = append(errs, errors.New("end date cannot be set without start date"))
	}

	if t.Start != nil && t.End != nil && t.Start.After(*t.End) {
		errs = append(errs, errors.New("wrong date range"))
	}
	return errs
}

type Projects []Project

func (p Projects) GetIDs() []string {
	ids := make([]string, 0)

	for _, project := range p {
		ids = append(ids, project.ID)
	}

	return ids
}

func (p Projects) HasID(projectID string) bool {
	for _, project := range p {
		if project.ID == projectID {
			return true
		}
	}

	return false
}
