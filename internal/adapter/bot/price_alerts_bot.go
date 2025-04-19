package bot

import (
	"log"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/gsouza97/my-bots/internal/constants"
)

type PriceAlertsBot struct {
	adapter *TelegramAdapter
	chatID  int64
}

func NewPriceAlertsBot(adapter *TelegramAdapter, chatID string) *PriceAlertsBot {
	chatIDInt, err := strconv.ParseInt(chatID, 10, 64)
	if err != nil {
		return nil
	}

	return &PriceAlertsBot{
		adapter: adapter,
		chatID:  chatIDInt,
	}
}

func (pab *PriceAlertsBot) Start() error {
	log.Println("Starting PriceAlertsBot")
	updates, err := pab.adapter.HandleUpdates()
	if err != nil {
		return err
	}

	for update := range updates {
		if update.Message.IsCommand() {
			pab.processCommand(update)
		} else {
			pab.adapter.SendMessage(update.Message.Chat.ID, "Comando inválido.")
		}
	}

	return nil
}

func (pab *PriceAlertsBot) SendMessage(message string) error {
	err := pab.adapter.SendMessage(pab.chatID, message)
	return err
}

func (pab *PriceAlertsBot) processCommand(update tgbotapi.Update) {
	command := update.Message.Command()
	chatID := update.Message.Chat.ID

	var response string

	switch command {
	case constants.StartCommand:
		response = "Olá! Eu sou o PriceAlertsBot, um bot para te enviar alertas."
	default:
		response = "Comando desconhecido."
	}

	pab.adapter.SendMessage(chatID, response)
}
