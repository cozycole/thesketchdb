BEGIN;
CREATE TABLE IF NOT EXISTS person (
    id SERIAL PRIMARY KEY,
    slug VARCHAR NOT NULL,
    first VARCHAR NOT NULL,
    last VARCHAR NOT NULL,
    description VARCHAR,
    birthdate DATE, 
    profile_img VARCHAR NOT NULL
);
COMMIT;
