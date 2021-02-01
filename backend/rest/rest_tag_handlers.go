package rest

import (
	"github.com/go-chi/render"
	"github.com/ts-dmitry/cronpad/backend/repository"
	"net/http"
)

type tagHandlers struct {
	tagStore     TagStore
	projectStore ProjectStore
	validator    *FormValidator
}

type TagStore interface {
	FindAllBaseAndProjectActiveTags(projectID []string) ([]repository.Tag, error)
	FindAllActive() ([]repository.Tag, error)
}

func (t *tagHandlers) findAll(writer http.ResponseWriter, request *http.Request) {
	user, err := GetUserInfo(request)
	if err != nil {
		SendAuthorizationErrorJSON(writer, request, err)
		return
	}

	tags, err := t.getTags(user)
	if err != nil {
		SendErrorJSON(writer, request, http.StatusBadRequest, err, "can't get tag", ErrInternal)
		return
	}

	render.Status(request, http.StatusOK)
	render.JSON(writer, request, tags)
}

func (t *tagHandlers) getTags(user UserInfo) ([]repository.Tag, error) {
	if user.hasRole(adminRole) {
		return t.tagStore.FindAllActive()
	}

	projects, err := t.projectStore.FindAllActiveProjectsByUser(user.ID)
	if err != nil {
		return nil, err
	}

	return t.tagStore.FindAllBaseAndProjectActiveTags(projects.GetIDs())
}
