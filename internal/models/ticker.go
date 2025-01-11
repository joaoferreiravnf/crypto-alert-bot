package models

import (
	"math"
	"time"
)

var rateLimit = 250

// Ticker represents a trading pair entity
type Ticker struct {
	Pair           string
	Currency       string  `json:"currency"`
	CurrentAsk     Float64 `json:"ask"`
	CurrentBid     Float64 `json:"bid"`
	PreviousAsk    Float64
	PreviousBid    Float64
	AskPriceChange float64
	AskPercChange  float64
	Config         TickerConfig
}

// TickerConfig represents the configuration settings for a ticker
type TickerConfig struct {
	RefreshRate     float64
	PercOscillation float64
	Lifetime        time.Duration
}

// NewTicker creates a new ticker entity
func NewTicker(pair string, refreshRate float64, percOscillation float64, lifetime time.Duration) *Ticker {
	return &Ticker{
		Pair: pair,
		Config: TickerConfig{
			RefreshRate:     refreshRate,
			PercOscillation: percOscillation,
			Lifetime:        lifetime,
		},
	}
}

// IsAbovePercOscillation checks if the current ask price is above the percentage oscillation threshold
func (t *Ticker) IsAbovePercOscillation() bool {
	t.setAskPriceChange()
	t.setAskPercChange()

	if t.PreviousAsk != 0 && t.AskPercChange >= t.Config.PercOscillation {
		return true
	}

	return false
}

// NormalizeValues resets the previous ask and bid prices to the current ask and bid prices for futures calculations
func (t *Ticker) setAskPriceChange() {
	previousAsk := t.PreviousAsk.Float64()
	currentAsk := t.CurrentAsk.Float64()

	t.AskPriceChange = math.Abs(previousAsk - currentAsk)
}

// setAskPercChange calculates the percentage change between the previous ask price and the current ask price
func (t *Ticker) setAskPercChange() {
	t.AskPercChange = t.AskPriceChange / t.PreviousAsk.Float64() * 100
}

// NormalizeValues resets the previous ask and bid prices to the current ask and bid prices for futures calculations
func (t *Ticker) NormalizeValues() {
	t.PreviousAsk = t.CurrentAsk
	t.PreviousBid = t.CurrentBid
}

type Tickers []*Ticker

func (ts *Tickers) IsAboveRateLimit() bool {
	totalCalls := 0

	for _, ticker := range *ts {
		call := int(math.Ceil(60 / ticker.Config.RefreshRate))

		totalCalls = totalCalls + call
	}

	if totalCalls > rateLimit {
		return true
	}

	return false
}
