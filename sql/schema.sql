BEGIN;

CREATE TABLE IF NOT EXISTS person (
    id SERIAL PRIMARY KEY,
    slug TEXT NOT NULL,
    first TEXT NOT NULL,
    last TEXT NOT NULL,
    description TEXT,
    birthdate DATE, 
    profile_img TEXT,
    search_vector tsvector,
    insert_timestamp TIMESTAMP DEFAULT now()
);

CREATE TABLE IF NOT EXISTS character (
    id SERIAL PRIMARY KEY, 
    slug TEXT NOT NULL,
    name TEXT NOT NULL, 
    description TEXT, 
    img_name TEXT,
    insert_timestamp TIMESTAMP DEFAULT now(),
    search_vector tsvector,
    person_id INT REFERENCES person(id)
);

CREATE TABLE IF NOT EXISTS creator (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    slug TEXT NOT NULL,
    page_url TEXT NOT NULL,
    description TEXT,
    profile_img TEXT, 
    date_established DATE,
    search_vector tsvector,
    insert_timestamp TIMESTAMP DEFAULT now()
);

CREATE TABLE IF NOT EXISTS show (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    profile_img TEXT,
    slug TEXT UNIQUE NOT NULL
);

-- CREATE TABLE IF NOT EXISTS show_creator (
--     id SERIAL PRIMARY KEY,
--     show_id INTEGER REFERENCES show(id),
--     creator_id INTEGER REFERENCES creator(id),
--     person_id INTEGER REFERENCES person(id),
--     CONSTRAINT unique_show_creator UNIQUE(show_id, creator_id),
--     CONSTRAINT unique_show_person UNIQUE(show_id, person_id),
--     CONSTRAINT at_least_one_creator CHECK(
-- 	creator_id IS NOT NULL OR person_id IS NOT NULL
--     )
-- );

-- CREATE TABLE IF NOT EXISTS network (
--     id SERIAL PRIMARY KEY,
--     name TEXT NOT NULL,
--     profile_img TEXT
-- );

-- CREATE TABLE IF NOT EXISTS show_network (
--     show_id INTEGER REFERENCES show(id),
--     network_id INTEGER REFERENCES network(id),
--     PRIMARY KEY(show_id, network_id)
-- );

CREATE TABLE IF NOT EXISTS season (
    id SERIAL PRIMARY KEY,
    show_id INTEGER REFERENCES show(id),
    season_number INTEGER NOT NULL,
    CONSTRAINT unique_show_season UNIQUE(show_id, season_number)
);

CREATE TABLE IF NOT EXISTS episode (
    id SERIAL PRIMARY KEY,
    season_id INTEGER REFERENCES season(id),
    title TEXT,
    episode_number INTEGER NOT NULL,
    thumbnail_name TEXT UNIQUE,
    air_date DATE,
    CONSTRAINT unique_season_episode UNIQUE(season_id, episode_number)
);

CREATE TABLE IF NOT EXISTS video (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    video_url TEXT,
    youtube_id TEXT, 
    slug TEXT NOT NULL,
    thumbnail_name TEXT,
    description TEXT,
    upload_date DATE,
    pg_rating rating,
    episode_id int REFERENCES episode(id),
    part_number int,
    sketch_number int,
    search_vector tsvector,
    insert_timestamp timestamp DEFAULT now()
);

CREATE TABLE IF NOT EXISTS cast_members (
    id SERIAL PRIMARY KEY,
    video_id INT references video(id) NOT NULL,
    person_id INT references person(id) NOT NULL,
    character_name text DEFAULT '',
    character_id INT references character(id),
    position INT,
    img_name TEXT,
    role TEXT CHECK (role IN ('host', 'cast', 'guest')),
    insert_timestamp timestamp DEFAULT now(),
    CONSTRAINT unique_cast_character UNIQUE(video_id, person_id, character_id)
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

CREATE TABLE IF NOT EXISTS categories (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP(0) with time zone NOT NULL DEFAULT NOW(),
    name TEXT NOT NULL,
    slug TEXT UNIQUE NOT NULL,
    parent_id INT REFERENCES categories(id) ON DELETE SET NULL
);

CREATE TABLE IF NOT EXISTS tags (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP(0) with time zone NOT NULL DEFAULT NOW(),
    name TEXT NOT NULL,
    slug TEXT NOT NULL,
    category_id INT REFERENCES categories(id) ON DELETE SET NULL
);

CREATE TABLE video_tags (
    video_id INT references video(id),
    tag_id INT references tags(id),
    PRIMARY KEY (video_id, tag_id)
);

CREATE TABLE IF NOT EXISTS sessions (
	token TEXT PRIMARY KEY,
	data BYTEA NOT NULL,
	expiry TIMESTAMPTZ NOT NULL
);

CREATE INDEX IF NOT EXISTS sessions_expiry_idx ON sessions (expiry);

COMMIT;
