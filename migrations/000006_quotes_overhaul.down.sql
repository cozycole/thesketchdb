-- restore transcription columns
ALTER TABLE transcription_lines
ADD COLUMN start_time interval,
ADD COLUMN end_time interval;

-- drop ms columns
ALTER TABLE transcription_lines
DROP COLUMN start_ms,
DROP COLUMN end_ms;

-- recreate moment table
CREATE TABLE moment (
  id SERIAL PRIMARY KEY,
  sketch_id INT references sketch(id) NOT NULL,
  timestamp INT,
  insert_timestamp timestamp DEFAULT now()
);

-- restore quote columns
ALTER TABLE quote
ADD COLUMN moment_id INT REFERENCES moment(id),
ADD COLUMN position INT,
ADD COLUMN cast_id INT;

-- drop relation table
DROP TABLE quote_cast_rel;

-- drop new columns
ALTER TABLE quote
DROP COLUMN start_time_ms,
DROP COLUMN end_time_ms;
