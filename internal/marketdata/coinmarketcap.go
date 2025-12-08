package marketdata

import (
	"github.com/go-resty/resty/v2"
)

type CoinMarketCapClient struct {
	http   *resty.Client
	apiKey string
}

func NewCoinMarketCapClient(apiKey string) *CoinMarketCapClient {
	return &CoinMarketCapClient{
		http:   resty.New().SetBaseURL("https://pro-api.coinmarketcap.com"),
		apiKey: apiKey,
	}
}

type cmcQuotesResponse struct {
	Data map[string]struct {
		Quote map[string]struct {
			Price float64 `json:"price"`
		} `json:"quote"`
	} `json:"data"`
}

// GetPrice fetches price from CoinMarketCap
// Symbol should be the coin ID (e.g., "tron" for TRX)
func (c *CoinMarketCapClient) GetPrice(symbol string) (float64, error) {
	var resp cmcQuotesResponse

	_, err := c.http.R().
		SetHeader("X-CMC_PRO_API_KEY", c.apiKey).
		SetQueryParam("symbol", symbol).
		SetQueryParam("convert", "USDT").
		SetResult(&resp).
		Get("/v2/cryptocurrency/quotes/latest")
	if err != nil {
		return 0, err
	}

	// Extract price from nested response
	if data, ok := resp.Data[symbol]; ok {
		if quote, ok := data.Quote["USDT"]; ok {
			return quote.Price, nil
		}
	}

	return 0, nil
}

// AlternativeGetPrice - if CMC API returns symbol instead of ID
// Tries to parse direct price response
func (c *CoinMarketCapClient) GetPriceByID(id string) (float64, error) {
	var resp cmcQuotesResponse

	_, err := c.http.R().
		SetHeader("X-CMC_PRO_API_KEY", c.apiKey).
		SetQueryParam("id", id).
		SetQueryParam("convert", "USDT").
		SetResult(&resp).
		Get("/v2/cryptocurrency/quotes/latest")
	if err != nil {
		return 0, err
	}

	// Extract price from nested response
	for _, data := range resp.Data {
		if quote, ok := data.Quote["USDT"]; ok {
			return quote.Price, nil
		}
	}

	return 0, nil
}
