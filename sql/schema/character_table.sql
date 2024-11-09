CREATE TABLE IF NOT EXISTS character (
    id SERIAL PRIMARY KEY, 
    name VARCHAR NOT NULL, 
    description VARCHAR, 
    img_name VARCHAR,
    person_id INT REFERENCES person(id)
);
