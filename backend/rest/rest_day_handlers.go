package rest

import (
	"errors"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/ts-dmitry/cronpad/backend/repository"
	"net/http"
	"net/url"
	"time"
)

const dateLayout = "2006-01-02"

type dayHandlers struct {
	store DayStore
}

type DayStore interface {
	FindByDateRange(from time.Time, to time.Time, userID string) ([]repository.Day, error)
	FindByDate(date time.Time, userID string) (repository.Day, error)
}

func (t *dayHandlers) findByDate(writer http.ResponseWriter, request *http.Request) {
	user, err := GetUserInfo(request)
	if err != nil {
		SendAuthorizationErrorJSON(writer, request, err)
		return
	}

	date, err := time.Parse(dateLayout, chi.URLParam(request, "date"))
	if err != nil {
		SendErrorJSON(writer, request, http.StatusBadRequest, err, "date is missing", ErrInternal)
		return
	}

	days, err := t.store.FindByDate(date, user.ID)
	if err != nil {
		SendErrorJSON(writer, request, http.StatusBadRequest, err, "can't get days", ErrInternal)
		return
	}

	render.Status(request, http.StatusOK)
	render.JSON(writer, request, days)
}

func (t *dayHandlers) findByDateRange(writer http.ResponseWriter, request *http.Request) {
	user, err := GetUserInfo(request)
	if err != nil {
		SendAuthorizationErrorJSON(writer, request, err)
		return
	}

	query := request.URL.Query()

	from, err := getDateQueryParam(query, "from")
	if err != nil {
		SendErrorJSON(writer, request, http.StatusBadRequest, err, "'from' date field is required", ErrInternal)
		return
	}

	to, err := getDateQueryParam(query, "to")
	if err != nil {
		SendErrorJSON(writer, request, http.StatusBadRequest, err, "'to' date field is required'", ErrInternal)
		return
	}

	days, err := t.store.FindByDateRange(from, to, user.ID)
	if err != nil {
		SendErrorJSON(writer, request, http.StatusBadRequest, err, "can't get days", ErrInternal)
		return
	}

	render.Status(request, http.StatusOK)
	render.JSON(writer, request, days)
}

func getDateQueryParam(query url.Values, name string) (time.Time, error) {
	fromString, fromPresent := query[name]

	if !fromPresent || len(fromString) == 0 {
		return time.Time{}, errors.New("date not present")
	}

	if len(fromString) > 1 {
		return time.Time{}, errors.New("wrong date")
	}

	return time.Parse(dateLayout, fromString[0])
}
