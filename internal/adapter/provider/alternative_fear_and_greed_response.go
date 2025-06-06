package provider

import "encoding/json"

type AlternativeFearAndGreedResponse struct {
	Data []AlternativeFearAndGreedDataResponse `json:"data"`
}

type AlternativeFearAndGreedDataResponse struct {
	Value               string `json:"value"`
	ValueClassification string `json:"value_classification"`
}

func ParseAlternativeFearAndGreedResponse(data []byte) (*AlternativeFearAndGreedResponse, error) {
	var response AlternativeFearAndGreedResponse
	err := json.Unmarshal(data, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}
