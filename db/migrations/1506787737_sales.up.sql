-- user sales to us
CREATE TABLE sales
(
  id                SERIAL                                 NOT NULL
    CONSTRAINT sales_pkey
    PRIMARY KEY,
  user_steam_id     BIGINT                                 NOT NULL
    CONSTRAINT sales_user_steam_id_foreign
    REFERENCES users,
  paid_to           VARCHAR(7)                             NOT NULL,
  created_at        TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL
);
CREATE INDEX sales_user_steam_id_index
  ON sales (user_steam_id);