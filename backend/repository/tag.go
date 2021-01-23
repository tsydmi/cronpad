package repository

type Tag struct {
	ID     string `json:"id,omitempty" bson:"_id"`
	Name   string `json:"name"`
	Color  string `json:"color"`
	Parent string `json:"parent"`
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