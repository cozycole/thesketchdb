-- This needs to be updated later to not execute on all rows
CREATE OR REPLACE TRIGGER character_search_update 
BEFORE INSERT OR UPDATE ON character
FOR EACH ROW EXECUTE FUNCTION tsvector_update_trigger(
  search_vector, 
  'pg_catalog.english', 
  name, description
);

CREATE OR REPLACE TRIGGER creator_search_update 
BEFORE INSERT OR UPDATE ON creator
FOR EACH ROW EXECUTE FUNCTION tsvector_update_trigger(
  search_vector, 
  'pg_catalog.english', 
  name, description
);

CREATE OR REPLACE TRIGGER person_search_update 
BEFORE INSERT OR UPDATE ON person
FOR EACH ROW EXECUTE FUNCTION tsvector_update_trigger(
  search_vector, 
  'pg_catalog.english', 
  first, last, description
);

CREATE OR REPLACE TRIGGER sketch_search_update 
BEFORE INSERT OR UPDATE ON sketch 
FOR EACH ROW EXECUTE FUNCTION tsvector_update_trigger(
  search_vector, 
  'pg_catalog.english', 
  title, description
);

