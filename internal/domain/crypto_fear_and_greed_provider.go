package domain

import (
	"github.com/gsouza97/my-bots/internal/adapter/provider"
)

type CryptoFearAndGreedProvider interface {
	GetFearAndGreedIndex() (provider.AlternativeFearAndGreedDataResponse, error)
}
