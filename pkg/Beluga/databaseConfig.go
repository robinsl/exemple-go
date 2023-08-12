package Beluga

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"os"
)

type DatabaseConfiguration struct {
	Url              string `envconfig:"MONGO_URL" required:"true"`
	Port             string `envconfig:"MONGO_PORT" required:"true"`
	Username         string `envconfig:"MONGO_USERNAME" required:"true"`
	Password         string `envconfig:"MONGO_PASSWORD" required:"true"`
	ConnectionString string
	DatabaseName     string `envconfig:"MONGO_DB" required:"true"`
}

func LoadDatabaseConfiguration(prefix string) (DatabaseConfiguration, error) {
	var cfg DatabaseConfiguration

	if err := envconfig.Process("", &cfg); err != nil {
		return cfg, err
	}

	url := os.Getenv(prefix + "_MONGO_URL")
	if url != "" {
		cfg.Url = url
	}

	port := os.Getenv(prefix + "_MONGO_PORT")
	if port != "" {
		cfg.Port = port
	}

	username := os.Getenv(prefix + "_MONGO_USERNAME")
	if username != "" {
		cfg.Username = username
	}

	password := os.Getenv(prefix + "_MONGO_PASSWORD")
	if password != "" {
		cfg.Password = password
	}

	db := os.Getenv(prefix + "_MONGO_DB")
	if db != "" {
		cfg.DatabaseName = db
	}

	cfg.ConnectionString = fmt.Sprintf("mongodb://%s:%s", cfg.Url, cfg.Port)

	return cfg, nil
}
