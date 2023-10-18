package calendar

import (
	"context"
	"time"

	"github.com/dianapovarnitsina/hw-test/hw12_13_14_15_calendar/interfaces"
	"github.com/dianapovarnitsina/hw-test/hw12_13_14_15_calendar/internal/storage"
)

type App struct {
	logger  interfaces.Logger
	storage interfaces.EventStorage
}

func NewApp(logger interfaces.Logger, storage interfaces.EventStorage) *App {
	return &App{
		logger:  logger,
		storage: storage,
	}
}

func (a *App) CreateEvent(ctx context.Context, event *storage.Event) error {
	return a.storage.CreateEvent(ctx, event)
}

func (a *App) UpdateEvent(ctx context.Context, event *storage.Event) error {
	return a.storage.UpdateEvent(ctx, event)
}

func (a *App) DeleteEvent(ctx context.Context, eventID string) error {
	return a.storage.DeleteEvent(ctx, eventID)
}

func (a *App) GetEvent(ctx context.Context, eventID string) (*storage.Event, error) {
	return a.storage.GetEvent(ctx, eventID)
}

func (a *App) ListEventsForDay(ctx context.Context, date time.Time) ([]*storage.Event, error) {
	return a.storage.ListEventsForDay(ctx, date)
}

func (a *App) ListEventsForWeek(ctx context.Context, date time.Time) ([]*storage.Event, error) {
	return a.storage.ListEventsForWeek(ctx, date)
}

func (a *App) ListEventsForMonth(ctx context.Context, date time.Time) ([]*storage.Event, error) {
	return a.storage.ListEventsForWeek(ctx, date)
}
