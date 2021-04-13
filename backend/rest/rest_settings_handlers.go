package rest

import (
	"encoding/json"
	"github.com/go-chi/render"
	R "github.com/go-pkgz/rest"
	"github.com/ts-dmitry/cronpad/backend/repository"
	"net/http"
)

type settingsHandlers struct {
	store     SettingsStore
	validator *FormValidator
}

type SettingsStore interface {
	FindByUser(userID string) (*repository.Settings, error)
	CreateOrUpdate(settings repository.Settings) error
}

func (s *settingsHandlers) get(writer http.ResponseWriter, request *http.Request) {
	user, err := GetUserInfo(request)
	if err != nil {
		SendAuthorizationErrorJSON(writer, request, err)
		return
	}

	settings, err := s.store.FindByUser(user.ID)
	if err != nil {
		defaultSettings := repository.CreateDefaultSettings()
		settings = &defaultSettings
	}

	settings.PrepareToSend()

	render.Status(request, http.StatusOK)
	render.JSON(writer, request, settings)
}

func (s *settingsHandlers) update(writer http.ResponseWriter, request *http.Request) {
	user, err := GetUserInfo(request)
	if err != nil {
		SendAuthorizationErrorJSON(writer, request, err)
		return
	}

	var settings repository.Settings
	err = json.NewDecoder(request.Body).Decode(&settings)
	if err != nil {
		SendErrorJSON(writer, request, http.StatusBadRequest, err, "can't parse json", ErrInternal)
		return
	}

	err = s.validator.validate(&settings)
	if err != nil {
		SendValidationErrorJSON(writer, request, err)
		return
	}

	settings.PrepareReceived()
	settings.UserID = user.ID

	err = s.store.CreateOrUpdate(settings)
	if err != nil {
		SendErrorJSON(writer, request, http.StatusBadRequest, err, "can't update settings", ErrInternal)
		return
	}

	render.Status(request, http.StatusCreated)
	render.JSON(writer, request, R.JSON{})
}
