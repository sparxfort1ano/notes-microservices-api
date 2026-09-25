package jwt

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type config struct {
	SecretKey              string        `envconfig:"SECRET_KEY" required:"true"`
	AccessTokenExpiration  time.Duration `envconfig:"ACCESS_TOKEN_EXPIRATION" default:"24h"`
	RefreshTokenExpiration time.Duration `envconfig:"REFRESH_TOKEN_EXPIRATION" default:"168h"`
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
