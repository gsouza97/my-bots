package domain

type CryptoPriceProvider interface {
	GetPrice(crypto string, others ...string) (float64, error)
	GetMultiplePrices(cryptos []string) (map[string]float64, error)
}
