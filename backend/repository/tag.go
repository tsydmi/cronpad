package repository

type Tag struct {
	ID          string `json:"id,omitempty" bson:"_id"`
	Name        string `json:"name" validate:"required,max=256"`
	Description string `json:"description" validate:"max=256"`
	Color       string `json:"color" validate:"required,max=16"`
}

func (t *Tag) PrepareReceived() {
	t.ID = ""
}

func (t *Tag) PrepareToSend() {
}

func (t *Tag) IsValid() []error {
	return nil
}
