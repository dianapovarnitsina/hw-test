package internalhttp

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/dianapovarnitsina/hw-test/hw12_13_14_15_calendar/interfaces"
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
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", s.helloHandler)
	mux.HandleFunc("/", s.helloHandler)

	handlerWithMiddleware := middleware(mux, s.logger)

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
	fmt.Fprintln(w, "Hello, World!")
}
