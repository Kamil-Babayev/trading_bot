# Copilot Instructions for Trading Bot

## Architecture Overview

The trading bot is a **concurrent streaming system** with three independent goroutine pipelines connected by channels:

```
MarketData → Price Channel → Strategies → Signal Channel → Trader → Execution + Logging
  (5m tick)      (float64)      (Multiple)    (map[str]Signal)  (Execute)
```

**Key Design Philosophy**: Loose coupling via channels; each component is independent and implements a well-defined interface. Multiple strategies can run in parallel and vote on trading signals.

## Core Components

### 1. **Market Data** (`internal/marketdata/`)
- **Interface**: `MarketData` with single method `GetPrice(symbol string) (float64, error)`
- **Implementations**: 
  - `BinanceClient`: Fetches live prices from Binance API (`/api/v3/ticker/price`)
  - `CoinMarketCapClient`: Fetches prices from CoinMarketCap Pro API (requires API key)
- **Pattern**: Factory constructors return pointer receivers; HTTP client baseURL set in constructor
- **Key Detail**: Symbol passed per-call, not stored in struct

### 2. **Strategies** (`internal/strategies/`)
- **Core Interface**: `Strategy` with `Check(price float64) Signal` method
- **Signal Type**: Enum-like constants (Hold=0, Buy=1, Sell=2)
- **Implementations**:
  - **MomentumStrategy**: Detects ±0.5% price changes; stateful
  - **SMAStrategy**: Simple Moving Average (default period=5); buys when price > SMA
  - **RSILikeStrategy**: RSI-based overbought/oversold (default period=14); buy < 30, sell > 70
- **Pattern**: Strategies maintain state; `Check()` updates internal state before returning signal
- **Convention**: First price call returns `Hold` (initialization phase)

### 3. **Trader** (`internal/trader/`)
- **Responsibility**: Executes signals and logs trades
- **Features**:
  - Prints timestamped signals to terminal with emoji + price + strategy name
  - Logs all trades to configurable file (`trades.log` by default)
  - `ExecuteWithDetails(signal, price, strategyName)` for rich output
  - Graceful close with `Close()` method
- **Pattern**: No persistence layer directly coupled; side-effect only

### 4. **Models & Storage** (`internal/models/`, `internal/storage/`)
- **Trade Model**: Timestamp, signal string, price snapshot
- **MemoryStorage**: Unused in current flow; available for future use

## Data Flow & Concurrency Patterns

**Entry Point**: `cmd/bot/main.go` with configurable parameters:

```go
// Three independent goroutine workers:
1. Ticker goroutine: Fetches price every N seconds/minutes → priceCh
2. Strategy goroutine: Evaluates all strategies, returns map[strategyName]Signal → signalCh
3. Trader goroutine: Loops through signals, executes each via Trader.ExecuteWithDetails()

// Graceful shutdown:
signal.Notify(sigChan) → close channels → wait → exit
```

**Important**: Channels use buffering (capacity=1) to prevent goroutine blocking. Each goroutine independent; one failure doesn't propagate.

## Configuration & Flags

Run bot with:
```bash
go run ./cmd/bot -interval 5m -api binance -strategy all -symbol TRXUSDT -log trades.log
```

Available flags:
- `-interval`: Duration between price fetches (default: 5m; e.g., `10s`, `5m`, `1h`)
- `-api`: `binance` or `coinmarketcap` (default: binance)
- `-cmc-key`: CoinMarketCap API key (required if `-api coinmarketcap`)
- `-strategy`: `momentum`, `sma`, `rsi`, or `all` (default: all)
- `-symbol`: Trading pair (default: TRXUSDT)
- `-log`: Log file path (default: trades.log)

## Critical Developer Conventions

### Naming & Packages
- **Package names**: lowercase, match directory (no underscores)
- **Types**: PascalCase (e.g., `BinanceClient`, `MomentumStrategy`)
- **Methods**: Pointer receivers standard (e.g., `func (c *BinanceClient)`)
- **Constructors**: Use `NewTypeName()` pattern

### Adding New Strategies
1. Create `internal/strategies/newstrategy.go`
2. Implement `Strategy` interface: `Check(price float64) Signal`
3. Maintain state as struct fields if needed
4. Register in `main.go` strategies slice (lines ~35-50)
5. No other changes needed (interface-based design)

**Example**:
```go
type CustomStrategy struct {
	state float64
}

func NewCustomStrategy() *CustomStrategy {
	return &CustomStrategy{}
}

func (s *CustomStrategy) Check(price float64) Signal {
	// implementation
	return Buy/Sell/Hold
}
```

### Adding New Market Data Providers
1. Create `internal/marketdata/newprovider.go`
2. Implement `MarketData` interface: `GetPrice(symbol string) (float64, error)`
3. Use `resty.Client` for HTTP consistency
4. Parse response into float64; handle errors gracefully
5. Swap constructor in main.go: `md := marketdata.NewXxxClient()`

### Error Handling
- **Strategy `Check()`**: Never errors; returns `Hold` on invalid state
- **MarketData `GetPrice()`**: Returns error; main.go logs and continues (line ~68)
- **No panics**: All error cases logged; goroutines survive individual failures
- **Graceful shutdown**: Ctrl+C closes channels, waits for goroutines, exits

## Build & Run

### Build
```bash
go build -o bot ./cmd/bot
```

### Run with defaults (5m interval, Binance, all strategies)
```bash
go run ./cmd/bot
```

### Run with custom settings (10s interval, all strategies)
```bash
go run ./cmd/bot -interval 10s
```

### Run with CoinMarketCap (requires API key from coinmarketcap.com/api/)
```bash
go run ./cmd/bot -api coinmarketcap -cmc-key YOUR_API_KEY
```

### Output
- **Terminal**: Timestamped signals with price and strategy name
- **File**: `trades.log` contains all trade history (append-only)

## Example Modifications

### Change tick interval to 10 minutes
```bash
go run ./cmd/bot -interval 10m
```

### Use only momentum strategy
```bash
go run ./cmd/bot -strategy momentum
```

### Change SMA period (modify main.go line ~46)
```go
strategies.NewSMAStrategy(10)  // 10-period instead of 5
```

### Change momentum threshold
File: `internal/strategies/momentum.go` lines 18-19
```go
if change > 1.0 {    // Buy threshold (1% instead of 0.5%)
  return Buy
}
if change < -1.0 {   // Sell threshold
  return Sell
}
```

## Dependencies
- `github.com/go-resty/resty/v2 v2.17.0`: HTTP client for API calls
- Standard library: `time`, `log`, `os`, `signal`, `flag`, `fmt`, `syscall`

## Testing & Debugging

### Test a strategy locally
```go
strat := strategies.NewSMAStrategy(5)
for i := 1.0; i <= 5.0; i += 0.1 {
	sig := strat.Check(i)
	// assert sig == Hold for first 5 checks, then Buy/Sell after
}
```

### Verify API connectivity
```bash
# Binance (no auth)
curl "https://api.binance.com/api/v3/ticker/price?symbol=TRXUSDT"

# CoinMarketCap (needs key)
curl -H "X-CMC_PRO_API_KEY: YOUR_KEY" \
  "https://pro-api.coinmarketcap.com/v2/cryptocurrency/quotes/latest?symbol=TRX&convert=USDT"
```

### Monitor trades in real-time
```bash
tail -f trades.log
```
