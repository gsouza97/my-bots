package domain

import "github.com/gsouza97/my-bots/internal/infrastructure/providers"

type CryptoAltcoinSeasonProvider interface {
	GetAltcoinSeasonIndex() (providers.CmcAltcoinSeasonIndexResponse, error)
}
