-- user withdrawals from balance
CREATE TABLE withdrawals
(
  id            SERIAL                                 NOT NULL
    CONSTRAINT withdrawals_pkey
    PRIMARY KEY,
  user_steam_id BIGINT                                 NOT NULL
    CONSTRAINT withdrawals_user_steam_id_foreign
    REFERENCES users,
  created_at    TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL
);
CREATE INDEX withdrawals_user_steam_id_index
  ON withdrawals (user_steam_id);