package config

import (
	"time"

	// Autoload .env file.
	_ "github.com/joho/godotenv/autoload"

	sc "github.com/fmotalleb/scrapper-go/config"
)

type Config struct {
	HTTPListenAddr string `mapstructure:"http_listen" default:"{{ env \"HTTP_LISTEN\" }}"`

	TelegramBotKey     string `mapstructure:"telegram_bot" default:"{{ env \"TELEGRAM_BOT\" }}" validate:"required"`
	DatabaseConnection string `mapstructure:"database" default:"{{ or (env \"DATABASE\") \"sqlite:///outage.db\" }}" validate:"required,uri"`

	CollectCycle    string             `mapstructure:"collect_cycle" default:"{{ or (env \"COLLECT_CRON\") \"0 0 * * * *\" }}" validate:"required,cron"`
	CollectTimeout  time.Duration      `mapstructure:"collect_timeout" default:"{{ or (env \"COLLECT_TIMEOUT\") \"1h\" | parseDuration }}"`
	CollectorConfig sc.ExecutionConfig `mapstructure:"collector"`

	CollectOnStart          *bool         `mapstructure:"collect_on_start" default:"{{ or (env \"COLLECT_ON_START\") \"true\" }}"`
	CollectOnStartThreshold time.Duration `mapstructure:"collect_on_start_threshold" default:"{{ or (env \"COLLECT_ON_START_THRESHOLD\") \"10m\" }}"`

	RotateAfter time.Duration `mapstructure:"max_age" default:"{{ or (env \"MAX_AGE\") \"1h\" | parseDuration }}"`
}
