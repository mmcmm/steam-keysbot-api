CREATE TABLE stats
(
  id                SERIAL                                 NOT NULL
    CONSTRAINT stats_pkey
    PRIMARY KEY,
  users_nr          INTEGER                                NOT NULL,
  money_transacted  NUMERIC(10, 2)                         NOT NULL,
  current_keys      INTEGER                                NOT NULL,
  users_stored      INTEGER                                NOT NULL
);

INSERT INTO stats (id, users_nr, money_transacted, current_keys, users_stored) 
VALUES (1, 0, 0, 0, 0) ON CONFLICT DO NOTHING;