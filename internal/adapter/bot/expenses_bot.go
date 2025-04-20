package bot

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/gsouza97/my-bots/internal/constants"
	"github.com/gsouza97/my-bots/internal/logger"
	"github.com/gsouza97/my-bots/internal/usecase"
)

type ExpensesBot struct {
	adapter               *TelegramAdapter
	saveUseCase           *usecase.SaveBill
	generateReportUseCase *usecase.GenerateReport
}

func NewExpensesBot(adapter *TelegramAdapter, saveUseCase *usecase.SaveBill, generateReportUseCase *usecase.GenerateReport) *ExpensesBot {
	return &ExpensesBot{
		adapter:               adapter,
		saveUseCase:           saveUseCase,
		generateReportUseCase: generateReportUseCase,
	}
}

func (eb *ExpensesBot) Start() error {
	logger.Log.Info("Starting ExpensesBot")
	updates, err := eb.adapter.HandleUpdates()
	if err != nil {
		return err
	}

	for update := range updates {
		if update.Message.IsCommand() {
			eb.processCommand(update)
		} else {
			eb.adapter.SendMessage(update.Message.Chat.ID, "Comando inválido.")
		}
	}

	return nil
}

func (eb *ExpensesBot) processCommand(update tgbotapi.Update) {
	command := update.Message.Command()
	args := update.Message.CommandArguments()
	chatID := update.Message.Chat.ID

	var response string

	switch command {
	case constants.StartCommand:
		response = "Olá! Eu sou o ExpensesBot, um bot para gerenciar suas contas."
	case constants.SaveCommand:
		response = eb.handleSave(args)
	case constants.ReportCommand:
		response = eb.handleReport(args)
	default:
		response = "Comando desconhecido."
	}

	eb.adapter.SendMessage(chatID, response)
}

func (eb *ExpensesBot) handleSave(msg string) string {
	ctx := context.Background()
	response, err := eb.saveUseCase.Execute(ctx, msg)
	if err != nil {
		return err.Error()
	}

	return response
}

func (eb *ExpensesBot) handleReport(msg string) string {
	ctx := context.Background()
	response, err := eb.generateReportUseCase.Execute(ctx, msg)
	if err != nil {
		return err.Error()
	}

	return response
}
