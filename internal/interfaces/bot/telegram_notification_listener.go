package bot

import (
	"context"
	"fmt"

	adaptorbot "github.com/gsouza97/my-bots/internal/adapter/bot"
	"github.com/gsouza97/my-bots/internal/domain"
	"github.com/gsouza97/my-bots/internal/domain/events"
	"github.com/gsouza97/my-bots/internal/logger"
)

// TelegramNotificationListener escuta eventos de domínio e envia notificações via Telegram
// Implementa EventHandler para cada tipo de evento que quer processar
type TelegramNotificationListener struct {
	telegramAdapter *adaptorbot.TelegramAdapter
	chatID          int64
}

func NewTelegramNotificationListener(adapter *adaptorbot.TelegramAdapter, chatID int64) *TelegramNotificationListener {
	return &TelegramNotificationListener{
		telegramAdapter: adapter,
		chatID:          chatID,
	}
}

func (tnl *TelegramNotificationListener) HandlePriceAlertTriggered(ctx context.Context, event events.Event) error {
	priceEvent, ok := event.(*events.PriceAlertTriggeredEvent)
	if !ok {
		return fmt.Errorf("invalid event type, expected PriceAlertTriggeredEvent")
	}

	emoji := "📈"
	if priceEvent.Status == "OVER_PRICE" {
		emoji = "🚀"
	} else if priceEvent.Status == "UNDER_PRICE" {
		emoji = "📉"
	}

	message := fmt.Sprintf(
		"%s *Alerta de Preço!*\n\n💰 *%s*\nAtual: $%.2f\nAlvo: $%.2f\nStatus: %s",
		emoji,
		priceEvent.Crypto,
		priceEvent.CurrentPrice,
		priceEvent.AlertPrice,
		priceEvent.Status,
	)

	message = getAlertPriceMessage(priceEvent.Crypto, priceEvent.AlertPrice, priceEvent.CurrentPrice, domain.AlertPriceStatus(priceEvent.Status))

	logger.Log.Infof("Sending price alert for %s to Telegram", priceEvent.Crypto)

	return tnl.telegramAdapter.SendMessage(tnl.chatID, message)
}

func (tnl *TelegramNotificationListener) HandleDailyAlertTriggered(ctx context.Context, event events.Event) error {
	dailyEvent, ok := event.(*events.DailyAlertTriggeredEvent)
	if !ok {
		return fmt.Errorf("invalid event type, expected DailyAlertTriggeredEvent")
	}

	logger.Log.Infof("Sending daily alert to Telegram")

	return tnl.telegramAdapter.SendMessage(tnl.chatID, dailyEvent.Message)
}

func (tnl *TelegramNotificationListener) HandlePoolAlertTriggered(ctx context.Context, event events.Event) error {
	poolEvent, ok := event.(*events.PoolAlertTriggeredEvent)
	if !ok {
		return fmt.Errorf("invalid event type, expected PoolAlertTriggeredEvent")
	}

	logger.Log.Infof("Sending pool alert for %s to Telegram", poolEvent.PoolID)

	return tnl.telegramAdapter.SendMessage(tnl.chatID, poolEvent.Message)
}

func (tnl *TelegramNotificationListener) HandleLoanAlertTriggered(ctx context.Context, event events.Event) error {
	loanEvent, ok := event.(*events.LoanAlertTriggeredEvent)
	if !ok {
		return fmt.Errorf("invalid event type, expected LoanAlertTriggeredEvent")
	}

	logger.Log.Infof("Sending loan alert for %s to Telegram", loanEvent.LoanID)

	return tnl.telegramAdapter.SendMessage(tnl.chatID, loanEvent.Message)
}

func (tnl *TelegramNotificationListener) HandleHomologacionAlertTriggered(ctx context.Context, event events.Event) error {
	homEvent, ok := event.(*events.HomologacionAlertTriggeredEvent)
	if !ok {
		return fmt.Errorf("invalid event type, expected HomologacionAlertTriggeredEvent")
	}

	logger.Log.Infof("Sending homologacion alert to Telegram")

	return tnl.telegramAdapter.SendMessage(tnl.chatID, homEvent.Message)
}

func getAlertPriceMessage(crypto string, alertPrice float64, currentPrice float64, newAlertStatus domain.AlertPriceStatus) string {
	statusStr := "ABAIXO"
	if newAlertStatus == domain.OverPrice {
		statusStr = "ACIMA"
	}
	return fmt.Sprintf("🚨 ALERTA: %s está %s de %f USD! Preço atual: %.4f USD.", crypto, statusStr, alertPrice, currentPrice)
}
