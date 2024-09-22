CREATE DATABASE sketch_data_test
    WITH
    OWNER = colet
    ENCODING = 'UTF8'
    CONNECTION LIMIT = -1
    IS_TEMPLATE = False;


CREATE TYPE IF NOT EXISTS rating AS ENUM ('PG', 'PG-13', 'R');