CREATE TABLE IF NOT EXISTS cast_tags_rel (
    cast_id INT NOT NULL,
    tag_id   INT NOT NULL,
    PRIMARY KEY (cast_id, tag_id),
    FOREIGN KEY (cast_id) REFERENCES cast_members(id) ON DELETE CASCADE,
    FOREIGN KEY (tag_id)   REFERENCES tags(id)   ON DELETE CASCADE
);
