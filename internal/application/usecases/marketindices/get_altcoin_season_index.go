package marketindices

import (
	"fmt"

	"github.com/gsouza97/my-bots/internal/domain"
	"github.com/gsouza97/my-bots/internal/infrastructure/providers"
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

func buildAltcoinSeasonMessage(data providers.CmcAltcoinSeasonIndexResponse) string {
	if data.Name == "" {
		return fmt.Sprintf(
			"Altcoin Season Index: %s",
			data.AltcoinIndex,
		)
	} else {
		return fmt.Sprintf(
			"Altcoin Season Index: %s → %s",
			data.AltcoinIndex, data.Name,
		)
	}
}
