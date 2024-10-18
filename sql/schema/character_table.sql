CREATE TABLE IF NOT EXISTS character (
    id SERIAL PRIMARY KEY, 
    first_name VARCHAR NOT NULL, 
    middle_initial VARCHAR, 
    last_name VARCHAR, 
    description VARCHAR, 
    img_name VARCHAR,
    person_id INT REFERENCES person(id)
);