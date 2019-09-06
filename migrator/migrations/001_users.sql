-- +migrate Up

CREATE TABLE users(
  id BIGSERIAL PRIMARY KEY,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  telegram_id BIGSERIAL NOT NULL,
  first_name VARCHAR(255) NOT NULL,
  last_name VARCHAR(255) NOT NULL,
  username VARCHAR(255) NOT NULL,
  odometer BIGINT NOT NULL DEFAULT 0,
  enable_monitoring BOOLEAN NOT NULL DEFAULT FALSE
);

INSERT INTO users (telegram_id, first_name, last_name, username, odometer, enable_monitoring)
VALUES
  (503868934, 'Руслан', 'Наэльток', 'bandit_hire5', 0, false);

-- +migrate Down

DROP TABLE IF EXISTS users;