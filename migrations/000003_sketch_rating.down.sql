DROP TRIGGER IF EXISTS sketch_rating_trigger;
DROP FUNCTION IF EXISTS update_sketch_rating;

ALTER TABLE sketch DROP COLUMN IF EXISTS rating;
ALTER TABLE sketch DROP COLUMN IF EXISTS total_ratings;

DROP TABLE IF EXISTS sketch_rating;
