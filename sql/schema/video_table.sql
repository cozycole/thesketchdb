BEGIN;
CREATE TYPE rating AS ENUM ('PG', 'PG-13', 'R');

CREATE TABLE video (
    id serial primary key,
    title VARCHAR NOT NULL,
    video_url VARCHAR NOT NULL,
    thumbnail_name VARCHAR,
    upload_date DATE,
    pg_rating rating,
    insert_timestamp timestamp DEFAULT now()
);

CREATE TABLE video_actor_rel (
    id serial primary key,
    actor_id int references actor(id),
    video_id int references video(id),
    UNIQUE (actor_id, video_id)
);

CREATE TABLE video_creator_rel (
    id serial primary key,
    creator_id int references creator(id),
    video_id int references video(id)
    UNIQUE (creator_id, video_id)
);

CREATE TABLE video_tag_rel (
    id serial primary key,
    video_id int references video(id),
    tag_id int references tag(id),
    UNIQUE (actor_id, video_id)
);
COMMIT;