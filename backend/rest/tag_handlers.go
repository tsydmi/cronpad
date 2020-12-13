package rest

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	R "github.com/go-pkgz/rest"
	"github.com/go-playground/validator/v10"
	"github.com/ts-dmitry/cronpad/backend/repository"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

type tagHandlers struct {
	store     TagStore
	validator *validator.Validate
}

type TagStore interface {
	Create(tag repository.Tag) (*mongo.InsertOneResult, error)
	FindAll(userID string) ([]repository.Tag, error)
	Delete(tagID string, userID string) error
}

func (t *tagHandlers) findAll(writer http.ResponseWriter, request *http.Request) {
	user, err := GetUserInfo(request)
	if err != nil {
		SendAuthorizationErrorJSON(writer, request, err)
		return
	}

	tags, err := t.store.FindAll(user.ID)
	if err != nil {
		SendErrorJSON(writer, request, http.StatusBadRequest, err, "can't get tag", ErrInternal)
		return
	}

	render.Status(request, http.StatusOK)
	render.JSON(writer, request, tags)
}

func (t *tagHandlers) create(writer http.ResponseWriter, request *http.Request) {
	user, err := GetUserInfo(request)
	if err != nil {
		SendAuthorizationErrorJSON(writer, request, err)
		return
	}

	var tag repository.Tag
	err = json.NewDecoder(request.Body).Decode(&tag)
	if err != nil {
		SendErrorJSON(writer, request, http.StatusBadRequest, err, "can't parse json", ErrInternal)
		return
	}

	err = t.validator.Struct(tag)
	if err != nil {
		SendValidationErrorJSON(writer, request, err)
		return
	}

	tag.PrepareReceived()
	tag.UserID = user.ID

	result, err := t.store.Create(tag)
	if err != nil || result == nil {
		SendErrorJSON(writer, request, http.StatusBadRequest, err, "can't insert tag", ErrInternal)
		return
	}

	render.Status(request, http.StatusCreated)
	render.JSON(writer, request, R.JSON{"id": result.InsertedID})
}

func (t *tagHandlers) delete(writer http.ResponseWriter, request *http.Request) {
	user, err := GetUserInfo(request)
	if err != nil {
		SendAuthorizationErrorJSON(writer, request, err)
		return
	}

	id := chi.URLParam(request, "id")

	err = t.store.Delete(id, user.ID)
	if err != nil {
		SendErrorJSON(writer, request, http.StatusInternalServerError, err, "can't delete tag", ErrInternal)
		return
	}

	render.Status(request, http.StatusOK)
	render.JSON(writer, request, R.JSON{"id": id})
}
