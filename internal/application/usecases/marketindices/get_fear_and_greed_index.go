package marketindices

import (
	"fmt"

	"github.com/gsouza97/my-bots/internal/domain"
	"github.com/gsouza97/my-bots/internal/infrastructure/providers"
)

type GetFearAndGreedIndex struct {
	fearAndGreedProvider domain.CryptoFearAndGreedProvider
}

func NewGetFearAndGreedIndex(fearAndGreedProvider domain.CryptoFearAndGreedProvider) *GetFearAndGreedIndex {
	return &GetFearAndGreedIndex{
		fearAndGreedProvider: fearAndGreedProvider,
	}
}

func (uc *GetFearAndGreedIndex) Execute() (string, error) {
	data, err := uc.fearAndGreedProvider.GetFearAndGreedIndex()
	if err != nil {
		return "", err
	}

	msg := buildFearAndGreedMessage(data)

	return msg, nil
}

func buildFearAndGreedMessage(data providers.AlternativeFearAndGreedDataResponse) string {
	return fmt.Sprintf(
		"Fear and Greed Index: %s → %s",
		data.Value, data.ValueClassification,
	)
}
