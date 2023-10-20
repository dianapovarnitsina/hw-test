package interfaces

import (
	"context"
	"time"

	"github.com/dianapovarnitsina/hw-test/hw12_13_14_15_calendar/internal/storage"
)

type EventStorage interface {
	Connect(ctx context.Context, dbPort int, dbHost, dbUser, dbPassword, dbName string) error
	Close() error
	Migrate(ctx context.Context, migrate string) error
	CreateEvent(ctx context.Context, event *storage.Event) error
	UpdateEvent(ctx context.Context, event *storage.Event) error
	DeleteEvent(ctx context.Context, eventID string) error
	GetEvent(ctx context.Context, eventID string) (*storage.Event, error)
	ListEventsForDay(ctx context.Context, startOfDay time.Time) ([]*storage.Event, error)
	ListEventsForWeek(ctx context.Context, startOfWeek time.Time) ([]*storage.Event, error)
	ListEventsForMonth(ctx context.Context, startOfMonth time.Time) ([]*storage.Event, error)
	SelectEventsForNotifications(ctx context.Context) ([]*storage.Event, error)
}
