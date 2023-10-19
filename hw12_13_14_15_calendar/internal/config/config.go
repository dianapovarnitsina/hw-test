package config

import (
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type Configure interface {
	Init(file string) error
}

type LoggerConf struct {
	Level string `json:"level"`
}

type StorageConf struct {
	Type      string `json:"type"`
	Migration string `json:"migration"`
}

type DataBaseConf struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Dbname   string `json:"dbname"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type HTTP struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

type GRPC struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

type RMQ struct {
	URI       string `json:"uri"`
	ReConnect struct {
		MaxElapsedTime  string  `json:"maxElapsedTime"`
		InitialInterval string  `json:"initialInterval"`
		Multiplier      float64 `json:"multiplier"`
		MaxInterval     string  `json:"maxInterval"`
	}
}

type Queue struct {
	ExchangeName string `json:"exchangeName"`
	ExchangeType string `json:"exchangeType"`
	QueueName    string `json:"queueName"`
	BindingKey   string `json:"bindingKey"` // Message routing rules
}

func Init(file string, c Configure) (Configure, error) {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.SetConfigFile(file)

	if err := viper.ReadInConfig(); err != nil {
		return nil, errors.Wrap(err, "open config file failed")
	}

	if err := viper.Unmarshal(c); err != nil {
		return nil, errors.Wrap(err, "unmarshal config file failed")
	}

	return c, nil
}
