package bot

import (
	"context"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/gsouza97/my-bots/internal/constants"
	"github.com/gsouza97/my-bots/internal/logger"
	"github.com/gsouza97/my-bots/internal/usecase"
)

type PoolsBot struct {
	adapter                *TelegramAdapter
	listActivePoolsUseCase *usecase.ListActivePools
	getPoolFeesUseCase     *usecase.GetPoolFees
	chatID                 int64
}

func NewPoolsBot(adapter *TelegramAdapter, listActivePoolsUseCase *usecase.ListActivePools, getPoolFeesUseCase *usecase.GetPoolFees, chatID string) *PoolsBot {
	chatIDInt, err := strconv.ParseInt(chatID, 10, 64)
	if err != nil {
		return nil
	}

	return &PoolsBot{
		adapter:                adapter,
		listActivePoolsUseCase: listActivePoolsUseCase,
		getPoolFeesUseCase:     getPoolFeesUseCase,
		chatID:                 chatIDInt,
	}
}

func (pb *PoolsBot) Start() error {
	logger.Log.Info("Starting PoolsBot")
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

func (pb *PoolsBot) processCommand(update tgbotapi.Update) {
	command := update.Message.Command()
	args := update.Message.CommandArguments()
	chatID := update.Message.Chat.ID

	var response string

	switch command {
	case constants.StartCommand:
		response = "Olá! Eu sou o PoolsBot, um bot para gerenciar suas contas."
	case constants.PoolsCommand:
		response = pb.handlePools(args)
	case constants.FeesCommand:
		response = pb.handleFees(args)
	default:
		response = "Comando desconhecido."
	}

	pb.adapter.SendMessage(chatID, response)
}

func (pb *PoolsBot) handlePools(msg string) string {
	ctx := context.Background()
	response, err := pb.listActivePoolsUseCase.Execute(ctx)
	if err != nil {
		return err.Error()
	}

	return response
}

func (pb *PoolsBot) handleFees(msg string) string {
	ctx := context.Background()
	response, err := pb.getPoolFeesUseCase.Execute(ctx)
	if err != nil {
		return err.Error()
	}

	return response
}

func (pb *PoolsBot) SendMessage(message string) error {
	logger.Log.Infof("Sending message to chat: %d", pb.chatID)
	err := pb.adapter.SendMessage(pb.chatID, message)
	return err
}
