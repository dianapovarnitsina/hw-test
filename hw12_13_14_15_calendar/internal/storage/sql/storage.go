package sqlstorage

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/dianapovarnitsina/hw-test/hw12_13_14_15_calendar/internal/storage"
	_ "github.com/lib/pq" // Blank import for side effects
	"github.com/pressly/goose/v3"
)

type Storage struct {
	db *sql.DB
}

func (s *Storage) Migrate(ctx context.Context, migrate string) (err error) {
	_ = ctx
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("cannot set dialect: %w", err)
	}

	if err := goose.Up(s.db, migrate); err != nil {
		return fmt.Errorf("cannot do up migration: %w", err)
	}

	return nil
}

func (s *Storage) Connect(ctx context.Context, dbPort int, dbHost, dbUser, dbPassword, dbName string) (err error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		dbUser, dbPassword, dbHost, dbPort, dbName)
	s.db, err = sql.Open("postgres", connStr)

	if err != nil {
		return fmt.Errorf("cannot open pgx driver: %w", err)
	}

	return s.db.PingContext(ctx)
}

func (s *Storage) Close() error {
	return s.db.Close()
}

func (s *Storage) CreateEvent(ctx context.Context, event *storage.Event) error {
	const query = `
		INSERT INTO events (id, title, date_time, duration, description, user_id, reminder)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := s.db.ExecContext(
		ctx,
		query,
		event.ID,
		event.Title,
		event.DateTime,
		event.Duration,
		event.Description,
		event.UserID,
		event.Reminder,
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) UpdateEvent(ctx context.Context, event *storage.Event) error {
	const query = `
		UPDATE events
		SET title = $2, date_time = $3, duration = $4, description = $5, user_id = $6, reminder = $7
		WHERE id = $1
	`

	_, err := s.db.ExecContext(
		ctx,
		query,
		event.ID,
		event.Title,
		event.DateTime,
		event.Duration,
		event.Description,
		event.UserID,
		event.Reminder,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) DeleteEvent(ctx context.Context, eventID string) error {
	const query = `DELETE FROM events WHERE id = $1`

	_, err := s.db.ExecContext(ctx, query, eventID)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) GetEvent(ctx context.Context, eventID string) (*storage.Event, error) {
	const query = `
		SELECT id, title, date_time, duration, description, user_id, reminder
		FROM events
		WHERE id = $1
	`

	row := s.db.QueryRowContext(ctx, query, eventID)

	event := &storage.Event{}

	// Scanning values from the query result to the event object
	err := row.Scan(
		&event.ID,
		&event.Title,
		&event.DateTime,
		&event.Duration,
		&event.Description,
		&event.UserID,
		&event.Reminder,
	)
	if err != nil {
		return nil, err
	}

	return event, nil
}

func (s *Storage) ListEventsForDay(ctx context.Context, day time.Time) ([]*storage.Event, error) {
	startOfDay := day.Add(-time.Second)
	endOfDay := startOfDay.Add(24 * time.Hour)

	const query = `
		SELECT id, title, date_time, duration, description, user_id, reminder
		FROM events
		WHERE date_time >= $1 AND date_time < $2
	`

	rows, err := s.db.QueryContext(ctx, query, startOfDay, endOfDay)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	events := []*storage.Event{}

	// Iterate on the results of the query and create event objects
	for rows.Next() {
		event := &storage.Event{}
		err := rows.Scan(
			&event.ID,
			&event.Title,
			&event.DateTime,
			&event.Duration,
			&event.Description,
			&event.UserID,
			&event.Reminder,
		)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	return events, nil
}

func (s *Storage) ListEventsForWeek(ctx context.Context, startOfWeek time.Time) ([]*storage.Event, error) {
	startOfWeek = startOfWeek.Add(-time.Second)
	endOfWeek := startOfWeek.AddDate(0, 0, 6)

	const query = `
		SELECT id, title, date_time, duration, description, user_id, reminder
		FROM events
		WHERE date_time >= $1 AND date_time < $2
	`

	rows, err := s.db.QueryContext(ctx, query, startOfWeek, endOfWeek)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	events := []*storage.Event{}

	// Iterate on the results of the query and create event objects
	for rows.Next() {
		event := &storage.Event{}
		err := rows.Scan(
			&event.ID,
			&event.Title,
			&event.DateTime,
			&event.Duration,
			&event.Description,
			&event.UserID,
			&event.Reminder,
		)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	return events, nil
}

func (s *Storage) ListEventsForMonth(ctx context.Context, startOfMonth time.Time) ([]*storage.Event, error) {
	startOfMonth = startOfMonth.Add(-time.Second)
	endOfMonth := startOfMonth.AddDate(0, 1, 0).Add(-time.Nanosecond)

	const query = `
		SELECT id, title, date_time, duration, description, user_id, reminder
		FROM events
		WHERE date_time >= $1 AND date_time < $2
	`

	rows, err := s.db.QueryContext(ctx, query, startOfMonth, endOfMonth)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	events := []*storage.Event{}

	// Iterate on the results of the query and create event objects
	for rows.Next() {
		event := &storage.Event{}
		err := rows.Scan(
			&event.ID,
			&event.Title,
			&event.DateTime,
			&event.Duration,
			&event.Description,
			&event.UserID,
			&event.Reminder,
		)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	return events, nil
}

func (s *Storage) SelectEventsForNotifications(ctx context.Context) ([]*storage.Event, error) {
	const query = `
		SELECT id, title, date_time, duration, description, user_id, reminder
		FROM events
		WHERE date_trunc('minute', date_time - (reminder * INTERVAL '1 minute')) = date_trunc('minute', NOW());
	`

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []*storage.Event

	// Итерируйтесь по результатам запроса и создайте объекты событий.
	for rows.Next() {
		event := &storage.Event{}
		err := rows.Scan(
			&event.ID,
			&event.Title,
			&event.DateTime,
			&event.Duration,
			&event.Description,
			&event.UserID,
			&event.Reminder,
		)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	return events, nil
}

func (s *Storage) DeleteOldEvents(ctx context.Context) error {
	const query = `DELETE FROM events WHERE date_time < NOW() - INTERVAL '1 year';`

	_, err := s.db.ExecContext(ctx, query)
	if err != nil {
		return err
	}

	return nil
}
