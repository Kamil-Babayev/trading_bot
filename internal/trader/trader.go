package trader

import (
	"fmt"
	"os"
	"time"
	"trading-bot/internal/models"
	"trading-bot/internal/strategies"
)

type Trader struct {
	logFile *os.File
}

func NewTrader(logFilePath string) *Trader {
	file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Printf("Warning: Could not open log file: %v\n", err)
	}
	return &Trader{logFile: file}
}

func (t *Trader) Execute(sig strategies.Signal) {
	t.ExecuteWithDetails(sig, 0, "")
}

func (t *Trader) ExecuteWithDetails(sig strategies.Signal, price float64, strategyName string) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")

	var action string
	var emoji string

	switch sig {
	case strategies.Buy:
		action = "BUY"
		emoji = "📈"
	case strategies.Sell:
		action = "SELL"
		emoji = "📉"
	default:
		return
	}

	// Print to terminal
	if price > 0 && strategyName != "" {
		fmt.Printf("%s [%s] %s SIGNAL at price: %.8f (Strategy: %s)\n",
			emoji, timestamp, action, price, strategyName)
	} else {
		fmt.Printf("%s [%s] %s SIGNAL\n", emoji, timestamp, action)
	}

	// Log to file
	if t.logFile != nil {
		logEntry := fmt.Sprintf("%s | %s | Price: %.8f | Strategy: %s\n",
			timestamp, action, price, strategyName)
		t.logFile.WriteString(logEntry)
		t.logFile.Sync() // Ensure data is written to disk
	}
}

func (t *Trader) LogTrade(trade models.Trade) {
	if t.logFile != nil {
		logEntry := fmt.Sprintf("%s | %s | Price: %.8f\n",
			trade.Time.Format("2006-01-02 15:04:05"), trade.Signal, trade.Price)
		t.logFile.WriteString(logEntry)
		t.logFile.Sync()
	}
}

func (t *Trader) Close() error {
	if t.logFile != nil {
		return t.logFile.Close()
	}
	return nil
}
