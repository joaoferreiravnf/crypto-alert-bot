CREATE SCHEMA IF NOT EXISTS crypto_alerts;

CREATE TABLE crypto_alerts.configs (
      id SERIAL PRIMARY KEY,
      refresh_rate NUMERIC(10, 5) NOT NULL,
      perc_oscillation NUMERIC(20, 10) NOT NULL
);

CREATE TABLE crypto_alerts.alerts (
      id SERIAL PRIMARY KEY,
      pair VARCHAR(20) NOT NULL,
      price_change NUMERIC(30, 20) NOT NULL,
      perc_change NUMERIC(30, 20) NOT NULL,
      final_price NUMERIC(30, 20) NOT NULL,
      config_id INT NOT NULL REFERENCES crypto_alerts.configs(id),
      timestamp TIMESTAMP NOT NULL
);