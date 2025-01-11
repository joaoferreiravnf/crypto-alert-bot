CREATE SCHEMA IF NOT EXISTS uphold_alerts;

CREATE TABLE uphold_alerts.alerts (
      id SERIAL PRIMARY KEY,
      pair VARCHAR(20) NOT NULL,
      price_change NUMERIC(30, 20) NOT NULL,
      perc_change NUMERIC(30, 20) NOT NULL,
      final_price numeric(30, 20) NOT NULL,
      config_refresh INT NOT NULL,
      config_perc_oscillation NUMERIC(20, 10) NOT NULL,
      timestamp TIMESTAMP NOT NULL
);