package usecase

import (
	"fmt"

	"github.com/gsouza97/my-bots/internal/adapter/provider"
	"github.com/gsouza97/my-bots/internal/domain"
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

func buildFearAndGreedMessage(data provider.AlternativeFearAndGreedDataResponse) string {
	return fmt.Sprintf(
		"Fear and Greed Index: %s → %s",
		data.Value, data.ValueClassification,
	)
}
