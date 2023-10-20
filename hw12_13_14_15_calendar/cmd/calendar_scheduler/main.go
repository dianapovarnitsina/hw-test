package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/dianapovarnitsina/hw-test/hw12_13_14_15_calendar/internal/app/scheduler"
	"github.com/dianapovarnitsina/hw-test/hw12_13_14_15_calendar/internal/config"
	"github.com/pkg/errors"
)

var schedulerConfigFile string

func init() {
	flag.StringVar(&schedulerConfigFile, "config", "scheduler_config.yaml", "Path to configuration file")
}

func main() {
	if err := mainImpl(); err != nil {
		log.Fatal(err)
	}
}

func mainImpl() error {
	ctx, cancel := context.WithCancel(context.TODO())
	// ctx, cancel := context.WithTimeout(context.TODO(), time.Minute)
	defer cancel()

	flag.Parse()

	if schedulerConfigFile == "" {
		return fmt.Errorf("please set: '--config=<Path to configuration file>'")
	}

	conf := new(config.SchedulerConfig)
	if err := conf.Init(schedulerConfigFile); err != nil {
		return errors.Wrap(err, "init config failed")
	}

	_, err := scheduler.NewSchedulerApp(ctx, conf)
	if err != nil {
		return fmt.Errorf("failed to create schedulerApp: %w", err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	// <-ctx.Done()

	return nil
}
