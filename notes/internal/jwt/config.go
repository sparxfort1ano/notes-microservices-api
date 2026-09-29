package jwt

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type config struct {
	SecretKey string `envconfig:"SECRET_KEY" required:"true"`
}

func newConfig() (config, error) {
	var cfg config
	if err := envconfig.Process("JWT", &cfg); err != nil {
		return config{}, fmt.Errorf("process envconfig: %w", err)
	}

	return cfg, nil
}

func NewConfigMust() config {
	cfg, err := newConfig()
	if err != nil {
		err = fmt.Errorf("get jwt config: %w", err)
		panic(err)
	}

	return cfg
}
