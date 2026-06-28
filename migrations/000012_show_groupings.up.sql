CREATE TABLE sketch_grouping (
    id SERIAL PRIMARY KEY,
    slug TEXT NOT NULL,
    title TEXT,
    description TEXT,
    position INT,
    show_id INT references show(id),
    creator_id INT references creator(id)
);

ALTER TABLE sketch ADD COLUMN grouping_id INT references sketch_grouping(id);
