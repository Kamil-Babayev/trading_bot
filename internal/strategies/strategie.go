package strategies

type Signal int

const (
	Hold Signal = iota
	Buy
	Sell
)

type Strategy interface {
	Check(price float64) Signal
}
