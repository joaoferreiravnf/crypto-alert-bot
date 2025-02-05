package prompt

import (
	"bufio"
	"crypto-alert-bot/internal/models"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

//go:generate mockgen -source=$GOFILE -destination=../mocks/mock_api/mock_$GOFILE
type ApiDataValidator interface {
	IsPairValid(pair string) (bool, error)
}

// AskUserInput prompts the user for various inputs and returns a slice of ticker structs
func AskUserInput(validator ApiDataValidator) *models.Tickers {
	var tickers models.Tickers

	for {
		pair := promptPair(validator)
		refreshRate := promptRefreshRate()
		percThreshold := promptPercThreshold()
		lifetime := promptLifetime()

		ticker := models.NewTicker(pair, refreshRate, percThreshold, lifetime)

		tickers = append(tickers, ticker)

		multiplePairs := promptMultiplePairs()

		if multiplePairs {
			continue
		}

		if tickers.IsAboveRateLimit() {
			fmt.Println("The current bot configuration would exceed the rate limit. Please adjust bot configurations")
			continue
		}

		return &tickers
	}
}

// promptPair prompts the user to choose a trading pair
func promptPair(validator ApiDataValidator) string {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("Chose a trading pair (e.g. BTCUSD, ETHEUR): ")

		input, _ := reader.ReadString('\n')

		pair := strings.TrimSpace(input)
		pair = strings.ToUpper(pair)

		isValid, err := validator.IsPairValid(pair)

		if !isValid {
			fmt.Println(err)
			continue
		} else {
			return pair
		}
	}
}

// promptRefreshRate prompts the user to choose a refresh rate
func promptRefreshRate() float64 {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("Define a refresh rate in seconds (e.g. 1, 20, 60): ")

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		refreshRate, err := strconv.ParseFloat(input, 64)

		if err != nil || refreshRate < 0 {
			fmt.Println("Invalid choice. Please enter a positive and valid integer")
			continue
		}

		return refreshRate
	}
}

// promptPercThreshold prompts the user to choose a percentage threshold
func promptPercThreshold() float64 {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("Define a percent change threshold (e.g. 0,02 or 2): ")

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		percThreshold, err := strconv.ParseFloat(input, 64)

		if err != nil || percThreshold < 0 {
			fmt.Println("Invalid choice. Please enter a positive and valid floating point")
			continue
		}

		return percThreshold
	}
}

// promptLifetime prompts the user to choose a lifetime for the ticker bot
func promptLifetime() time.Duration {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("Set for how long, in seconds, the ticker should run (e.g. 30, 100, 5000) or just hit enter if forever): ")

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "" {
			return 0
		}
		refreshRate, err := strconv.Atoi(input)

		if err != nil || refreshRate < 0 {
			fmt.Println("Invalid choice. Please enter a positive and valid integer")
			continue
		}

		return time.Duration(refreshRate)
	}
}

// promptMultiplePairs prompts the user if they want to track more pairs
func promptMultiplePairs() bool {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Do you want to track more pairs? If yes write y, otherwise write any other letter or just hit enter: ")

	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	input = strings.ToUpper(input)

	if input == "Y" {
		return true
	}

	return false
}
