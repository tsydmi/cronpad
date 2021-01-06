package rest

import (
	"github.com/go-chi/render"
	"github.com/ts-dmitry/cronpad/backend/service"
	"net/http"
)

type userHandlers struct {
	service UserService
}

type UserService interface {
	FindAll(authorizationHeader string) ([]service.User, error)
	FindByProject(token string, projectID string) ([]service.User, error)
}

func (t *userHandlers) findAll(writer http.ResponseWriter, request *http.Request) {
	token, err := GetAuthTokenFromHeader(request)
	if err != nil {
		SendErrorJSON(writer, request, http.StatusBadRequest, err, "can't get authentication header", ErrInternal)
		return
	}

	users, err := t.service.FindAll(token)
	if err != nil {
		SendErrorJSON(writer, request, http.StatusBadRequest, err, "can't get users", ErrInternal)
		return
	}

	render.Status(request, http.StatusOK)
	render.JSON(writer, request, users)
}
