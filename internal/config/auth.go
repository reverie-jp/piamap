package config

import "time"

type AuthConfig struct {
	JWTSecretKey      string        `env:"AUTH_JWT_SECRET_KEY"`
	AccessExpiration  time.Duration `env:"AUTH_ACCESS_EXPIRATION" envDefault:"15m"`
	RefreshExpiration time.Duration `env:"AUTH_REFRESH_EXPIRATION" envDefault:"720h"`
}
