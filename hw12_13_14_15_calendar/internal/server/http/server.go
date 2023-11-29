package internalhttp

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/dianapovarnitsina/hw-test/hw12_13_14_15_calendar/interfaces"
	"github.com/dianapovarnitsina/hw-test/hw12_13_14_15_calendar/internal/storage"
	"github.com/gorilla/mux"
)

type Server struct {
	host    string
	port    int
	logger  interfaces.Logger
	storage interfaces.EventStorage
	server  *http.Server
}

func NewServer(host string, port int, logger interfaces.Logger, storage interfaces.EventStorage) *Server {
	return &Server{
		host:    host,
		port:    port,
		logger:  logger,
		storage: storage,
	}
}

func (s *Server) Start(ctx context.Context) error {
	router := mux.NewRouter()

	router.HandleFunc("/hello", s.helloHandler).Methods(http.MethodGet)
	router.HandleFunc("/", s.helloHandler).Methods(http.MethodGet)
	router.HandleFunc("/event/create", s.createEventHandler).Methods(http.MethodPost)
	router.HandleFunc("/event/{id}", s.getEventHandler).Methods(http.MethodGet)
	router.HandleFunc("/event/{id}", s.updateEventHandler).Methods(http.MethodPut)
	router.HandleFunc("/event/{id}", s.deleteEventHandler).Methods(http.MethodDelete)
	router.HandleFunc("/events", s.listEventsHandler).Methods(http.MethodGet)

	handlerWithMiddleware := middleware(router, s.logger)

	s.server = &http.Server{
		Addr:              fmt.Sprintf("%s:%d", s.host, s.port),
		Handler:           handlerWithMiddleware,
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		<-ctx.Done()
		s.Stop(ctx)
	}()

	s.logger.Info(fmt.Sprintf("Server started at %s:%d", s.host, s.port))
	return s.server.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	if s.server != nil {
		s.logger.Info("Shutting down the server...")
		return s.server.Shutdown(ctx)
	}
	return nil
}

func (s *Server) helloHandler(w http.ResponseWriter, r *http.Request) {
	_ = r
	fmt.Fprintln(w, "Hello, World!")
}

func (s *Server) createEventHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var event storage.Event
	if err := decoder.Decode(&event); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	err := s.storage.CreateEvent(r.Context(), &event)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(event)
}

func (s *Server) getEventHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	eventID, ok := vars["id"]
	if !ok {
		http.Error(w, "Event ID not provided", http.StatusBadRequest)
		return
	}

	event, err := s.storage.GetEvent(r.Context(), eventID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(event)
}

func (s *Server) updateEventHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	eventID, ok := vars["id"]
	if !ok {
		http.Error(w, "Event ID not provided", http.StatusBadRequest)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var updatedEvent storage.Event

	if err := decoder.Decode(&updatedEvent); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if eventID == updatedEvent.ID {
		s.storage.UpdateEvent(r.Context(), &updatedEvent)
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, "Event updated successfully")
	} else {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
}

func (s *Server) deleteEventHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	eventID, ok := vars["id"]
	if !ok {
		http.Error(w, "Event ID not provided", http.StatusBadRequest)
		return
	}

	err := s.storage.DeleteEvent(r.Context(), eventID)
	if err != nil {
		if errors.Is(err, storage.ErrEventNotFound) {
			http.Error(w, "Event not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to delete event", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintln(w, "Event deleted successfully")
}

func (s *Server) listEventsHandler(w http.ResponseWriter, r *http.Request) {
	var events []*storage.Event
	var err error

	queryParams := r.URL.Query()

	switch {
	case len(queryParams["day"]) > 0:
		dayStr := queryParams["day"][0]
		day, parseErr := time.Parse("2006-01-02", dayStr)
		if parseErr != nil {
			http.Error(w, "Invalid date format", http.StatusBadRequest)
			return
		}
		events, err = s.storage.ListEventsForDay(r.Context(), day)

	case len(queryParams["week"]) > 0:
		weekStr := queryParams["week"][0]
		week, parseErr := time.Parse("2006-01-02", weekStr)
		if parseErr != nil {
			http.Error(w, "Invalid date format", http.StatusBadRequest)
			return
		}
		events, err = s.storage.ListEventsForWeek(r.Context(), week)

	case len(queryParams["month"]) > 0:
		monthStr := queryParams["month"][0]
		month, parseErr := time.Parse("2006-01-02", monthStr)
		if parseErr != nil {
			http.Error(w, "Invalid date format", http.StatusBadRequest)
			return
		}
		events, err = s.storage.ListEventsForMonth(r.Context(), month)

	default:
		http.Error(w, "Invalid query", http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(events)
}
