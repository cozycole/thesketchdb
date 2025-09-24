CREATE TABLE IF NOT EXISTS sketch_rating (
    sketch_id INT references sketch(id) NOT NULL,
    user_id INT references users(id) NOT NULL,
    rating INT NOT NULL,
    created_at TIMESTAMP(0) with time zone NOT NULL DEFAULT NOW(),
    PRIMARY KEY (sketch_id, user_id)
);

ALTER TABLE sketch ADD COLUMN IF NOT EXISTS rating REAL DEFAULT 0;
ALTER TABLE sketch ADD COLUMN IF NOT EXISTS total_ratings INT DEFAULT 0;

CREATE OR REPLACE FUNCTION update_sketch_rating()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' OR TG_OP = 'UPDATE' THEN
        UPDATE sketch 
        SET 
            rating = (
                SELECT AVG(rating::REAL) 
                FROM sketch_rating 
                WHERE sketch_id = NEW.sketch_id
            ),
            total_ratings = (
                SELECT COUNT(*) 
                FROM sketch_rating 
                WHERE sketch_id = NEW.sketch_id
            )
        WHERE id = NEW.sketch_id;
        RETURN NEW;
    END IF;
    
    IF TG_OP = 'DELETE' THEN
        UPDATE sketch 
        SET 
            rating = COALESCE((
                SELECT AVG(rating::REAL) 
                FROM sketch_rating 
                WHERE sketch_id = OLD.sketch_id
            ), 0),
            total_ratings = COALESCE((
                SELECT COUNT(*) 
                FROM sketch_rating 
                WHERE sketch_id = OLD.sketch_id
            ), 0)
        WHERE id = OLD.sketch_id;
        RETURN OLD;
    END IF;
    
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER sketch_rating_trigger
    AFTER INSERT OR UPDATE OR DELETE ON sketch_rating
    FOR EACH ROW
    EXECUTE FUNCTION update_sketch_rating();
