package usecase

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/gsouza97/my-bots/internal/domain"
	"github.com/gsouza97/my-bots/internal/logger"
	"github.com/gsouza97/my-bots/internal/repository"
)

type CheckPriceAlert struct {
	alertRepository repository.PriceAlertRepository
	priceProvider   domain.CryptoPriceProvider
	notifier        domain.Notifier
}

func NewCheckPriceAlert(alertRepository repository.PriceAlertRepository, priceProvider domain.CryptoPriceProvider, notifier domain.Notifier) *CheckPriceAlert {
	return &CheckPriceAlert{
		alertRepository: alertRepository,
		priceProvider:   priceProvider,
		notifier:        notifier,
	}
}

func (cpa *CheckPriceAlert) Execute() error {
	ctx := context.Background()
	alerts, err := cpa.alertRepository.FindAllByActiveIsTrue(ctx)
	if err != nil {
		return err
	}

	numWorkers := 8
	alertChannel := make(chan *domain.PriceAlert, numWorkers)
	errChannel := make(chan error, 1)
	t := time.Now()

	var wg sync.WaitGroup

	for _, alert := range alerts {
		wg.Add(1)

		alertChannel <- alert

		go func(alert *domain.PriceAlert) {
			defer wg.Done()
			defer func() { <-alertChannel }()

			err := cpa.processAlert(ctx, alert)
			if err != nil {
				select {
				case errChannel <- err:
				default:
				}
			}
		}(alert)

	}
	wg.Wait()
	t2 := time.Now()
	logger.Log.Infof("Tempo total para check_price_alert: %s", t2.Sub(t))

	select {
	case err := <-errChannel:
		return err
	default:
		return nil
	}
}

func (cpa *CheckPriceAlert) processAlert(ctx context.Context, alert *domain.PriceAlert) error {
	price, err := cpa.priceProvider.GetPrice(alert.Crypto)
	if err != nil {
		return fmt.Errorf("error getting price for %s: %w", alert.Crypto, err)
	}
	logger.Log.Infof("price for %s: %f", alert.Crypto, price)

	newAlertStatus := domain.UnderPrice
	if price >= alert.AlertPrice {
		newAlertStatus = domain.OverPrice
	}

	if newAlertStatus != alert.PriceStatus {
		message := cpa.getAlertPriceMessage(alert, price, newAlertStatus)
		err := cpa.notifier.SendMessage(message)
		if err != nil {
			return fmt.Errorf("error sending message: %w", err)
		}
		alert.PriceStatus = newAlertStatus
		err = cpa.updatePriceAlert(ctx, alert, newAlertStatus)
		if err != nil {
			return fmt.Errorf("error updating alert: %w", err)
		}
	}
	return nil
}

func (cpa *CheckPriceAlert) updatePriceAlert(ctx context.Context, alert *domain.PriceAlert, alertStatus domain.AlertPriceStatus) error {
	alert.PriceStatus = alertStatus
	err := cpa.alertRepository.Update(ctx, alert)
	return err
}

func (cpa *CheckPriceAlert) getAlertPriceMessage(alert *domain.PriceAlert, price float64, newAlertStatus domain.AlertPriceStatus) string {
	statusStr := "ABAIXO"
	if newAlertStatus == domain.OverPrice {
		statusStr = "ACIMA"
	}
	return fmt.Sprintf("🚨 ALERTA: %s está %s de %.2f USD! Preço atual: %.2f USD 🚀", alert.Crypto, statusStr, alert.AlertPrice, price)
}
