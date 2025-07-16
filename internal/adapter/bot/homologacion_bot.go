package bot

import (
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/gsouza97/my-bots/internal/constants"
	"github.com/gsouza97/my-bots/internal/logger"
)

type HomologacionBot struct {
	adapter *TelegramAdapter
	chatID  int64
}

func NewHomologacionBot(adapter *TelegramAdapter, chatID string) *HomologacionBot {
	chatIDInt, err := strconv.ParseInt(chatID, 10, 64)
	if err != nil {
		return nil
	}

	return &HomologacionBot{
		adapter: adapter,
		chatID:  chatIDInt,
	}
}

func (pb *HomologacionBot) Start() error {
	logger.Log.Info("Starting HomologacionBot")
	updates, err := pb.adapter.HandleUpdates()
	if err != nil {
		return err
	}

	for update := range updates {
		if update.Message.IsCommand() {
			pb.processCommand(update)
		} else {
			pb.adapter.SendMessage(update.Message.Chat.ID, "Comando inválido.")
		}
	}

	return nil
}

func (pb *HomologacionBot) processCommand(update tgbotapi.Update) {
	command := update.Message.Command()
	// args := update.Message.CommandArguments()
	chatID := update.Message.Chat.ID

	var response string

	switch command {
	case constants.StartCommand:
		response = "Olá! Eu sou o HomologacionBot, um bot para consultar sua homologação."
	default:
		response = "Comando desconhecido."
	}

	pb.adapter.SendMessage(chatID, response)
}

func (pb *HomologacionBot) SendMessage(message string) error {
	logger.Log.Infof("Sending message to chat: %d", pb.chatID)
	err := pb.adapter.SendMessage(pb.chatID, message)
	return err
}
