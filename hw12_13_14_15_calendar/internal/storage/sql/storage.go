package sqlstorage

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/dianapovarnitsina/hw-test/hw12_13_14_15_calendar/internal/config"
	"github.com/dianapovarnitsina/hw-test/hw12_13_14_15_calendar/internal/storage"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

type Storage struct {
	db *sql.DB
}

func (r *Storage) Migrate(ctx context.Context, migrate string) (err error) {

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("cannot set dialect: %w", err)
	}

	if err := goose.Up(r.db, migrate); err != nil {
		return fmt.Errorf("cannot do up migration: %w", err)
	}

	return nil
}

func (r *Storage) Connect(ctx context.Context, conf *config.Config) (err error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		conf.Database.Host, conf.Database.Port, conf.Database.Username, conf.Database.Password, conf.Database.Dbname)
	r.db, err = sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("cannot open pgx driver: %w", err)
	}

	return r.db.PingContext(ctx)
}

func (s *Storage) Close() error {
	return s.db.Close()
}

func (s *Storage) CreateEvent(ctx context.Context, event *storage.Event) error {
	query := `
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
	query := `
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
	query := `
		DELETE FROM events
		WHERE id = $1
	`

	_, err := s.db.ExecContext(ctx, query, eventID)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) GetEvent(ctx context.Context, eventID string) (*storage.Event, error) {
	query := `
		SELECT id, title, date_time, duration, description, user_id, reminder
		FROM events
		WHERE id = $1
	`

	row := s.db.QueryRowContext(ctx, query, eventID)

	event := &storage.Event{}

	// Сканируем значения из результата запроса в объект события
	err := row.Scan(&event.ID, &event.Title, &event.DateTime, &event.Duration, &event.Description, &event.UserID, &event.Reminder)
	if err != nil {
		return nil, err
	}

	return event, nil
}

func (s *Storage) ListEventsForDay(ctx context.Context, date time.Time) ([]*storage.Event, error) {
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	query := `
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

	// Итерируемся по результатам запроса и создаем объекты событий
	for rows.Next() {
		event := &storage.Event{}
		err := rows.Scan(&event.ID, &event.Title, &event.DateTime, &event.Duration, &event.Description, &event.UserID, &event.Reminder)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	return events, nil
}

func (s *Storage) ListEventsForWeek(ctx context.Context, startOfWeek time.Time) ([]*storage.Event, error) {
	startOfWeek = time.Date(startOfWeek.Year(), startOfWeek.Month(), startOfWeek.Day(), 0, 0, 0, 0, startOfWeek.Location())
	endOfWeek := startOfWeek.AddDate(0, 0, 6)

	query := `
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

	// Итерируемся по результатам запроса и создаем объекты событий
	for rows.Next() {
		event := &storage.Event{}
		err := rows.Scan(&event.ID, &event.Title, &event.DateTime, &event.Duration, &event.Description, &event.UserID, &event.Reminder)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	return events, nil
}

func (s *Storage) ListEventsForMonth(ctx context.Context, startOfMonth time.Time) ([]*storage.Event, error) {
	startOfMonth = time.Date(startOfMonth.Year(), startOfMonth.Month(), startOfMonth.Day(), 0, 0, 0, 0, startOfMonth.Location())
	endOfMonth := startOfMonth.AddDate(0, 1, 0).Add(-time.Nanosecond)

	query := `
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

	// Итерируемся по результатам запроса и создаем объекты событий
	for rows.Next() {
		event := &storage.Event{}
		err := rows.Scan(&event.ID, &event.Title, &event.DateTime, &event.Duration, &event.Description, &event.UserID, &event.Reminder)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	return events, nil
}
