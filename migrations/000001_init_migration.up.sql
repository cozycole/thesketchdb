CREATE EXTENSION IF NOT EXISTS citext;
CREATE TYPE cast_role AS ENUM ('cast', 'guest', 'host', '');
CREATE TYPE character_type AS ENUM ('original', 'impression', 'fictional_impression', 'generic');

CREATE TABLE IF NOT EXISTS person (
    id SERIAL PRIMARY KEY,
    slug TEXT NOT NULL,
    first TEXT NOT NULL,
    last TEXT NOT NULL,
    aliases TEXT,
    professions TEXT NOT NULL,
    description TEXT,
    birthdate DATE, 
    profile_img TEXT,
    search_vector tsvector,
    popularity_score REAL DEFAULT 0,
    insert_timestamp TIMESTAMP DEFAULT now()
);

CREATE TABLE IF NOT EXISTS character (
    id SERIAL PRIMARY KEY, 
    slug TEXT NOT NULL,
    name TEXT NOT NULL, 
    aliases TEXT,
    character_type character_type NOT NULL,
    description TEXT, 
    img_name TEXT,
    search_vector tsvector,
    person_id INT REFERENCES person(id),
    popularity_score REAL DEFAULT 0,
    insert_timestamp TIMESTAMP DEFAULT now()
);

CREATE TABLE IF NOT EXISTS creator (
    id SERIAL PRIMARY KEY,
    slug TEXT NOT NULL,
    name TEXT NOT NULL,
    aliases TEXT,
    page_url TEXT NOT NULL,
    description TEXT,
    profile_img TEXT, 
    date_established DATE,
    search_vector tsvector,
    popularity_score REAL DEFAULT 0,
    insert_timestamp TIMESTAMP DEFAULT now()
);

CREATE TABLE IF NOT EXISTS show (
    id SERIAL PRIMARY KEY,
    slug TEXT NOT NULL,
    name TEXT NOT NULL,
    aliases TEXT,
    profile_img TEXT,
    popularity_score REAL DEFAULT 0,
    search_vector tsvector,
    insert_timestamp TIMESTAMP DEFAULT now()
);

CREATE TABLE IF NOT EXISTS series (
    id SERIAL PRIMARY KEY,
    slug TEXT NOT NULL,
    title TEXT NOT NULL,
    description TEXT,
    thumbnail_name TEXT,
    insert_timestamp TIMESTAMP DEFAULT now()
);

CREATE TABLE IF NOT EXISTS recurring (
    id SERIAL PRIMARY KEY,
    slug TEXT NOT NULL,
    title TEXT NOT NULL,
    description TEXT,
    thumbnail_name TEXT,
    insert_timestamp TIMESTAMP DEFAULT now()
);

CREATE TABLE IF NOT EXISTS season (
    id SERIAL PRIMARY KEY,
    slug TEXT,
    show_id INTEGER REFERENCES show(id),
    season_number INTEGER NOT NULL,
    air_date DATE,
    CONSTRAINT unique_show_season UNIQUE(show_id, season_number)
);

CREATE TABLE IF NOT EXISTS episode (
    id SERIAL PRIMARY KEY,
    slug TEXT,
    season_id INTEGER REFERENCES season(id),
    title TEXT,
    episode_number INTEGER NOT NULL,
    url VARCHAR,
    youtube_id VARCHAR,
    thumbnail_name TEXT,
    air_date DATE,
    CONSTRAINT unique_season_episode UNIQUE(season_id, episode_number)
);

CREATE TABLE IF NOT EXISTS sketch (
    id SERIAL PRIMARY KEY,
    slug TEXT NOT NULL,
    title TEXT NOT NULL,
    sketch_url TEXT,
    youtube_id TEXT, 
    thumbnail_name TEXT,
    description TEXT,
    transcript TEXT, 
    diarization TEXT,
    upload_date DATE,
    duration INT,
    episode_id INT REFERENCES episode(id),
    episode_start INT, 
    sketch_number INT,
    series_id INT REFERENCES series(id), 
    part_number INT,
    recurring_id INT REFERENCES recurring(id),
    popularity_score REAL DEFAULT 0,
    search_vector tsvector,
    insert_timestamp timestamp DEFAULT now()
);

CREATE TABLE IF NOT EXISTS cast_members (
    id SERIAL PRIMARY KEY,
    sketch_id INT references sketch(id) NOT NULL,
    person_id INT references person(id),
    character_name text DEFAULT '',
    character_id INT references character(id),
    position INT,
    thumbnail_name TEXT,
    profile_img TEXT,
    role cast_role,
    minor bool,
    insert_timestamp timestamp DEFAULT now(),
    CONSTRAINT unique_cast_character UNIQUE(sketch_id, person_id, character_id)
);

CREATE TABLE IF NOT EXISTS sketch_creator_rel (
    creator_id INT references creator(id),
    sketch_id INT references sketch(id),
    position INT,
    insert_timestamp timestamp DEFAULT now(),
    PRIMARY KEY (creator_id, sketch_id)
);

CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP(0) with time zone NOT NULL DEFAULT NOW(),
    username VARCHAR(20) UNIQUE NOT NULL,
    email CITEXT UNIQUE NOT NULL,
    password_hash BYTEA NOT NULL,
    activated BOOL NOT NULL,
    role TEXT NOT NULL DEFAULT 'viewer' CHECK (role IN ('admin', 'editor', 'viewer')),
    profile_image TEXT DEFAULT 'missing-profile.jpg'
);

CREATE TABLE IF NOT EXISTS moment (
    id SERIAL PRIMARY KEY,
    sketch_id INT references sketch(id) NOT NULL,
    timestamp INT NOT NULL,
    insert_timestamp timestamp DEFAULT now()
);

CREATE TABLE IF NOT EXISTS quote (
    id SERIAL PRIMARY KEY,
    moment_id INT references moment(id) NOT NULL,
    cast_id INT references cast_members(id) NOT NULL,
    text TEXT,
    type TEXT,
    funny TEXT,
    position INT
);

CREATE TABLE IF NOT EXISTS categories (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP(0) with time zone NOT NULL DEFAULT NOW(),
    name TEXT NOT NULL,
    slug TEXT NOT NULL,
    parent_id INT REFERENCES categories(id) ON DELETE SET NULL
);

CREATE TABLE IF NOT EXISTS tags (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP(0) with time zone NOT NULL DEFAULT NOW(),
    name TEXT NOT NULL,
    slug TEXT NOT NULL,
    type TEXT,
    category_id INT REFERENCES categories(id) ON DELETE SET NULL
);

CREATE TABLE IF NOT EXISTS quote_tags_rel (
    quote_id INT NOT NULL,
    tag_id   INT NOT NULL,
    PRIMARY KEY (quote_id, tag_id),
    FOREIGN KEY (quote_id) REFERENCES quote(id) ON DELETE CASCADE,
    FOREIGN KEY (tag_id)   REFERENCES tags(id)   ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS likes (
    created_at TIMESTAMP(0) with time zone NOT NULL DEFAULT NOW(),
    user_id INT references users(id) NOT NULL,
    sketch_id INT references sketch(id) NOT NULL,
    PRIMARY KEY (user_id, sketch_id)
);

CREATE TABLE IF NOT EXISTS sketch_tags (
    sketch_id INT references sketch(id),
    tag_id INT references tags(id),
    PRIMARY KEY (sketch_id, tag_id)
);

CREATE TABLE IF NOT EXISTS sessions (
	token TEXT PRIMARY KEY,
	data BYTEA NOT NULL,
	expiry TIMESTAMPTZ NOT NULL
);

CREATE INDEX IF NOT EXISTS sessions_expiry_idx ON sessions (expiry);

CREATE OR REPLACE TRIGGER character_search_update 
BEFORE INSERT OR UPDATE ON character
FOR EACH ROW EXECUTE FUNCTION tsvector_update_trigger(
  search_vector, 
  'pg_catalog.english', 
  name, description, alias
);

CREATE OR REPLACE TRIGGER show_search_update 
BEFORE INSERT OR UPDATE ON show
FOR EACH ROW EXECUTE FUNCTION tsvector_update_trigger(
  search_vector, 
  'pg_catalog.english', 
  name, alias
);

CREATE OR REPLACE TRIGGER creator_search_update 
BEFORE INSERT OR UPDATE ON creator
FOR EACH ROW EXECUTE FUNCTION tsvector_update_trigger(
  search_vector, 
  'pg_catalog.english', 
  name, description, alias
);

CREATE OR REPLACE TRIGGER person_search_update 
BEFORE INSERT OR UPDATE ON person
FOR EACH ROW EXECUTE FUNCTION tsvector_update_trigger(
  search_vector, 
  'pg_catalog.english', 
  first, last, description, alias
);

CREATE OR REPLACE TRIGGER sketch_search_update 
BEFORE INSERT OR UPDATE ON sketch 
FOR EACH ROW EXECUTE FUNCTION tsvector_update_trigger(
  search_vector, 
  'pg_catalog.english', 
  title, description
);
