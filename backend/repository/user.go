package repository

type User struct {
	ID   string `json:"id" bson:"_id"`
	Name string `json:"name"`
}
