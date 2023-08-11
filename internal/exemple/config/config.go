package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"time"
)

type Database struct {
	Url              string `envconfig:"MONGO_URL" required:"true"`
	Port             string `envconfig:"MONGO_PORT" required:"true"`
	Username         string `envconfig:"MONGO_USERNAME" required:"true"`
	Password         string `envconfig:"MONGO_PASSWORD" required:"true"`
	ConnectionString string
	DatabaseName     string `envconfig:"MONGO_DB" required:"true"`
}

type HttpServer struct {
	Port         int           `envconfig:"HTTP_PORT" required:"true" default:"80"`
	IdleTimeout  time.Duration `envconfig:"HTTP_IDLE_TIMEOUT" required:"true" default:"2m"`
	ReadTimeout  time.Duration `envconfig:"HTTP_READ_TIMEOUT" required:"true" default:"5s"`
	WriteTimeout time.Duration `envconfig:"HTTP_WRITE_TIMEOUT" required:"true" default:"10s"`
}

type Configurations struct {
	Database
	HttpServer
}

func Load() (Configurations, error) {
	var cfg Configurations
	if err := godotenv.Load("configs/.env"); err != nil {
		return cfg, err
	}

	if err := envconfig.Process("", &cfg); err != nil {
		return cfg, err
	}

	cfg.Database.ConnectionString = fmt.Sprintf("mongodb://%s:%s", cfg.Database.Url, cfg.Database.Port)

	return cfg, nil
}
