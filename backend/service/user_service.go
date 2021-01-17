package service

import (
	"encoding/json"
	"fmt"
	"github.com/ts-dmitry/cronpad/backend/repository"
	"net/http"
	"time"
)

const allUsersPath = "/auth/admin/realms/cronpad/users?briefRepresentation=true"

type UserService struct {
	keycloakUrl  string
	projectStore ProjectStore
}

func CreateUserService(keycloakUrl string, projectStore ProjectStore) *UserService {
	return &UserService{keycloakUrl: keycloakUrl, projectStore: projectStore}
}

type ProjectStore interface {
	GetProjectWithUsersByID(projectID string) (repository.Project, error)
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

func (t *UserService) FindByProject(token string, projectID string) ([]User, error) {
	allUsers, err := t.FindAll(token)
	if err != nil {
		return nil, err
	}

	project, err := t.projectStore.GetProjectWithUsersByID(projectID)
	if err != nil {
		return nil, err
	}

	result := make([]User, 0)
	for i := range allUsers {
		user := allUsers[i]
		for j := range project.Users {
			projectUserID := project.Users[j]
			if user.Id == projectUserID {
				result = append(result, user)
			}
		}
	}

	return result, nil
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
