package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dianapovarnitsina/hw-test/hw12_13_14_15_calendar/interfaces"
	"github.com/dianapovarnitsina/hw-test/hw12_13_14_15_calendar/internal/app"
	config "github.com/dianapovarnitsina/hw-test/hw12_13_14_15_calendar/internal/config"
	"github.com/dianapovarnitsina/hw-test/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/dianapovarnitsina/hw-test/hw12_13_14_15_calendar/internal/server/http"
	memorystorage "github.com/dianapovarnitsina/hw-test/hw12_13_14_15_calendar/internal/storage/memory"
	sql "github.com/dianapovarnitsina/hw-test/hw12_13_14_15_calendar/internal/storage/sql"
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
		return fmt.Errorf("Please set: '--config=<Path to configuration file>'")
	}

	conf, err := config.ReadConfig(configFile)
	if err != nil {
		return fmt.Errorf("cannot read config: %v", err)
	}

	var storage interfaces.EventStorage

	if conf.Storage.Type == "postgres" {
		storage = new(sql.Storage)
		if err := storage.Connect(ctx, conf); err != nil {
			return fmt.Errorf("cannot connect to psql: %v", err)
		}

		err := storage.Migrate(ctx, conf.Storage.Migration) //TODO разобраться почему миграция не накатывается
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

	//логер
	logg := logger.New(conf.Logger.Level, os.Stdout)

	//приложение
	calendar := app.NewApp(logg, storage)
	server := internalhttp.NewServer(conf.Http.Host, conf.Http.Port, logg, calendar)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

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
		cancel()
		os.Exit(1) //nolint:gocritic
	}

	return nil
}

func main() {
	if err := mainImpl(); err != nil {
		log.Fatal(err)
	}
}
