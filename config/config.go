package config

import (
	"fmt"
	"path"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		HTTP     `yaml:"http"`
		Database `yaml:"database"`
		Log      `yaml:"log"`
	}

	HTTP struct {
		Port            string        `env-required:"true" yaml:"port"`
		Host            string        `env-required:"true" yaml:"host"`
		Timeout         time.Duration `env-required:"true" yaml:"timeout" env-default:"4s"`
		ShutdownTimeout time.Duration `env-required:"true" yaml:"shutdown_timeout" env-default:"2s"`
	}

	Database struct {
		Conn         string        `env-required:"true" env:"POSTGRES_CONN"`
		MaxPoolSize  int           `env-required:"true" yaml:"max_pool_size" env-default:"1"`
		ConnAttempts int           `env-required:"true" yaml:"conn_attempts" env-default:"5"`
		ConnTimeout  time.Duration `env-required:"true" yaml:"conn_timeout" env-default:"3s"`
	}

	Log struct {
		Level string `env-required:"true" yaml:"level" env-default:"local"`
	}
)

func NewConfig(configPath string) (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadConfig(path.Join("./", configPath), cfg)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	err = cleanenv.UpdateEnv(cfg)
	if err != nil {
		return nil, fmt.Errorf("error updating env: %w", err)
	}

	return cfg, nil
}
