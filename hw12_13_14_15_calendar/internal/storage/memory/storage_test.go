package memorystorage

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"strconv"
	"sync"
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

func TestConcurrentCreateEvent(t *testing.T) {
	storage := New()
	var wg sync.WaitGroup
	numOperations := 10

	for i := 0; i < numOperations; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			event := &stor.Event{
				ID:       strconv.Itoa(i),
				Title:    fmt.Sprintf("Event %d", i),
				DateTime: time.Now(),
			}
			err := storage.CreateEvent(context.Background(), event)
			assert.NoError(t, err)
		}(i)
	}

	wg.Wait()

	// Check that all events were created
	for i := 0; i < numOperations; i++ {
		eventID := strconv.Itoa(i)
		event, err := storage.GetEvent(context.Background(), eventID)
		assert.NoError(t, err)
		assert.NotNil(t, event)
	}
}

func TestConcurrentUpdateEvent(t *testing.T) {
	storage := New()
	event := &stor.Event{
		ID:       "1",
		Title:    "Initial Title",
		DateTime: time.Now(),
	}
	err := storage.CreateEvent(context.Background(), event)
	assert.NoError(t, err)

	var wg sync.WaitGroup
	numOperations := 10

	for i := 0; i < numOperations; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			event := &stor.Event{
				ID:       "1",
				Title:    fmt.Sprintf("Updated Title %d", i),
				DateTime: time.Now(),
			}
			err := storage.UpdateEvent(context.Background(), event)
			assert.NoError(t, err)
		}(i)
	}

	wg.Wait()

	// Check that the event was updated
	updatedEvent, err := storage.GetEvent(context.Background(), "1")
	assert.NoError(t, err)
	assert.NotNil(t, updatedEvent)
	assert.Contains(t, updatedEvent.Title, "Updated Title")
}

func TestConcurrentDeleteEvent(t *testing.T) {
	storage := New()
	event := &stor.Event{
		ID:       "1",
		Title:    "Event to delete",
		DateTime: time.Now(),
	}
	err := storage.CreateEvent(context.Background(), event)
	assert.NoError(t, err)

	// Try to delete the event concurrently
	var wg sync.WaitGroup
	numOperations := 10

	// Channel to collect errors
	errCh := make(chan error, numOperations)

	for i := 0; i < numOperations; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := storage.DeleteEvent(context.Background(), "1")
			errCh <- err
		}()
	}

	wg.Wait()
	close(errCh)

	// Count successful deletions
	successCount := 0

	for err := range errCh {
		if err == nil {
			successCount++
		}
	}

	// Only one deletion should succeed, the others should return error
	assert.Equal(t, 1, successCount)

	// Check that the event was deleted only once
	deletedEvent, err := storage.GetEvent(context.Background(), "1")
	assert.Error(t, err) // Event should not exist
	assert.Nil(t, deletedEvent)
}

func TestConcurrentListEventsForDay(t *testing.T) {
	storage := New()

	// Create some events
	event1 := &stor.Event{ID: "1", DateTime: time.Now()}
	event2 := &stor.Event{ID: "2", DateTime: time.Now()}
	_ = storage.CreateEvent(context.Background(), event1)
	_ = storage.CreateEvent(context.Background(), event2)

	// Concurrently list events for the day
	var wg sync.WaitGroup
	numOperations := 10

	// Channel to collect events
	eventsCh := make(chan []*stor.Event, numOperations)

	for i := 0; i < numOperations; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			events, _ := storage.ListEventsForDay(context.Background(), time.Now())
			eventsCh <- events
		}()
	}

	wg.Wait()
	close(eventsCh)

	// Collect events from the channel
	eventsList := [][]*stor.Event{}
	for events := range eventsCh {
		eventsList = append(eventsList, events)
	}

	// Compare the events
	for i := 1; i < len(eventsList); i++ {
		assert.ElementsMatch(t, eventsList[0], eventsList[i])
	}
}

func TestConcurrentListEventsForWeek(t *testing.T) {
	storage := New()

	// Create some events
	event1 := &stor.Event{ID: "1", DateTime: time.Now()}
	event2 := &stor.Event{ID: "2", DateTime: time.Now()}
	_ = storage.CreateEvent(context.Background(), event1)
	_ = storage.CreateEvent(context.Background(), event2)

	// Concurrently list events for the week
	var wg sync.WaitGroup
	numOperations := 10

	// Create a slice to store the results
	results := make([][]*stor.Event, numOperations)

	for i := 0; i < numOperations; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			events, _ := storage.ListEventsForWeek(context.Background(), time.Now())
			results[index] = events
		}(i)
	}

	wg.Wait()

	// Compare the events
	for i := 1; i < len(results); i++ {
		assert.ElementsMatch(t, results[0], results[i])
	}
}

func TestConcurrentListEventsForMonth(t *testing.T) {
	storage := New()

	// Create some events
	event1 := &stor.Event{ID: "1", DateTime: time.Now()}
	event2 := &stor.Event{ID: "2", DateTime: time.Now()}
	_ = storage.CreateEvent(context.Background(), event1)
	_ = storage.CreateEvent(context.Background(), event2)

	// Concurrently list events for the month
	var wg sync.WaitGroup
	numOperations := 10

	// Create a slice to store the results
	results := make([][]*stor.Event, numOperations)

	for i := 0; i < numOperations; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			events, _ := storage.ListEventsForMonth(context.Background(), time.Now())
			results[index] = events
		}(i)
	}

	wg.Wait()

	// Compare the events
	for i := 1; i < len(results); i++ {
		assert.ElementsMatch(t, results[0], results[i])
	}
}
