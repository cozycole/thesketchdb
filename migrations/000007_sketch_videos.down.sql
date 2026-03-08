ALTER TABLE pipeline_jobs DROP COLUMN IF EXISTS video_id;
ALTER TABLE pipeline_jobs ADD COLUMN video_filename TEXT;
ALTER TABLE pipeline_jobs ADD COLUMN sketch_id INT REFERENCES sketch(id);

DROP TABLE IF EXISTS sketch_videos ;
