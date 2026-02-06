package provider

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gsouza97/my-bots/internal/constants"
)

type BinancePriceProvider struct {
}

func NewBinancePriceProvider() *BinancePriceProvider {
	return &BinancePriceProvider{}
}

func (bp *BinancePriceProvider) GetPrice(crypto string, others ...string) (float64, error) {
	var pair string
	if len(others) > 0 {
		pair = fmt.Sprintf("%s%s", crypto, others[0])
	} else {
		pair = fmt.Sprintf("%sUSDT", crypto)
	}

	url := fmt.Sprintf("%s%s", constants.BinanceAPI, pair)
	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, errors.New("erro ao buscar preço na Binance para o par " + pair)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	response, err := ParseBinancePriceResponse(body)
	if err != nil {
		return 0, err
	}

	price, err := strconv.ParseFloat(response.Price, 64)
	if err != nil {
		return 0, err
	}

	return price, nil
}

func (bp *BinancePriceProvider) GetMultiplePrices(cryptos []string) (map[string]float64, error) {
	prices := make(map[string]float64)

	for _, crypto := range cryptos {
		price, err := bp.GetPrice(crypto)
		if err != nil {
			return nil, err
		}
		prices[crypto] = price
	}

	return prices, nil
}
