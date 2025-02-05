# Crypto Alert Bot

This project consists of a bot that continuously checks specified trading pairs and triggers an alert whenever a given price change threshold is met. All thresholds met are saved into a Postgresql database

### How It Works
1. User Prompts: Upon starting, the bot asks for:
- A Trading Pair (e.g., BTCUSD, ETHEUR)
- Refresh Interval (in seconds) for API data querying
- Percentage Threshold for price oscillation
- Lifetime (in seconds) the bot should run: if no value is provided, it runs indefinitely
- If you want to monitor more than one pair: if yes, just enter "Y" and the bot will aks for the next pair

2. Data Fetching: The bot periodically queries the API _api.uphold.com/v0/ticker/:pair_ to retrieve up-to-date bid/ask prices for your chosen trading pairs


3. Alert Logic:
- It compares current ask prices with previous ask prices (bot developed from the buyer's perspective)
- If the percentage change exceeds your specified threshold, an alert is logged and the event is stored in the database
4. Database: 
- It uses Flyway to manage schema migrations, ensuring the database table structure is set up before the bot starts

### Prerequisites
- Before starting, make sure you have installed:
1. Docker
2. Docker Compose
3. Make (for ease of use)

### Getting Started
1. Clone the Repository:
- Open your terminal and enter:
```
git clone https://github.com/joaoferreiravnf/crypto-alert-bot.git
cd crypto-alert-bot
```
2. Run the Bot Locally:
- To start the bot, run:
```
make run-local
``` 
- Alternatively, if by any reason you can't use Make, you can run:
```
docker-compose build
docker-compose up -d db
docker-compose up flyway
docker-compose run --rm bot
```

This / these command spins up:
- A Postgresql database
- The Flyway migration container, applying database migrations
- The Alert Bot container

3. Interact with the Bot:
- You’ll see prompts in the terminal asking for your input
- Enter your desired trading pairs, refresh intervals, thresholds, and lifetimes

4. Stop the Bot:
- Press Ctrl + C in your terminal or run:
```
docker compose down
```

5. Query the database:
- At any point, before or after stopping the bot, you can check the database for the stored alerts:
```
docker exec -it crypto_alert_db psql -U postgres -d crypto_alert_db
SELECT * FROM crypto_alerts.alerts;
```

### Project Structure
- cmd/main.go: Entry point for the application


- internal:
  - api: Responsible for connecting and retrieving data from API
  - models: Defines the domain entities (e.g. Ticker) and related logic
  - prompt: Handles all user input prompts
  - repository: Manages saving ticker events to the Postgres database
  - services: Holds core functionality as scheduling and alerts publishing
  - config: Database connection and configuration loading logic
  - migrations: SQL migration scripts run by Flyway


- Dockerfile: Multi-stage build for a minimal container image


- docker-compose.yml: Orchestrates services (Postgres, Flyway and the Bot) 

