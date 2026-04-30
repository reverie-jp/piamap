package config

import "time"

type DatabaseConfig struct {
	DSN             string        `env:"DATABASE_DSN"`
	MaxConns        int32         `env:"DATABASE_MAX_CONNS" envDefault:"10"`
	MinConns        int32         `env:"DATABASE_MIN_CONNS" envDefault:"1"`
	MaxConnLifetime time.Duration `env:"DATABASE_MAX_CONN_LIFETIME" envDefault:"30m"`
	MaxConnIdleTime time.Duration `env:"DATABASE_MAX_CONN_IDLE_TIME" envDefault:"5m"`
}
