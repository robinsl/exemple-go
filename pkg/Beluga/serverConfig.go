package Beluga

import (
	"github.com/kelseyhightower/envconfig"
	"time"
)

type HttpServerConfiguration struct {
	Port         int           `envconfig:"HTTP_PORT" required:"true" default:"80"`
	IdleTimeout  time.Duration `envconfig:"HTTP_IDLE_TIMEOUT" required:"true" default:"2m"`
	ReadTimeout  time.Duration `envconfig:"HTTP_READ_TIMEOUT" required:"true" default:"5s"`
	WriteTimeout time.Duration `envconfig:"HTTP_WRITE_TIMEOUT" required:"true" default:"10s"`
}

func LoadHttpServerConfiguration() (HttpServerConfiguration, error) {
	var cfg HttpServerConfiguration

	if err := envconfig.Process("", &cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}
