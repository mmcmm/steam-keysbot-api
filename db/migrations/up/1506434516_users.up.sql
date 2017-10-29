CREATE TABLE users
(
  steam_id        BIGINT             NOT NULL
    CONSTRAINT users_pkey
    PRIMARY KEY,
  trade_link_url  VARCHAR(120),
  bitcoin_balance NUMERIC(10, 8) DEFAULT 0  NOT NULL,
  display_name    VARCHAR(40)               NOT NULL,
  avatar          VARCHAR(255)              NOT NULL,
  added_at        TIMESTAMP WITH TIME ZONE DEFAULT now()
);

CREATE OR REPLACE FUNCTION stats_users_nr_increment()
  RETURNS TRIGGER AS
$func$
BEGIN
  UPDATE stats SET users_nr = users_nr + 1;
  RETURN NEW;
END;
$func$
LANGUAGE plpgsql;

CREATE TRIGGER stats_users_nr_increment
AFTER INSERT ON users FOR EACH ROW
EXECUTE PROCEDURE stats_users_nr_increment();

CREATE OR REPLACE FUNCTION stats_users_stored_increment()
  RETURNS TRIGGER AS
$func$
BEGIN
IF NEW.bitcoin_balance <> OLD.bitcoin_balance THEN
  UPDATE stats SET users_stored = users_stored + NEW.bitcoin_balance - OLD.bitcoin_balance;
END IF;  
  RETURN NEW;
END;
$func$
LANGUAGE plpgsql;

CREATE TRIGGER stats_users_stored_increment
AFTER UPDATE ON users FOR EACH ROW
EXECUTE PROCEDURE stats_users_stored_increment();