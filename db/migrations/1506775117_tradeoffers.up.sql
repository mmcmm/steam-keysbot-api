CREATE TABLE tradeoffers
(
  id                SERIAL                                 NOT NULL
    CONSTRAINT tradeoffers_pkey
    PRIMARY KEY,
  user_steam_id     BIGINT                                 NOT NULL
    CONSTRAINT tradeoffers_user_steam_id_foreign
    REFERENCES users,
  type              INTEGER                                NOT NULL,
  status            INTEGER                                NOT NULL,
  merchant_steam_id BIGINT,
  steam_offer_id    BIGINT,
  failure_details   VARCHAR(255),
  user_items        VARCHAR(255),
  bot_items         VARCHAR(255),
  amount            INTEGER                                NOT NULL,
  created_at        TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL
);
CREATE INDEX tradeoffers_user_steam_id_index
  ON tradeoffers (user_steam_id);


CREATE OR REPLACE FUNCTION stats_tradeoffers_daily_increment()
  RETURNS TRIGGER AS
$func$
BEGIN
IF NEW.status = 2 THEN -- ACCEPTED
  UPDATE stats SET daily_tradeoffers = daily_tradeoffers + 1;
END IF;  
  RETURN NEW;
END;
$func$
LANGUAGE plpgsql;

CREATE TRIGGER stats_tradeoffers_daily_increment
AFTER UPDATE ON tradeoffers FOR EACH ROW
EXECUTE PROCEDURE stats_tradeoffers_daily_increment()