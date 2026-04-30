package config

type ModerationConfig struct {
	AutoHideThreshold int `env:"MODERATION_AUTO_HIDE_THRESHOLD" envDefault:"5"`
}
