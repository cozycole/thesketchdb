BEGIN;

CREATE TABLE IF NOT EXISTS video (
    id SERIAL PRIMARY KEY,
    title VARCHAR NOT NULL,
    video_url VARCHAR NOT NULL,
    slug VARCHAR NOT NULL,
    thumbnail_name VARCHAR,
    description VARCHAR,
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
    CONSTRAINT unique_video_person_character UNIQUE(video_id, person_id, character_id)
);

CREATE TABLE IF NOT EXISTS video_creator_rel (
    creator_id int references creator(id),
    video_id int references video(id),
    position int,
    PRIMARY KEY (creator_id, video_id)
);

-- CREATE TABLE video_tag_rel (
--     id serial primary key,
--     video_id int references video(id),
--     tag_id int references tag(id),
--     UNIQUE (video_id, tag_id)
-- );
COMMIT;
