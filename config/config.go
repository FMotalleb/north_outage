package config

import (
	"time"

	// Autoload .env file.
	_ "github.com/joho/godotenv/autoload"

	sc "github.com/fmotalleb/scrapper-go/config"
)

type Config struct {
	TelegramBotID      string             `mapstructure:"telegram_bot" default:"{{ env \"TELEGRAM_BOT\" }}" validate:"required"`
	DatabaseConnection string             `mapstructure:"database" default:"{{ or (env \"DATABASE\") \"sqlite:///outage.db\" }}" validate:"required,uri"`
	CollectCycle       string             `mapstructure:"collect_cycle" default:"{{ or (env \"COLLECT_CRON\") \"0 0 * * * *\" }}" validate:"required,cron"`
	RotateAfter        time.Duration      `mapstructure:"max_age" default:"{{ or (env \"MAX_AGE\") \"1h\" | parseDuration }}"`
	CollectorConfig    sc.ExecutionConfig `mapstructure:"collector"`
}
