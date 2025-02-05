package main

import (
	"context"
	"crypto-alert-bot/config"
	"crypto-alert-bot/internal/adapters/api"
	"crypto-alert-bot/internal/adapters/logger"
	"crypto-alert-bot/internal/adapters/postgres"
	"crypto-alert-bot/internal/adapters/prompt"
	"crypto-alert-bot/internal/models"
	"crypto-alert-bot/internal/services"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	loadDbConfigs := config.LoadDatabaseConfig()

	db, err := config.ConnectToDatabase(loadDbConfigs)
	if err != nil {
		log.Fatal("error on initializing db connection", err)
	}
	defer db.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	repo := postgres.NewPostgres(db, loadDbConfigs.Schema, loadDbConfigs.TableConfigs, loadDbConfigs.TableAlerts)

	upholdApi := api.NewUpholdApi(nil)

	publisher := logger.NewTickerPublisher()

	tickers := prompt.AskUserInput(upholdApi)

	var wg sync.WaitGroup

	wg.Add(len(*tickers))

	fmt.Println("Starting bot")

	for _, t := range *tickers {
		go runSchedulerBot(ctx, &wg, *t, upholdApi, repo, publisher)
	}

	go gracefulShutdown(cancel)

	wg.Wait()
}

func runSchedulerBot(ctx context.Context, wg *sync.WaitGroup, ticker models.Ticker, upholdApi *api.UpholdApi, repo *postgres.Postgres, publisher services.Publisher) {
	defer wg.Done()

	tickerScheduler := services.NewTickerScheduler(upholdApi, &ticker, repo, publisher)

	tickerScheduler.SchedulerStart(ctx)

	if ticker.Config.Lifetime > 0 {
		select {
		case <-time.After(ticker.Config.Lifetime * time.Second):

			tickerScheduler.SchedulerStop()

			fmt.Printf("Scheduler for %s completed", ticker.Pair)
		case <-ctx.Done():
			fmt.Println("Shutting down scheduler for", ticker.Pair)

			tickerScheduler.SchedulerStop()
		}
	} else {
		select {
		case <-ctx.Done():
			fmt.Println("Shutting down scheduler for", ticker.Pair)

			tickerScheduler.SchedulerStop()

			return
		}
	}
}

func gracefulShutdown(cancel context.CancelFunc) {
	sigChan := make(chan os.Signal, 1)

	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	<-sigChan

	fmt.Println("Shutting down...")

	cancel()
}
