# Trading Bot - TRX/USDT University Project

A concurrent streaming cryptocurrency trading bot that monitors TRX/USDT prices and generates buy/sell signals based on multiple strategies.

## Features

✨ **Three Trading Strategies:**
- **Momentum**: Detects ±0.5% price changes
- **Simple Moving Average (SMA)**: Buys when price > 5-period average
- **RSI-like**: Detects overbought (>70) and oversold (<30) conditions

✨ **Market Data Sources:**
- Binance API (free, no auth required)
- CoinMarketCap Pro API (requires API key)

✨ **Configuration:**
- Configurable price fetch interval (5 seconds to hours)
- Run all strategies or select specific ones
- Trade history logged to file with timestamps

## Quick Start

### Build
```bash
go build -o bot ./cmd/bot
```

### Run with Defaults (5-minute interval, all strategies, Binance)
```bash
./bot
```

### Run Examples

**10-second interval for testing:**
```bash
./bot -interval 10s
```

**Only momentum strategy, 1-minute interval:**
```bash
./bot -strategy momentum -interval 1m
```

**Only SMA strategy:**
```bash
./bot -strategy sma
```

**Using CoinMarketCap (requires API key):**
```bash
./bot -api coinmarketcap -cmc-key YOUR_API_KEY
```

**Custom symbol and log file:**
```bash
./bot -symbol ETHUSDT -log my_trades.log
```

## Output

Terminal output shows real-time signals:
```
✅ Trading bot running (press Ctrl+C to exit)...

📈 [2025-12-08 14:23:45] BUY SIGNAL at price: 0.25948273 (Strategy: Momentum)
📉 [2025-12-08 14:28:52] SELL SIGNAL at price: 0.25867491 (Strategy: SMA(5))
```

Trade history is saved to `trades.log` (or your custom log file):
```
2025-12-08 14:23:45 | BUY | Price: 0.25948273 | Strategy: Momentum
2025-12-08 14:28:52 | SELL | Price: 0.25867491 | Strategy: SMA(5)
```

## All Flags

```
-interval duration
    Price fetch interval (default: 5m, examples: 10s, 1m, 5m, 1h)

-api string
    Market data provider: "binance" or "coinmarketcap" (default: binance)

-cmc-key string
    CoinMarketCap API key (required if -api coinmarketcap)

-strategy string
    Strategy to use: "momentum", "sma", "rsi", or "all" (default: all)

-symbol string
    Trading pair symbol (default: TRXUSDT)

-log string
    File path for trade history (default: trades.log)
```

## Architecture

Three concurrent goroutines communicate via channels:

```
MarketData API
     ↓
 Price Channel (float64)
     ↓
 Strategy Goroutine (evaluates all active strategies)
     ↓
 Signal Channel (map[strategyName]Signal)
     ↓
 Trader Goroutine (prints + logs)
     ↓
 Output + trades.log
```

## Tech Stack

- **Language**: Go 1.23+
- **HTTP Client**: `github.com/go-resty/resty/v2`
- **Concurrency**: Goroutines + channels
- **Design Pattern**: Interface-based architecture for easy extension

## For Your University Project

This bot demonstrates:
- ✅ Concurrent programming with goroutines and channels
- ✅ Real API integration (Binance & CoinMarketCap)
- ✅ Strategy pattern for trading algorithms
- ✅ Configuration management with flags
- ✅ File I/O for trade history
- ✅ Error handling and graceful shutdown
- ✅ Multiple independent data sources

Run with different strategies and intervals to observe different trading behaviors!
