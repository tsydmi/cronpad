package repository

type Tag struct {
	ID     string `json:"id,omitempty" bson:"_id"`
	Name   string `json:"name" validate:"required,max=256"`
	Color  string `json:"color" validate:"required,max=16"`
	Parent string `json:"parent" validate:"max=36"`
	UserID string `json:"user,omitempty"`
}

func (t *Tag) PrepareReceived() {
	t.ID = ""
	t.UserID = ""
}

func (t *Tag) PrepareToSend() {
	t.UserID = ""
}

func (t *Tag) IsValid() []error {
	return nil
}
