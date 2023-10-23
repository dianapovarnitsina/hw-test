package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"sync"

	"github.com/dianapovarnitsina/hw-test/hw12_13_14_15_calendar/internal/app/calendar"
	"github.com/dianapovarnitsina/hw-test/hw12_13_14_15_calendar/internal/config"
	"github.com/pkg/errors"
)

var (
	calendarConfigFile string
	wg                 sync.WaitGroup
)

func init() {
	flag.StringVar(&calendarConfigFile, "config", "calendar_config.toml", "Path to configuration file")
}

func main() {
	if err := mainImpl(); err != nil {
		log.Fatal(err)
	}
}

func mainImpl() error {
	ctx := context.TODO()

	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return fmt.Errorf("version")
	}

	if calendarConfigFile == "" {
		return fmt.Errorf("please set: '--config=<Path to configuration file>'")
	}

	conf := new(config.CalendarConfig)
	if err := conf.Init(calendarConfigFile); err != nil {
		return errors.Wrap(err, "init config failed")
	}

	// Создание и инициализация приложения
	app, err := calendar.NewApp(ctx, conf)
	if err != nil {
		return fmt.Errorf("failed to create calendarApp: %w", err)
	}

	// Увеличиваем счетчик WaitGroup на 2, так как у нас два сервера
	wg.Add(2)

	// Отдельные горутины для запуска серверов
	go func() {
		defer wg.Done()
		// Ожидаем завершения работы сервера HTTP
		<-app.GetHTTPServerShutdownSignal()
	}()

	go func() {
		defer wg.Done()
		// Ожидаем завершения работы сервера gRPC
		<-app.GetGrpcServerShutdownSignal()
	}()

	// Ожидаем завершения работы всех серверов
	wg.Wait()

	return nil
}
