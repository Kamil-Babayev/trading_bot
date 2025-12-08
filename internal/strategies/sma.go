package strategies

// SMAStrategy implements a Simple Moving Average strategy
// Buys when current price > SMA, sells when current price < SMA
type SMAStrategy struct {
	prices []float64
	period int
}

func NewSMAStrategy(period int) *SMAStrategy {
	if period < 2 {
		period = 5 // default
	}
	return &SMAStrategy{
		prices: make([]float64, 0, period),
		period: period,
	}
}

func (s *SMAStrategy) Check(price float64) Signal {
	s.prices = append(s.prices, price)
	
	// Keep only the last 'period' prices
	if len(s.prices) > s.period {
		s.prices = s.prices[1:]
	}
	
	// Need minimum data points
	if len(s.prices) < s.period {
		return Hold
	}
	
	// Calculate SMA
	sum := 0.0
	for _, p := range s.prices {
		sum += p
	}
	sma := sum / float64(s.period)
	
	// Compare current price with SMA
	if price > sma*1.001 { // 0.1% above SMA
		return Buy
	}
	if price < sma*0.999 { // 0.1% below SMA
		return Sell
	}
	return Hold
}
