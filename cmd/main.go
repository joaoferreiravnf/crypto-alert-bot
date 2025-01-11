package main

import (
	"fmt"
	"log"
	"sync"
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

	for _, t := range *tickers {
		go runSchedulerBot(&wg, *t, apiResponse, repo)
	}

	wg.Wait()

	fmt.Println("All tickers finished")
}

func runSchedulerBot(wg *sync.WaitGroup, tickerValue models.Ticker, apiResponse *api.ApiResponse, repo *postgres.Postgres) {
	defer wg.Done()

	tickerPublisher := logger.NewTickerPublisher()

	tickerScheduler := services.NewTickerScheduler(apiResponse, &tickerValue, repo, tickerPublisher)
	tickerScheduler.SchedulerStart()

	if tickerValue.Config.Lifetime > 0 {
		time.Sleep(tickerValue.Config.Lifetime * time.Second)

		tickerScheduler.SchedulerStop()
	} else {
		select {}
	}
}
