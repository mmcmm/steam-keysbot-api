CREATE TABLE settings
(
  id              SERIAL                                 NOT NULL
    CONSTRAINT settings_pkey
    PRIMARY KEY,
  btc_to_usd_rate NUMERIC(10, 2)                         NOT NULL,
  updated_at      TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL
);

INSERT INTO stats (id, btc_to_usd_rate) 
VALUES (1, 0) ON CONFLICT DO NOTHING;