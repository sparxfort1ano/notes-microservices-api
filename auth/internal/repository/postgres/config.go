package postgres

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type config struct {
	User     string        `envconfig:"USER" required:"true"`
	Password string        `envconfig:"PASSWORD" required:"true"`
	Host     string        `envconfig:"HOST" required:"true"`
	Port     string        `envconfig:"PORT" default:"5432"`
	DB       string        `envconfig:"DB" required:"true"`
	SSLMode  string        `envconfig:"SSL_MODE" default:"disable"`
	Timeout  time.Duration `envconfig:"TIMEOUT" default:"5s"`
}

func newConfig() (config, error) {
	var cfg config
	if err := envconfig.Process("POSTGRES", &cfg); err != nil {
		return config{}, fmt.Errorf("process envconfig: %w", err)
	}

	return cfg, nil
}

func NewConfigMust() config {
	cfg, err := newConfig()
	if err != nil {
		err = fmt.Errorf("get postgres config: %w", err)
		panic(err)
	}

	return cfg
}
