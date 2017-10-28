-- user withdrawals from balance
CREATE TABLE withdrawals
(
  id              SERIAL                                 NOT NULL
    CONSTRAINT withdrawals_pkey
    PRIMARY KEY,
  user_steam_id   BIGINT                                 NOT NULL
    CONSTRAINT withdrawals_user_steam_id_foreign
    REFERENCES users,
  status          INTEGER DEFAULT 9 /* PENDING */        NOT NULL,
  payment_address VARCHAR(60)                            NOT NULL,
  usd_rate        NUMERIC(10, 2)                         NOT NULL,
  currency        INTEGER                                NOT NULL,
  usd_total       NUMERIC(10, 2)                         NOT NULL,
  crypto_total    NUMERIC(10, 8)                         NOT NULL,
  txhash          VARCHAR(80),                            
  created_at      TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL
);
CREATE INDEX withdrawals_user_steam_id_index
  ON withdrawals (user_steam_id);  
CREATE INDEX withdrawals_status_index
  ON withdrawals (status);     