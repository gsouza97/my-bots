package usecase

import (
	"fmt"

	"github.com/gsouza97/my-bots/internal/adapter/provider"
	"github.com/gsouza97/my-bots/internal/domain"
)

type GetAltcoinSeasonIndex struct {
	altcoinSeasonProvider domain.CryptoAltcoinSeasonProvider
}

func NewGetAltcoinSeasonIndex(altcoinSeasonProvider domain.CryptoAltcoinSeasonProvider) *GetAltcoinSeasonIndex {
	return &GetAltcoinSeasonIndex{
		altcoinSeasonProvider: altcoinSeasonProvider,
	}
}

func (uc *GetAltcoinSeasonIndex) Execute() (string, error) {
	data, err := uc.altcoinSeasonProvider.GetAltcoinSeasonIndex()
	if err != nil {
		return "", err
	}

	msg := buildAltcoinSeasonMessage(data)
	return msg, nil
}

func buildAltcoinSeasonMessage(data provider.CmcAltcoinSeasonIndexResponse) string {
	return fmt.Sprintf(
		"Altcoin Season Index: %s → %s",
		data.AltcoinIndex, data.Name,
	)
}
