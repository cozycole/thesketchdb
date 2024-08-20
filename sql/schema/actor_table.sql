BEGIN;
CREATE TABLE actor (
    id serial primary key,
    givenname VARCHAR NOT NULL,
    surname VARCHAR NOT NULL,
    birthdate DATE NOT NULL, 
    profile_img VARCHAR NOT NULL
);

CREATE TABLE actor_creator_rel (
    id serial primary key,
    actor_id int references actor(id),
    creator_id int references creator(id),
    UNIQUE (actor_id, creator_id)
);
COMMIT;