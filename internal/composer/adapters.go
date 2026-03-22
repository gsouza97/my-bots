package composer

import (
	"github.com/gsouza97/my-bots/config"
	"github.com/gsouza97/my-bots/internal/adapter/bot"
	"github.com/gsouza97/my-bots/internal/logger"
)

type AdaptersComposer struct {
	TelegramExpensesAdapter     *bot.TelegramAdapter
	TelegramPriceAlertsAdapter  *bot.TelegramAdapter
	TelegramPoolsAdapter        *bot.TelegramAdapter
	TelegramHomologacionAdapter *bot.TelegramAdapter
	TelegramLoansAdapter        *bot.TelegramAdapter
}

func NewAdaptersComposer(cfg *config.Config) (*AdaptersComposer, error) {
	var err error
	ac := &AdaptersComposer{}

	// Inicializa cada adapter, capturando erros
	ac.TelegramExpensesAdapter, err = bot.NewTelegramAdapter(cfg.ExpensesBotToken)
	if err != nil {
		logger.Log.Errorf("Erro ao inicializar Telegram Expenses Adapter: %v", err)
		return nil, err
	}

	ac.TelegramPriceAlertsAdapter, err = bot.NewTelegramAdapter(cfg.PriceAlertsBotToken)
	if err != nil {
		logger.Log.Errorf("Erro ao inicializar Telegram Price Alerts Adapter: %v", err)
		return nil, err
	}

	ac.TelegramPoolsAdapter, err = bot.NewTelegramAdapter(cfg.PoolsBotToken)
	if err != nil {
		logger.Log.Errorf("Erro ao inicializar Telegram Pools Adapter: %v", err)
		return nil, err
	}

	ac.TelegramHomologacionAdapter, err = bot.NewTelegramAdapter(cfg.HomologacionBotToken)
	if err != nil {
		logger.Log.Errorf("Erro ao inicializar Telegram Homologacion Adapter: %v", err)
		return nil, err
	}

	ac.TelegramLoansAdapter, err = bot.NewTelegramAdapter(cfg.LoansBotToken)
	if err != nil {
		logger.Log.Errorf("Erro ao inicializar Telegram Loans Adapter: %v", err)
		return nil, err
	}

	return ac, nil
}
