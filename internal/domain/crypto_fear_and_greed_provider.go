package domain

import "github.com/gsouza97/my-bots/internal/infrastructure/providers"

type CryptoFearAndGreedProvider interface {
	GetFearAndGreedIndex() (providers.AlternativeFearAndGreedDataResponse, error)
}
