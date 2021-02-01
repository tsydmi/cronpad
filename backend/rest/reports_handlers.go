package rest

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/ts-dmitry/cronpad/backend/repository"
	"github.com/ts-dmitry/cronpad/backend/service/report"
	"net/http"
)

type reportsHandlers struct {
	projectStore ProjectStore
	service      ReportService
}

type ReportService interface {
	CalculateUserReport(form repository.DaySearchForm) (report.UserReport, error)
	CalculateProjectReport(projectID string) (report.ProjectReport, error)
}

func (t *reportsHandlers) userReport(writer http.ResponseWriter, request *http.Request) {
	var form repository.DaySearchForm
	err := json.NewDecoder(request.Body).Decode(&form)
	if err != nil {
		SendErrorJSON(writer, request, http.StatusBadRequest, err, "can't parse json", ErrInternal)
		return
	}

	userReport, err := t.service.CalculateUserReport(form)
	if err != nil {
		SendErrorJSON(writer, request, http.StatusBadRequest, err, "can't get report", ErrInternal)
		return
	}

	render.Status(request, http.StatusOK)
	render.JSON(writer, request, userReport)
}

func (t *reportsHandlers) projectReport(writer http.ResponseWriter, request *http.Request) {
	id := chi.URLParam(request, "id")
	if len(id) == 0 {
		SendErrorJSON(writer, request, http.StatusBadRequest, errors.New("id can't be empty"), "", ErrInternal)
		return
	}

	err := t.verifyRightsForSelectedProject(request, id)
	if err != nil {
		SendAuthorizationErrorJSON(writer, request, err)
		return
	}

	projectReport, err := t.service.CalculateProjectReport(id)
	if err != nil {
		SendErrorJSON(writer, request, http.StatusBadRequest, err, "can't get report", ErrInternal)
		return
	}

	render.Status(request, http.StatusOK)
	render.JSON(writer, request, projectReport)
}

func (t *reportsHandlers) verifyRightsForSelectedProject(request *http.Request, projectID string) error {
	userInfo, err := GetUserInfo(request)
	if err != nil {
		return err
	}

	if userInfo.hasRole(adminRole) {
		return nil
	}

	projects, err := t.projectStore.FindAllActiveProjectsByUser(userInfo.ID)
	if err != nil {
		return err
	}

	if !projects.HasID(projectID) {
		return errors.New("user does not have rights for selected project")
	}

	return nil
}
