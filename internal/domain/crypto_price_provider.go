package domain

type CryptoPriceProvider interface {
	GetPrice(crypto string, others ...string) (float64, error)
}
