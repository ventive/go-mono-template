package v1

import (
	configparser "github.com/ventive/go-mono-template/pkg/config-parser"
	"github.com/ventive/go-mono-template/pkg/logger"
)

const (
	defaultAppEnv      = "staging"
	defaultLoggerLevel = "debug"
)

type queuesConfig struct {
	Publish struct {
		Default string `mapstructure:"default"`
		Errors  string `mapstructure:"errors"`
	} `mapstructure:"publish"`
	Subscribe struct {
		Queue string `mapstructure:"queue"`
		Group string `mapstructure:"group"`
	} `mapstructure:"subscribe"`
}

type natsConfig struct {
	URL  string `mapstructure:"url"`
	Name string `mapstructure:"name"`
	User string `mapstructure:"user"`
	Pass string `mapstructure:"pass"`
	TLS  struct {
		Enabled bool   `mapstructure:"enabled"`
		Cert    string `mapstructure:"cert"`
		Key     string `mapstructure:"key"`
		CA      string `mapstructure:"ca"`
	} `mapstructure:"tls"`
}

type appConfig struct {
	Env    string       `mapstructure:"env"`
	Nats   natsConfig   `mapstructure:"nats"`
	Queues queuesConfig `mapstructure:"queues"`
}

type config struct {
	App    appConfig     `mapstructure:"app"`
	Logger logger.Config `mapstructure:"logger"`
}

func newConfig() (config, error) {
	cfg := config{}

	defaults := map[string]interface{}{
		"app.env":      defaultAppEnv,
		"logger.level": defaultLoggerLevel,
	}

	if err := configparser.Parse(configFile, &cfg, defaults); err != nil {
		return cfg, err
	}

	return cfg, nil
}
