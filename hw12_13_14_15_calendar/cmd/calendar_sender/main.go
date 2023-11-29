package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/dianapovarnitsina/hw-test/hw12_13_14_15_calendar/internal/app/sender"
	"github.com/dianapovarnitsina/hw-test/hw12_13_14_15_calendar/internal/config"
	"github.com/pkg/errors"
)

var senderConfigFile string

func init() {
	flag.StringVar(&senderConfigFile, "config", "sender_config.yaml", "Path to configuration file")
}

func main() {
	if err := mainImpl(); err != nil {
		log.Fatal(err)
	}
}

func mainImpl() error {
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	flag.Parse()

	if senderConfigFile == "" {
		return fmt.Errorf("please set: '--config=<Path to configuration file>'")
	}

	conf := new(config.SenderConfig)
	if err := conf.Init(senderConfigFile); err != nil {
		return errors.Wrap(err, "init config failed")
	}

	_, err := sender.NewSenderApp(ctx, conf)
	if err != nil {
		return fmt.Errorf("failed to create senderApp: %w", err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	return nil
}
