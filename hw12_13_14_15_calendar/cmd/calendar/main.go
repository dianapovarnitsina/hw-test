package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"sync"

	"github.com/dianapovarnitsina/hw-test/hw12_13_14_15_calendar/internal/app/calendar"
	config "github.com/dianapovarnitsina/hw-test/hw12_13_14_15_calendar/internal/config"
)

var calendarConfigFile string

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
		return fmt.Errorf("please set: '--calendar_config=<Path to configuration file>'")
	}
	conf, err := config.ReadConfig(calendarConfigFile)
	if err != nil {
		return fmt.Errorf("cannot read config: %w", err)
	}

	var wg sync.WaitGroup

	// Создание и инициализация приложения
	app, err := calendar.NewApp(ctx, conf)
	if err != nil {
		return fmt.Errorf("failed to create app: %w", err)
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
