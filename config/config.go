package config

import (
	"os"

	"github.com/gsouza97/my-bots/internal/logger"
	"github.com/joho/godotenv"
)

type Config struct {
	ExpensesBotToken        string
	PriceAlertsBotToken     string
	PoolsBotToken           string
	HomologacionBotToken    string
	MongoURI                string
	BotChatID               string
	NotificationCooldown    string
	MongoDBName             string
	DailyAlertCron          string
	AlertMonitorCron        string
	PoolsMonitorCron        string
	HomologacionMonitorCron string
}

func LoadConfig() *Config {
	err := godotenv.Load("../.env")
	if err != nil {
		logger.Log.Info("No .env file found. Using production environment variables.")
	}
	return &Config{
		ExpensesBotToken:        os.Getenv("TELEGRAM_EXPENSES_BOT_TOKEN"),
		PriceAlertsBotToken:     os.Getenv("TELEGRAM_PRICE_ALERTS_BOT_TOKEN"),
		PoolsBotToken:           os.Getenv("TELEGRAM_POOLS_BOT_TOKEN"),
		HomologacionBotToken:    os.Getenv("TELEGRAM_HOMOLOGACION_BOT_TOKEN"),
		MongoURI:                os.Getenv("MONGO_URI"),
		BotChatID:               os.Getenv("TELEGRAM_CHAT_ID"),
		NotificationCooldown:    os.Getenv("TELEGRAM_NOTIFICATION_COOLDOWN_SECONDS"),
		MongoDBName:             os.Getenv("MONGO_DB_NAME"),
		DailyAlertCron:          os.Getenv("DAILY_ALERT_CRON"),
		AlertMonitorCron:        os.Getenv("ALERT_MONITOR_CRON"),
		PoolsMonitorCron:        os.Getenv("POOLS_MONITOR_CRON"),
		HomologacionMonitorCron: os.Getenv("HOMOLOGACION_MONITOR_CRON"),
	}
}
