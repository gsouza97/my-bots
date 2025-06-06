package provider

import (
	"errors"
	"io"
	"net/http"

	"github.com/gsouza97/my-bots/internal/constants"
)

type AlternativeFearAndGreedProvider struct {
}

func NewAlternativeFearAndGreedProvider() *AlternativeFearAndGreedProvider {
	return &AlternativeFearAndGreedProvider{}
}

func (p *AlternativeFearAndGreedProvider) GetFearAndGreedIndex() (AlternativeFearAndGreedDataResponse, error) {
	resp, err := http.Get(constants.AlternativeAPI)
	if err != nil {
		return AlternativeFearAndGreedDataResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return AlternativeFearAndGreedDataResponse{}, errors.New("erro ao buscar Fear and Greed Index")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return AlternativeFearAndGreedDataResponse{}, err
	}

	response, err := ParseAlternativeFearAndGreedResponse(body)
	if err != nil {
		return AlternativeFearAndGreedDataResponse{}, err
	}

	return response.Data[0], nil
}
