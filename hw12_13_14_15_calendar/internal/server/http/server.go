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

	s.server = &http.Server{
		Addr:              fmt.Sprintf("%s:%d", s.host, s.port),
		Handler:           mux,
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
	startTime := time.Now()
	qValue := r.URL.Query().Get("q")
	latency := time.Since(startTime)

	s.logger.Info(fmt.Sprintf(
		"%s %s [%s] %s %s %s %d %d \"%s\"",
		r.RemoteAddr,
		qValue,
		startTime.Format("02/Jan/2006:15:04:05 -0700"),
		r.Method,
		r.URL.Path,
		r.Proto,
		http.StatusOK,
		latency.Milliseconds(),
		r.Header.Get("User-Agent"),
	))

	fmt.Fprintln(w, "Hello, World!")
}
