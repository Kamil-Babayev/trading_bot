package storage

import "trading-bot/internal/models"

type MemoryStorage struct {
	Trades []models.Trade
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{Trades: make([]models.Trade, 0)}
}

func (s *MemoryStorage) SaveTrade(t models.Trade) {
	s.Trades = append(s.Trades, t)
}
