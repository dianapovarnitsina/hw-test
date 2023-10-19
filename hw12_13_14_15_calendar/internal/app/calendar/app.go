package calendar

import (
	"context"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/dianapovarnitsina/hw-test/hw12_13_14_15_calendar/interfaces"
	"github.com/dianapovarnitsina/hw-test/hw12_13_14_15_calendar/internal/config"
	"github.com/dianapovarnitsina/hw-test/hw12_13_14_15_calendar/internal/logger"
	internalgrpc "github.com/dianapovarnitsina/hw-test/hw12_13_14_15_calendar/internal/server/grpc"
	internalhttp "github.com/dianapovarnitsina/hw-test/hw12_13_14_15_calendar/internal/server/http"
	"github.com/dianapovarnitsina/hw-test/hw12_13_14_15_calendar/internal/server/pb"
	"github.com/dianapovarnitsina/hw-test/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/dianapovarnitsina/hw-test/hw12_13_14_15_calendar/internal/storage/memory"
	sql "github.com/dianapovarnitsina/hw-test/hw12_13_14_15_calendar/internal/storage/sql"
	"google.golang.org/grpc"
)

type App struct {
	logger       interfaces.Logger
	storage      interfaces.EventStorage
	serverGRPC   *grpc.Server
	serverHTTP   *internalhttp.Server
	grpcShutdown chan struct{} // Канал для сигнала завершения gRPC сервера
	httpShutdown chan struct{} // Канал для сигнала завершения HTTP сервера
}

func NewApp(ctx context.Context, conf *config.CalendarConfig) (*App, error) {
	app := &App{}

	// Инициализация логгера.
	logger := logger.New(conf.Logger.Level, os.Stdout)
	app.logger = logger

	// Инициализация хранилища данных.
	var eventStorage interfaces.EventStorage
	if conf.Storage.Type == "postgres" {
		psqlStorage := new(sql.Storage)
		if err := psqlStorage.Connect(ctx, conf); err != nil {
			return nil, fmt.Errorf("cannot connect to PostgreSQL: %w", err)
		}
		err := psqlStorage.Migrate(ctx, conf.Storage.Migration)
		if err != nil {
			return nil, fmt.Errorf("migration did not work out: %w", err)
		}
		eventStorage = psqlStorage
	} else {
		eventStorage = memorystorage.New()
	}
	app.storage = eventStorage

	// Инициализация HTTP-сервера.
	app.serverHTTP = internalhttp.NewServer(conf.HTTP.Host, conf.HTTP.Port, logger, app.storage)
	go func() {
		logger.Info("Starting HTTP server on port %d", conf.HTTP.Port)
		if err := app.serverHTTP.Start(ctx); err != nil {
			logger.Error("HTTP server failed: %v", err)
		}
		close(app.httpShutdown) // Отправляем сигнал о завершении работы HTTP сервера.
	}()

	// Инициализация gRPC-сервера
	app.serverGRPC = grpc.NewServer(
		grpc.UnaryInterceptor(internalgrpc.NewLoggingInterceptor(logger).UnaryServerInterceptor),
	)

	api := internalgrpc.NewEventServiceServer(app.storage)
	pb.RegisterCalendarServiceServer(app.serverGRPC, api)

	grpcListener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", conf.GRPC.Host, conf.GRPC.Port))
	if err != nil {
		logger.Error("Failed to listen: %v", err)
	}
	go func() {
		logger.Info("Starting gRPC server on port %s", fmt.Sprintf("%s:%d", conf.GRPC.Host, conf.GRPC.Port))
		if err := app.serverGRPC.Serve(grpcListener); err != nil {
			logger.Error("gRPC server failed: %v", err)
		}
		close(app.grpcShutdown) // Отправляем сигнал о завершении работы gRPC сервера.
	}()

	return app, nil
}

// GetGrpcServerShutdownSignal метод получения сигналов о завершении работы сервера.
func (a *App) GetGrpcServerShutdownSignal() <-chan struct{} {
	return a.grpcShutdown
}

// GetHTTPServerShutdownSignal метод получения сигналов о завершении работы сервера.
func (a *App) GetHTTPServerShutdownSignal() <-chan struct{} {
	return a.httpShutdown
}

func (a *App) CreateEvent(ctx context.Context, event *storage.Event) error {
	return a.storage.CreateEvent(ctx, event)
}

func (a *App) UpdateEvent(ctx context.Context, event *storage.Event) error {
	return a.storage.UpdateEvent(ctx, event)
}

func (a *App) DeleteEvent(ctx context.Context, eventID string) error {
	return a.storage.DeleteEvent(ctx, eventID)
}

func (a *App) GetEvent(ctx context.Context, eventID string) (*storage.Event, error) {
	return a.storage.GetEvent(ctx, eventID)
}

func (a *App) ListEventsForDay(ctx context.Context, date time.Time) ([]*storage.Event, error) {
	return a.storage.ListEventsForDay(ctx, date)
}

func (a *App) ListEventsForWeek(ctx context.Context, date time.Time) ([]*storage.Event, error) {
	return a.storage.ListEventsForWeek(ctx, date)
}

func (a *App) ListEventsForMonth(ctx context.Context, date time.Time) ([]*storage.Event, error) {
	return a.storage.ListEventsForWeek(ctx, date)
}
