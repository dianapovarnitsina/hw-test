package main

import (
	"flag"
	"fmt"
	"log"

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
	//ctx := context.TODO()
	flag.Parse()

	if senderConfigFile == "" {
		return fmt.Errorf("please set: '--config=<Path to configuration file>'")
	}

	conf := new(config.SenderConfig)
	if err := conf.Init(senderConfigFile); err != nil {
		return errors.Wrap(err, "init config failed")
	}

	return nil
}
