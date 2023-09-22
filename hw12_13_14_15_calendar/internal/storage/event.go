package storage

import (
	"errors"
	"time"
)

var (
	ErrEventAlreadyExists = errors.New("err event already exists")
	ErrEventNotFound      = errors.New("err event not found")
)

type Event struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	DateTime    time.Time `json:"date_time"`
	Duration    int64     `json:"duration"`
	Description string    `json:"description"`
	UserID      string    `json:"user_id"`
	Reminder    int64     `json:"reminder"`
}

type Notification struct {
	EventID  string    `json:"event_id"`
	Title    string    `json:"title"`
	DateTime time.Time `json:"date_time"`
	UserID   string    `json:"user_id"`
}
