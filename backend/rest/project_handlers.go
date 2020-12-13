package rest

import (
	"github.com/go-chi/render"
	"github.com/ts-dmitry/cronpad/backend/repository"
	"net/http"
)

type projectHandlers struct {
	store ProjectStore
}

type ProjectStore interface {
	FindAllProjectsByUser(userID string) ([]repository.Project, error)
}

func (t *projectHandlers) findAllByUser(writer http.ResponseWriter, request *http.Request) {
	user, err := GetUserInfo(request)
	if err != nil {
		SendAuthorizationErrorJSON(writer, request, err)
		return
	}

	projects, err := t.store.FindAllProjectsByUser(user.ID)
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
