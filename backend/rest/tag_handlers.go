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

type tagHandlers struct {
	store     TagStore
	validator *FormValidator
}

type TagStore interface {
	Create(tag repository.Tag) (*mongo.InsertOneResult, error)
	FindAllActive() ([]repository.Tag, error)
	Update(tag repository.Tag) (string, error)
	Delete(tagID string) error
}

func (t *tagHandlers) create(writer http.ResponseWriter, request *http.Request) {
	var tag repository.Tag
	err := json.NewDecoder(request.Body).Decode(&tag)
	if err != nil {
		SendErrorJSON(writer, request, http.StatusBadRequest, err, "can't parse json", ErrInternal)
		return
	}

	err = t.validator.validate(&tag)
	if err != nil {
		SendValidationErrorJSON(writer, request, err)
		return
	}

	tag.PrepareReceived()

	result, err := t.store.Create(tag)
	if err != nil || result == nil {
		SendErrorJSON(writer, request, http.StatusBadRequest, err, "can't insert tag", ErrInternal)
		return
	}

	render.Status(request, http.StatusCreated)
	render.JSON(writer, request, R.JSON{"id": result.InsertedID})
}

func (t *tagHandlers) findAll(writer http.ResponseWriter, request *http.Request) {
	tags, err := t.store.FindAllActive()
	if err != nil {
		SendErrorJSON(writer, request, http.StatusBadRequest, err, "can't get tag", ErrInternal)
		return
	}

	render.Status(request, http.StatusOK)
	render.JSON(writer, request, tags)
}

func (t *tagHandlers) update(writer http.ResponseWriter, request *http.Request) {
	id := chi.URLParam(request, "id")
	if len(id) == 0 {
		SendErrorJSON(writer, request, http.StatusBadRequest, errors.New("id can't be empty"), "", ErrInternal)
		return
	}

	var tag repository.Tag
	err := json.NewDecoder(request.Body).Decode(&tag)
	if err != nil {
		SendErrorJSON(writer, request, http.StatusBadRequest, err, "can't parse json", ErrInternal)
		return
	}

	err = t.validator.validate(&tag)
	if err != nil {
		SendValidationErrorJSON(writer, request, err)
		return
	}

	tag.PrepareReceived()
	tag.ID = id

	id, err = t.store.Update(tag)
	if err != nil {
		SendErrorJSON(writer, request, http.StatusBadRequest, err, "can't update tag", ErrInternal)
		return
	}

	render.Status(request, http.StatusCreated)
	render.JSON(writer, request, R.JSON{"id": id})
}

func (t *tagHandlers) delete(writer http.ResponseWriter, request *http.Request) {
	id := chi.URLParam(request, "id")

	err := t.store.Delete(id)
	if err != nil {
		SendErrorJSON(writer, request, http.StatusInternalServerError, err, "can't delete tag", ErrInternal)
		return
	}

	render.Status(request, http.StatusOK)
	render.JSON(writer, request, R.JSON{"id": id})
}
