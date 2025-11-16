ALTER TABLE person ADD COLUMN IF NOT EXISTS wiki_page TEXT;
ALTER TABLE person ADD COLUMN IF NOT EXISTS imdb_id TEXT;
ALTER TABLE person ADD COLUMN IF NOT EXISTS tmdb_id TEXT;

ALTER TABLE creator ADD COLUMN IF NOT EXISTS aliases TEXT;

ALTER TABLE cast_members alter person_id DROP NOT NULL;
ALTER TABLE cast_members DROP CONSTRAINT IF EXISTS unique_cast_character;

CREATE OR REPLACE TRIGGER person_search_update 
BEFORE INSERT OR UPDATE ON person
FOR EACH ROW EXECUTE FUNCTION tsvector_update_trigger(
  search_vector, 
  'pg_catalog.english', 
  first, last, description, aliases
);

CREATE OR REPLACE TRIGGER character_search_update 
BEFORE INSERT OR UPDATE ON character
FOR EACH ROW EXECUTE FUNCTION tsvector_update_trigger(
  search_vector, 
  'pg_catalog.english', 
  name, description, aliases
);

CREATE OR REPLACE TRIGGER creator_search_update 
BEFORE INSERT OR UPDATE ON creator
FOR EACH ROW EXECUTE FUNCTION tsvector_update_trigger(
  search_vector, 
  'pg_catalog.english', 
  name, description, aliases
);

CREATE OR REPLACE TRIGGER show_search_update 
BEFORE INSERT OR UPDATE ON show
FOR EACH ROW EXECUTE FUNCTION tsvector_update_trigger(
  search_vector, 
  'pg_catalog.english', 
  name, aliases
);
