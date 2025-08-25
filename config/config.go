package config

import (
	// Autoload .env file.
	"time"

	_ "github.com/joho/godotenv/autoload"
)

type Config struct {
	TelegramBotID      string        `default:"{{ env \"TELEGRAM_BOT\" }}" validate:"required"`
	DatabaseConnection string        `default:"{{ or (env \"DATABASE\") \"sqlite:///outage.db\" }}" validate:"required,uri"`
	CollectCycle       string        `default:"{{ or (env \"COLLECT_CRON\") \"0 0 * * * *\" }}"`
	RotateAfter        time.Duration `default:"{{ or (env \"ROTATE_AGE\") \"1h\" | parseDuration }}"`
}
