package provider

import "encoding/json"

type BinancePriceResponse struct {
	Price string `json:"price"`
}

func ParseBinancePriceResponse(data []byte) (*BinancePriceResponse, error) {
	var response BinancePriceResponse
	err := json.Unmarshal(data, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}
