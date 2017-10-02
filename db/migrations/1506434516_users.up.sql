CREATE TABLE users
(
  steam_id  BIGINT NOT NULL
    CONSTRAINT users_pkey
    PRIMARY KEY,
  trade_link_url VARCHAR(120),
  bitcoin_balance INTEGER DEFAULT 0  NOT NULL,
  csgokey_balance INTEGER DEFAULT 0  NOT NULL,
  display_name  VARCHAR(40) NOT NULL,
  avatar        VARCHAR(255) NOT NULL,
  added_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);