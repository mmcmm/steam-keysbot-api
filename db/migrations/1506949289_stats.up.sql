CREATE TABLE stats
(
  id                SERIAL                                 NOT NULL
    CONSTRAINT stats_pkey
    PRIMARY KEY,
  users_nr          INTEGER                                NOT NULL,
  daily_tradeoffers INTEGER                                NOT NULL,
  money_transacted  NUMERIC(10, 2)                         NOT NULL,
  current_value     NUMERIC(10, 2)                         NOT NULL,
  users_stored      NUMERIC(10, 2)                         NOT NULL,
  updated_at        TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL
);

INSERT INTO stats (id, users_nr, daily_tradeoffers, money_transacted, current_value, users_stored) 
VALUES (1, 0, 0, 0, 0, 0) ON CONFLICT DO NOTHING;