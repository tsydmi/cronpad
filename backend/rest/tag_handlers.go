package rest

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	R "github.com/go-pkgz/rest"
	"github.com/ts-dmitry/cronpad/backend/repository"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

type tagHandlers struct {
	store TagStore
}

type TagStore interface {
	Create(record repository.Tag) (*mongo.InsertOneResult, error)
	FindAll(userID string) ([]repository.Tag, error)
	Delete(timeID string, userID string) error
}

func (t *tagHandlers) findAll(writer http.ResponseWriter, request *http.Request) {
	user, err := GetUserInfo(request)
	if err != nil {
		SendErrorJSON(writer, request, http.StatusForbidden, err, "user should be logged in", ErrInternal)
		return
	}

	timeRecords, err := t.store.FindAll(user.ID)
	if err != nil {
		SendErrorJSON(writer, request, http.StatusBadRequest, err, "can't get time record", ErrInternal)
		return
	}

	render.Status(request, http.StatusOK)
	render.JSON(writer, request, timeRecords)
}

func (t *tagHandlers) create(writer http.ResponseWriter, request *http.Request) {
	user, err := GetUserInfo(request)
	if err != nil {
		SendErrorJSON(writer, request, http.StatusForbidden, err, "user should be logged in", ErrInternal)
		return
	}

	var tag repository.Tag
	err = json.NewDecoder(request.Body).Decode(&tag)
	if err != nil {
		SendErrorJSON(writer, request, http.StatusBadRequest, err, "can't parse json", ErrInternal)
		return
	}

	tag.PrepareReceived()
	tag.UserID = user.ID

	result, err := t.store.Create(tag)
	if err != nil || result == nil {
		SendErrorJSON(writer, request, http.StatusBadRequest, err, "can't insert time record", ErrInternal)
		return
	}

	render.Status(request, http.StatusCreated)
	render.JSON(writer, request, R.JSON{"id": result.InsertedID})
}

func (t *tagHandlers) delete(writer http.ResponseWriter, request *http.Request) {
	user, err := GetUserInfo(request)
	if err != nil {
		SendErrorJSON(writer, request, http.StatusForbidden, err, "user should be logged in", ErrInternal)
		return
	}

	id := chi.URLParam(request, "id")

	err = t.store.Delete(id, user.ID)
	if err != nil {
		SendErrorJSON(writer, request, http.StatusInternalServerError, err, "can't delete comment", ErrInternal)
		return
	}

	render.Status(request, http.StatusOK)
	render.JSON(writer, request, R.JSON{"id": id})
}
