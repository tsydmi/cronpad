package rest

import (
	"errors"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/ts-dmitry/cronpad/backend/repository"
	"net/http"
)

type projectHandlers struct {
	store       ProjectStore
	userService UserService
}

type ProjectStore interface {
	FindAllActiveProjectsByUser(userID string) (repository.Projects, error)
}

func (t *projectHandlers) findAllByUser(writer http.ResponseWriter, request *http.Request) {
	user, err := GetUserInfo(request)
	if err != nil {
		SendAuthorizationErrorJSON(writer, request, err)
		return
	}

	projects, err := t.store.FindAllActiveProjectsByUser(user.ID)
	if err != nil {
		SendErrorJSON(writer, request, http.StatusBadRequest, err, "can't get project", ErrInternal)
		return
	}

	for i := range projects {
		projects[i].PrepareToSend()
	}

	render.Status(request, http.StatusOK)
	render.JSON(writer, request, projects)
}

func (t *projectHandlers) users(writer http.ResponseWriter, request *http.Request) {
	projectID := chi.URLParam(request, "id")

	token, err := GetAuthTokenFromHeader(request)
	if err != nil {
		SendErrorJSON(writer, request, http.StatusBadRequest, err, "can't get authentication header", ErrInternal)
		return
	}

	err = t.verifyRightsForSelectedProject(request, projectID)
	if err != nil {
		SendAuthorizationErrorJSON(writer, request, err)
		return
	}

	report, err := t.userService.FindByProject(token, projectID)
	if err != nil {
		SendErrorJSON(writer, request, http.StatusBadRequest, err, "can't get users", ErrInternal)
		return
	}

	render.Status(request, http.StatusOK)
	render.JSON(writer, request, report)
}

func (t *projectHandlers) verifyRightsForSelectedProject(request *http.Request, projectID string) error {
	userInfo, err := GetUserInfo(request)
	if err != nil {
		return err
	}

	if userInfo.hasRole(adminRole) {
		return nil
	}

	projects, err := t.store.FindAllActiveProjectsByUser(userInfo.ID)
	if err != nil {
		return err
	}

	if !projects.HasID(projectID) {
		return errors.New("user does not have rights for selected project")
	}

	return nil
}
