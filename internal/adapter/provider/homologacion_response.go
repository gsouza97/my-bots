package provider

import "encoding/json"

type HomologacionResponse []struct {
	NumeroExpediente string `json:"numeroExpediente"`
	Estado           string `json:"estado"`
	CodAcceda        string `json:"codAcceda"`
	Fecha            string `json:"fecha"`
}

func ParseHomologacionResponse(data []byte) (*HomologacionResponse, error) {
	var response HomologacionResponse
	err := json.Unmarshal(data, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}
