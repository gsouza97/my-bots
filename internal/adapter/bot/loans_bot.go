package bot

import (
	"context"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/gsouza97/my-bots/internal/constants"
	"github.com/gsouza97/my-bots/internal/logger"
	"github.com/gsouza97/my-bots/internal/usecase"
)

type LoansBot struct {
	adapter         *TelegramAdapter
	getLoansUseCase *usecase.GetLoans
	chatID          int64
}

func NewLoansBot(adapter *TelegramAdapter, getLoansUseCase *usecase.GetLoans, chatID string) *LoansBot {
	chatIDInt, err := strconv.ParseInt(chatID, 10, 64)
	if err != nil {
		return nil
	}

	return &LoansBot{
		adapter:         adapter,
		getLoansUseCase: getLoansUseCase,
		chatID:          chatIDInt,
	}
}

func (lb *LoansBot) Start() error {
	logger.Log.Info("Starting LoansBot")
	updates, err := lb.adapter.HandleUpdates()
	if err != nil {
		return err
	}

	for update := range updates {
		if update.Message.IsCommand() {
			lb.processCommand(update)
		} else {
			lb.adapter.SendMessage(update.Message.Chat.ID, "Comando inválido.")
		}
	}

	return nil
}

func (lb *LoansBot) processCommand(update tgbotapi.Update) {
	command := update.Message.Command()
	args := update.Message.CommandArguments()
	chatID := update.Message.Chat.ID

	var response string

	switch command {
	case constants.StartCommand:
		response = "Olá! Eu sou o LoansBot, um bot para te enviar alertas sobre os empréstimos."
	case constants.LoansCommand:
		response = lb.handleLoans(args)
	default:
		response = "Comando desconhecido."
	}

	lb.adapter.SendMessage(chatID, response)
}

func (lb *LoansBot) handleLoans(msg string) string {
	ctx := context.Background()
	response, err := lb.getLoansUseCase.Execute(ctx)
	if err != nil {
		return err.Error()
	}

	return response
}

func (lb *LoansBot) SendMessage(message string) error {
	logger.Log.Infof("Sending message to chat: %d", lb.chatID)
	err := lb.adapter.SendMessage(lb.chatID, message)
	return err
}
