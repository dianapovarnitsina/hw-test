package internalhttp

import (
	"net/http"
	"time"

	"github.com/dianapovarnitsina/hw-test/hw12_13_14_15_calendar/interfaces"
)

func middleware(wrappedHandler http.Handler, logger interfaces.Logger) http.Handler { //nolint:unused
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		StartAt := time.Now()
		lrw := NewLoggingResponseWriter(w)
		wrappedHandler.ServeHTTP(lrw, r)
		a := struct {
			ClientIPAddress string
			StartAt         time.Time
			HTTPMethod      string
			HTTPVersion     string
			Query           string
			StatusCode      int
			UserAgent       string
			Latency         time.Duration
		}{
			ClientIPAddress: r.RemoteAddr,
			StartAt:         StartAt,
			HTTPMethod:      r.Method,
			HTTPVersion:     r.Proto,
			Query:           r.URL.Query().Get("q"),
			StatusCode:      lrw.StatusCode,
			UserAgent:       r.UserAgent(),
			Latency:         time.Since(StartAt),
		}
		logger.Info("%+v", a)
	})
}

type LoggingResponseWriter struct {
	http.ResponseWriter
	StatusCode int
}

func NewLoggingResponseWriter(writer http.ResponseWriter) *LoggingResponseWriter {
	return &LoggingResponseWriter{writer, 0}
}

func (lrw *LoggingResponseWriter) WriteHeader(code int) {
	lrw.StatusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}
