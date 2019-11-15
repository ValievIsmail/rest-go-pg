package main

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

type config struct {
	API struct {
		Port string `default:":8080" envconfig:"API_PORT"`
	}
	DB struct {
		Host         string        `default:"localhost" envconfig:"DB_HOST"`
		Name         string        `default:"restdb" envconfig:"DB_NAME"`
		User         string        `default:"dbuser" envconfig:"DB_USER"`
		Password     string        `default:"qwerty123" envconfig:"DB_PASSWORD"`
		Port         int           `default:"54320" envconfig:"DB_PORT"`
		PoolSize     int           `default:"5" envconfig:"DB_POOLSIZE"`
		MaxIdleConns int           `default:"3" envconfig:"DB_MAX_IDLE"`
		ConnLifetime time.Duration `default:"10m" envconfig:"DB_CONNLIFETIME"`
		Tmpl         string        `default:"host=%s port=%d dbname=%s user=%s password=%s sslmode=disable application_name=%s"`
	}
}

func parseConfig(app string) (cfg config, err error) {
	if err := envconfig.Process(app, &cfg); err != nil {
		if err := envconfig.Usage(app, &cfg); err != nil {
			return cfg, err
		}
		return cfg, err
	}
	return cfg, nil
}
