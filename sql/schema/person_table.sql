BEGIN;
CREATE TABLE IF NOT EXISTS person (
    id serial primary key,
    first VARCHAR NOT NULL,
    last VARCHAR NOT NULL,
    description VARCHAR,
    birthdate DATE, 
    profile_img VARCHAR NOT NULL
);
COMMIT;
