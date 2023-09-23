package memorystorage

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"

	stor "github.com/dianapovarnitsina/hw-test/hw12_13_14_15_calendar/internal/storage"
)

func TestCreateEvent(t *testing.T) {
	storage := New()

	event := &stor.Event{
		ID:       "event1",
		Title:    "Test Event",
		DateTime: time.Now(),
	}

	// Create an event
	err := storage.CreateEvent(context.Background(), event)
	assert.NoError(t, err)

	// Re-creating the event again, an error message should return
	err = storage.CreateEvent(context.Background(), event)
	assert.Equal(t, stor.ErrEventAlreadyExists, err)
}

func TestUpdateEvent(t *testing.T) {
	storage := New()
	event := &stor.Event{
		ID:       "event1",
		Title:    "Test Event",
		DateTime: time.Now(),
	}

	// Create an event
	err := storage.CreateEvent(context.Background(), event)
	assert.NoError(t, err)

	// Update the event
	event.Title = "Updated Event"
	err = storage.UpdateEvent(context.Background(), event)
	assert.NoError(t, err)

	// Event should be updated
	updatedEvent, err := storage.GetEvent(context.Background(), event.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Updated Event", updatedEvent.Title)

	// Try updating a non-existing event
	err = storage.UpdateEvent(context.Background(), &stor.Event{ID: "nonexistent"})
	assert.Equal(t, stor.ErrEventNotFound, err)
}

func TestDeleteEvent(t *testing.T) {
	storage := New()
	eventID := "event1"
	event := &stor.Event{
		ID:       eventID,
		Title:    "Test Event",
		DateTime: time.Now(),
	}

	// Create an event
	err := storage.CreateEvent(context.Background(), event)
	assert.NoError(t, err)

	// Delete the event
	err = storage.DeleteEvent(context.Background(), eventID)
	assert.NoError(t, err)

	// Try deleting a non-existing event
	err = storage.DeleteEvent(context.Background(), "nonexistent")
	assert.Equal(t, stor.ErrEventNotFound, err)
}

func TestGetEvent(t *testing.T) {
	storage := New()
	eventID := "event1"
	event := &stor.Event{
		ID:       eventID,
		Title:    "Test Event",
		DateTime: time.Now(),
	}

	// Create an event
	err := storage.CreateEvent(context.Background(), event)
	assert.NoError(t, err)

	// Get the event
	retrievedEvent, err := storage.GetEvent(context.Background(), eventID)
	assert.NoError(t, err)
	assert.Equal(t, event.Title, retrievedEvent.Title)

	// Try getting a non-existing event
	_, err = storage.GetEvent(context.Background(), "nonexistent")
	assert.Equal(t, stor.ErrEventNotFound, err)
}

func TestListEventsForDay(t *testing.T) {
	storage := New()
	day := time.Now()
	event1 := &stor.Event{
		ID:       "event1",
		Title:    "Event 1",
		DateTime: day,
	}
	event2 := &stor.Event{
		ID:       "event2",
		Title:    "Event 2",
		DateTime: day.Add(time.Hour * 2),
	}

	// Create events
	err := storage.CreateEvent(context.Background(), event1)
	assert.NoError(t, err)
	err = storage.CreateEvent(context.Background(), event2)
	assert.NoError(t, err)

	// List events for the given day
	events, err := storage.ListEventsForDay(context.Background(), day)
	assert.NoError(t, err)
	assert.Len(t, events, 2)
	assert.Equal(t, event1.ID, events[0].ID)

	// Try listing events for a day with no events
	noEventsDay := day.AddDate(0, 0, 1)
	noEvents, err := storage.ListEventsForDay(context.Background(), noEventsDay)
	assert.NoError(t, err)
	assert.Len(t, noEvents, 0)
}

func TestListEventsForWeek(t *testing.T) {
	storage := New()
	startOfWeek := time.Now().AddDate(0, 0, 1)
	event1 := &stor.Event{
		ID:       "event1",
		Title:    "Event 1",
		DateTime: startOfWeek.AddDate(0, 0, 3),
	}
	event2 := &stor.Event{
		ID:       "event2",
		Title:    "Event 2",
		DateTime: startOfWeek.AddDate(0, 0, 5),
	}

	err := storage.CreateEvent(context.Background(), event1)
	assert.NoError(t, err)
	err = storage.CreateEvent(context.Background(), event2)
	assert.NoError(t, err)

	// List events for the week
	events, err := storage.ListEventsForWeek(context.Background(), startOfWeek)
	assert.NoError(t, err)
	assert.Len(t, events, 2)

	// Try listing events for a week with no events
	noEventsWeek := startOfWeek.AddDate(0, 0, 6)
	noEvents, err := storage.ListEventsForWeek(context.Background(), noEventsWeek)
	assert.NoError(t, err)
	assert.Len(t, noEvents, 0)
}

func TestListEventsForMonth(t *testing.T) {
	storage := New()
	startOfMonth := time.Now()
	event1 := &stor.Event{
		ID:       "event1",
		Title:    "Event 1",
		DateTime: startOfMonth.AddDate(0, 0, 10),
	}
	event2 := &stor.Event{
		ID:       "event2",
		Title:    "Event 2",
		DateTime: startOfMonth.AddDate(0, 0, 20),
	}

	err := storage.CreateEvent(context.Background(), event1)
	assert.NoError(t, err)
	err = storage.CreateEvent(context.Background(), event2)
	assert.NoError(t, err)

	// List events for the month
	events, err := storage.ListEventsForMonth(context.Background(), startOfMonth)
	assert.NoError(t, err)
	assert.Len(t, events, 2)

	// Try listing events for a month with no events
	noEventsMonth := startOfMonth.AddDate(0, 0, 34)
	noEvents, err := storage.ListEventsForMonth(context.Background(), noEventsMonth)
	assert.NoError(t, err)
	assert.Len(t, noEvents, 0)
	fmt.Println()
}
