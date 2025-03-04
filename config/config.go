package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	ApiServerHost string `env:"APISERVER_HOST"`
	ApiServerPort string `env:"APISERVER_PORT"`
}

func New() (*Config, error) {
	cfg, err := env.ParseAs[Config]()
	if err != nil {
		return nil, fmt.Errorf("error parsing config: %s", err)
	}
	return &cfg, nil
}
