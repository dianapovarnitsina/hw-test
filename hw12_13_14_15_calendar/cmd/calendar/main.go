package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dianapovarnitsina/hw-test/hw12_13_14_15_calendar/interfaces"
	"github.com/dianapovarnitsina/hw-test/hw12_13_14_15_calendar/internal/app"
	config "github.com/dianapovarnitsina/hw-test/hw12_13_14_15_calendar/internal/config"
	"github.com/dianapovarnitsina/hw-test/hw12_13_14_15_calendar/internal/logger"
	internalgrpc "github.com/dianapovarnitsina/hw-test/hw12_13_14_15_calendar/internal/server/grpc"
	internalhttp "github.com/dianapovarnitsina/hw-test/hw12_13_14_15_calendar/internal/server/http"
	"github.com/dianapovarnitsina/hw-test/hw12_13_14_15_calendar/internal/server/pb"
	memorystorage "github.com/dianapovarnitsina/hw-test/hw12_13_14_15_calendar/internal/storage/memory"
	sql "github.com/dianapovarnitsina/hw-test/hw12_13_14_15_calendar/internal/storage/sql"
	"google.golang.org/grpc"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "", "Path to configuration file")
}

func mainImpl() error {
	ctx := context.TODO()

	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return fmt.Errorf("version")
	}
	if configFile == "" {
		return fmt.Errorf("please set: '--config=<Path to configuration file>'")
	}

	conf, err := config.ReadConfig(configFile)
	if err != nil {
		return fmt.Errorf("cannot read config: %w", err)
	}

	var storage interfaces.EventStorage

	if conf.Storage.Type == "postgres" {
		storage = new(sql.Storage)
		if err := storage.Connect(ctx, conf); err != nil {
			return fmt.Errorf("cannot connect to psql: %w", err)
		}

		err := storage.Migrate(ctx, conf.Storage.Migration)
		if err != nil {
			return fmt.Errorf("migration did not work out")
		}

		defer func() {
			if err := storage.Close(); err != nil {
				log.Println("cannot close psql connection", err)
			}
		}()
	} else {
		storage = memorystorage.New()
	}

	logg := logger.New(conf.Logger.Level, os.Stdout)
	calendar := app.NewApp(logg, storage)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	if conf.Net.API == "grpc" {
		err = startGRPC(ctx, *conf, logg, storage)
	} else {
		err = startHTTP(ctx, *conf, logg, calendar)
	}

	if err != nil {
		return err
	}

	return nil
}

func startGRPC(ctx context.Context, conf config.Config, logg *logger.Logger, storage interfaces.EventStorage) error {
	server := grpc.NewServer(
		grpc.UnaryInterceptor(internalgrpc.NewLoggingInterceptor(logg).UnaryServerInterceptor),
	)

	go func() {
		<-ctx.Done()

		_, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()
		server.Stop()
	}()

	api := internalgrpc.NewEventServiceServer(storage)
	pb.RegisterCalendarServiceServer(server, api)

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", conf.Net.Host, conf.Net.Port))
	if err != nil {
		return err
	}

	logg.Info("calendar is running...")

	if err := server.Serve(listener); err != nil {
		logg.Error("failed to start grpc server: " + err.Error())
		return err
	}
	return nil
}

func startHTTP(ctx context.Context, conf config.Config, logg *logger.Logger, calendar *app.App) error {
	server := internalhttp.NewServer(conf.Net.Host, conf.Net.Port, logg, calendar)

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()

	logg.Info("calendar is running...")

	if err := server.Start(ctx); err != nil {
		logg.Error("failed to start http server: " + err.Error())
		return err
	}
	return nil
}

func main() {
	if err := mainImpl(); err != nil {
		log.Fatal(err)
	}
}
