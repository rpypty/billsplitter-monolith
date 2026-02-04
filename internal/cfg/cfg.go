package cfg

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
	"github.com/subosito/gotenv"
)

var (
	configPath = "."
	debug      = false
)

type Telegram struct {
	BotToken string `mapstructure:"bot_token"`
}

type Server struct {
	Http Http `mapstructure:"http"`
}

type Http struct {
	Port string `mapstructure:"port"`
}

type Storage struct {
	Postgres Postgres `mapstructure:"postgres"`
}

type Postgres struct {
	Host           string `mapstructure:"host"`
	Port           int    `mapstructure:"port"`
	Database       string `mapstructure:"database"`
	User           string `mapstructure:"user"`
	Password       string `mapstructure:"password"`
	SSLMode        string `mapstructure:"sslmode"`
	ConnectTimeout int    `mapstructure:"connect_timeout"`
	MaxOpenConns   int    `mapstructure:"max_open_conns"`
	MaxIdleConns   int    `mapstructure:"max_idle_conns"`
	DSN            string `mapstructure:"dsn"`
}

type Config struct {
	Debug    bool     `mapstructure:"debug"`
	Server   Server   `mapstructure:"server"`
	Telegram Telegram `mapstructure:"telegram"`
	Storage  Storage  `mapstructure:"storage"`
}

func IsDebug() bool {
	return debug
}

func LoadConfig(file, extension string) (Config, error) {
	v := viper.New()
	v.SetConfigName(file)
	v.SetConfigType(extension)
	v.AddConfigPath(configPath)

	if err := v.ReadInConfig(); err != nil {
		return Config{}, fmt.Errorf("failed to read config: %w", err)
	}

	if err := gotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return Config{}, fmt.Errorf("failed to load .env: %w", err)
	}

	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()
	bindEnvKeys(v)

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return Config{}, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	cfg.Storage.Postgres.DSN = cfg.Storage.Postgres.DSNValue()
	debug = cfg.Debug

	return cfg, nil
}

func bindEnvKeys(v *viper.Viper) {
	_ = v.BindEnv("debug", "DEBUG")
	_ = v.BindEnv("server.http.port", "HTTP_PORT")
	_ = v.BindEnv("telegram.bot_token", "TELEGRAM_BOT_TOKEN")
	_ = v.BindEnv("storage.postgres.host", "PG_HOST")
	_ = v.BindEnv("storage.postgres.port", "PG_PORT")
	_ = v.BindEnv("storage.postgres.database", "PG_DATABASE")
	_ = v.BindEnv("storage.postgres.user", "PG_USER")
	_ = v.BindEnv("storage.postgres.password", "PG_PASSWORD")
	_ = v.BindEnv("storage.postgres.sslmode", "PG_SSLMODE")
	_ = v.BindEnv("storage.postgres.connect_timeout", "PG_CONNECT_TIMEOUT")
	_ = v.BindEnv("storage.postgres.max_open_conns", "PG_MAX_OPEN_CONNS")
	_ = v.BindEnv("storage.postgres.max_idle_conns", "PG_MAX_IDLE_CONNS")
	_ = v.BindEnv("storage.postgres.dsn", "PG_DSN")
}

func (p Postgres) DSNValue() string {
	if p.DSN != "" {
		return p.DSN
	}

	parts := make([]string, 0, 8)
	if p.Host != "" {
		parts = append(parts, fmt.Sprintf("host=%s", p.Host))
	}
	if p.Port != 0 {
		parts = append(parts, fmt.Sprintf("port=%d", p.Port))
	}
	if p.User != "" {
		parts = append(parts, fmt.Sprintf("user=%s", p.User))
	}
	if p.Password != "" {
		parts = append(parts, fmt.Sprintf("password=%s", p.Password))
	}
	if p.Database != "" {
		parts = append(parts, fmt.Sprintf("dbname=%s", p.Database))
	}
	if p.SSLMode != "" {
		parts = append(parts, fmt.Sprintf("sslmode=%s", p.SSLMode))
	}
	if p.ConnectTimeout != 0 {
		parts = append(parts, fmt.Sprintf("connect_timeout=%d", p.ConnectTimeout))
	}

	return strings.Join(parts, " ")
}
