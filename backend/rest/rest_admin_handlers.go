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

type adminHandlers struct {
	validator     *FormValidator
	projectStore  AdminProjectStore
	baseTagStore  BaseTagStore
	userService   *service.UserService
	reportService *service.ReportService
}

type AdminProjectStore interface {
	Create(record repository.Project) (*mongo.InsertOneResult, error)
	Search(form repository.ProjectSearchForm) ([]repository.Project, error)
	Update(record repository.Project) (string, error)
	Delete(projectID string) error
}

type BaseTagStore interface {
	Create(tag repository.Tag) (*mongo.InsertOneResult, error)
	UpdateBaseTag(tag repository.Tag) (string, error)
	DeleteBaseTag(tagID string) error
}

// Project CRUD
func (h *adminHandlers) createProject(writer http.ResponseWriter, request *http.Request) {
	var project repository.Project
	err := json.NewDecoder(request.Body).Decode(&project)
	if err != nil {
		SendErrorJSON(writer, request, http.StatusBadRequest, err, "can't parse json", ErrInternal)
		return
	}

	err = h.validator.validate(&project)
	if err != nil {
		SendValidationErrorJSON(writer, request, err)
		return
	}

	project.PrepareReceived()

	result, err := h.projectStore.Create(project)
	if err != nil || result == nil {
		SendErrorJSON(writer, request, http.StatusBadRequest, err, "can't insert project", ErrInternal)
		return
	}

	render.Status(request, http.StatusCreated)
	render.JSON(writer, request, R.JSON{"id": result.InsertedID})
}

func (h *adminHandlers) updateProject(writer http.ResponseWriter, request *http.Request) {
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

	err = h.validator.validate(&project)
	if err != nil {
		SendValidationErrorJSON(writer, request, err)
		return
	}

	project.PrepareReceived()
	project.ID = id

	id, err = h.projectStore.Update(project)
	if err != nil {
		SendErrorJSON(writer, request, http.StatusBadRequest, err, "can't update project", ErrInternal)
		return
	}

	render.Status(request, http.StatusCreated)
	render.JSON(writer, request, R.JSON{"id": id})
}

func (h *adminHandlers) deleteProject(writer http.ResponseWriter, request *http.Request) {
	id := chi.URLParam(request, "id")

	err := h.projectStore.Delete(id)
	if err != nil {
		SendErrorJSON(writer, request, http.StatusInternalServerError, err, "can't delete project", ErrInternal)
		return
	}

	render.Status(request, http.StatusOK)
	render.JSON(writer, request, R.JSON{"id": id})
}

func (h *adminHandlers) findProject(writer http.ResponseWriter, request *http.Request) {
	var form repository.ProjectSearchForm
	err := json.NewDecoder(request.Body).Decode(&form)
	if err != nil {
		SendErrorJSON(writer, request, http.StatusBadRequest, err, "can't parse json", ErrInternal)
		return
	}

	report, err := h.projectStore.Search(form)
	if err != nil {
		SendErrorJSON(writer, request, http.StatusBadRequest, err, "can't get projects", ErrInternal)
		return
	}

	render.Status(request, http.StatusOK)
	render.JSON(writer, request, report)
}

// Base tag CRUD
func (h *adminHandlers) createBaseTag(writer http.ResponseWriter, request *http.Request) {
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

	tag.PrepareReceivedBaseTag()

	result, err := h.baseTagStore.Create(tag)
	if err != nil || result == nil {
		SendErrorJSON(writer, request, http.StatusBadRequest, err, "can't insert tag", ErrInternal)
		return
	}

	render.Status(request, http.StatusCreated)
	render.JSON(writer, request, R.JSON{"id": result.InsertedID})
}

func (h *adminHandlers) updateBaseTag(writer http.ResponseWriter, request *http.Request) {
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

	tag.PrepareReceivedBaseTag()
	tag.ID = id

	id, err = h.baseTagStore.UpdateBaseTag(tag)
	if err != nil {
		SendErrorJSON(writer, request, http.StatusBadRequest, err, "can't update tag", ErrInternal)
		return
	}

	render.Status(request, http.StatusCreated)
	render.JSON(writer, request, R.JSON{"id": id})
}

func (h *adminHandlers) deleteBaseTag(writer http.ResponseWriter, request *http.Request) {
	id := chi.URLParam(request, "id")

	err := h.baseTagStore.DeleteBaseTag(id)
	if err != nil {
		SendErrorJSON(writer, request, http.StatusInternalServerError, err, "can't delete tag", ErrInternal)
		return
	}

	render.Status(request, http.StatusOK)
	render.JSON(writer, request, R.JSON{"id": id})
}

func (h *adminHandlers) findAllUser(writer http.ResponseWriter, request *http.Request) {
	token, err := GetAuthTokenFromHeader(request)
	if err != nil {
		SendErrorJSON(writer, request, http.StatusBadRequest, err, "can't get authentication header", ErrInternal)
		return
	}

	users, err := h.userService.FindAll(token)
	if err != nil {
		SendErrorJSON(writer, request, http.StatusBadRequest, err, "can't get users", ErrInternal)
		return
	}

	render.Status(request, http.StatusOK)
	render.JSON(writer, request, users)
}

func (t *adminHandlers) userReport(writer http.ResponseWriter, request *http.Request) {
	var form repository.DaySearchForm
	err := json.NewDecoder(request.Body).Decode(&form)
	if err != nil {
		SendErrorJSON(writer, request, http.StatusBadRequest, err, "can't parse json", ErrInternal)
		return
	}

	userReport, err := t.reportService.CalculateUserReport(form)
	if err != nil {
		SendErrorJSON(writer, request, http.StatusBadRequest, err, "can't get report", ErrInternal)
		return
	}

	render.Status(request, http.StatusOK)
	render.JSON(writer, request, userReport)
}
