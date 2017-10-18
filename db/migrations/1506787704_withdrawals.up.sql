-- user withdrawals from balance
CREATE TABLE withdrawals
(
  id              SERIAL                                 NOT NULL
    CONSTRAINT withdrawals_pkey
    PRIMARY KEY,
  user_steam_id   BIGINT                                 NOT NULL
    CONSTRAINT withdrawals_user_steam_id_foreign
    REFERENCES users,
  tradeoffer_id   BIGINT
    CONSTRAINT withdrawals_tradeoffer_id_foreign
    REFERENCES tradeoffers,
  status          INTEGER DEFAULT 9 /* PENDING */        NOT NULL,
  payment_address VARCHAR(60)                            NOT NULL,
  usd_rate        NUMERIC(10, 2)                         NOT NULL,
  amount          NUMERIC(10, 8)                         NOT NULL,
  txhash          VARCHAR(80)                            NOT NULL,
  created_at      TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL
);
CREATE INDEX withdrawals_user_steam_id_index
  ON withdrawals (user_steam_id);
CREATE INDEX withdrawals_tradeoffer_id_index
  ON withdrawals (tradeoffer_id);  