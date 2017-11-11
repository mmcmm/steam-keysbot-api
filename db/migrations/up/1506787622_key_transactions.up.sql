-- user key_transactions from our bots
CREATE TABLE key_transactions
(
  id               SERIAL                                 NOT NULL
    CONSTRAINT key_transactions_pkey
    PRIMARY KEY,
  user_steam_id    BIGINT                                 NOT NULL
    CONSTRAINT key_transactions_user_steam_id_foreign
    REFERENCES users,
  tradeoffer_id    BIGINT                                 NOT NULL
    CONSTRAINT key_transactions_tradeoffer_id_foreign
    REFERENCES tradeoffers
    CONSTRAINT key_transactions_tradeoffer_id_unique
    UNIQUE,
  status           INTEGER DEFAULT 9 /* PENDING */        NOT NULL,
  type             INTEGER                                NOT NULL,
  transaction_type INTEGER                                NOT NULL,
  amount           INTEGER                                NOT NULL,
  unit_price       NUMERIC(10, 2)                         NOT NULL,
  payment_address  VARCHAR(60),
  usd_rate         NUMERIC(10, 2)                         NOT NULL,
  currency         INTEGER                                NOT NULL,
  usd_total        NUMERIC(10, 2)                         NOT NULL,
  crypto_total     NUMERIC(10, 8)                         NOT NULL,
  app_id           INTEGER                                NOT NULL,
  created_at       TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL
);
CREATE INDEX key_transactions_user_steam_id_index
  ON key_transactions (user_steam_id);

CREATE OR REPLACE FUNCTION stats_money_transacted_increment()
  RETURNS TRIGGER AS
$func$
BEGIN
  IF NEW.status = 7 /* COMPLETED */
  THEN
    UPDATE stats
    SET money_transacted = money_transacted + NEW.usd_total;
  END IF;
  RETURN NEW;
END;
$func$
LANGUAGE plpgsql;

CREATE TRIGGER stats_money_transacted_increment
AFTER UPDATE ON key_transactions
FOR EACH ROW
EXECUTE PROCEDURE stats_money_transacted_increment()