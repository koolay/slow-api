CREATE TABLE public.slow_sql (
  created_on TIMESTAMP WITHOUT TIME ZONE,
  "user" CHARACTER VARYING(200),
  host CHARACTER VARYING(100),
  query_time DOUBLE PRECISION,
  lock_time DOUBLE PRECISION,
  rows_sent INTEGER,
  rows_examined INTEGER,
  sql TEXT,
  id INTEGER PRIMARY KEY NOT NULL DEFAULT nextval('slow_sql_id_seq'::regclass)
);
CREATE UNIQUE INDEX slow_sql_id_uindex ON slow_sql USING BTREE (id);
