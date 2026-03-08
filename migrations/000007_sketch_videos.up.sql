CREATE TABLE IF NOT EXISTS sketch_video (
    id SERIAL PRIMARY KEY,
    sketch_id INT REFERENCES sketch(id),
    hot_s3_key TEXT,
    cold_s3_key TEXT,
    archived_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT now()
);

ALTER TABLE pipeline_jobs DROP COLUMN IF EXISTS video_filename;
ALTER TABLE pipeline_jobs DROP COLUMN IF EXISTS sketch_id;
ALTER TABLE pipeline_jobs ADD COLUMN video_id INT REFERENCES sketch_video(id) ON DELETE CASCADE;
