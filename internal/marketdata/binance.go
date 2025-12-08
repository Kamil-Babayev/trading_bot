package marketdata

import (
	"strconv"

	"github.com/go-resty/resty/v2"
)

type BinanceClient struct {
	http *resty.Client
}

func NewBinanceClient() *BinanceClient {
	return &BinanceClient{
		http: resty.New().SetBaseURL("https://api.binance.com"),
	}
}

type binancePriceResponse struct {
	Price string `json:"price"`
}

func (c *BinanceClient) GetPrice(symbol string) (float64, error) {
	var resp binancePriceResponse

	_, err := c.http.R().
		SetQueryParam("symbol", symbol).
		SetResult(&resp).
		Get("/api/v3/ticker/price")
	if err != nil {
		return 0, err
	}

	return strconv.ParseFloat(resp.Price, 64)
}
