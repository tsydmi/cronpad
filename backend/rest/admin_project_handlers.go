package rest

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	R "github.com/go-pkgz/rest"
	"github.com/ts-dmitry/cronpad/backend/repository"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

type adminProjectHandlers struct {
	validator   *FormValidator
	store       AdminProjectStore
	userService UserService
}

type AdminProjectStore interface {
	Create(record repository.Project) (*mongo.InsertOneResult, error)
	Search(form repository.ProjectSearchForm) ([]repository.Project, error)
	Update(record repository.Project) (string, error)
	Delete(projectID string) error
}

func (t *adminProjectHandlers) create(writer http.ResponseWriter, request *http.Request) {
	var project repository.Project
	err := json.NewDecoder(request.Body).Decode(&project)
	if err != nil {
		SendErrorJSON(writer, request, http.StatusBadRequest, err, "can't parse json", ErrInternal)
		return
	}

	err = t.validator.validate(&project)
	if err != nil {
		SendValidationErrorJSON(writer, request, err)
		return
	}

	project.PrepareReceived()

	result, err := t.store.Create(project)
	if err != nil || result == nil {
		SendErrorJSON(writer, request, http.StatusBadRequest, err, "can't insert project", ErrInternal)
		return
	}

	render.Status(request, http.StatusCreated)
	render.JSON(writer, request, R.JSON{"id": result.InsertedID})
}

func (t *adminProjectHandlers) update(writer http.ResponseWriter, request *http.Request) {
	id := chi.URLParam(request, "id")
	if len(id) == 0 {
		SendErrorJSON(writer, request, http.StatusBadRequest, errors.New("id can't be empty"), "", ErrInternal)
		return
	}

	var project repository.Project
	err := json.NewDecoder(request.Body).Decode(&project)
	if err != nil {
		SendErrorJSON(writer, request, http.StatusBadRequest, err, "can't parse json", ErrInternal)
		return
	}

	err = t.validator.validate(&project)
	if err != nil {
		SendValidationErrorJSON(writer, request, err)
		return
	}

	project.PrepareReceived()
	project.ID = id

	id, err = t.store.Update(project)
	if err != nil {
		SendErrorJSON(writer, request, http.StatusBadRequest, err, "can't update project", ErrInternal)
		return
	}

	render.Status(request, http.StatusCreated)
	render.JSON(writer, request, R.JSON{"id": id})
}

func (t *adminProjectHandlers) delete(writer http.ResponseWriter, request *http.Request) {
	id := chi.URLParam(request, "id")

	err := t.store.Delete(id)
	if err != nil {
		SendErrorJSON(writer, request, http.StatusInternalServerError, err, "can't delete project", ErrInternal)
		return
	}

	render.Status(request, http.StatusOK)
	render.JSON(writer, request, R.JSON{"id": id})
}

func (t *adminProjectHandlers) search(writer http.ResponseWriter, request *http.Request) {
	var form repository.ProjectSearchForm
	err := json.NewDecoder(request.Body).Decode(&form)
	if err != nil {
		SendErrorJSON(writer, request, http.StatusBadRequest, err, "can't parse json", ErrInternal)
		return
	}

	report, err := t.store.Search(form)
	if err != nil {
		SendErrorJSON(writer, request, http.StatusBadRequest, err, "can't get projects", ErrInternal)
		return
	}

	render.Status(request, http.StatusOK)
	render.JSON(writer, request, report)
}
