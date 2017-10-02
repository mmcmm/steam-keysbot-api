CREATE TABLE stats
(
  id                SERIAL                                 NOT NULL,
    CONSTRAINT stats_pkey
    PRIMARY KEY,
  users             INTEGER                       NOT NULL,                    
  daily_trades      INTEGER                       NOT NULL,                    
  money_transacted  NUMERIC(10, 2)                NOT NULL,                    
  current_value     NUMERIC(10, 2)                NOT NULL,        
  users_stored      NUMERIC(10, 2)                NOT NULL,       
  updated_at        TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL
);