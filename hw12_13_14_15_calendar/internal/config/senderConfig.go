package config

import (
	"github.com/pkg/errors"
)

var _ Configure = (*SenderConfig)(nil)

type SenderConfig struct {
	Logger   LoggerConf   `json:"logger"`
	Database DataBaseConf `json:"database"`
	RMQ      RMQ          `json:"rmq"`

	Queues struct {
		Events  Queue `json:"events"`
		Senders Queue `json:"senders"`
	}

	Consumer struct {
		ConsumerTag      string `json:"consumerTag"` // Consumer tag name
		QosPrefetchCount int    `json:"qosPrefetchCount"`
		Threads          int    `json:"threads"` // Count threads for reading queue messages
	}
}

func (c *SenderConfig) Init(file string) error {
	cfg, err := Init(file, c)

	_, ok := cfg.(*SenderConfig)
	if !ok {
		return errors.Wrap(err, "init config failed")
	}

	return nil
}
