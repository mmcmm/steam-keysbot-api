CREATE TABLE steam_bots (
  steam_id        BIGINT             NOT NULL
    CONSTRAINT steam_bots_pkey
    PRIMARY KEY,
  ip_address      VARCHAR(45)        NOT NULL
    CONSTRAINT steam_bots_ip_address_unique
    UNIQUE,
  trade_link_url  VARCHAR(120),
  csgokey_balance INTEGER DEFAULT 0  NOT NULL,
  display_name    VARCHAR(40)        NOT NULL,
  updated_at      TIMESTAMP WITH TIME ZONE DEFAULT now()
);

CREATE OR REPLACE FUNCTION stats_steam_bots_keys_balance()
  RETURNS TRIGGER AS
$func$
BEGIN
  IF NEW.csgokey_balance <> OLD.csgokey_balance
  THEN
    UPDATE stats
    SET current_keys = current_keys + NEW.csgokey_balance - OLD.csgokey_balance;
  END IF;
  RETURN NEW;
END;
$func$
LANGUAGE plpgsql;

CREATE TRIGGER stats_steam_bots_keys_balance
AFTER UPDATE ON steam_bots FOR EACH ROW
EXECUTE PROCEDURE stats_steam_bots_keys_balance();