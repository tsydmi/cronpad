package repository

type Project struct {
	ID          string   `json:"id,omitempty" bson:"_id"`
	Name        string   `json:"name" validate:"required,max=256"`
	Description string   `json:"description" validate:"required,max=256"`
	Users       []string `json:"users"`
}

func (t *Project) PrepareReceived() {
	t.ID = ""
}

func (t *Project) PrepareToSend() {
	t.Users = nil
}

func (t *Project) IsValid() []error {
	return nil
}
