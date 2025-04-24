package config

import (
	"fmt"
	"github.com/pelletier/go-toml/v2"
)

func ParseAndValidate(filename string) (Config, error) {
	cfg := Config{}

	err := toml.Unmarshal([]byte(filename), &cfg)
	if err != nil {
		panic(err)
	}

	// FIXME: 2) Валидируем валидатором из internal/validator.

	if err := cfg.Validate(); err != nil {
		return cfg, fmt.Errorf("cannot vaidate options: %w", err)
	}

	return Config{}, nil
}
