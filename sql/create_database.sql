CREATE DATABASE sketch_data_test
    WITH
    OWNER = colet
    ENCODING = 'UTF8'
    CONNECTION LIMIT = -1
    IS_TEMPLATE = False;

CREATE EXTENSION IF NOT EXISTS citext;
CREATE TYPE IF NOT EXISTS vid_role AS ENUM ('Cast', 'Guest', 'Host');
