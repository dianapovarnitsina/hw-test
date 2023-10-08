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
	host   string
	port   uint16
	logger interfaces.Logger
	app    interfaces.Application
	server *http.Server
}

func NewServer(host string, port uint16, logger interfaces.Logger, app interfaces.Application) *Server {
	return &Server{
		host:   host,
		port:   port,
		logger: logger,
		app:    app,
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

	err := s.app.CreateEvent(r.Context(), &event)
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

	event, err := s.app.GetEvent(r.Context(), eventID)
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
		s.app.UpdateEvent(r.Context(), &updatedEvent)
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

	err := s.app.DeleteEvent(r.Context(), eventID)
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
	_ = r
	fmt.Fprintln(w, "list events")
}
