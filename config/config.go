package config

type Config struct {
	TelegramBotID      string `default:"{{ env \"TELEGRAM_BOT\" }}"`
	DatabaseConnection string `default:"{{ env \"DATABASE\" }}"`
	CollectCycle       string `default:"{{ env \"COLLECT_CRON\" }}"`
	RotateAfter        string `default:"{{ env \"ROTATE_AGE\" }}"`
}
