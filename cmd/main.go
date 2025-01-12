package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
	"uphold-alert-bot/config"
	"uphold-alert-bot/internal/adapters/api"
	"uphold-alert-bot/internal/adapters/logger"
	"uphold-alert-bot/internal/adapters/postgres"
	"uphold-alert-bot/internal/adapters/prompt"
	"uphold-alert-bot/internal/models"
	"uphold-alert-bot/internal/services"
)

func main() {
	loadDbConfigs := config.LoadDatabaseConfig()

	db, err := config.ConnectToDatabase(loadDbConfigs)
	if err != nil {
		log.Fatal("error on initializing db connection", err)
	}
	defer db.Close()

	repo := postgres.NewPostgres(db, loadDbConfigs.Schema, loadDbConfigs.Table)

	apiResponse := api.NewAPIResponse(nil)

	tickers := prompt.AskUserInput(apiResponse)

	var wg sync.WaitGroup

	wg.Add(len(*tickers))

	fmt.Println("Starting bot")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
		<-sigChan
		fmt.Println("Shutting down...")
		cancel()
	}()

	for _, t := range *tickers {
		go runSchedulerBot(ctx, &wg, *t, apiResponse, repo)
	}

	wg.Wait()

	fmt.Println("All tickers finished")
}

func runSchedulerBot(ctx context.Context, wg *sync.WaitGroup, tickerValue models.Ticker, apiResponse *api.ApiResponse, repo *postgres.Postgres) {
	defer wg.Done()

	tickerPublisher := logger.NewTickerPublisher()

	tickerScheduler := services.NewTickerScheduler(apiResponse, &tickerValue, repo, tickerPublisher)
	tickerScheduler.SchedulerStart(ctx)

	if tickerValue.Config.Lifetime > 0 {
		select {
		case <-time.After(tickerValue.Config.Lifetime * time.Second):
			tickerScheduler.SchedulerStop()
		case <-ctx.Done():
			fmt.Println("Shutting down scheduler for", tickerValue.Pair)
			tickerScheduler.SchedulerStop()
		}
	} else {
		select {
		case <-ctx.Done():
			fmt.Println("Shutting down scheduler for", tickerValue.Pair)
			tickerScheduler.SchedulerStop()
			return
		}
	}
}
