BEGIN;

CREATE TABLE IF NOT EXISTS person (
    id SERIAL PRIMARY KEY,
    slug VARCHAR NOT NULL,
    first VARCHAR NOT NULL,
    last VARCHAR NOT NULL,
    description VARCHAR,
    birthdate DATE, 
    profile_img VARCHAR NOT NULL,
    search_vector tsvector,
    insert_timestamp timestamp DEFAULT now()
);

CREATE TABLE IF NOT EXISTS character (
    id SERIAL PRIMARY KEY, 
    slug VARCHAR NOT NULL,
    name VARCHAR NOT NULL, 
    description VARCHAR, 
    img_name VARCHAR,
    insert_timestamp timestamp DEFAULT now(),
    search_vector tsvector,
    person_id INT REFERENCES person(id)
);

CREATE TABLE IF NOT EXISTS creator (
    id serial primary key,
    name VARCHAR NOT NULL,
    slug VARCHAR NOT NULL,
    page_url VARCHAR NOT NULL,
    description VARCHAR,
    profile_img VARCHAR, 
    date_established DATE,
    search_vector tsvector,
    insert_timestamp timestamp DEFAULT now()
);


CREATE TABLE IF NOT EXISTS video (
    id SERIAL PRIMARY KEY,
    title VARCHAR NOT NULL,
    video_url TEXT,
    youtube_id VARCHAR, 
    slug TEXT NOT NULL,
    thumbnail_name TEXT,
    description TEXT,
    upload_date DATE,
    pg_rating rating,
    search_vector tsvector,
    insert_timestamp timestamp DEFAULT now()
);

CREATE TABLE IF NOT EXISTS video_person_rel (
    id SERIAL PRIMARY KEY,
    video_id INT references video(id) NOT NULL,
    person_id INT references person(id) NOT NULL,
    character_id INT references character(id),
    position INT,
    img VARCHAR,
    insert_timestamp timestamp DEFAULT now(),
    CONSTRAINT unique_video_person_character UNIQUE(video_id, person_id, character_id)
);

CREATE TABLE IF NOT EXISTS video_creator_rel (
    creator_id INT references creator(id),
    video_id INT references video(id),
    position INT,
    insert_timestamp timestamp DEFAULT now(),
    PRIMARY KEY (creator_id, video_id)
);

CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP(0) with time zone NOT NULL DEFAULT NOW(),
    username VARCHAR(20) UNIQUE NOT NULL,
    email CITEXT UNIQUE NOT NULL,
    password_hash BYTEA NOT NULL,
    activated BOOL NOT NULL,
    role TEXT NOT NULL DEFAULT 'viewer' CHECK (role IN ('admin', 'editor', 'viewer'))
);

CREATE TABLE IF NOT EXISTS likes (
    created_at TIMESTAMP(0) with time zone NOT NULL DEFAULT NOW(),
    user_id INT references users(id) NOT NULL,
    video_id INT references video(id) NOT NULL,
    PRIMARY KEY (user_id, video_id)
);

--CREATE TABLE IF NOT EXISTS tag (
--    id SERIAL PRIMARY KEY,
--    tag VARCHAR UNIQUE
--);

-- CREATE TABLE video_tag_rel (
--     id serial PRIMARY KEY,
--     video_id INT references video(id),
--     tag_id INT references tag(id),
--     UNIQUE (video_id, tag_id)
-- );

CREATE TABLE IF NOT EXISTS sessions (
	token TEXT PRIMARY KEY,
	data BYTEA NOT NULL,
	expiry TIMESTAMPTZ NOT NULL
);

CREATE INDEX IF NOT EXISTS sessions_expiry_idx ON sessions (expiry);

COMMIT;
