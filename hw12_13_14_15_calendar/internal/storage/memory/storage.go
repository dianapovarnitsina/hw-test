package memorystorage

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/dianapovarnitsina/hw-test/hw12_13_14_15_calendar/internal/config"
	"github.com/dianapovarnitsina/hw-test/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct {
	events map[string]*storage.Event
	mu     sync.RWMutex
}

func New() *Storage {
	return &Storage{
		events: make(map[string]*storage.Event),
	}
}

func (s *Storage) Connect(ctx context.Context, conf *config.Config) error {
	_, _ = ctx, conf
	return nil
}

func (s *Storage) Close() error {
	return nil
}

func (s *Storage) Migrate(ctx context.Context, migrate string) (err error) {
	_, _ = ctx, migrate
	return nil
}

func (s *Storage) CreateEvent(ctx context.Context, event *storage.Event) error {
	_ = ctx
	s.mu.Lock()
	defer s.mu.Unlock()

	// Проверка наличия события с таким ID
	if _, found := s.events[event.ID]; found {
		return storage.ErrEventAlreadyExists
	}

	// Добавляем событие
	s.events[event.ID] = event
	return nil
}

func (s *Storage) UpdateEvent(ctx context.Context, event *storage.Event) error {
	_ = ctx
	s.mu.Lock()
	defer s.mu.Unlock()

	// Проверка наличия события с таким ID
	if _, found := s.events[event.ID]; !found {
		return storage.ErrEventNotFound
	}

	// Обновляем событие
	s.events[event.ID] = event
	return nil
}

func (s *Storage) DeleteEvent(ctx context.Context, eventID string) error {
	_ = ctx
	s.mu.Lock()
	defer s.mu.Unlock()

	// Проверка наличия события с таким ID
	if _, found := s.events[eventID]; !found {
		return storage.ErrEventNotFound
	}

	// Удаляем событие
	delete(s.events, eventID)
	return nil
}

func (s *Storage) GetEvent(ctx context.Context, eventID string) (*storage.Event, error) {
	_ = ctx
	s.mu.RLock()
	defer s.mu.RUnlock()

	event, found := s.events[eventID]
	if !found {
		return nil, errors.New("err event not found")
	}

	return event, nil
}

func (s *Storage) ListEventsForDay(ctx context.Context, day time.Time) ([]*storage.Event, error) {
	_ = ctx
	s.mu.RLock()
	defer s.mu.RUnlock()

	var eventsForDay []*storage.Event

	for _, event := range s.events {
		if event.DateTime.Year() == day.Year() &&
			event.DateTime.Month() == day.Month() &&
			event.DateTime.Day() == day.Day() {
			eventsForDay = append(eventsForDay, event)
		}
	}

	return eventsForDay, nil
}

func (s *Storage) ListEventsForWeek(ctx context.Context, startOfWeek time.Time) ([]*storage.Event, error) {
	_ = ctx
	s.mu.RLock()
	defer s.mu.RUnlock()

	var eventsForWeek []*storage.Event

	for _, event := range s.events {
		if event.DateTime.After(startOfWeek) &&
			event.DateTime.Before(startOfWeek.AddDate(0, 0, 7)) {
			eventsForWeek = append(eventsForWeek, event)
		}
	}

	return eventsForWeek, nil
}

func (s *Storage) ListEventsForMonth(ctx context.Context, startOfMonth time.Time) ([]*storage.Event, error) {
	_ = ctx
	s.mu.RLock()
	defer s.mu.RUnlock()

	var eventsForMonth []*storage.Event

	for _, event := range s.events {
		if event.DateTime.Month() == startOfMonth.Month() {
			eventsForMonth = append(eventsForMonth, event)
		}
	}

	return eventsForMonth, nil
}
