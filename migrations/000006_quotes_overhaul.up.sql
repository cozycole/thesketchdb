-- update quote table
ALTER TABLE quote 
ADD COLUMN start_time_ms INT, 
ADD COLUMN end_time_ms INT;

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

ALTER TABLE quote DROP COLUMN moment_id;
ALTER TABLE quote DROP COLUMN postion;

CREATE TABLE IF NOT EXISTS quote_cast_rel (
	quote_id INT REFERENCES quote(id),
	cast_id INT REFERENCES cast_members(id),
	PRIMARY KEY (quote_id, cast_id)
);

INSERT INTO quote_cast_rel (quote_id, cast_id)
SELECT id, cast_id
FROM quote
WHERE cast_id IS NOT NULL
ON CONFLICT DO NOTHING;

ALTER TABLE quote DROP COLUMN cast_id;

-- delete moments table
DROP TABLE moment;

-- update transcript table
ALTER TABLE transcription_lines 
ADD COLUMN start_ms INT,
ADD COLUMN end_ms INT;

UPDATE transcription_lines
SET start_ms = (EXTRACT(EPOCH FROM start_time) * 1000)::integer,
    end_ms   = (EXTRACT(EPOCH FROM end_time)   * 1000)::integer;

ALTER TABLE transcription_lines 
DROP COLUMN start_time, 
DROP COLUMN end_time;
