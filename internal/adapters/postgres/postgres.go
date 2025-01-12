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
	DB       *sql.DB
	DbSchema string
	DbTable  string
}

// NewPostgres returns a new instance of Postgres
func NewPostgres(db *sql.DB, dbSchema, dbTable string) *Postgres {
	return &Postgres{
		DB:       db,
		DbSchema: dbSchema,
		DbTable:  dbTable,
	}
}

// Save saves the ticker to the database
func (p *Postgres) Save(ctx context.Context, timestamp time.Time, ticker *models.Ticker) error {
	query := fmt.Sprintf("INSERT INTO %s.%s (pair, price_change, perc_change, final_price, config_refresh, config_perc_oscillation, timestamp) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		p.DbSchema, p.DbTable)

	_, err := p.DB.ExecContext(ctx, query, ticker.Pair, ticker.AskPriceChange, ticker.AskPercChange, ticker.CurrentAsk, ticker.Config.RefreshRate, ticker.Config.PercOscillation, timestamp)
	if err != nil {
		return errors.Wrap(err, "failed to save ticker to postgres")
	}

	return nil
}
