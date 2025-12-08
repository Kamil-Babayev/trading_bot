package strategies

// RSILikeStrategy implements a simple RSI-like strategy
// Tracks gains and losses over a period to detect overbought/oversold conditions
type RSILikeStrategy struct {
	prices []float64
	period int
}

func NewRSILikeStrategy(period int) *RSILikeStrategy {
	if period < 2 {
		period = 14 // RSI default
	}
	return &RSILikeStrategy{
		prices: make([]float64, 0, period+1),
		period: period,
	}
}

func (s *RSILikeStrategy) Check(price float64) Signal {
	s.prices = append(s.prices, price)

	// Keep only the last 'period+1' prices (need one extra for comparison)
	if len(s.prices) > s.period+1 {
		s.prices = s.prices[1:]
	}

	// Need minimum data points
	if len(s.prices) < s.period+1 {
		return Hold
	}

	// Calculate gains and losses
	gains := 0.0
	losses := 0.0

	for i := 1; i < len(s.prices); i++ {
		change := s.prices[i] - s.prices[i-1]
		if change > 0 {
			gains += change
		} else {
			losses += -change
		}
	}

	avgGain := gains / float64(s.period)
	avgLoss := losses / float64(s.period)

	// Avoid division by zero
	if avgLoss == 0 {
		if avgGain > 0 {
			return Buy
		}
		return Hold
	}

	rs := avgGain / avgLoss
	rsi := 100 - (100 / (1 + rs))

	// Simple thresholds: oversold < 30, overbought > 70
	if rsi < 30 {
		return Buy
	}
	if rsi > 70 {
		return Sell
	}
	return Hold
}
