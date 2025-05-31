package usecase

import (
	"fmt"

	"github.com/gsouza97/my-bots/internal/domain"
	"github.com/gsouza97/my-bots/internal/repository"
	"github.com/gsouza97/my-bots/pkg/parser"
)

type CheckPrice struct {
	priceAlertRepository repository.PriceAlertRepository
	priceProvider        domain.CryptoPriceProvider
}

func NewCheckPrice(priceAlertRepository repository.PriceAlertRepository, priceProvider domain.CryptoPriceProvider) *CheckPrice {
	return &CheckPrice{
		priceAlertRepository: priceAlertRepository,
		priceProvider:        priceProvider,
	}
}

func (uc *CheckPrice) Execute(message string) (string, error) {
	crypto, err := parser.ParseCheckPriceMesage(message)
	if err != nil {
		return "", err
	}
	price, err := uc.priceProvider.GetPrice(crypto)
	if err != nil {
		return "Erro ao buscar preço", fmt.Errorf("error checking price: %w", err)
	}
	return fmt.Sprintf(
		"%s: %f USD",
		crypto, price,
	), nil
}
