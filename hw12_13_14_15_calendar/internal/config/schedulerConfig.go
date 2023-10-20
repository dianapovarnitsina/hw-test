package config

import (
	"github.com/pkg/errors"
)

var _ Configure = (*SchedulerConfig)(nil)

type SchedulerConfig struct {
	Logger   LoggerConf   `json:"logger"`
	Storage  StorageConf  `json:"storage"`
	Database DataBaseConf `json:"database"`
	RMQ      RMQ          `json:"rmq"`

	Queues struct {
		Events Queue
	}
}

func (c *SchedulerConfig) Init(file string) error {
	cfg, err := Init(file, c)

	_, ok := cfg.(*SchedulerConfig)
	if !ok {
		return errors.Wrap(err, "init config failed")
	}

	return nil
}
