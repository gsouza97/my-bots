package provider

import "encoding/json"

type CmcAltcoinSeasonResponse struct {
	Data CmcAltcoinSeasonDataResponse `json:"data"`
}

type CmcAltcoinSeasonDataResponse struct {
	HistoricalValues CmcAltcoinSeasonHistoricalValuesResponse `json:"historicalValues"`
}

type CmcAltcoinSeasonHistoricalValuesResponse struct {
	Now CmcAltcoinSeasonIndexResponse `json:"now"`
}

type CmcAltcoinSeasonIndexResponse struct {
	Name         string `json:"name"`
	AltcoinIndex string `json:"altcoinIndex"`
}

func ParseCmcAltcoinSeasonResponse(data []byte) (*CmcAltcoinSeasonIndexResponse, error) {
	var response CmcAltcoinSeasonResponse
	err := json.Unmarshal(data, &response)
	if err != nil {
		return nil, err
	}
	return &response.Data.HistoricalValues.Now, nil
}
