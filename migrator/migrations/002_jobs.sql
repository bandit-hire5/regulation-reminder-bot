-- +migrate Up

CREATE TABLE jobs(
  id                BIGSERIAL PRIMARY KEY,
  user_id           BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  name              VARCHAR(255) NOT NULL,
  regulation        BIGINT NOT NULL,
  last_updated_at    DATE NOT NULL,
  last_odometer      BIGINT NOT NULL,
  next_odometer      BIGINT NOT NULL,
  left_odometer      BIGINT NOT NULL
);

INSERT INTO jobs (user_id, name, regulation, last_updated_at, last_odometer, next_odometer, left_odometer)
VALUES
  (1, 'Масло ДВС, фильтр масляный', 10000, '2019-08-15', 119000, 129000, 9300),
  (1, 'Регулировка клапанов', 30000, '2009-08-15', 0, 129000, 9300),
  (1, 'Воздушный фильтр', 20000, '2019-08-15', 119000, 139000, 9300),
  (1, 'Фильтр салона', 20000, '2019-08-15', 119000, 139000, 9300),
  (1, 'Масло АКПП', 45000, '2009-08-15', 0, 129000, 9300),
  (1, 'Задний редуктор', 15000, '2019-08-15', 119000, 134000, 9300),
  (1, 'Раздатка', 45000, '2009-08-15', 0, 129000, 9300),
  (1, 'Тормозная жидкость', 45000, '2009-08-15', 0, 129000, 9300),
  (1, 'Антифриз', 100000, '2009-08-15', 0, 129000, 9300),
  (1, 'Жидкость ГУР', 45000, '2009-08-15', 0, 129000, 9300),
  (1, 'Свечи', 100000, '2009-08-15', 0, 129000, 9300),
  (1, 'Помпа', 100000, '2009-08-15', 0, 200000, 9300),
  (1, 'Ремень ГРМ', 120000, '2019-08-23', 119000, 240000, 9300),
  (1, 'Ремень навесного', 100000, '2009-08-15', 0, 200000, 9300),
  (1, 'Газовые фильтра', 10000, '2019-08-15', 119700, 129000, 9300),
  (1, 'Топливный фильтр', 90000, '2009-09-03', 119000, 129000, 9300);

-- +migrate Down

DROP TABLE IF EXISTS jobs;