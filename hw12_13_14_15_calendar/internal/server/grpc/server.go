package internalgrpc

import (
	"context"

	"github.com/dianapovarnitsina/hw-test/hw12_13_14_15_calendar/interfaces"
	pb "github.com/dianapovarnitsina/hw-test/hw12_13_14_15_calendar/internal/server/pb"
	"github.com/dianapovarnitsina/hw-test/hw12_13_14_15_calendar/internal/storage"
	"github.com/golang/protobuf/ptypes/timestamp"
)

type EventServiceServer struct {
	storage interfaces.EventStorage
	pb.UnimplementedCalendarServiceServer
}

func NewEventServiceServer(storage interfaces.EventStorage) *EventServiceServer {
	return &EventServiceServer{
		storage: storage,
	}
}

//nolint:dupl
func (s *EventServiceServer) CreateEvent(ctx context.Context, req *pb.EventRequest) (*pb.EventResponse, error) {
	// Преобразуем pb.EventRequest в структуру storage.Event
	event := &storage.Event{
		ID:          req.Event.Id,
		Title:       req.Event.Title,
		Description: req.Event.Description,
		UserID:      req.Event.UserId,
		Duration:    req.Event.Duration,
		Reminder:    req.Event.Reminder,
		DateTime:    req.Event.DateTime.AsTime(),
	}

	// Вызываем хранилище для создания события
	err := s.storage.CreateEvent(ctx, event)
	if err != nil {
		return nil, err
	}

	// Преобразуем storage.Event в pb.EventResponse
	eventResp := &pb.EventResponse{
		Event: &pb.Event{
			Id:          event.ID,
			Title:       event.Title,
			Description: event.Description,
			UserId:      event.UserID,
			Duration:    event.Duration,
			Reminder:    event.Reminder,
			DateTime:    &timestamp.Timestamp{Seconds: event.DateTime.Unix()},
		},
	}
	return eventResp, nil
}

//nolint:dupl
func (s *EventServiceServer) UpdateEvent(ctx context.Context, req *pb.EventRequest) (*pb.EventResponse, error) {
	// Преобразуем pb.EventRequest в структуру storage.Event
	event := &storage.Event{
		ID:          req.Event.Id,
		Title:       req.Event.Title,
		Description: req.Event.Description,
		UserID:      req.Event.UserId,
		Duration:    req.Event.Duration,
		Reminder:    req.Event.Reminder,
		DateTime:    req.Event.DateTime.AsTime(),
	}

	// Вызываем хранилище для обновления события
	err := s.storage.UpdateEvent(ctx, event)
	if err != nil {
		return nil, err
	}

	// Преобразуем storage.Event в pb.EventResponse
	eventResp := &pb.EventResponse{
		Event: &pb.Event{
			Id:          event.ID,
			Title:       event.Title,
			Description: event.Description,
			UserId:      event.UserID,
			Duration:    event.Duration,
			Reminder:    event.Reminder,
			DateTime:    &timestamp.Timestamp{Seconds: event.DateTime.Unix()},
		},
	}
	return eventResp, nil
}

func (s *EventServiceServer) DeleteEvent(ctx context.Context, req *pb.DeleteEventRequest) (*pb.EventResponse, error) {
	// Получаем идентификатор события для удаления
	eventID := req.EventId

	// Вызываем хранилище для удаления события по его идентификатору
	err := s.storage.DeleteEvent(ctx, eventID)
	if err != nil {
		return nil, err
	}

	// Возвращаем пустой объект ответа, так как событие успешно удалено
	eventResp := &pb.EventResponse{}
	return eventResp, nil
}

func (s *EventServiceServer) GetEvent(ctx context.Context, req *pb.GetEventRequest) (*pb.EventResponse, error) {
	eventID := req.EventId

	event, err := s.storage.GetEvent(ctx, eventID)
	if err != nil {
		return nil, err
	}

	// Преобразуем поля event в соответствующие поля pb.Event
	pbEvent := &pb.Event{
		Id:          event.ID,
		Title:       event.Title,
		Description: event.Description,
		UserId:      event.UserID,
		Duration:    event.Duration,
		Reminder:    event.Reminder,
		DateTime:    &timestamp.Timestamp{Seconds: event.DateTime.Unix()},
	}

	eventResponse := &pb.EventResponse{
		Event: pbEvent,
	}
	return eventResponse, nil
}

func (s *EventServiceServer) ListEventsForDay(
	ctx context.Context, req *pb.ListEventsRequest,
) (*pb.ListEventsResponse, error) {
	date := req.Date

	events, err := s.storage.ListEventsForDay(ctx, date.AsTime())
	if err != nil {
		return nil, err
	}

	return &pb.ListEventsResponse{Events: eventsToPBEvents(events)}, nil
}

func (s *EventServiceServer) ListEventsForWeek(
	ctx context.Context, req *pb.ListEventsRequest,
) (*pb.ListEventsResponse, error) {
	date := req.Date

	events, err := s.storage.ListEventsForWeek(ctx, date.AsTime())
	if err != nil {
		return nil, err
	}

	return &pb.ListEventsResponse{Events: eventsToPBEvents(events)}, nil
}

func (s *EventServiceServer) ListEventsForMonth(
	ctx context.Context, req *pb.ListEventsRequest,
) (*pb.ListEventsResponse, error) {
	date := req.Date

	events, err := s.storage.ListEventsForMonth(ctx, date.AsTime())
	if err != nil {
		return nil, err
	}

	return &pb.ListEventsResponse{Events: eventsToPBEvents(events)}, nil
}

func eventsToPBEvents(events []*storage.Event) []*pb.Event {
	pbEvents := make([]*pb.Event, len(events))
	for i, event := range events {
		pbEvents[i] = &pb.Event{
			Id:          event.ID,
			Title:       event.Title,
			Description: event.Description,
			UserId:      event.UserID,
			Duration:    event.Duration,
			Reminder:    event.Reminder,
			DateTime:    &timestamp.Timestamp{Seconds: event.DateTime.Unix(), Nanos: int32(event.DateTime.Nanosecond())},
		}
	}
	return pbEvents
}
