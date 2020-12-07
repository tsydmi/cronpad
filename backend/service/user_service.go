package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const allUsersPath = "/auth/admin/realms/cronpad/users?briefRepresentation=true"

type UserService struct {
	keycloakUrl string
}

func CreateUserService(keycloakUrl string) *UserService {
	return &UserService{keycloakUrl: keycloakUrl}
}

func (t *UserService) FindAll(token string) ([]User, error) {
	client := http.Client{
		Timeout: time.Duration(5 * time.Second),
	}

	usersURL := t.keycloakUrl + allUsersPath
	request, err := http.NewRequest("GET", usersURL, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Add("Authorization", "Bearer "+token)

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	users := make([]keycloakUser, 0)
	if err := json.NewDecoder(response.Body).Decode(&users); err != nil {
		return nil, fmt.Errorf("unable to read key %s", err)
	}

	return convertToUser(users), nil
}

func convertToUser(keycloakUsers []keycloakUser) []User {
	result := make([]User, 0)
	for i := range keycloakUsers {
		user := keycloakUsers[i]
		result = append(result, User{
			Id:        user.Id,
			FirstName: user.FirstName,
			LastName:  user.LastName,
		})
	}

	return result
}

type User struct {
	Id        string `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type keycloakUser struct {
	Id               string `json:"id"`
	CreatedTimestamp int64  `json:"createdTimestamp"`
	Username         string `json:"username"`
	Enabled          bool   `json:"enabled"`
	EmailVerified    bool   `json:"emailVerified"`
	FirstName        string `json:"firstName"`
	LastName         string `json:"lastName"`
	Access           access `json:"access"`
}

type access struct {
	ManageGroupMembership bool `json:"manageGroupMembership"`
	View                  bool `json:"view"`
	MapRoles              bool `json:"mapRoles"`
	Impersonate           bool `json:"impersonate"`
	Manage                bool `json:"manage"`
}
