package interfaces

import (
	"context"
	"time"

	"github.com/dianapovarnitsina/hw-test/hw12_13_14_15_calendar/internal/config"
	"github.com/dianapovarnitsina/hw-test/hw12_13_14_15_calendar/internal/storage"
)

type EventStorage interface {
	Connect(ctx context.Context, conf *config.Config) error
	Close() error
	Migrate(ctx context.Context, migrate string) error
	CreateEvent(ctx context.Context, event *storage.Event) error
	UpdateEvent(ctx context.Context, event *storage.Event) error
	DeleteEvent(ctx context.Context, eventID string) error
	GetEvent(ctx context.Context, eventID string) (*storage.Event, error)
	ListEventsForDay(ctx context.Context, date time.Time) ([]*storage.Event, error)
	ListEventsForWeek(ctx context.Context, startOfWeek time.Time) ([]*storage.Event, error)
	ListEventsForMonth(ctx context.Context, startOfMonth time.Time) ([]*storage.Event, error)
}