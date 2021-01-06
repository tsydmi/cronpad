package rest

import (
	"encoding/json"
	"github.com/go-chi/render"
	"github.com/ts-dmitry/cronpad/backend/repository"
	"github.com/ts-dmitry/cronpad/backend/service"
	"net/http"
)

type reportsHandlers struct {
	service ReportService
}

type ReportService interface {
	Search(form repository.DaySearchForm) (service.ChartReport, error)
}

func (t *reportsHandlers) search(writer http.ResponseWriter, request *http.Request) {
	var form repository.DaySearchForm
	err := json.NewDecoder(request.Body).Decode(&form)
	if err != nil {
		SendErrorJSON(writer, request, http.StatusBadRequest, err, "can't parse json", ErrInternal)
		return
	}

	report, err := t.service.Search(form)
	if err != nil {
		SendErrorJSON(writer, request, http.StatusBadRequest, err, "can't get report", ErrInternal)
		return
	}

	render.Status(request, http.StatusOK)
	render.JSON(writer, request, report)
}
