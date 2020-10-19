package finnhub

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"

	"github.com/dbeliakov/stocks-bot/stocks"
)

const (
	authHeader = "X-Finnhub-Token"
)

type Provider struct {
	apiURL     string
	httpClient *resty.Client
}

var _ stocks.Provider = &Provider{}

func NewProvider(apiURL, apiKey string) *Provider {
	return &Provider{
		httpClient: resty.New().SetHeader(authHeader, apiKey),
		apiURL:     apiURL,
	}
}

func (p *Provider) CurrentPrice(symbol string) (float64, error) {
	resp, err := p.httpClient.R().
		SetQueryParam("symbol", symbol).Get(p.apiURL)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch info: %w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return 0, fmt.Errorf("unexpected status code: %d", resp.StatusCode())
	}

	var c currentPriceResponse
	if err := json.Unmarshal(resp.Body(), &c); err != nil {
		return 0, fmt.Errorf("failed to unmarshal data: %w", err)
	}

	if c.CurrentPrice == 0 {
		return 0, errors.New("symbol not found")
	}

	return c.CurrentPrice, nil
}

type currentPriceResponse struct {
	CurrentPrice float64 `json:"c"`
}
