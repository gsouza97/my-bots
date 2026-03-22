package helper

import "github.com/gsouza97/my-bots/internal/domain"

func ExtractPriceAlertsAssets(priceAlerts []*domain.PriceAlert) []string {
	assetsSet := make(map[string]bool)

	for _, priceAlert := range priceAlerts {
		assetsSet[priceAlert.Crypto] = true
	}

	// Converter map em slice
	var assetsList []string
	for asset := range assetsSet {
		assetsList = append(assetsList, asset)
	}

	return assetsList
}
