package provider

import "encoding/json"

type RevertPoolResponse struct {
	Success bool                   `json:"success"`
	Data    RevertPoolDataResponse `json:"data"`
}

type RevertPoolDataResponse struct {
	UncollectedFees0 string                             `json:"uncollected_fees0"`
	UncollectedFees1 string                             `json:"uncollected_fees1"`
	Token0           string                             `json:"token0"`
	Token1           string                             `json:"token1"`
	Tokens           map[string]RevertPoolTokenResponse `json:"tokens"`
}

type RevertPoolTokenResponse struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
}

func ParseRevertPoolResponse(data []byte) (*RevertPoolResponse, error) {
	var response RevertPoolResponse
	err := json.Unmarshal(data, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}
