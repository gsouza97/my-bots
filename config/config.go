package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ExpensesBotToken     string
	PriceAlertsBotToken  string
	PoolsBotToken        string
	MongoURI             string
	BotChatID            string
	NotificationCooldown string
	MongoDBName          string
}

func LoadConfig() *Config {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Println("No .env file found. Using production environment variables.")
	}
	return &Config{
		ExpensesBotToken:     os.Getenv("TELEGRAM_EXPENSES_BOT_TOKEN"),
		PriceAlertsBotToken:  os.Getenv("TELEGRAM_PRICE_ALERTS_BOT_TOKEN"),
		PoolsBotToken:        os.Getenv("TELEGRAM_POOLS_BOT_TOKEN"),
		MongoURI:             os.Getenv("MONGO_URI"),
		BotChatID:            os.Getenv("TELEGRAM_CHAT_ID"),
		NotificationCooldown: os.Getenv("TELEGRAM_NOTIFICATION_COOLDOWN_SECONDS"),
		MongoDBName:          os.Getenv("MONGO_DB_NAME"),
	}
}
