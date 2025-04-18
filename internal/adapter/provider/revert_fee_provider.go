package provider

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gsouza97/my-bots/internal/constants"
)

type RevertFeeProvider struct {
}

func NewRevertFeeProvider() *RevertFeeProvider {
	return &RevertFeeProvider{}
}

func (rp *RevertFeeProvider) GetFees(ctx context.Context, chain string, nftId string) (RevertPoolDataResponse, error) {
	url := fmt.Sprintf("%s%s/uniswapv3/%s", constants.RevertAPI, chain, nftId)
	resp, err := http.Get(url)
	if err != nil {
		return RevertPoolDataResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return RevertPoolDataResponse{}, errors.New("erro ao buscar preço no Revert")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return RevertPoolDataResponse{}, err
	}

	response, err := ParseRevertPoolResponse(body)
	if err != nil {
		return RevertPoolDataResponse{}, err
	}

	if !response.Success {
		return RevertPoolDataResponse{}, fmt.Errorf("error: invalid response")
	}

	log.Println("Revert response:", response)

	return response.Data, nil
}
