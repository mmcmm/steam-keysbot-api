CREATE TABLE errors (
    id              SERIAL                                 NOT NULL
    CONSTRAINT errors_pkey
    PRIMARY KEY,
    source          VARCHAR(255)                           NOT NULL,
    message         VARCHAR(255)                           NOT NULL,
    created_at      TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL
);


CREATE OR REPLACE FUNCTION errors_delete_old()
  RETURNS TRIGGER AS
$func$
BEGIN
  DELETE FROM errors WHERE created_at <= NOW() - interval '30 days';
  RETURN NEW;
END;
$func$
LANGUAGE plpgsql;

CREATE TRIGGER errors_delete_old
AFTER INSERT ON errors FOR EACH ROW
EXECUTE PROCEDURE errors_delete_old();

