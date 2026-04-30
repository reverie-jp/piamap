package config

import (
	"github.com/caarlos0/env/v11"
)

type Env string

const (
	EnvDevelopment Env = "development"
	EnvStaging     Env = "staging"
	EnvProduction  Env = "production"
)

type Config struct {
	Env        Env `env:"ENVIRONMENT"`
	Auth       AuthConfig
	Database   DatabaseConfig
	Google     GoogleConfig
	Log        LogConfig
	Server     ServerConfig
	Moderation ModerationConfig
}

func New() *Config {
	return &Config{}
}

func (c *Config) LoadFromEnv() error {
	for _, target := range []any{c, &c.Auth, &c.Database, &c.Google, &c.Log, &c.Server, &c.Moderation} {
		if err := env.Parse(target); err != nil {
			return err
		}
	}
	return nil
}
