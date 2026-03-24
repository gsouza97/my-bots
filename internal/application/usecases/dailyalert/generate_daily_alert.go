package dailyalert

import (
	"context"
	"fmt"
	"strings"

	"github.com/gsouza97/my-bots/internal/application/usecases/marketindices"
	"github.com/gsouza97/my-bots/internal/application/usecases/pools"
	"github.com/gsouza97/my-bots/internal/domain"
	"github.com/gsouza97/my-bots/internal/domain/events"
	"github.com/gsouza97/my-bots/internal/infrastructure/repository"
	"github.com/gsouza97/my-bots/internal/logger"
)

type GenerateDailyAlert struct {
	getPoolFees      *pools.GetPoolFees
	getFearAndGreed  *marketindices.GetFearAndGreedIndex
	getAltcoinSeason *marketindices.GetAltcoinSeasonIndex
	alertRepository  repository.PriceAlertRepository
	priceProvider    domain.CryptoPriceProvider
	eventPublisher   domain.EventPublisher
}

func NewGenerateDailyAlert(getPoolFees *pools.GetPoolFees, getFearAndGreed *marketindices.GetFearAndGreedIndex, getAltcoinSeason *marketindices.GetAltcoinSeasonIndex, alertRepository repository.PriceAlertRepository, priceProvider domain.CryptoPriceProvider, eventPublisher domain.EventPublisher) *GenerateDailyAlert {
	return &GenerateDailyAlert{
		getPoolFees:      getPoolFees,
		getFearAndGreed:  getFearAndGreed,
		getAltcoinSeason: getAltcoinSeason,
		alertRepository:  alertRepository,
		priceProvider:    priceProvider,
		eventPublisher:   eventPublisher,
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

	fearAndGreedMsg, err := gda.getFearAndGreed.Execute()
	if err != nil {
		return fmt.Errorf("error getting fear and greed index: %w", err)
	}

	altcoinSeasonMsg, err := gda.getAltcoinSeason.Execute()
	if err != nil {
		return fmt.Errorf("error getting altcoin season index: %w", err)
	}

	dailyAlertMsg = append(dailyAlertMsg, alertsPriceMsgs...)
	dailyAlertMsg = append(dailyAlertMsg, "\n")
	dailyAlertMsg = append(dailyAlertMsg, poolsMsg)
	dailyAlertMsg = append(dailyAlertMsg, "\n")
	dailyAlertMsg = append(dailyAlertMsg, "📌 "+fearAndGreedMsg)
	dailyAlertMsg = append(dailyAlertMsg, "\n")
	dailyAlertMsg = append(dailyAlertMsg, "📌 "+altcoinSeasonMsg)

	event := events.NewDailyAlertTriggeredEvent(strings.Join(dailyAlertMsg, "\n"))

	err = gda.eventPublisher.Publish(ctx, event)
	if err != nil {
		return fmt.Errorf("error publishing event: %w", err)
	}
	return nil
}
