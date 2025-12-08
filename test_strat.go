package main

import (
"fmt"
"log"
"trading-bot/internal/marketdata"
"trading-bot/internal/strategies"
)

func main() {
	md := marketdata.NewBinanceClient()
	strat := strategies.NewMomentumStrategy()
	
	for i := 0; i < 10; i++ {
		price, err := md.GetPrice("TRXUSDT")
		if err != nil {
			log.Println("Error:", err)
			continue
		}
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
