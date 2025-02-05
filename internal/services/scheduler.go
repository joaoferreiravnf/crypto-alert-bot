package services

import (
	"context"
	"log/slog"
	"time"
	"crypto-alert-bot/internal/models"
)

var apiTimeout = 5 * time.Second
var dbTimeout = 5 * time.Second

//go:generate mockgen -source=$GOFILE -destination=../mocks/mock_scheduler/mock_$GOFILE
type DataRetriever interface {
	FetchPairData(context.Context, *models.Ticker) error
}

//go:generate mockgen -source=$GOFILE -destination=../mocks/mock_scheduler/mock_$GOFILE
type Publisher interface {
	Publish(time.Time, *models.Ticker)
}

//go:generate mockgen -source=$GOFILE -destination=../mocks/mock_scheduler/mock_$GOFILE
type Recorder interface {
	Save(context.Context, time.Time, *models.Ticker) error
}

// TickerScheduler represents the scheduler for the ticker, orchestrating the fetching of data and publishing of alerts
type TickerScheduler struct {
	api       DataRetriever
	ticker    *models.Ticker
	publisher Publisher
	repo      Recorder
	stop      chan struct{}
}

// NewTickerScheduler returns a new instance of TickerScheduler
func NewTickerScheduler(apiResponse DataRetriever, ticker *models.Ticker, repo Recorder, publisher Publisher) *TickerScheduler {
	return &TickerScheduler{
		api:       apiResponse,
		ticker:    ticker,
		publisher: publisher,
		repo:      repo,
		stop:      make(chan struct{}),
	}
}

// SchedulerStart starts the scheduler
func (ts *TickerScheduler) SchedulerStart(ctx context.Context) {
	interval := time.Duration((ts.ticker.Config.RefreshRate) * float64(time.Second))
	timeTicker := time.NewTicker(interval)

	go func() {
		for {
			select {
			case <-timeTicker.C:
				apiCtx, cancel := context.WithTimeout(ctx, apiTimeout*time.Second)
				defer cancel()

				err := ts.api.FetchPairData(apiCtx, ts.ticker)
				if err != nil {
					slog.Error("error fetching data", "error", err)
					break
				}

				if !ts.ticker.IsAbovePercOscillation() {
					break
				}

				dbCtx, dbCancel := context.WithTimeout(ctx, dbTimeout*time.Second)
				defer dbCancel()

				timestamp := time.Now().UTC()

				ts.publisher.Publish(timestamp, ts.ticker)

				err = ts.repo.Save(dbCtx, timestamp, ts.ticker)
				if err != nil {
					slog.Error("error saving to database", "error", err)
				}

				ts.ticker.NormalizeValues()

			case <-ctx.Done():
				slog.Info("scheduler canceled by context")
				timeTicker.Stop()
				return
			}
		}
	}()
}

// SchedulerStop stops the scheduler
func (ts *TickerScheduler) SchedulerStop() {
	close(ts.stop)
}
