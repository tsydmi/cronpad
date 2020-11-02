package service

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/ts-dmitry/cronpad/backend/repository"
	"testing"
	"time"
)

type mockedUuidProvider struct {
}

func (m mockedUuidProvider) New() string {
	return "new-uuid"
}

func TestSortEventsByStartDate(t *testing.T) {
	currentTime := time.Now()

	t.Run("Empty array", func(t *testing.T) {
		notSortedArray := []repository.Event{}
		sortedArray := []repository.Event{}

		result := SortEventsByStartDate(notSortedArray)

		assert.Len(t, result, 0, "should have the same size")
		for i := range result {
			assert.True(t, result[i].Start.Equal(sortedArray[i].Start), "should sort elements")
		}
	})

	t.Run("Array with a single element", func(t *testing.T) {
		notSortedArray := []repository.Event{repository.Event{Start: currentTime}}
		sortedArray := []repository.Event{repository.Event{Start: currentTime}}

		result := SortEventsByStartDate(notSortedArray)

		assert.Len(t, result, 1, "should have the same size")
		for i := range result {
			assert.True(t, result[i].Start.Equal(sortedArray[i].Start), "should sort elements")
		}
	})

	t.Run("Array with multiple element", func(t *testing.T) {
		notSortedArray := []repository.Event{
			repository.Event{Start: currentTime.Add(3 * time.Minute)},
			repository.Event{Start: currentTime.Add(1 * time.Minute)},
			repository.Event{Start: currentTime.Add(2 * time.Minute)},
		}
		sortedArray := []repository.Event{
			repository.Event{Start: currentTime.Add(1 * time.Minute)},
			repository.Event{Start: currentTime.Add(2 * time.Minute)},
			repository.Event{Start: currentTime.Add(3 * time.Minute)},
		}

		result := SortEventsByStartDate(notSortedArray)

		assert.Len(t, result, 3, "should have the same size")
		for i := range result {
			assert.Equal(t, sortedArray[i].Start, result[i].Start, "should sort elements")
		}
	})
}

func TestNormalizeEvents(t *testing.T) {
	currentTime := time.Now()

	t.Run("Empty array", func(t *testing.T) {
		event := repository.Event{
			Start: currentTime,
			End:   currentTime.Add(5 * time.Hour),
			Name:  "New Event",
		}

		emptyArray := []repository.Event{}

		result := AddEventProperly(event, emptyArray, mockedUuidProvider{})

		assert.Len(t, result, 1, "should have a single element")
		for i := range result {
			assert.True(t, result[i].Start.Equal(event.Start), "should keep start date unchanged")
			assert.True(t, result[i].End.Equal(event.End), "should keep end date unchanged")
			assert.True(t, result[i].Name == event.Name, "should keep name unchanged")
		}
	})

	t.Run("Existed events doesn't have conflicts with a new one", func(t *testing.T) {
		event := repository.Event{
			ID:    uuid.New().String(),
			Start: currentTime.Add(4 * time.Hour),
			End:   currentTime.Add(5 * time.Hour),
			Name:  "New Event",
		}

		after := repository.Event{
			ID:    uuid.New().String(),
			Start: currentTime.Add(6 * time.Hour),
			End:   currentTime.Add(7 * time.Hour),
			Name:  "After",
		}

		before := repository.Event{
			ID:    uuid.New().String(),
			Start: currentTime.Add(1 * time.Hour),
			End:   currentTime.Add(2 * time.Hour),
			Name:  "Before",
		}

		result := AddEventProperly(event, []repository.Event{after, before}, mockedUuidProvider{})

		assert.Len(t, result, 3, "should have all elements")

		expectedResults := []repository.Event{before, event, after}
		for i := range result {
			assert.Equal(t, expectedResults[i].Start, result[i].Start, "should keep start date unchanged")
			assert.Equal(t, expectedResults[i].End, result[i].End, "should keep end date unchanged")
			assert.Equal(t, expectedResults[i].Name, result[i].Name, "should keep name unchanged")
		}
	})

	t.Run("All existed events have conflicts with a new one", func(t *testing.T) {
		// scheme: [--(--]--[--]--[--][--)--]

		event := repository.Event{
			ID:    uuid.New().String(),
			Name:  "New Event",
			Start: currentTime.Add(2 * time.Hour),
			End:   currentTime.Add(8 * time.Hour),
		}

		endInside := repository.Event{
			ID:    uuid.New().String(),
			Name:  "End Inside",
			Start: currentTime.Add(1 * time.Hour),
			End:   currentTime.Add(3 * time.Hour),
		}

		wholeEventInside := repository.Event{
			ID:    uuid.New().String(),
			Name:  "Whole Event Inside",
			Start: currentTime.Add(4 * time.Hour),
			End:   currentTime.Add(5 * time.Hour),
		}

		connectedToNext := repository.Event{
			ID:    uuid.New().String(),
			Name:  "Connected to the next",
			Start: currentTime.Add(6 * time.Hour),
			End:   currentTime.Add(7 * time.Hour),
		}

		startInside := repository.Event{
			ID:    uuid.New().String(),
			Name:  "Start Inside",
			Start: currentTime.Add(7 * time.Hour),
			End:   currentTime.Add(10 * time.Hour),
		}

		result := AddEventProperly(event, []repository.Event{wholeEventInside, endInside, connectedToNext, startInside}, mockedUuidProvider{})

		assert.Len(t, result, 6, "should have all elements")

		expectedResults := []repository.Event{
			endInside,
			repository.Event{
				ID:    event.ID,
				Name:  event.Name,
				Start: currentTime.Add(3 * time.Hour),
				End:   currentTime.Add(4 * time.Hour),
			},
			wholeEventInside,
			repository.Event{
				ID:    "new-uuid",
				Name:  event.Name,
				Start: currentTime.Add(5 * time.Hour),
				End:   currentTime.Add(6 * time.Hour),
			},
			connectedToNext,
			startInside,
		}
		for i := range result {
			assert.Exactly(t, expectedResults[i], result[i], "should return correct result")
		}
	})
}

func TestCreateEventBlocks(t *testing.T) {
	t.Run("Single event", func(t *testing.T) {
		currentTime := time.Now()

		event := repository.Event{
			Start: currentTime.Add(1 * time.Hour),
			End:   currentTime.Add(1 * time.Hour),
			Name:  "Event",
		}

		result := createEventBlocks([]repository.Event{event}, repository.Event{})

		assert.Len(t, result, 1, "should have all elements")

		eventBlock := result[0]
		assert.Equal(t, event.Start, eventBlock.start, "should have the same start date as event")
		assert.Equal(t, event.End, eventBlock.end, "should have the same end date as event")
		assert.Len(t, eventBlock.events, 1, "should have all events")
		assert.Exactly(t, event, eventBlock.events[0], "should have event inside event block")
	})

	t.Run("Not connected events", func(t *testing.T) {
		currentTime := time.Now()

		after := repository.Event{
			ID:    uuid.New().String(),
			Start: currentTime.Add(1 * time.Hour),
			End:   currentTime.Add(1 * time.Hour),
			Name:  "After",
		}

		before := repository.Event{
			ID:    uuid.New().String(),
			Start: currentTime.Add(-2 * time.Hour),
			End:   currentTime.Add(-2 * time.Hour),
			Name:  "Before",
		}

		result := createEventBlocks([]repository.Event{after, before}, repository.Event{})

		assert.Len(t, result, 2, "should have all elements")

		firstEventBlock := result[0]
		assert.Equal(t, before.Start, firstEventBlock.start, "should have the same start date as event")
		assert.Equal(t, before.End, firstEventBlock.end, "should have the same end date as event")
		assert.Len(t, firstEventBlock.events, 1, "should have all events")
		assert.Exactly(t, before, firstEventBlock.events[0], "should have event inside event block")

		secondEventBlock := result[1]
		assert.Equal(t, after.Start, secondEventBlock.start, "should have the same start date as event")
		assert.Equal(t, after.End, secondEventBlock.end, "should have the same end date as event")
		assert.Len(t, secondEventBlock.events, 1, "should have all events")
		assert.Exactly(t, after, secondEventBlock.events[0], "should have event inside event block")
	})

	t.Run("Multiple connected events and a single disconnected", func(t *testing.T) {
		currentTime := time.Now()

		firstConnected := repository.Event{
			ID:    uuid.New().String(),
			Start: currentTime.Add(1 * time.Hour),
			End:   currentTime.Add(2 * time.Hour),
			Name:  "First",
		}

		secondConnected := repository.Event{
			ID:    uuid.New().String(),
			Start: currentTime.Add(2 * time.Hour),
			End:   currentTime.Add(5 * time.Hour),
			Name:  "Second",
		}

		thirdConnected := repository.Event{
			ID:    uuid.New().String(),
			Start: currentTime.Add(5 * time.Hour),
			End:   currentTime.Add(8 * time.Hour),
			Name:  "Third",
		}

		disconnectedEvent := repository.Event{
			ID:    uuid.New().String(),
			Start: currentTime.Add(8 * time.Hour).Add(1 * time.Second),
			End:   currentTime.Add(9 * time.Hour),
			Name:  "Disconnected",
		}

		result := createEventBlocks([]repository.Event{thirdConnected, disconnectedEvent, firstConnected, secondConnected}, repository.Event{})

		assert.Len(t, result, 2, "should have all elements")

		firstBlock := result[0]
		assert.Equal(t, firstConnected.Start, firstBlock.start, "should have the same start date as firstConnected event")
		assert.Equal(t, thirdConnected.End, firstBlock.end, "should have the same end date as thirdConnected event")
		assert.Len(t, firstBlock.events, 3, "should have all events")
		assert.Exactly(t, firstConnected, firstBlock.events[0], "should have firstConnected event inside event block")
		assert.Exactly(t, secondConnected, firstBlock.events[1], "should have secondConnected event inside event block")
		assert.Exactly(t, thirdConnected, firstBlock.events[2], "should have thirdConnected event inside event block")

		secondBlock := result[1]
		assert.Equal(t, disconnectedEvent.Start, secondBlock.start, "should have the same start date as firstConnected event")
		assert.Equal(t, disconnectedEvent.End, secondBlock.end, "should have the same end date as thirdConnected event")
		assert.Len(t, secondBlock.events, 1, "should have all events")
		assert.Exactly(t, disconnectedEvent, secondBlock.events[0], "should have disconnectedEvent event inside event block")
	})
}
