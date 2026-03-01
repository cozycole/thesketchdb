-- update quote table
ALTER TABLE quote 
ADD COLUMN IF NOT EXISTS start_time_ms INT, 
ADD COLUMN IF NOT EXISTS end_time_ms INT,
ADD COLUMN IF NOT EXISTS sketch_id INT REFERENCES sketch(id);

UPDATE quote AS q
SET start_time_ms = m.timestamp * 1000
FROM moment AS m 
WHERE q.moment_id = m.id;

-- extract the sketch_id from the moment table before
-- removing the moment_id 
UPDATE quote 
SET sketch_id=m.sketch_id 
FROM moment AS m 
WHERE moment_id = m.id;

ALTER TABLE quote DROP COLUMN IF EXISTS moment_id;
ALTER TABLE quote DROP COLUMN IF EXISTS position;

CREATE TABLE IF NOT EXISTS quote_cast_rel (
	quote_id INT REFERENCES quote(id) ON DELETE CASCADE,
	cast_id INT REFERENCES cast_members(id),
	PRIMARY KEY (quote_id, cast_id)
);

INSERT INTO quote_cast_rel (quote_id, cast_id)
SELECT id, cast_id
FROM quote
WHERE cast_id IS NOT NULL
ON CONFLICT DO NOTHING;

ALTER TABLE quote DROP COLUMN IF EXISTS cast_id;

-- alter quote_tag table to allow quote deletes
ALTER TABLE quote_tag_rel
  DROP CONSTRAINT IF EXISTS quote_tag_rel_quote_id_fkey,
  ADD CONSTRAINT quote_tag_rel_quote_id_fkey
    FOREIGN KEY (quote_id) REFERENCES quote(id) ON DELETE CASCADE;

-- delete moments table
DROP TABLE moment;

-- update transcript table
ALTER TABLE transcription_lines 
ADD COLUMN IF NOT EXISTS start_ms INT,
ADD COLUMN IF NOT EXISTS end_ms INT;

UPDATE transcription_lines
SET start_ms = (EXTRACT(EPOCH FROM start_time) * 1000)::integer,
    end_ms   = (EXTRACT(EPOCH FROM end_time)   * 1000)::integer;

ALTER TABLE transcription_lines 
DROP COLUMN IF EXISTS start_time, 
DROP COLUMN IF EXISTS end_time;
