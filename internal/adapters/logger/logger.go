package logger

import (
	"log/slog"
	"time"
	"uphold-alert-bot/internal/models"
)

// TickerPublisher is a struct that implements the Publisher
type TickerPublisher struct{}

// NewTickerPublisher returns a new instance of TickerPublisher
func NewTickerPublisher() *TickerPublisher {
	return &TickerPublisher{}
}

// Publish publishes the ticker
func (tp *TickerPublisher) Publish(timestamp time.Time, ticker *models.Ticker) {
	slog.Info(
		"Above threshold alert:", "pair", ticker.Pair,
		"percent_change:", ticker.AskPercChange,
		"price_change:", ticker.AskPriceChange,
		"time:", timestamp)
}
