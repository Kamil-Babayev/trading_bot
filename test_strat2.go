package main

import (
"fmt"
"trading-bot/internal/strategies"
)

func main() {
	strat := strategies.NewMomentumStrategy()
	
	// Simulate price changes
	prices := []float64{0.28630000, 0.28640000, 0.28635000, 0.28650000, 0.28628000}
	
	for _, price := range prices {
		sig := strat.Check(price)
		sigName := "HOLD"
		if sig == strategies.Buy {
			sigName = "BUY"
		} else if sig == strategies.Sell {
			sigName = "SELL"
		}
		fmt.Printf("Price: %.8f → Signal: %s\n", price, sigName)
	}
}
