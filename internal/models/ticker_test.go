package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewTicker(t *testing.T) {
	pair := "BTC-USD"
	refreshRate := 60.0
	percOscillation := 5.0
	lifetime := 24 * time.Hour

	ticker := NewTicker(pair, refreshRate, percOscillation, lifetime)

	assert.Equal(t, pair, ticker.Pair, "Pair should match the input.")
	assert.Equal(t, refreshRate, ticker.Config.RefreshRate, "RefreshRate should match the input.")
	assert.Equal(t, percOscillation, ticker.Config.PercOscillation, "PercOscillation should match the input.")
	assert.Equal(t, lifetime, ticker.Config.Lifetime, "Lifetime should match the input.")
}

func TestIsAbovePercOscillation(t *testing.T) {
	tests := []struct {
		name            string
		previousAsk     Float64
		currentAsk      Float64
		percOscillation float64
		wantIsAbove     bool
	}{
		{
			name:            "Previous ask = 0",
			previousAsk:     0,
			currentAsk:      100.0,
			percOscillation: 5.0,
			wantIsAbove:     false,
		},
		{
			name:            "Ask change within threshold",
			previousAsk:     100.0,
			currentAsk:      103.0,
			percOscillation: 5.0,
			wantIsAbove:     false,
		},
		{
			name:            "Ask change above threshold",
			previousAsk:     100.0,
			currentAsk:      110.0,
			percOscillation: 5.0,
			wantIsAbove:     true,
		},
		{
			name:            "Ask change exactly threshold",
			previousAsk:     100.0,
			currentAsk:      105.0,
			percOscillation: 5.0,
			wantIsAbove:     true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			ticker := &Ticker{
				PreviousAsk: tt.previousAsk,
				CurrentAsk:  tt.currentAsk,
				Config: TickerConfig{
					PercOscillation: tt.percOscillation,
				},
			}

			gotIsAbove := ticker.IsAbovePercOscillation()
			assert.Equal(t, tt.wantIsAbove, gotIsAbove)
		})
	}
}

func TestNormalizeValues(t *testing.T) {
	ticker := &Ticker{
		CurrentAsk:  120.0,
		CurrentBid:  118.0,
		PreviousAsk: 100.0,
		PreviousBid: 99.0,
	}

	ticker.NormalizeValues()

	assert.Equal(t, float64(ticker.CurrentAsk), float64(ticker.PreviousAsk),
		"PreviousAsk should be updated to CurrentAsk.")
	assert.Equal(t, float64(ticker.CurrentBid), float64(ticker.PreviousBid),
		"PreviousBid should be updated to CurrentBid.")
}
