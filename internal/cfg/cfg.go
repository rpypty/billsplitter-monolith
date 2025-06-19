package cfg

import (
	"fmt"

	"github.com/spf13/viper"
)

var (
	configPath = "."
	debug      = false
)

type Telegram struct {
	BotToken string
}

type Server struct {
	Telegram Telegram
	Http     Http
}

type Http struct {
	Port string
}

type Storage struct {
	Postgres Postgres
}

type Postgres struct {
	DSN string
}

type Config struct {
	Debug   bool
	Server  Server
	Storage Storage
}

func IsDebug() bool {
	return debug
}

func LoadConfig() (Config, error) {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(configPath)

	if err := v.ReadInConfig(); err != nil {
		return Config{}, fmt.Errorf("failed to read config: %w", err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return Config{}, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	debug = cfg.Debug

	return cfg, nil
}
