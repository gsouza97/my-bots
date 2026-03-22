package composer

import (
	"github.com/gsouza97/my-bots/internal/domain"
	"github.com/gsouza97/my-bots/internal/infrastructure/event_publishing"
	"github.com/gsouza97/my-bots/internal/interfaces/bot"
	"github.com/gsouza97/my-bots/internal/logger"
)

// EventPublishingComposer gerencia a inicialização do EventPublisher
// e o registro de todos os subscribers (listeners)
type EventPublishingComposer struct {
	EventPublisher domain.EventPublisher
}

func NewEventPublishingComposer(adapters *AdaptersComposer, chatID int64) (*EventPublishingComposer, error) {

	// 1. Criar o publisher local
	publisher := event_publishing.NewLocalEventPublisher()

	// 2. Criar listeners (subscribers)
	priceAlertListener := bot.NewTelegramNotificationListener(
		adapters.TelegramPriceAlertsAdapter,
		chatID,
	)

	dailyAlertListener := bot.NewTelegramNotificationListener(
		adapters.TelegramPriceAlertsAdapter,
		chatID,
	)

	poolAlertListener := bot.NewTelegramNotificationListener(
		adapters.TelegramPoolsAdapter,
		chatID,
	)

	loanAlertListener := bot.NewTelegramNotificationListener(
		adapters.TelegramLoansAdapter,
		chatID,
	)

	homologacionAlertListener := bot.NewTelegramNotificationListener(
		adapters.TelegramHomologacionAdapter,
		chatID,
	)

	// 3. Registrar handlers no publisher

	// Price Alerts
	if err := publisher.Subscribe("PriceAlertTriggered", priceAlertListener.HandlePriceAlertTriggered); err != nil {
		logger.Log.Errorf("Error subscribing to PriceAlertTriggered: %v", err)
		return nil, err
	}

	// Daily Alerts
	if err := publisher.Subscribe("DailyAlertTriggered", dailyAlertListener.HandleDailyAlertTriggered); err != nil {
		logger.Log.Errorf("Error subscribing to DailyAlertTriggered: %v", err)
		return nil, err
	}

	// Pool Alerts
	if err := publisher.Subscribe("PoolAlertTriggered", poolAlertListener.HandlePoolAlertTriggered); err != nil {
		logger.Log.Errorf("Error subscribing to PoolAlertTriggered: %v", err)
		return nil, err
	}

	// Loan Alerts
	if err := publisher.Subscribe("LoanAlertTriggered", loanAlertListener.HandleLoanAlertTriggered); err != nil {
		logger.Log.Errorf("Error subscribing to LoanAlertTriggered: %v", err)
		return nil, err
	}

	// Homologacion Alerts
	if err := publisher.Subscribe("HomologacionAlertTriggered", homologacionAlertListener.HandleHomologacionAlertTriggered); err != nil {
		logger.Log.Errorf("Error subscribing to HomologacionAlertTriggered: %v", err)
		return nil, err
	}

	logger.Log.Info("EventPublisher initialized with all subscribers registered")

	return &EventPublishingComposer{
		EventPublisher: publisher,
	}, nil
}
