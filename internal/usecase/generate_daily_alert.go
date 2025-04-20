package usecase

import (
	"context"
	"fmt"
	"strings"

	"github.com/gsouza97/my-bots/internal/domain"
	"github.com/gsouza97/my-bots/internal/logger"
	"github.com/gsouza97/my-bots/internal/repository"
)

type GenerateDailyAlert struct {
	getPoolFees     *GetPoolFees
	alertRepository repository.PriceAlertRepository
	priceProvider   domain.CryptoPriceProvider
	notifier        domain.Notifier
}

func NewGenerateDailyAlert(getPoolFees *GetPoolFees, alertRepository repository.PriceAlertRepository, priceProvider domain.CryptoPriceProvider, notifier domain.Notifier) *GenerateDailyAlert {
	return &GenerateDailyAlert{
		getPoolFees:     getPoolFees,
		alertRepository: alertRepository,
		priceProvider:   priceProvider,
		notifier:        notifier,
	}
}

func (gda *GenerateDailyAlert) Execute() error {
	ctx := context.Background()
	dailyAlertMsg := []string{"📌 Alerta Diário\n"}

	alertsPriceMsgs := []string{}
	alerts, err := gda.alertRepository.FindAllByActiveIsTrue(ctx)
	if err != nil {
		return fmt.Errorf("error getting active alerts: %w", err)
	}
	alertsUnique := make(map[string]bool)
	var uniqueList []string
	for _, alert := range alerts {
		if !alertsUnique[alert.Crypto] {
			alertsUnique[alert.Crypto] = true
			uniqueList = append(uniqueList, alert.Crypto)
		}
	}
	logger.Log.Infof("uniqueList: %v", uniqueList)

	for _, crypto := range uniqueList {
		price, err := gda.priceProvider.GetPrice(crypto)
		if err != nil {
			return fmt.Errorf("error getting price for %s: %w", crypto, err)
		}
		alertsPriceMsgs = append(alertsPriceMsgs, fmt.Sprintf("📈 %s: %.4f", crypto, price))
	}

	poolsMsg, err := gda.getPoolFees.ExecuteAndUpdateLastFees(ctx)
	if err != nil {
		return err
	}
	dailyAlertMsg = append(dailyAlertMsg, alertsPriceMsgs...)
	dailyAlertMsg = append(dailyAlertMsg, "\n")
	dailyAlertMsg = append(dailyAlertMsg, poolsMsg)

	gda.notifier.SendMessage(strings.Join(dailyAlertMsg, "\n"))
	if err != nil {
		return fmt.Errorf("error sending message: %w", err)
	}
	return nil
}
