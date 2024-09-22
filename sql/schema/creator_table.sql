CREATE TABLE IF NOT EXISTS creator (
    id serial primary key,
    name VARCHAR NOT NULL,
    profile_img_path VARCHAR, 
    date_established DATE
);