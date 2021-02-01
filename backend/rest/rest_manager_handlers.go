package rest

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	R "github.com/go-pkgz/rest"
	"github.com/ts-dmitry/cronpad/backend/repository"
	"github.com/ts-dmitry/cronpad/backend/service"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

type managerHandlers struct {
	validator     *FormValidator
	projectStore  ProjectStore
	tagStore      ManagerTagStore
	userService   *service.UserService
	reportService *service.ReportService
}

type ManagerTagStore interface {
	Create(tag repository.Tag) (*mongo.InsertOneResult, error)
	Update(tag repository.Tag) (string, error)
	DeleteByProjectID(tagID string, projectIDs []string) error
	Delete(tagID string) error
}

func (h *managerHandlers) getProjectReport(writer http.ResponseWriter, request *http.Request) {
	id := chi.URLParam(request, "id")
	if len(id) == 0 {
		SendErrorJSON(writer, request, http.StatusBadRequest, errors.New("id can't be empty"), "", ErrInternal)
		return
	}

	err := h.verifyRightsForSelectedProject(request, id)
	if err != nil {
		SendAuthorizationErrorJSON(writer, request, err)
		return
	}

	projectReport, err := h.reportService.CalculateProjectReport(id)
	if err != nil {
		SendErrorJSON(writer, request, http.StatusBadRequest, err, "can't get report", ErrInternal)
		return
	}

	render.Status(request, http.StatusOK)
	render.JSON(writer, request, projectReport)
}

func (h *managerHandlers) getProjectUsers(writer http.ResponseWriter, request *http.Request) {
	projectID := chi.URLParam(request, "id")

	token, err := GetAuthTokenFromHeader(request)
	if err != nil {
		SendErrorJSON(writer, request, http.StatusBadRequest, err, "can't get authentication header", ErrInternal)
		return
	}

	err = h.verifyRightsForSelectedProject(request, projectID)
	if err != nil {
		SendAuthorizationErrorJSON(writer, request, err)
		return
	}

	report, err := h.userService.FindByProject(token, projectID)
	if err != nil {
		SendErrorJSON(writer, request, http.StatusBadRequest, err, "can't get users", ErrInternal)
		return
	}

	render.Status(request, http.StatusOK)
	render.JSON(writer, request, report)
}

func (h *managerHandlers) createTag(writer http.ResponseWriter, request *http.Request) {
	var tag repository.Tag
	err := json.NewDecoder(request.Body).Decode(&tag)
	if err != nil {
		SendErrorJSON(writer, request, http.StatusBadRequest, err, "can't parse json", ErrInternal)
		return
	}

	err = h.validator.validate(&tag)
	if err != nil {
		SendValidationErrorJSON(writer, request, err)
		return
	}

	err = h.verifyRightsForSelectedProject(request, tag.ProjectID)
	if err != nil {
		SendAuthorizationErrorJSON(writer, request, err)
		return
	}

	tag.PrepareReceivedProjectTag()

	result, err := h.tagStore.Create(tag)
	if err != nil || result == nil {
		SendErrorJSON(writer, request, http.StatusBadRequest, err, "can't insert tag", ErrInternal)
		return
	}

	render.Status(request, http.StatusCreated)
	render.JSON(writer, request, R.JSON{"id": result.InsertedID})
}

func (h *managerHandlers) updateTag(writer http.ResponseWriter, request *http.Request) {
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

	err = h.validator.validate(&tag)
	if err != nil {
		SendValidationErrorJSON(writer, request, err)
		return
	}

	err = h.verifyRightsForSelectedProject(request, tag.ProjectID)
	if err != nil {
		SendAuthorizationErrorJSON(writer, request, err)
		return
	}

	tag.PrepareReceivedProjectTag()
	tag.ID = id

	id, err = h.tagStore.Update(tag)
	if err != nil {
		SendErrorJSON(writer, request, http.StatusBadRequest, err, "can't update tag", ErrInternal)
		return
	}

	render.Status(request, http.StatusCreated)
	render.JSON(writer, request, R.JSON{"id": id})
}

func (h *managerHandlers) deleteTag(writer http.ResponseWriter, request *http.Request) {
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

	projects, err := h.projectStore.FindAllActiveProjectsByUser(userInfo.ID)
	if err != nil {
		SendAuthorizationErrorJSON(writer, request, err)
		return
	}

	err = h.tagStore.DeleteByProjectID(id, projects.GetIDs())
	if err != nil {
		SendErrorJSON(writer, request, http.StatusInternalServerError, err, "can't delete tag", ErrInternal)
		return
	}

	render.Status(request, http.StatusOK)
	render.JSON(writer, request, R.JSON{"id": id})
}

func (h *managerHandlers) verifyRightsForSelectedProject(request *http.Request, projectID string) error {
	userInfo, err := GetUserInfo(request)
	if err != nil {
		return err
	}

	if userInfo.hasRole(adminRole) {
		return nil
	}

	projects, err := h.projectStore.FindAllActiveProjectsByUser(userInfo.ID)
	if err != nil {
		return err
	}

	if !projects.HasID(projectID) {
		return errors.New("user does not have rights for selected project")
	}

	return nil
}
