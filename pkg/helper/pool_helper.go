package helper

import (
	"fmt"
	"strconv"

	"github.com/gsouza97/my-bots/internal/adapter/provider"
	"github.com/gsouza97/my-bots/internal/domain"
)

func BuildPoolResponseMessage(pools []*domain.Pool) string {
	message := "📌 *Pools Ativas:* \n"
	for _, pool := range pools {
		poolStatus := "DENTRO"
		if pool.OutOfRange {
			poolStatus = "FORA"
		}
		message += fmt.Sprintf(
			"- *%s/%s*:\n  Pool: %s\n  Range: %f - %f  ➡  *%s*\n\n",
			pool.Crypto1, pool.Crypto2, pool.Description, pool.MinPrice, pool.MaxPrice, poolStatus,
		)
	}
	return message
}

func BuildFeesToCollectMessage(pool *domain.Pool, totalFeesToCollect float64, feesToCollectVariation float64) string {
	var arrow string
	switch {
	case feesToCollectVariation > 0:
		arrow = "↑"
	case feesToCollectVariation < 0:
		arrow = "↓"
	default:
		arrow = ""
	}

	poolRangeStatus := "DENTRO"
	if pool.OutOfRange {
		poolRangeStatus = "FORA"
	}

	message := fmt.Sprintf(
		"\n- Pool: %s\n Status: %s\n Fees para coletar: %.2f USD (%s%.2f)",
		pool.Description, poolRangeStatus, totalFeesToCollect, arrow, feesToCollectVariation,
	)

	return message
}

func BuildRangeMessage(pool *domain.Pool, outOfRange bool, price float64) string {
	message := fmt.Sprintf("ATENÇÃO: POOL SAIU DO RANGE!\n Pool: %s\n %s/%s\n Preço: %.2f\n Preço fora do range configurado.", pool.Description, pool.Crypto1, pool.Crypto2, price)
	if !outOfRange {
		message = fmt.Sprintf("ATENÇÃO: POOL VOLTOU AO RANGE!\n Pool: %s\n %s/%s\n Preço: %.2f\n Preço voltou ao range configurado.", pool.Description, pool.Crypto1, pool.Crypto2, price)
	}
	return message
}

func BuildWarningMessage(pool *domain.Pool, price float64, percent float64, diff float64, maxWarningMessage bool) string {
	message := fmt.Sprintf("ATENÇÃO: PREÇO PRÓXIMO AO TOPO!\n %s\n %s\n Preço: %.2f\n Margem de Risco: %.2f %%\n Precisa subir %.2f (%.6f %s) para atingir o topo.", pool.Description, pool.Crypto1, price, pool.RiskRate*100, percent*100, diff, pool.Crypto2)
	if !maxWarningMessage {
		message = fmt.Sprintf("ATENÇÃO: PREÇO PRÓXIMO AO MÍNIMO!\n %s\n %s\n Preço: %.2f\n Margem de Risco: %.2f %%\n Precisa baixar %.2f (%.6f %s) para atingir o mínimo.", pool.Description, pool.Crypto1, price, pool.RiskRate*100, percent*100, diff, pool.Crypto2)
	}
	return message
}

func CalculateFeesToCollect(data provider.RevertPoolDataResponse) (float64, error) {
	token1Price, err := strconv.ParseFloat(data.Tokens[data.Token0].Price, 64)
	if err != nil {
		return 0, fmt.Errorf("erro ao converter preço do token 1: %w", err)
	}

	token2Price, err := strconv.ParseFloat(data.Tokens[data.Token1].Price, 64)
	if err != nil {
		return 0, fmt.Errorf("erro ao converter preço do token 2: %w", err)
	}

	token1Uncollected, err := strconv.ParseFloat(data.UncollectedFees0, 64)
	if err != nil {
		return 0, fmt.Errorf("erro ao converter fees do token 1: %w", err)
	}

	token2Uncollected, err := strconv.ParseFloat(data.UncollectedFees1, 64)
	if err != nil {
		return 0, fmt.Errorf("erro ao converter fees do token 2: %w", err)
	}

	feesToken1 := token1Uncollected * token1Price
	feesToken2 := token2Uncollected * token2Price
	totalFees := feesToken1 + feesToken2

	return totalFees, nil
}

func CalculatePercentToMaxPrice(price float64, maxPrice float64) float64 {
	diffToMaxPrice := maxPrice - price
	return diffToMaxPrice / price
}

func CalculatePercentToMinPrice(price float64, minPrice float64) float64 {
	diffToMinPrice := price - minPrice
	return diffToMinPrice / price
}
