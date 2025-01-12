package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/pkg/errors"
	"time"
	"uphold-alert-bot/internal/models"
)

// Postgres represents the postgres postgres
type Postgres struct {
	DB             *sql.DB
	DbSchema       string
	DbTableConfigs string
	DbTableAlerts  string
}

// NewPostgres returns a new instance of Postgres
func NewPostgres(db *sql.DB, dbSchema, dbTableConfigs, dbTableAlerts string) *Postgres {
	return &Postgres{
		DB:             db,
		DbSchema:       dbSchema,
		DbTableConfigs: dbTableConfigs,
		DbTableAlerts:  dbTableAlerts,
	}
}

// Save saves the ticker to the database
func (p *Postgres) Save(ctx context.Context, timestamp time.Time, ticker *models.Ticker) error {
	tx, err := p.DB.BeginTx(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "failed to start transaction")
	}
	defer tx.Rollback()

	configQuery := fmt.Sprintf("INSERT INTO %s.%s (refresh_rate, perc_oscillation) VALUES ($1, $2) RETURNING ID", p.DbSchema, p.DbTableConfigs)

	var configID int
	err = tx.QueryRowContext(ctx, configQuery, ticker.Config.RefreshRate, ticker.Config.PercOscillation).Scan(&configID)
	if err != nil {
		return errors.Wrap(err, "failed to save ticker configs into configs table")
	}

	query := fmt.Sprintf("INSERT INTO %s.%s (pair, price_change, perc_change, final_price, config_id, timestamp) VALUES ($1, $2, $3, $4, $5, $6)",
		p.DbSchema, p.DbTableAlerts)

	_, err = tx.ExecContext(ctx, query, ticker.Pair, ticker.AskPriceChange, ticker.AskPercChange, ticker.CurrentAsk, configID, timestamp)
	if err != nil {
		return errors.Wrap(err, "failed to save ticker into alerts table")
	}

	err = tx.Commit()
	if err != nil {
		return errors.Wrap(err, "failed to commit transaction")
	}

	return nil
}
