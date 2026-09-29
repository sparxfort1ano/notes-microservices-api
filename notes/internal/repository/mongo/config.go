package mongo

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type config struct {
	User       string        `envconfig:"ROOT_USERNAME" required:"true"`
	Password   string        `envconfig:"ROOT_PASSWORD" required:"true"`
	Host       string        `envconfig:"HOST" required:"true"`
	Port       string        `envconfig:"PORT" default:"27017"`
	DB         string        `envconfig:"DB" required:"true"`
	Collection string        `envconfig:"COLLECTION" default:"notes"`
	Timeout    time.Duration `envconfig:"TIMEOUT" default:"5s"`
	SSLMode    bool          `envconfig:"SSL_MODE" default:"false"`
}

func newConfig() (config, error) {
	var cfg config
	if err := envconfig.Process("MONGO_INITDB", &cfg); err != nil {
		return config{}, fmt.Errorf("process envconfig: %w", err)
	}

	return cfg, nil
}

func NewConfigMust() config {
	cfg, err := newConfig()
	if err != nil {
		err = fmt.Errorf("get mongo config: %w", err)
		panic(err)
	}

	return cfg
}

func Collection() string {
	return NewConfigMust().Collection
}
