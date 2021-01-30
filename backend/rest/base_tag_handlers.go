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

type baseTagHandlers struct {
	store     BaseTagStore
	validator *FormValidator
}

type BaseTagStore interface {
	Create(tag repository.Tag) (*mongo.InsertOneResult, error)
	UpdateBaseTag(tag repository.Tag) (string, error)
	DeleteBaseTag(tagID string) error
}

func (t *baseTagHandlers) create(writer http.ResponseWriter, request *http.Request) {
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

	tag.PrepareReceivedBaseTag()

	result, err := t.store.Create(tag)
	if err != nil || result == nil {
		SendErrorJSON(writer, request, http.StatusBadRequest, err, "can't insert tag", ErrInternal)
		return
	}

	render.Status(request, http.StatusCreated)
	render.JSON(writer, request, R.JSON{"id": result.InsertedID})
}

func (t *baseTagHandlers) update(writer http.ResponseWriter, request *http.Request) {
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

	tag.PrepareReceivedBaseTag()
	tag.ID = id

	id, err = t.store.UpdateBaseTag(tag)
	if err != nil {
		SendErrorJSON(writer, request, http.StatusBadRequest, err, "can't update tag", ErrInternal)
		return
	}

	render.Status(request, http.StatusCreated)
	render.JSON(writer, request, R.JSON{"id": id})
}

func (t *baseTagHandlers) delete(writer http.ResponseWriter, request *http.Request) {
	id := chi.URLParam(request, "id")

	err := t.store.DeleteBaseTag(id)
	if err != nil {
		SendErrorJSON(writer, request, http.StatusInternalServerError, err, "can't delete tag", ErrInternal)
		return
	}

	render.Status(request, http.StatusOK)
	render.JSON(writer, request, R.JSON{"id": id})
}
