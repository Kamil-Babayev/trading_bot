package marketdata

type MarketData interface {
	GetPrice(symbol string) (float64, error)
}
