package strategies

type MomentumStrategy struct {
	last float64
}

func NewMomentumStrategy() *MomentumStrategy {
	return &MomentumStrategy{}
}

// Простая стратегия: если изменение >0.01% → Buy, если <−0.01% → Sell
func (s *MomentumStrategy) Check(price float64) Signal {
	if s.last == 0 {
		s.last = price
		return Hold
	}

	change := (price - s.last) / s.last * 100
	s.last = price

	if change > 0.01 {
		return Buy
	}
	if change < -0.01 {
		return Sell
	}
	return Hold
}
