package config

import (
	"github.com/pkg/errors"
)

var _ Configure = (*CalendarConfig)(nil)

type CalendarConfig struct {
	Logger   LoggerConf   `json:"logger"`
	FilePath string       `json:"file_path"` //nolint:tagliatelle
	Database DataBaseConf `json:"database"`
	HTTP     HTTP         `json:"http"`
	GRPC     GRPC         `json:"grpc"`
	Storage  StorageConf  `json:"storage"`
}

func (c *CalendarConfig) Init(file string) error {
	cfg, err := Init(file, c)

	_, ok := cfg.(*CalendarConfig)
	if !ok {
		return errors.Wrap(err, "init config failed")
	}

	return nil
}
