-- user purchases from our bots
CREATE TABLE purchases
(
  id              SERIAL                                 NOT NULL
    CONSTRAINT purchases_pkey
    PRIMARY KEY,
  user_steam_id   BIGINT                                 NOT NULL
    CONSTRAINT purchases_user_steam_id_foreign
    REFERENCES users,
  tradeoffer_id   BIGINT                                 NOT NULL
    CONSTRAINT purchases_tradeoffer_id_foreign
    REFERENCES tradeoffers,
  status          VARCHAR(25)                            NOT NULL,
  type            VARCHAR(25)                            NOT NULL,
  amount          INTEGER                                NOT NULL,
  unit_price      NUMERIC(10, 2)                         NOT NULL,
  payment_address VARCHAR(60),
  usd_rate        NUMERIC(10, 2)                         NOT NULL,
  currency        VARCHAR(3)                             NOT NULL,
  usd_total       NUMERIC(10, 2)                         NOT NULL,
  crypto_total    NUMERIC(10, 2)                         NOT NULL,
  created_at      TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL
);
CREATE INDEX purchases_user_steam_id_index
  ON purchases (user_steam_id);
CREATE INDEX purchases_tradeoffer_id_index
  ON purchases (tradeoffer_id);
CREATE INDEX purchases_payment_address_index
  ON purchases (payment_address);  