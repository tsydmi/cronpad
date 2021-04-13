package repository

import "errors"

type Settings struct {
	TimeRange TimeRange `json:"timeRange" validate:"required"`
	UserID    string    `json:"user,omitempty" bson:"user"`
}

type TimeRange struct {
	Min int8 `json:"min" validate:"required,min=0,max=24"`
	Max int8 `json:"max" validate:"required,min=0,max=24"`
}

func (s *Settings) PrepareReceived() {
	s.UserID = ""
}

func (s *Settings) PrepareToSend() {
	s.UserID = ""
}

func (s *Settings) IsValid() []error {
	var errs []error

	if s.TimeRange.Min > s.TimeRange.Max {
		errs = append(errs, errors.New("wrong time range"))
	}

	return errs
}

func CreateDefaultSettings() Settings {
	return Settings{
		TimeRange: TimeRange{
			Min: 6,
			Max: 20,
		},
	}
}