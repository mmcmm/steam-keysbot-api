CREATE TABLE trade_offers
(
  id                SERIAL                                 NOT NULL
    CONSTRAINT trade_offers_pkey
    PRIMARY KEY,
  user_steam_id     BIGINT                                 NOT NULL
    CONSTRAINT trade_offers_user_steam_id_foreign
    REFERENCES users,
  merchant_steam_id BIGINT,
  steam_offer_id    BIGINT,
  failure_details   VARCHAR(255),
  amount            INTEGER,
  type              VARCHAR(20),
  created_at        TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
  updated_at        TIMESTAMP WITH TIME ZONE
);
CREATE INDEX trade_offers_user_steam_id_index
  ON trade_offers (user_steam_id);