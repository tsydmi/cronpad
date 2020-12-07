package repository

type Project struct {
	ID          string   `json:"id,omitempty" bson:"_id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Users       []string `json:"users,omitempty"`
}

func (t *Project) PrepareReceived() {
	t.ID = ""
}

func (t *Project) PrepareToSend() {
	t.Users = nil
}
