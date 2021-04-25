package service

import (
	"github.com/ts-dmitry/cronpad/backend/repository"
	"github.com/ts-dmitry/cronpad/backend/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type EventService struct {
	dayStore     DayStore
	uuidProvider utils.UuidProvider
}

type DayStore interface {
	Create(day repository.Day) (*mongo.InsertOneResult, error)
	FindByEventID(eventID string, userID string) (repository.Day, error)
	FindByDate(date time.Time, userID string) (repository.Day, error)
	FindByDateRange(from time.Time, to time.Time, userID string) ([]repository.Day, error)
	Update(day repository.Day) (repository.Day, error)
	GetUsedNames(userID string, tagID string, from time.Time, to time.Time) ([]string, error)
}

func CreateEventService(dayStore DayStore, uuidProvider utils.UuidProvider) *EventService {
	return &EventService{dayStore: dayStore, uuidProvider: uuidProvider}
}

func (t *EventService) Create(event repository.Event, userID string) (string, error) {
	event.ID = t.uuidProvider.New()

	day, err := t.dayStore.FindByDate(event.Start, userID)
	if err != nil {
		day = repository.Day{
			Date:   time.Date(event.Start.Year(), event.Start.Month(), event.Start.Day(), 0, 0, 0, 0, time.UTC),
			UserID: userID,
			Events: []repository.Event{event},
		}

		_, err = t.dayStore.Create(day)
		return event.ID, err
	}

	if len(day.Events) > 0 {
		day.Events = AddEventProperly(event, day.Events, t.uuidProvider)
	} else {
		day.Events = append(day.Events, event)
	}

	_, err = t.dayStore.Update(day)
	return event.ID, err
}

func (t *EventService) Update(event repository.Event, userID string) (string, error) {
	day, err := t.dayStore.FindByDate(event.Start, userID)
	if err != nil {
		return "", err
	}

	events := deleteEventFromSlice(day.Events, event.ID)

	if len(day.Events) > 0 {
		day.Events = AddEventProperly(event, events, t.uuidProvider)
	} else {
		day.Events = append(day.Events, event)
	}

	_, err = t.dayStore.Update(day)
	return event.ID, err
}

func (t *EventService) Delete(eventID string, userID string) error {
	day, err := t.dayStore.FindByEventID(eventID, userID)
	if err != nil {
		return err
	}

	day.Events = deleteEventFromSlice(day.Events, eventID)
	_, err = t.dayStore.Update(day)
	return err
}

func (e *EventService) GetUsedNames(userID string, tagID string, from time.Time, to time.Time) ([]string, error) {
	return e.dayStore.GetUsedNames(userID, tagID, from, to)
}

func deleteEventFromSlice(events []repository.Event, eventID string) []repository.Event {
	result := make([]repository.Event, 0)
	for i := range events {
		if events[i].ID != eventID {
			result = append(result, events[i])
		}
	}

	return result
}
