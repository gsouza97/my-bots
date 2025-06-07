package provider

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/gsouza97/my-bots/internal/constants"
)

type CmcAltcoinSeasonProvider struct {
}

func NewCmcAltcoinSeasonProvider() *CmcAltcoinSeasonProvider {
	return &CmcAltcoinSeasonProvider{}
}

func (p *CmcAltcoinSeasonProvider) GetAltcoinSeasonIndex() (CmcAltcoinSeasonIndexResponse, error) {
	now := time.Now().UTC()
	startDate := now.AddDate(0, 0, -2)
	start := strconv.FormatInt(startDate.Unix(), 10)
	end := strconv.FormatInt(now.Unix(), 10)

	url, err := url.Parse(constants.CmcAltcoinSeasonAPI)
	if err != nil {
		return CmcAltcoinSeasonIndexResponse{}, errors.New("erro ao parsear URL")
	}

	query := url.Query()
	query.Set("start", start)
	query.Set("end", end)
	url.RawQuery = query.Encode()

	resp, err := http.Get(url.String())
	if err != nil {
		return CmcAltcoinSeasonIndexResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return CmcAltcoinSeasonIndexResponse{}, errors.New("erro ao buscar Altcoin Season Index")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return CmcAltcoinSeasonIndexResponse{}, err
	}

	response, err := ParseCmcAltcoinSeasonResponse(body)
	if err != nil {
		return CmcAltcoinSeasonIndexResponse{}, err
	}

	return *response, nil
}
