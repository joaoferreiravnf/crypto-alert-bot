package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"uphold-alert-bot/internal/models"
)

var PublicURLTicker = "https://api.uphold.com/v0/ticker"

// ApiResponse represents the API response
type ApiResponse struct {
	client *http.Client
}

// NewAPIResponse returns an new instance of ApiResponse
func NewAPIResponse(client *http.Client) *ApiResponse {
	if client == nil {
		client = http.DefaultClient
	}

	return &ApiResponse{
		client: client,
	}
}

// FetchPairData fetches the data for a given pair
func (a *ApiResponse) FetchPairData(ticker *models.Ticker) error {
	pairUrl := fmt.Sprintf(PublicURLTicker+"/%s", ticker.Pair)

	resp, err := a.client.Get(pairUrl)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	err = a.ParseAPIData(resp, ticker)
	if err != nil {
		return errors.Wrap(err, "error parsing API response")
	}

	return nil
}

// ParseAPIData parses the API response
func (a *ApiResponse) ParseAPIData(response *http.Response, ticker *models.Ticker) error {
	respBody, err := io.ReadAll(response.Body)
	if err != nil {
		return errors.Wrap(err, "error reading api response")
	}

	err = json.Unmarshal(respBody, ticker)
	if err != nil {
		return errors.Wrap(err, "error unmarshalling api response")
	}

	return nil
}

// IsPairValid checks if the pair exists and if it's a single pair
func (a *ApiResponse) IsPairValid(pair string) (bool, error) {
	pairUrl := fmt.Sprintf(PublicURLTicker+"/%s", pair)

	resp, err := a.client.Get(pairUrl)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return false, errors.Errorf("pair doesn't exist, try again")
	}

	if resp.StatusCode != http.StatusOK {
		return false, errors.Errorf("can't validate pair, unexpected status code: %d", resp.StatusCode)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, errors.Wrap(err, "error reading api response, try again")
	}

	pairsCount := bytes.Count(respBody, []byte(pair))
	if pairsCount > 1 {
		return false, errors.Errorf("cant use %s because it returns more than one pair. Please specify a ticker for a single pair", pair)
	}

	return true, nil
}
