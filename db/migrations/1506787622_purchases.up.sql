-- user purchases from our bots
CREATE TABLE purchases
(
  id                SERIAL                                 NOT NULL
    CONSTRAINT purchases_pkey
    PRIMARY KEY,
  user_steam_id     BIGINT                                 NOT NULL
    CONSTRAINT purchases_user_steam_id_foreign
    REFERENCES users,
  created_at        TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL
);
CREATE INDEX purchases_user_steam_id_index
  ON purchases (user_steam_id);