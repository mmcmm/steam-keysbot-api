CREATE TABLE settings
(
  id              SERIAL                                        NOT NULL
    CONSTRAINT settings_pkey
    PRIMARY KEY,
  btc_to_usd_rate         NUMERIC(10, 2)                        NOT NULL,
  buy_csgokey_price       NUMERIC(10, 2)                        NOT NULL,
  sell_csgokey_price      NUMERIC(10, 2)                        NOT NULL,
  updated_at      TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL
);

INSERT INTO settings (id, btc_to_usd_rate, buy_csgokey_price, sell_csgokey_price) 
VALUES (1, 0, 0, 0) ON CONFLICT DO NOTHING;