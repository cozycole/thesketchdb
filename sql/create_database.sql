CREATE DATABASE sketch_data_test
    WITH
    OWNER = colet
    ENCODING = 'UTF8'
    CONNECTION LIMIT = -1
    IS_TEMPLATE = False;

CREATE EXTENSION IF NOT EXISTS citext;
CREATE TYPE cast_role AS ENUM ('cast', 'guest', 'host', 'minor');
CREATE TYPE character_type AS ENUM ('original', 'impression', 'generic');
