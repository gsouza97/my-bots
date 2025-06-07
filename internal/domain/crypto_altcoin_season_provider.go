package domain

import (
	"github.com/gsouza97/my-bots/internal/adapter/provider"
)

type CryptoAltcoinSeasonProvider interface {
	GetAltcoinSeasonIndex() (provider.CmcAltcoinSeasonIndexResponse, error)
}
