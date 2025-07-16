package provider

import (
	"errors"
	"io"
	"net/http"

	"github.com/gsouza97/my-bots/internal/constants"
)

type HomologacionProvider struct {
}

func NewHomologacionProvider() *HomologacionProvider {
	return &HomologacionProvider{}
}

func (p *HomologacionProvider) GetHomologacionStatus(params string) (HomologacionResponse, error) {
	resp, err := http.Get(buildHomologacionAPIUrl(params))
	if err != nil {
		return HomologacionResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return HomologacionResponse{}, errors.New("erro ao buscar resposta da homologação")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return HomologacionResponse{}, err
	}

	response, err := ParseHomologacionResponse(body)
	if err != nil {
		return HomologacionResponse{}, err
	}

	return *response, nil
}

func buildHomologacionAPIUrl(params string) string {
	return constants.HomologacionAPI + params
}
