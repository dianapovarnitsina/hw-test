package config

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

func ReadConfig(file string) (*Config, error) {
	ext := filepath.Ext(file)
	switch strings.ToLower(ext) {
	case ".json":
		return readJSON(file)
	case ".yaml", ".yml":
		return readYAML(file)
	case ".toml":
		return readTOML(file)
	default:
		return nil, errors.New("unsupported config file format")
	}
}

func readJSON(file string) (*Config, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read config file")
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal config")
	}

	return &config, nil
}

func readYAML(file string) (*Config, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read config file")
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal YAML config")
	}

	return &config, nil
}

func readTOML(file string) (*Config, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read config file")
	}

	var config Config
	if err := toml.Unmarshal(data, &config); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal TOML config")
	}

	return &config, nil
}
