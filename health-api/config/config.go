package config

import (
	"fmt"
	"path"

	"github.com/sergicanet9/scv-go-tools/v4/api/utils"
)

type Async struct {
	Run      bool
	Interval utils.Duration
}

type Config struct {
	// set in flags
	Version     string
	Environment string
	HTTPPort    int
	URLs        string

	// set in json config files
	config
}

type config struct {
	Timeout utils.Duration
	Async   Async
}

// ReadConfig from the project´s JSON config files.
// Default values are specified in the default configuration file, config/config.json
// and can be overrided with values specified in the environment configuration files, config/config.{env}.json.
func ReadConfig(version, env string, httpPort int, urls, configPath string) (Config, error) {
	var c Config
	c.Version = version
	c.Environment = env
	c.HTTPPort = httpPort
	c.URLs = urls

	var cfg config

	if err := utils.LoadJSON(path.Join(configPath, "config.json"), &cfg); err != nil {
		return c, fmt.Errorf("error parsing configuration, %s", err)
	}

	if err := utils.LoadJSON(path.Join(configPath, "config."+env+".json"), &cfg); err != nil {
		return c, fmt.Errorf("error parsing environment configuration, %s", err)
	}

	c.config = cfg

	return c, nil
}
