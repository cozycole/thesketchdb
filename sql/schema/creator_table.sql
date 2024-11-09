CREATE TABLE IF NOT EXISTS creator (
    id serial primary key,
    name VARCHAR NOT NULL,
    slug VARCHAR NOT NULL,
    page_url VARCHAR NOT NULL,
    description VARCHAR,
    profile_img VARCHAR, 
    date_established DATE
);
