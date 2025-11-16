CREATE TABLE IF NOT EXISTS pipeline_jobs (
    id SERIAL PRIMARY KEY,
    sketch_id INT references sketch(id) NOT NULL,
    video_filename TEXT,
    status TEXT DEFAULT 'pending',
    error TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    edited_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS cast_auto_screenshots (
    id SERIAL PRIMARY KEY,
    sketch_id INT references sketch(id) NOT NULL,
    cluster_number INT,
    image_number INT, -- img number within a cluster
    thumbnail_img TEXT,
    profile_img TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS transcription_lines (
    id SERIAL PRIMARY KEY,
    sketch_id INT references sketch(id) NOT NULL,
    line_number INT,
    start_time INTERVAL,
    end_time INTERVAL,
    text TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE FUNCTION set_edited_at()
RETURNS trigger AS $$
BEGIN
  NEW.edited_at := NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_edited_at
BEFORE UPDATE ON pipeline_jobs
FOR EACH ROW
EXECUTE FUNCTION set_edited_at();
