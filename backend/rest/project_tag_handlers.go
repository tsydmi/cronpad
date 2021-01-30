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

type projectTagHandlers struct {
	tagStore     ProjectTagStore
	projectStore ProjectStore
	validator    *FormValidator
}

type ProjectTagStore interface {
	Create(tag repository.Tag) (*mongo.InsertOneResult, error)
	Update(tag repository.Tag) (string, error)
	DeleteByProjectID(tagID string, projectIDs []string) error
	Delete(tagID string) error
}

func (t *projectTagHandlers) create(writer http.ResponseWriter, request *http.Request) {
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

	err = t.verifyRightsForSelectedProject(request, tag.ProjectID)
	if err != nil {
		SendAuthorizationErrorJSON(writer, request, err)
		return
	}

	tag.PrepareReceivedProjectTag()

	result, err := t.tagStore.Create(tag)
	if err != nil || result == nil {
		SendErrorJSON(writer, request, http.StatusBadRequest, err, "can't insert tag", ErrInternal)
		return
	}

	render.Status(request, http.StatusCreated)
	render.JSON(writer, request, R.JSON{"id": result.InsertedID})
}

func (t *projectTagHandlers) update(writer http.ResponseWriter, request *http.Request) {
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

	err = t.verifyRightsForSelectedProject(request, tag.ProjectID)
	if err != nil {
		SendAuthorizationErrorJSON(writer, request, err)
		return
	}

	tag.PrepareReceivedProjectTag()
	tag.ID = id

	id, err = t.tagStore.Update(tag)
	if err != nil {
		SendErrorJSON(writer, request, http.StatusBadRequest, err, "can't update tag", ErrInternal)
		return
	}

	render.Status(request, http.StatusCreated)
	render.JSON(writer, request, R.JSON{"id": id})
}

func (t *projectTagHandlers) delete(writer http.ResponseWriter, request *http.Request) {
	id := chi.URLParam(request, "id")
	if len(id) == 0 {
		SendErrorJSON(writer, request, http.StatusBadRequest, errors.New("id can't be empty"), "", ErrInternal)
		return
	}

	userInfo, err := GetUserInfo(request)
	if err != nil {
		SendAuthorizationErrorJSON(writer, request, err)
		return
	}

	projects, err := t.projectStore.FindAllActiveProjectsByUser(userInfo.ID)
	if err != nil {
		SendAuthorizationErrorJSON(writer, request, err)
		return
	}

	err = t.tagStore.DeleteByProjectID(id, projects.GetIDs())
	if err != nil {
		SendErrorJSON(writer, request, http.StatusInternalServerError, err, "can't delete tag", ErrInternal)
		return
	}

	render.Status(request, http.StatusOK)
	render.JSON(writer, request, R.JSON{"id": id})
}

func (t *projectTagHandlers) verifyRightsForSelectedProject(request *http.Request, projectID string) error {
	userInfo, err := GetUserInfo(request)
	if err != nil {
		return err
	}

	projects, err := t.projectStore.FindAllActiveProjectsByUser(userInfo.ID)
	if err != nil {
		return err
	}

	if projects.HasID(projectID) {
		return errors.New("user does not have rights for selected project")
	}

	return nil
}
