-- +migrate Up

CREATE TYPE types AS ENUM ('odometer', 'date');

CREATE TABLE jobs(
  id                 BIGSERIAL PRIMARY KEY,
  user_id            BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  type               types,
  name               VARCHAR(255) NOT NULL,
  per_days           BIGINT NOT NULL,
  per_odometer       BIGINT NOT NULL,
  last_date          DATE,
  last_odometer      BIGINT NOT NULL,
  periodically       BOOLEAN DEFAULT true
);

INSERT INTO jobs (user_id, name, type, per_days, per_odometer, last_date, last_odometer, periodically)
VALUES
  (1, 'Масло ДВС, фильтр масляный', 'odometer', 0, 10000, '2019-08-15', 119000, true),
  (1, 'Регулировка клапанов', 'odometer', 0, 30000, '2009-08-15', 0, true),
  (1, 'Воздушный фильтр', 'odometer', 0, 20000, '2019-08-15', 119000, true),
  (1, 'Фильтр салона', 'odometer', 0, 20000, '2019-08-15', 119000, true),
  (1, 'Масло АКПП', 'odometer', 0, 45000, '2009-08-15', 0, true),
  (1, 'Задний редуктор', 'odometer', 0, 15000, '2019-08-15', 119000, true),
  (1, 'Раздатка', 'odometer', 0, 45000, '2009-08-15', 0, true),
  (1, 'Тормозная жидкость', 'odometer', 0, 45000, '2009-08-15', 0, true),
  (1, 'Антифриз', 'odometer', 0, 100000, '2009-08-15', 0, true),
  (1, 'Жидкость ГУР', 'odometer', 0, 45000, '2009-08-15', 0, true),
  (1, 'Свечи', 'odometer', 0, 100000, '2009-08-15', 0, true),
  (1, 'Помпа', 'odometer', 0, 100000, '2009-08-15', 0, true),
  (1, 'Ремень ГРМ', 'odometer', 0, 120000, '2019-08-23', 119000, true),
  (1, 'Ремень навесного', 'odometer', 0, 100000, '2009-08-15', 0, true),
  (1, 'Газовые фильтра', 'odometer', 0, 10000, '2019-08-15', 119700, true),
  (1, 'Топливный фильтр', 'odometer', 0, 90000, '2009-09-03', 119000, true),
  (1, 'Лифт/Проставки', 'date', 40, 0, NULL, 119000, false);

-- +migrate Down

DROP TABLE IF EXISTS jobs;
DROP TYPE IF EXISTS types;