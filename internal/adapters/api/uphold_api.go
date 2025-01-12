package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"uphold-alert-bot/internal/models"
)

var PublicURLTicker = "https://api.uphold.com/v0/ticker"

// UpholdApi represents the API response
type UpholdApi struct {
	client *http.Client
}

// NewAPIResponse returns an new instance of UpholdApi
func NewAPIResponse(client *http.Client) *UpholdApi {
	if client == nil {
		client = http.DefaultClient
	}

	return &UpholdApi{
		client: client,
	}
}

// FetchPairData fetches the data for a given pair
func (a *UpholdApi) FetchPairData(ctx context.Context, ticker *models.Ticker) error {
	pairUrl := fmt.Sprintf(PublicURLTicker+"/%s", ticker.Pair)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, pairUrl, nil)
	if err != nil {
		return errors.Wrap(err, "error creating api request")
	}

	resp, err := a.client.Do(req)
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
func (a *UpholdApi) ParseAPIData(response *http.Response, ticker *models.Ticker) error {
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
func (a *UpholdApi) IsPairValid(pair string) (bool, error) {
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
