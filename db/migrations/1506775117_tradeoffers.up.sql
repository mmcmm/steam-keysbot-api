CREATE TABLE tradeoffers
(
  id                SERIAL                                 NOT NULL
    CONSTRAINT tradeoffers_pkey
    PRIMARY KEY,
  user_steam_id     BIGINT                                 NOT NULL
    CONSTRAINT tradeoffers_user_steam_id_foreign
    REFERENCES users,
  type              VARCHAR(25)                            NOT NULL,
  status            VARCHAR(25)                            NOT NULL,
  merchant_steam_id BIGINT,                        
  steam_offer_id    BIGINT,
  failure_details   VARCHAR(255),
  amount            INTEGER                                NOT NULL,
  created_at        TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL
);
CREATE INDEX tradeoffers_user_steam_id_index
  ON tradeoffers (user_steam_id);