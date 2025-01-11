package api

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"uphold-alert-bot/internal/models"
)

func TestFetchPairData(t *testing.T) {
	tests := []struct {
		name         string
		statusCode   int
		responseBody string
		wantErr      bool
		errContains  string
	}{
		{
			name:         "Success - 200 OK with valid JSON",
			statusCode:   http.StatusOK,
			responseBody: `{"ask":"123.45","bid":"120.00"}`,
			wantErr:      false,
		},
		{
			name:         "Error - non-OK status code (404)",
			statusCode:   http.StatusNotFound,
			responseBody: `{"error":"not found"}`,
			wantErr:      true,
			errContains:  "unexpected status code: 404",
		},
		{
			name:         "Error - invalid JSON",
			statusCode:   http.StatusOK,
			responseBody: `{"ask":"badjson"`,
			wantErr:      true,
			errContains:  "error unmarshalling api response",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
				_, _ = w.Write([]byte(tt.responseBody))
			}))
			defer server.Close()

			oldURL := PublicURLTicker
			PublicURLTicker = server.URL
			defer func() { PublicURLTicker = oldURL }()

			a := NewAPIResponse(nil)

			ticker := &models.Ticker{Pair: "BTC-USD"}

			err := a.FetchPairData(ticker)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				assert.NoError(t, err)
				if err == nil {
					assert.Equal(t, 123.45, ticker.CurrentAsk.Float64())
					assert.Equal(t, 120.00, ticker.CurrentBid.Float64())
				}
			}
		})
	}
}

func TestParseAPIData(t *testing.T) {
	testCases := []struct {
		name         string
		responseBody string
		wantErr      bool
		errContains  string
	}{
		{
			name:         "Valid JSON",
			responseBody: `{"ask":"150.75","bid":"149.00"}`,
			wantErr:      false,
		},
		{
			name:         "Invalid JSON",
			responseBody: `{"ask":"150.75"`,
			wantErr:      true,
			errContains:  "error unmarshalling api response",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			resp := &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewBufferString(tc.responseBody)),
			}
			defer resp.Body.Close()

			a := NewAPIResponse(nil)
			ticker := &models.Ticker{}

			err := a.ParseAPIData(resp, ticker)

			if tc.wantErr {
				assert.Error(t, err)
				if tc.errContains != "" {
					assert.Contains(t, err.Error(), tc.errContains)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, 150.75, ticker.CurrentAsk.Float64())
				assert.Equal(t, 149.00, ticker.CurrentBid.Float64())
			}
		})
	}
}

func TestIsPairValid(t *testing.T) {
	tests := []struct {
		name         string
		statusCode   int
		responseBody string
		wantValid    bool
		wantErr      bool
		errContains  string
	}{
		{
			name:        "Pair not found (404)",
			statusCode:  http.StatusNotFound,
			wantValid:   false,
			wantErr:     true,
			errContains: "pair doesn't exist",
		},
		{
			name:        "Unexpected status code (500)",
			statusCode:  http.StatusInternalServerError,
			wantValid:   false,
			wantErr:     true,
			errContains: "unexpected status code: 500",
		},
		{
			name:         "Valid (200)",
			statusCode:   http.StatusOK,
			responseBody: `{"pairs": ["FAKE-PAIR"]}`,
			wantValid:    true,
			wantErr:      false,
		},
		{
			name:         "Multiple pairs returned",
			statusCode:   http.StatusOK,
			responseBody: `{"pairs": ["BTCUSDT", "BTCUEUR"]}`,
			wantValid:    false,
			wantErr:      true,
			errContains:  "cant use BTC because it returns more than one pair. Please specify a ticker for a single pair",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
				if tt.responseBody != "" {
					w.Write([]byte(tt.responseBody))
				}
			}))
			defer server.Close()

			oldURL := PublicURLTicker
			PublicURLTicker = server.URL
			defer func() { PublicURLTicker = oldURL }()

			a := NewAPIResponse(nil)
			valid, err := a.IsPairValid("BTC")

			assert.Equal(t, tt.wantValid, valid)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
