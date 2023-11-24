//go:build integration
// +build integration

package test_test

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/dianapovarnitsina/hw-test/hw12_13_14_15_calendar/internal/server/pb"
	"github.com/golang/protobuf/ptypes/timestamp"
	_ "github.com/lib/pq" // Blank import for side effects
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"os"
	"testing"
	"time"
)

const (
	querySelectById = `
		SELECT id, title, date_time, duration, description, user_id, reminder
		FROM events
		WHERE id = $1
	`
	queryInsert = `
		INSERT INTO events (id, title, date_time, duration, description, user_id, reminder)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
)

type CalendarSuite struct {
	suite.Suite
	ctx          context.Context
	calendarConn *grpc.ClientConn
	client       pb.CalendarServiceClient
	db           *sql.DB
}

func (s *CalendarSuite) SetupSuite() { //общая настройка для всего сьюта
	s.ctx = context.TODO()

	host := os.Getenv("GRPC_HOST")
	port := os.Getenv("GRPC_PORT")
	calendarHost := host + ":" + port
	//calendarHost := ""

	if calendarHost == "" {
		calendarHost = "127.0.0.1:8082"
	}
	var err error
	s.calendarConn, err = grpc.Dial(calendarHost, grpc.WithInsecure())
	s.Require().NoError(err)
	s.client = pb.NewCalendarServiceClient(s.calendarConn)

	connectionString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		//"postgres", "postgres", "localhost", 5432, "postgres")
		"postgres", "postgres", os.Getenv("POSTGRES_HOST"), 5432, "postgres")
	s.db, err = sql.Open("postgres", connectionString)
	s.Require().NoError(err)
}

func (s *CalendarSuite) SetupTest() { // настройка для конкретного теста

}

func (s *CalendarSuite) TearDownTest() { //очистка для конкретного теста
	query := `DELETE FROM events`
	_, err := s.db.Exec(query)
	s.Require().NoError(err)
}

func (s *CalendarSuite) TearDownSuite() { // общая очистка после завершения сьюта
	defer s.db.Close()
}

func TestCalendarPost(t *testing.T) {
	suite.Run(t, new(CalendarSuite))
}

func (s *CalendarSuite) TestCalendar_GetEvent() {
	_, err := s.db.Exec(queryInsert, "1", "Event 1", time.Now().Format("2006-01-02 15:04:05"), 1, "Description 1", "1", 1)
	s.Require().NoError(err)

	request := &pb.GetEventRequest{EventId: "1"}
	response, err := s.client.GetEvent(s.ctx, request)
	s.Require().NoError(err)
	s.Require().NotNil(response, "Expected a non-nil createdEvent")
	s.Equal(response.Event.Id, "1", "Event ID mismatch")
	s.Equal(response.Event.Title, "Event 1", "Event Title mismatch")
	s.Equal(response.Event.Duration, int64(1), "Event Duration mismatch")
	s.Equal(response.Event.Description, "Description 1", "Event Description mismatch")
	s.Equal(response.Event.UserId, "1", "Event UserId mismatch")
	s.Equal(response.Event.Reminder, int64(1), "Event Reminder mismatch")
	s.Nil(err, "Expected no error, but got: %v", err)
}

func (s *CalendarSuite) TestCalendar_CreateEvent() {
	request := &pb.EventRequest{
		Event: &pb.Event{
			Id:          "1",
			Title:       "Sample Event",
			DateTime:    &timestamp.Timestamp{Seconds: time.Now().Unix()},
			Duration:    60,
			Description: "Description of the event",
			UserId:      "1",
			Reminder:    15,
		},
	}
	response, err := s.client.CreateEvent(s.ctx, request)
	s.Require().NoError(err)
	s.Require().NotNil(response, "Expected a non-nil createdEvent")
	s.Equal(response.Event.Id, request.Event.Id, "Event ID mismatch")
	s.Equal(response.Event.Title, request.Event.Title, "Event Title mismatch")
	s.Equal(response.Event.Duration, request.Event.Duration, "Event Duration mismatch")
	s.Equal(response.Event.Description, request.Event.Description, "Event Description mismatch")
	s.Equal(response.Event.UserId, request.Event.UserId, "Event UserId mismatch")
	s.Equal(response.Event.Reminder, request.Event.Reminder, "Event Reminder mismatch")
	s.Nil(err, "Expected no error, but got: %v", err)

	row := s.db.QueryRow(querySelectById, request.Event.Id)
	event := &pb.Event{}
	var datetime time.Time
	err = row.Scan(&event.Id, &event.Title, &datetime, &event.Duration, &event.Description, &event.UserId, &event.Reminder)
	s.Require().NoError(err)
	event.DateTime = &timestamp.Timestamp{
		Seconds: datetime.Unix(),
		Nanos:   int32(datetime.Nanosecond()),
	}
	s.Require().NoError(err)
	s.Require().NotNil(response, "Expected a non-nil createdEvent")
	s.Equal(response.Event.Id, event.Id, "Event ID mismatch")
	s.Equal(response.Event.Title, event.Title, "Event Title mismatch")
	s.Equal(response.Event.Duration, event.Duration, "Event Duration mismatch")
	s.Equal(response.Event.Description, event.Description, "Event Description mismatch")
	s.Equal(response.Event.UserId, event.UserId, "Event UserId mismatch")
	s.Equal(response.Event.Reminder, event.Reminder, "Event Reminder mismatch")
	s.Nil(err, "Expected no error, but got: %v", err)
}

func (s *CalendarSuite) TestCalendar_CreateEvent_Error() {
	// Попытка создать событие с недопустимыми данными
	request := &pb.EventRequest{
		Event: &pb.Event{
			Id:          "2",
			Title:       "", // Пустой заголовок
			DateTime:    &timestamp.Timestamp{Seconds: time.Now().Unix()},
			Duration:    60,
			Description: "Description of the event",
			UserId:      "1",
			Reminder:    15,
		},
	}
	response, err := s.client.CreateEvent(s.ctx, request)

	// Проверка на наличие ошибки
	s.Error(err)
	s.Nil(response)
}

func (s *CalendarSuite) TestCalendar_UpdateEvent() {
	request := &pb.EventRequest{
		Event: &pb.Event{
			Id:          "1",
			Title:       "Event 1",
			DateTime:    &timestamp.Timestamp{Seconds: time.Now().Unix()},
			Duration:    1,
			Description: "Description 1",
			UserId:      "1",
			Reminder:    1,
		},
	}
	eventTimeStr := time.Unix(request.Event.DateTime.Seconds, 0).Format("2006-01-02 15:04:05")

	_, err := s.db.Exec(
		queryInsert,
		request.Event.Id,
		request.Event.Title,
		eventTimeStr,
		request.Event.Duration,
		request.Event.Description,
		request.Event.UserId,
		request.Event.Reminder,
	)
	s.Require().NoError(err)

	request.Event.Title = "Event 2"

	response, err := s.client.UpdateEvent(s.ctx, request)
	s.Require().NoError(err)
	s.Require().NotNil(response, "Expected a non-nil createdEvent")
	s.Equal(response.Event.Id, response.Event.Id, "Event ID mismatch")
	s.Equal(response.Event.Title, response.Event.Title, "Event Title mismatch")
	s.Equal(response.Event.Duration, response.Event.Duration, "Event Duration mismatch")
	s.Equal(response.Event.Description, response.Event.Description, "Event Description mismatch")
	s.Equal(response.Event.UserId, response.Event.UserId, "Event UserId mismatch")
	s.Equal(response.Event.Reminder, response.Event.Reminder, "Event Reminder mismatch")
	s.Nil(err, "Expected no error, but got: %v", err)

	//сделать запрос БД, получить запись и сравнить
	row := s.db.QueryRow(querySelectById, request.Event.Id)
	updatedEvent := &pb.Event{}
	var datetime time.Time
	err = row.Scan(
		&updatedEvent.Id,
		&updatedEvent.Title,
		&datetime,
		&updatedEvent.Duration,
		&updatedEvent.Description,
		&updatedEvent.UserId,
		&updatedEvent.Reminder,
	)
	s.Require().NoError(err)
	updatedEvent.DateTime = &timestamp.Timestamp{
		Seconds: datetime.Unix(),
		Nanos:   int32(datetime.Nanosecond()),
	}
	s.Require().NoError(err)
	s.Require().NotNil(response, "Expected a non-nil createdEvent")
	s.Equal(response.Event.Id, updatedEvent.Id, "Event ID mismatch")
	s.Equal(response.Event.Title, updatedEvent.Title, "Event Title mismatch")
	s.Equal(response.Event.Duration, updatedEvent.Duration, "Event Duration mismatch")
	s.Equal(response.Event.Description, updatedEvent.Description, "Event Description mismatch")
	s.Equal(response.Event.UserId, updatedEvent.UserId, "Event UserId mismatch")
	s.Equal(response.Event.Reminder, updatedEvent.Reminder, "Event Reminder mismatch")
	s.Nil(err, "Expected no error, but got: %v", err)
}

func (s *CalendarSuite) TestCalendar_DeleteEvent() {
	_, err := s.db.Exec(queryInsert, "1", "Event 1", time.Now().Format("2006-01-02 15:04:05"), 1, "Description 1", "1", 1)
	s.Require().NoError(err)

	request := &pb.DeleteEventRequest{EventId: "1"}
	response, err := s.client.DeleteEvent(s.ctx, request)

	s.Require().NoError(err)
	s.Require().NotNil(response, "Expected a non-nil createdEvent")
	s.Require().Nil(response.Event, "Expected a non-nil createdEvent")
	s.Nil(err, "Expected no error, but got: %v", err)
}

func (s *CalendarSuite) TestCalendar_DeleteNonExistingEvent() {
	// Попытка удалить несуществующее событие
	request := &pb.DeleteEventRequest{EventId: "999"} // Несуществующий ID
	response, err := s.client.DeleteEvent(s.ctx, request)

	// Проверка на отсутствие ошибки и ожидаемый ответ
	s.NoError(err)
	s.NotNil(response)
	s.Nil(response.Event)
}

func (s *CalendarSuite) TestCalendar_ListEventsForDay() {
	now := time.Now()
	endOfDay := now.Add(24 * time.Hour)

	// Запрос на получение событий за день
	request := &pb.ListEventsRequest{
		Date: &timestamp.Timestamp{Seconds: now.Unix()},
	}
	response, err := s.client.ListEventsForDay(s.ctx, request)

	s.Require().NoError(err)
	s.Require().NotNil(response)

	// Проверка списка событий за день
	for _, event := range response.Events {
		eventDateTime := time.Unix(event.DateTime.Seconds, int64(event.DateTime.Nanos))
		s.True(eventDateTime.After(now) && eventDateTime.Before(endOfDay),
			"Event does not belong to the specified day")
	}
}

func (s *CalendarSuite) TestCalendar_ListEventsForWeek() {
	now := time.Now()
	endOfWeek := now.Add(7 * 24 * time.Hour)

	// Запрос на получение событий за неделю
	request := &pb.ListEventsRequest{
		Date: &timestamp.Timestamp{Seconds: now.Unix()},
	}
	response, err := s.client.ListEventsForWeek(s.ctx, request)

	s.Require().NoError(err)
	s.Require().NotNil(response)

	// Проверка списка событий за неделю
	for _, event := range response.Events {
		eventDateTime := time.Unix(event.DateTime.Seconds, int64(event.DateTime.Nanos))
		s.True(eventDateTime.After(now) && eventDateTime.Before(endOfWeek),
			"Event does not belong to the specified week")
	}
}

func (s *CalendarSuite) TestCalendar_ListEventsForMonth() {
	now := time.Now()
	firstDayOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local)
	lastDayOfMonth := firstDayOfMonth.AddDate(0, 1, -1)

	// Запрос на получение событий за месяц
	request := &pb.ListEventsRequest{
		Date: &timestamp.Timestamp{Seconds: now.Unix()},
	}
	response, err := s.client.ListEventsForMonth(s.ctx, request)

	s.Require().NoError(err)
	s.Require().NotNil(response)

	// Проверка списка событий за месяц
	for _, event := range response.Events {
		eventDateTime := time.Unix(event.DateTime.Seconds, int64(event.DateTime.Nanos))
		s.True(eventDateTime.After(firstDayOfMonth) && eventDateTime.Before(lastDayOfMonth),
			"Event does not belong to the specified month")
	}
}
