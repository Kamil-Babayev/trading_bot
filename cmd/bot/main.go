package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	"trading-bot/internal/marketdata"
	"trading-bot/internal/strategies"
	"trading-bot/internal/trader"
)

func main() {
	// Configuration flags
	tickInterval := flag.Duration("interval", 5*time.Minute, "Price fetch interval (e.g., 5m, 10s)")
	apiProvider := flag.String("api", "binance", "Market data provider: binance or coinmarketcap")
	cmcAPIKey := flag.String("cmc-key", "", "CoinMarketCap API key (required if api=coinmarketcap)")
	strategyName := flag.String("strategy", "all", "Strategy to use: momentum, sma, rsi, or all")
	symbol := flag.String("symbol", "TRXUSDT", "Trading pair symbol (e.g., TRXUSDT for Binance)")
	logFile := flag.String("log", "trades.log", "File path for trade history")
	flag.Parse()

	log.Printf("🤖 Trading Bot Started")
	log.Printf("  API: %s, Symbol: %s, Interval: %v, Strategy: %s",
		*apiProvider, *symbol, *tickInterval, *strategyName)

	// Initialize market data
	var md marketdata.MarketData
	if *apiProvider == "coinmarketcap" {
		if *cmcAPIKey == "" {
			log.Fatal("CoinMarketCap API key required: use -cmc-key flag")
		}
		md = marketdata.NewCoinMarketCapClient(*cmcAPIKey)
	} else {
		md = marketdata.NewBinanceClient()
	}

	// Initialize strategies
	type strategyEntry struct {
		name     string
		instance strategies.Strategy
	}
	var strats []strategyEntry

	if *strategyName == "momentum" || *strategyName == "all" {
		strats = append(strats, strategyEntry{"Momentum", strategies.NewMomentumStrategy()})
	}
	if *strategyName == "sma" || *strategyName == "all" {
		strats = append(strats, strategyEntry{"SMA(5)", strategies.NewSMAStrategy(5)})
	}
	if *strategyName == "rsi" || *strategyName == "all" {
		strats = append(strats, strategyEntry{"RSI(14)", strategies.NewRSILikeStrategy(14)})
	}

	// Initialize trader
	tr := trader.NewTrader(*logFile)
	defer tr.Close()

	// Channels
	type priceWithSignals struct {
		price   float64
		signals map[string]strategies.Signal
	}
	priceCh := make(chan float64, 1)
	signalCh := make(chan priceWithSignals, 1)

	// Goroutine 1: Fetch price at interval
	go func() {
		ticker := time.NewTicker(*tickInterval)
		defer ticker.Stop()
		for range ticker.C {
			price, err := md.GetPrice(*symbol)
			if err != nil {
				log.Printf("❌ Market data error: %v", err)
				continue
			}
			select {
			case priceCh <- price:
			default:
			}
		}
	}()

	// Goroutine 2: Evaluate all strategies
	go func() {
		for price := range priceCh {
			signals := make(map[string]strategies.Signal)
			for _, strat := range strats {
				sig := strat.instance.Check(price)
				signals[strat.name] = sig
			}
			select {
			case signalCh <- priceWithSignals{price, signals}:
			default:
			}
		}
	}()

	// Goroutine 3: Execute trades
	go func() {
		for priceSignals := range signalCh {
			for stratName, sig := range priceSignals.signals {
				if sig != strategies.Hold {
					tr.ExecuteWithDetails(sig, priceSignals.price, stratName)
				}
			}
		}
	}()

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	fmt.Println("\n✅ Trading bot running (press Ctrl+C to exit)...\n")

	<-sigChan
	log.Println("\n🛑 Shutting down...")
	close(priceCh)
	close(signalCh)
	time.Sleep(100 * time.Millisecond)
	log.Println("✅ Bot stopped gracefully")
}
