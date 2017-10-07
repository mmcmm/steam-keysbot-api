CREATE TABLE settings
(
  id              SERIAL                                 NOT NULL
    CONSTRAINT settings_pkey
    PRIMARY KEY,
  btc_to_usd_price NUMERIC(10, 2)                         NOT NULL,
  updated_at      TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL
);

INSERT INTO settings (id, btc_to_usd_price) 
VALUES (1, 0) ON CONFLICT DO NOTHING;