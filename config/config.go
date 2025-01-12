package config

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"os"
)

// DatabaseConfig holds the configuration for the database connection
type DatabaseConfig struct {
	Host         string
	User         string
	Password     string
	Port         string
	Name         string
	Schema       string
	TableAlerts  string
	TableConfigs string
}

// LoadDatabaseConfig loads the database configuration from the environment variables defined on docker-compose.yml
func LoadDatabaseConfig() *DatabaseConfig {
	return &DatabaseConfig{
		User:         os.Getenv("USER"),
		Host:         os.Getenv("HOST"),
		Password:     os.Getenv("PASSWORD"),
		Port:         os.Getenv("PORT"),
		Name:         os.Getenv("NAME"),
		Schema:       os.Getenv("SCHEMA"),
		TableAlerts:  os.Getenv("TABLE_ALERTS"),
		TableConfigs: os.Getenv("TABLE_CONFIGS"),
	}
}

// ConnectToDatabase connects to the database using the provided configuration
func ConnectToDatabase(config *DatabaseConfig) (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s user=%s password=%s port=%s dbname=%s sslmode=disable",
		config.Host, config.User, config.Password, config.Port, config.Name)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, errors.Wrap(err, "error opening connection to db")
	}

	if err = db.Ping(); err != nil {
		return nil, errors.Wrap(err, "error pinging connection to db")
	}

	return db, nil
}
